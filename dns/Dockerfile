FROM ubuntu:bionic

RUN apt-get update \
  && apt-get install -y \
  bind9 \
  bind9utils \
  bind9-doc

# Enable IPv4
RUN sed -i 's/OPTIONS=.*/OPTIONS="-4 -u bind"/' /etc/default/bind9

# Copy configuration files
COPY etc/bind/named.conf.options /etc/bind/
COPY etc/bind/named.conf.local /etc/bind/
COPY etc/bind/zones/db.test-neferpitool.com /etc/bind/zones/
COPY etc/bind/zones/db.neferpitool /etc/bind/zones/
COPY etc/bind/zones/db2.neferpitool /etc/bind/zones/
# Run eternal loop
CMD ["/bin/bash", "-c", "while :; do sleep 10; done"]
