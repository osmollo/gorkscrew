FROM debian:buster

ENV DEBIAN_FRONTEND noninteractive

EXPOSE 3128

RUN apt-get -qq update && \
    apt-get -qq install locales squid vim apache2-utils procps net-tools && \
    apt-get -qq clean && \
    locale-gen "en_US.UTF-8" && \
    echo "LC_ALL=\"en_US.UTF-8\"" >> /etc/default/locale && \
    mkdir /var/cache/squid && \
    chmod 777 /var/cache/squid

COPY entrypoint.sh /
ENTRYPOINT [ "/entrypoint.sh"]
CMD ["/usr/sbin/squid", "-f", "/etc/squid/squid.conf", "--foreground", "-YCd", "1"]
