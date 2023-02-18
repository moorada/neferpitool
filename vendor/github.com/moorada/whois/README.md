# whois
WHOIS client for Go-lang

```
go get github.com/moorada/whois
```

## Example
```
result, _ := whois.Whois("google.com")
fmt.Println(result)
	
```

This library store the WHOIS servers locally in a json file for each TLD and use the following policy:


![alt policy diagram](https://github.com/moorada/whois/blob/master/policy.png)


