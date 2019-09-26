# neferpitool
Domain name permutation engine for detecting typo squatting and monitoring domains



## Install
```
go get github.com/jackal00/neferpitool

```

## Usage
in cmd/domainmonitoring path
```
// add a new domain in the monitor-zone
go run main.go example.domain
go run main.go -debug -log=logx.txt google.com
// manage the monitor-zone
go run main.go
go run main.go -debug -log=logx.txt
```


## Dependencies

```
go get -u github.com/golang/dep/cmd/dep
dep ensure -v
```
