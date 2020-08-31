#!/bin/bash

echo "==================================================================================="
echo "==== Kerberos Client =============================================================="
echo "==================================================================================="
KADMIN_PRINCIPAL_FULL=$KADMIN_PRINCIPAL@$REALM

echo "REALM: $REALM"
echo "KADMIN_PRINCIPAL_FULL: $KADMIN_PRINCIPAL_FULL"
echo "KADMIN_PASSWORD: $KADMIN_PASSWORD"
echo ""

function kadminCommand {
    kadmin -p $KADMIN_PRINCIPAL_FULL -w $KADMIN_PASSWORD -q "$1"
}


rm -fr /tmp/keytabs/*.keytab
# create princ/keytab for squid
kadminCommand "addprinc -randkey HTTP/$(hostname -f)@${REALM}"
kadminCommand "ktadd -k /tmp/keytabs/$(hostname -f).keytab HTTP/$(hostname -f)@${REALM}"

# create princ/keytab for client
kadminCommand "addprinc -randkey client@${REALM}"
kadminCommand "ktadd -k /tmp/keytabs/client.keytab client@${REALM}"

chmod 777 /tmp/keytabs/*.keytab


/usr/sbin/squid -f /etc/squid/squid.conf --foreground -YCd 1
