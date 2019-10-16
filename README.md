<p align="center">
  <img alt="Neferpitool" src="https://raw.githubusercontent.com/moorada/neferpitool/master/logo.png" width="40%" />
</p>

A tool that combines DNS and WHOIS to automatically monitor domain name changes.

```
go get github.com/moorada/neferpitool/cmd
GO111MODULE=on go mod vendor
```

## usage
Make config file
```
./Neferpitool -mc

```
Add domain to monitoring
```
./Neferpitool github.com

```
Active background scan
```
./Neferpitool github.com -bg

```
