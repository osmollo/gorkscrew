#!/bin/bash

htpasswd -bc /var/tmp/squid_users test test1234

exec $@
