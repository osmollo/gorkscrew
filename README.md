# GORKSCREW

- [GORKSCREW](#gorkscrew)
  - [Dependencies](#dependencies)
  - [Arguments](#arguments)
  - [Build gorkscrew](#build-gorkscrew)
  - [Execute Gorkscrew](#execute-gorkscrew)
  - [Testing](#testing)
    - [No authentication](#no-authentication)
    - [Basic authentication](#basic-authentication)
    - [Kerberos authentication](#kerberos-authentication)

## Dependencies

For **GO** installation:

```bash
git clone git@github.com:ohermosa/my_workstation.git
cd my_workstation/ansible
ansible-playbook install.yml -t go
```

After clone the repository, there will be to define the following environment variables:

```bash
mkdir $HOME/go

export GOROOT=/usr/local/go
export GOPATH=$HOME/go
```

And install the following external **GO** modules:

```bash
go get github.com/jcmturner/gokrb5/v8/client
go get github.com/jcmturner/gokrb5/v8/config
go get github.com/jcmturner/gokrb5/v8/credentials
go get github.com/jcmturner/gokrb5/v8/spnego
```

## Arguments

`gorkscrew` can receive the following arguments:

| NAME | DESCRIPTION | DEFAULT |
|--|--|--|
| proxy_host | proxy hostname/IP | `squid` |
| proxy_port | proxy port | `3128` |
| proxy_timeout | proxy timeout connection | `5` |
| dest_host | destination host | `foo_bar.com` |
| dest_port | destination port | `22` |
| krb_auth | enable kerberos authentication | `false` |
| krb5conf | path to `krb5.conf` file | `/etc/krb5.conf` |
| krb_spn | Kerberos SPN for kerberos authentication with proxy | `HTTP/squid-samuel` |
| basic_auth | enable basic authenticacion | `false` |
| creds_file | path to file with proxy credentials | `/foo/bar` |
| log | enable logging | `false` |
| log_file | path to log file | `/tmp/gorkscrew_$TIMESTAMP.log` |
| version | show gorkscrew version | false |

## Build gorkscrew

With this command, the `./gorkscrew` binary will be builded:

```bash
go build -ldflags "-X 'main.GorkscrewVersion=$(jq -r .version release.json)' -X 'main.GoVersion=$(jq -r .go_version release.json)'" gorkscrew.go
```

## Execute Gorkscrew

`Gorkscrew` can receive the following arguments:

```bash
./gorkscrew -h
Usage of gorkscrew:
  -basic_auth
        Use basic authentication for proxy users
  -creds_file string
        Filepath of proxy credentials (default "/foo/bar")
  -dest_host string
        Destination Host (default "foo_bar.com")
  -dest_port int
        Destination Port (default 22)
  -krb5conf string
        Path to Kerberos Config (default "/etc/krb5.conf")
  -krb_auth
        Use Kerberos authentication for proxy users
  -proxy_host string
        Proxy Host (default "squid")
  -proxy_port int
        Proxy Port (default 3128)
  -proxy_timeout int
        Proxy Timeout Connection (default 3)
  -krb_spn string
        Kerberos Service Principal Name for proxy authentication (default "HTTP/squid-samuel")
  -log
        enable logging
  -log_file string
        Save log execution to file (default "/foo/bar.log")
  -version
        Show gorkscrew version
```

According to the type of proxy authentication, we will need to use an argument or another:

```text
Host foo_bar.com
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host squid.internal.domain --proxy_port 3128 --dest_host %h --dest_port %p
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host squid.internal.domain --proxy_port 3128 --dest_host %h --dest_port %p --basic_auth --creds_file /tmp/userpass.txt
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host squid.internal.domain --proxy_port 3128 --dest_host %h --dest_port %p --krb_auth --krb_spn HTTP/my_squid
```

## Testing

```bash
cd tests
```

### No authentication

```bash
cd no_auth
docker-compose up -d
curl  -x 172.23.0.3:3128 https://www.google.com -vvv
```

This section must be present in `~/.ssh/config`:

```text
Host github.com
  LogLevel DEBUG3
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host 172.23.0.3 --proxy_port 3128 --dest_host %h --dest_port %p
```

We can clone any github repository:

```bash
git clone git@github.com:ohermosa/gorkscrew.git /tmp/gorkscrew
```

### Basic authentication

```bash
cd basic_auth
docker-compose up -d
curl -x test:test1234@172.21.0.3:3128 https://www.google.com -vvv
```

This section must be present in `~/.ssh/config`:

```text
Host github.com
  LogLevel DEBUG3
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host 172.21.0.3 --proxy_port 3128 --dest_host %h --dest_port %p --basic_auth
```

We must define the environment variable `GORKSCREW_AUTH` with proxy credentials:

```bash
export GORKSCREW_AUTH="test:test1234"
```

Finally, we can clone any github repo:

```bash
git clone git@github.com:ohermosa/gorkscrew.git /tmp/gorkscrew
```

### Kerberos authentication

```bash
cd krb_auth
chmod 777 squid/keytabs
docker-compose up -d
```

We must install kerberos client package:

```bash
sudo apt install krb5-user

tee /etc/krb5.conf <<EOF
[libdefaults]
  default_realm = EXAMPLE.COM
[realms]
  EXAMPLE.COM = {
    kdc = 172.22.0.2
    admin_server = 172.22.0.2
  }
EOF
```

We can check if proxy is working with kerberos authenticacion:

```bash
kinit -kt squid/keytabs/client.keytab client
curl --proxy-negotiate -u : -x 172.22.0.3:3128 https://www.google.com -vvv
```

This section must be present in `~/.ssh/config`:

```text
Host github.com
  LogLevel DEBUG3
  ProxyCommand /usr/local/bin/gorkscrew --proxy_host 172.22.0.3 --proxy_port 3128 --dest_host %h --dest_port %p --krb_auth --krb_spn HTTP/squid
```

Finally, we can clone any github repo:

```bash
git clone git@github.com:ohermosa/gorkscrew.git /tmp/gorkscrew
```
