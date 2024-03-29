# DNS testing


This part of the project allows to configure a DNS server for testing purposes.


## Create a docker network
```
sudo docker network create --subnet=172.20.0.0/16 neferpitool-net
```

## Build the Docker image
```
sudo docker build -t bind9 .

```

## Run the DNS server container in background

Using the same IP as in the db.test-neferpitool.com file and the same Docker network created:

```
sudo docker run -d --rm --name=dns-server --net=neferpitool-net --ip=172.20.0.2 bind9 

```
Run bind9
```
sudo docker exec -d dns-server /etc/init.d/bind9 start
```

## Change your default nameserver

Neferpitool reads the resolv.conf file inside the config folder

```
nameserver 172.20.0.2

```

## Build and Run the neferpitool container TO DO

Run these commands in the root folder of neferpitool

```
cd ..
sudo docker build -t neferpitool .
sudo docker run -d --rm --name=neferpitool --net=neferpitool-net --ip=172.20.0.3 --dns=172.20.0.2 neferpitool /bin/bash -c "while :; do sleep 10; done"

```
## Bind9 check configuration files
```
named-checkconf /etc/bind/named.conf.local
named-checkzone neferpitool /etc/bind/zones/db.neferpitool
named-checkzone test-neferpitool.com /etc/bind/zones/db.test-neferpitool.com

```