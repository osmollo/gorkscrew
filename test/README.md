# TEST

- [TEST](#test)
  - [No authentication](#no-authentication)
  - [Basic authentication](#basic-authentication)
  - [Kerberos authentication](#kerberos-authentication)

## No authentication

```shell
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

```shell
git clone git@github.com:osmollo/gorkscrew.git /tmp/gorkscrew
```

## Basic authentication

```shell
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

```shell
export GORKSCREW_AUTH="test:test1234"
```

Finally, we can clone any github repo:

```shell
git clone git@github.com:osmollo/gorkscrew.git /tmp/gorkscrew
```

## Kerberos authentication

```shell
cd krb_auth
chmod 777 squid/keytabs
docker-compose up -d
```

We must install kerberos client package:

```shell
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

```shell
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

```shell
git clone git@github.com:osmollo/gorkscrew.git /tmp/gorkscrew
```
