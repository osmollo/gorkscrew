package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

const bufsize = 4096

type Progress struct {
	bytes uint64
}

// this function redirects data between socket, stdin and stdout
func FeelTheMagic(con net.Conn) {
	c := make(chan Progress)

	// Read from Reader and write to Writer until EOF
	copy := func(r io.ReadCloser, w io.WriteCloser) {
		defer func() {
			r.Close()
			w.Close()
		}()
		n, _ := io.Copy(w, r)
		c <- Progress{bytes: uint64(n)}
	}

	go copy(con, os.Stdout)
	go copy(os.Stdin, con)

	p := <-c
	log.Printf("[%s]: Connection has been closed by remote peer, %d bytes has been received\n", con.RemoteAddr(), p.bytes)
	p = <-c
	log.Printf("[%s]: Local peer has been stopped, %d bytes has been sent\n", con.RemoteAddr(), p.bytes)
}

// Returns the URI for proxy connection for kerberos authenticacion
func GetURIKerberosAuth(krb5conf *string, spn *string, desthost *string, destport *int) string {
	var (
		ccpath          string = GetCredentialsCachePath(*krb5conf)
		repository      string
		b64tgs          string
		parentComand    string
		gitShellCommand string
		cl              *client.Client
	)

	if !FileExists(*krb5conf) {
		log.Printf("File %s doesn't exist\n", *krb5conf)
		os.Exit(1)
	}
	cfg, err := config.Load(*krb5conf)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	if FileExists(ccpath) {
		ccache, _ := credentials.LoadCCache(ccpath)
		cl, _ = client.NewFromCCache(ccache, cfg)
	} else {
		log.Println("User must be kerberized")
		fmt.Println("User must be kerberized")
		os.Exit(1)
	}

	cl.Login()

	s := spnego.SPNEGOClient(cl, *spn)
	err = s.AcquireCred()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
	st, _ := s.InitSecContext()
	nb, _ := st.Marshal()
	b64tgs = b64.StdEncoding.EncodeToString(nb)

	log.Println("Authenticated User:", cl.Credentials.UserName())
	parentComand = GetParentCommand()
	log.Print("Parent command:", parentComand)
	repository, gitShellCommand = GetRepositoryFromCommand(parentComand)
	log.Print("Repository:", repository)
	log.Print("Git Shell Command: ", gitShellCommand)
	if repository == "" {
		log.Println("Error: repository not found")
		os.Exit(1)
	}
	return "CONNECT " + *desthost + ":" + strconv.Itoa(*destport) + " HTTP/1.0\nHost: " + *desthost + ":" + strconv.Itoa(*destport) + "\nProxy-Authorization: Negotiate " + b64tgs + "\nRepository: " + repository + "\nGitShellCommand: " + gitShellCommand + "\nUserAuth: " + cl.Credentials.UserName() + "\r\n\r\n"
}

// returns the URI for proxy connections with basic authentication
func GetURIBasicAuth(credsfilename *string, desthost *string, destport *int) string {
	var (
		repository      string
		gitShellCommand string
		proxycreds      string
		b64creds        string
		parentComand    string
	)
	gvalue, gpresent := os.LookupEnv("GORKSCREW_AUTH")
	cvalue, cpresent := os.LookupEnv("CORKSCREW_AUTH")

	if *credsfilename != "" && FileExists(*credsfilename) {
		content, err := ioutil.ReadFile(*credsfilename)
		if err != nil {
			log.Printf("Error reading proxy credentials file '%s'\n", *credsfilename)
			os.Exit(5)
		}
		proxycreds = string(content)
	} else if gpresent {
		proxycreds = gvalue
	} else if cpresent {
		proxycreds = cvalue
	} else {
		log.Println("Proxy credentials not found")
		os.Exit(5)
	}

	log.Println("Authenticated User:", strings.Split(proxycreds, ":")[0])
	b64creds = b64.StdEncoding.EncodeToString([]byte(proxycreds))
	parentComand = GetParentCommand()
	log.Print("Parent command:", parentComand)
	repository, gitShellCommand = GetRepositoryFromCommand(parentComand)
	log.Print("Repository:", repository)
	log.Print("Git Shell Command: ", gitShellCommand)
	return "CONNECT " + *desthost + ":" + strconv.Itoa(*destport) + " HTTP/1.0\nHost: " + *desthost + ":" + strconv.Itoa(*destport) + "\nProxy-Authorization: Basic " + b64creds + "\nRepository: " + repository + "\nGitShellCommand: " + gitShellCommand + "\nUserAuth: " + strings.Split(proxycreds, ":")[0] + "\r\n\r\n"
}

// returns URI for proxy connection when there's no authentication

// returns the path to kerberos credentials cache if its set in krb5.conf or default value
func GetCredentialsCachePath(krb5conf string) string {
	data, _ := ioutil.ReadFile(krb5conf)
	file := string(data)
	line := 0
	temp := strings.Split(file, "\n")
	u, _ := user.Current()

	for _, item := range temp {
		if strings.Contains(item, "default_ccache_name") {
			return strings.Replace(strings.Split(item, " ")[len(strings.Split(item, " "))-1], "%{uid}", u.Uid, -1)
		}
		line++
	}

	return "/tmp/krb5cc_" + u.Uid
}

// Returns the command that executes the current program
func GetParentCommand() string {
	ppid := os.Getppid()
	fileName := "/proc/" + strconv.Itoa(ppid) + "/cmdline"
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Can't read file '%s'\n", fileName)
		return ""
	}
	output := ""
	for i, elem := range dat {
		if elem == byte(0) {
			output += " "
		} else {
			output += string(dat[i])
		}
	}
	return output
}

