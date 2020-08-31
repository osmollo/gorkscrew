FROM debian:buster

ENV DEBIAN_FRONTEND noninteractive
ENV REALM ${REALM:-EXAMPLE.COM}
ENV KADMIN_PRINCIPAL ${KADMIN_PRINCIPAL:-kadmin/admin}
ENV KADMIN_PASSWORD ${KADMIN_PASSWORD:-MITiys4K5}
ENV KRB5_KTNAME /tmp/keytabs/squid.keytab

EXPOSE 3128

RUN apt-get -qq update && \
    apt-get -qq install locales squid vim krb5-user procps net-tools && \
    apt-get -qq clean && \
    locale-gen "en_US.UTF-8" && \
    echo "LC_ALL=\"en_US.UTF-8\"" >> /etc/default/locale && \
    mkdir /var/cache/squid && \
    chmod 777 /var/cache/squid

COPY entrypoint.sh /
ENTRYPOINT [ "/entrypoint.sh"]
