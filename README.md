# neferpitool
The tool have the function of storing and showing information about the domains and typo domains of interest. It will also be able to identify, monitor and notify changes about typo domains via email.



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