// Parse COMMAND and returns in format 'ssh://'
func GetRepositoryFromCommand(command string) (string, string) {
	var (
		destination     string
		port            string = "22"
		repository      string
		gitShellCommand string
		arg             []string = strings.Split(command, " ")
	)
	/// "usr/bin/ssh git@github.com git-upload-pack 'bitexploder/timmy.git'"
	// "usr/bin/ssh -p 7999 git@globaldevtools.bbva.com git-receive-pack '/uqnwi/bitbucket_lifecycle.git'"

	for i, s := range arg {
		if strings.HasPrefix(s, "-p") && len(s) == 2 {
			port = arg[i+1]
		} else if strings.HasPrefix(s, "-p") && len(s) > 2 {
			port = s[2:]
		}
		if strings.HasPrefix(s, "git@") {
			destination = s
		}
		if strings.HasPrefix(s, "'") && strings.HasSuffix(s, ".git'") {
			repository = s[1 : len(s)-1]
		} else if strings.HasPrefix(s, "/") && strings.HasSuffix(s, ".git") {
			repository = s
		}
		if strings.HasPrefix(s, "git-") {
			gitShellCommand = s
		}
	}
	if !(strings.HasPrefix(repository, "/")) {
		repository = "/" + repository
	}
	if len(destination) > 0 && len(port) > 0 && len(repository) > 0 {
		return "ssh://" + destination + ":" + port + repository, gitShellCommand
	}
	return "", ""
}

// check if FILENAME exists and is a file
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// returns socket connection with proxyaddr (string with format `host:port`)
func CreateNetSocket(host *string, port *int, timeout *int) net.Conn {
	proxyaddr := *host + ":" + strconv.Itoa(*port)
	conn, err := net.DialTimeout("tcp", proxyaddr, time.Duration(*timeout)*time.Second)
	if err != nil {
		log.Printf("No connection to Proxy '%s:%d'\n", *host, *port)
		os.Exit(1)
	}
	return conn
}

func main() {
	log.SetFlags(0)
	var (
		uri               string
		buffer            = make([]byte, bufsize)
		read              int
		write             int
		setup             int = 0
		proxyhost             = flag.String("proxy_host", "squid", "Proxy Host")
		proxyport             = flag.Int("proxy_port", 3128, "Proxy Port")
		proxytimeout          = flag.Int("proxy_timeout", 3, "Proxy Timeout Connection")
		desthost              = flag.String("dest_host", "foo_bar.com", "Destination Host")
		destport              = flag.Int("dest_port", 22, "Destination Port")
		krb5conf              = flag.String("krb5conf", "/etc/krb5.conf", "Path to Kerberos Config")
		krbspn                = flag.String("krb_spn", "HTTP/squid-samuel", "Kerberos Service Principal Name for proxy authentication")
		krbauth               = flag.Bool("krb_auth", false, "Use Kerberos authentication for proxy users")
		basicauth             = flag.Bool("basic_auth", false, "Use basic authentication for proxy users")
		basicauthcredfile     = flag.String("creds_file", "/foo/bar", "Filepath of proxy credentials")
	)

	flag.Parse()

	logfile := "/tmp/gorkscrew_" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetOutput(file)

	log.Println("Proxy Host:", *proxyhost)
	log.Println("Proxy Port:", *proxyport)
	log.Println("Proxy Timeout:", *proxytimeout)
	log.Println("Destination Host:", *desthost)
	log.Println("Destination Port:", *destport)
	log.Println("Basic authentication:", strconv.FormatBool(*basicauth))
	if *basicauth {
		if *basicauthcredfile != "/foo/bar" {
			log.Println("Credentials will be loaded from environment variables 'GORKSCREW_AUTH' or 'CORKSCREW_AUTH'")
		} else {
			log.Printf("BasicAuth Credentials file: '%s'\n", *basicauthcredfile)
		}
	}

	log.Println("Kerberos authentication:", strconv.FormatBool(*krbauth))
	if *krbauth {
		log.Println("Kerberos Config:", *krb5conf)
		log.Println("Kerberos SPN:", *krbspn)
	}

	if !FileExists(*krb5conf) {
		log.Printf("kerberos configfile '%s' not exists", *krb5conf)
		os.Exit(10)
	}

	if *krbauth {
		uri = GetURIKerberosAuth(krb5conf, krbspn, desthost, destport)
	} else if *basicauth {
		uri = GetURIBasicAuth(basicauthcredfile, desthost, destport)
	} else {
		uri = GetURINoAuth(desthost, destport)
	}

	conn := CreateNetSocket(proxyhost, proxyport, proxytimeout)
	defer conn.Close()
	log.Println("")

	for {
		if setup == 0 {
			write, _ = conn.Write([]byte(uri))
			if write <= 0 {
				break
			}
			read, _ = conn.Read(buffer)
			if read <= 0 {
				break
			}
			statusStr := strings.Split(string(buffer[:]), " ")[1]
			statusCode, _ := strconv.Atoi(statusStr)
			if statusCode >= 200 && statusCode < 300 {
				log.Printf("Connection stablished. STATUS CODE: %d\n", statusCode)
				setup = 1
			} else if statusCode >= 407 {
				log.Printf("Proxy could not open connection. STATUS CODE: %d\n", statusCode)
				os.Exit(1)
			}
		} else {
			FeelTheMagic(conn)
		}
	}
}
