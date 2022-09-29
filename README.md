# GORKSCREW

- [GORKSCREW](#gorkscrew)
  - [GO version](#go-version)
  - [Build the binary](#build-the-binary)
    - [Dependencies](#dependencies)
    - [Build gorkscrew](#build-gorkscrew)
  - [How to use](#how-to-use)
    - [Execute Gorkscrew](#execute-gorkscrew)
  - [Testing](#testing)
  - [Buy me a coffee](#buy-me-a-coffee)

## GO version

![Go version](https://img.shields.io/badge/Go-1.18-brightgreen.svg)

## Build the binary
### Dependencies

For **GO** installation, download the desired version from [official download page](https://go.dev/dl/):

```shell
tar xvzf go1.18.6.linux-amd64.tar.gz
sudo mv go /usr/local

mkdir $HOME/go

export GOROOT=/usr/local/go
export GOPATH=$HOME/go
```

### Build gorkscrew

With this command, the `./gorkscrew` binary will be builded:

```shell
go build -ldflags "-X 'main.GorkscrewVersion=$(jq -r .version release.json)' -X 'main.GoVersion=$(jq -r .go_version release.json)'" gorkscrew.go
```

## How to use
### Execute Gorkscrew

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

You can see the usage help using the `-h` argument:

```shell
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

Please, read the [TESTING README](test/README.md)

## Buy me a coffee

If this repository has been helpful to you, but especially if you feel like it, you can invite me to a coffee

[![buy me a coffee](https://camo.githubusercontent.com/c3f856bacd5b09669157ed4774f80fb9d8622dd45ce8fdf2990d3552db99bd27/68747470733a2f2f7777772e6275796d6561636f666665652e636f6d2f6173736574732f696d672f637573746f6d5f696d616765732f6f72616e67655f696d672e706e67)](https://www.buymeacoffee.com/osmollo)
