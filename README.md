<p align="center">
  <img alt="Neferpitool" src="https://raw.githubusercontent.com/moorada/neferpitool/master/logo.png" width="40%" />

</p>


  <p align="center"><a href="https://github.com/moorada/neferpitool/blob/master/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a></p>

A tool that combines DNS and WHOIS to automatically monitor domain name changes.

## Features

* Generation and monitoring of domain variations (typo-domains)
* Identification and storage of changes regarding the status of typo-domains (status about DNS records and WHOIS info)
* Background monitoring with user notification via email
* Configuration of scan and display parameters

## Install
```
go get github.com/moorada/neferpitool/cmd
//in neferpitool path
GO111MODULE=on go mod vendor
//in cmd path
go build
./cmd
```

## usage
Add to monitoring or manage a domain
```
./Neferpitool github.com

```

Active background scan
```
./Neferpitool -bg

```

Active background scan with log
```
./Neferpitool -bg -logs

```

Make config file (if you want to configure the scan and display parameters)
```
./Neferpitool -mc

```
