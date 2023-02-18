<p align="center">
  <img alt="Neferpitool" src="https://raw.githubusercontent.com/moorada/neferpitool/master/logo.png" width="40%" />

</p>


  <p align="center"><a href="https://github.com/moorada/neferpitool/blob/master/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a></p>

Neferpitool is a tool that combines DNS and WHOIS to automatically monitor domain name changes.

The faster, unrestricted DNS protocol is used on all domains to verify their status by analyzing DNS records. Meanwhile, the WHOIS protocol is used only on registered domains and thus on a more limited number of domains.
This will take advantage of the DNS protocol to filter the domains to be monitored and then use the protocol WOIS which is more limited, to extract information about domains.

## Features

* Generation and monitoring of domain variations (typo-domains)
* Identification and storage of changes regarding the status of typo-domains (status about DNS records and WHOIS info)
* Background monitoring with notification via email
* Configuration of scan and display parameters

## Install

```
go get github.com/moorada/neferpitool/cmd
```

## usage

Add to monitoring or manage a domain

```
./Neferpitool github.com

```
### Flags

```
  -bg
    	Active monitoring in background
  -it string -p string
    	Import Typos from file : -main domain -path of the file
  -logs
    	Avtive logs on file
  -mc
    	Make config file
  -pd
    	Check if domains are present
  -td string
    	Manage one typodomain

```


### Configuration file fields

```
    "TYPOSALGHORITM": List of typosquatting alghoritms
    "EXPIRATIONTIME": Number of days remaining until the domain expires for which to be notified by email
    "MAXATTEMPTSWHOIS": number of attempted WHOIS requests if there are failures
    "MAXATTEMPTSSOA": number of attempted DNS requests if there are failures
    "TIMETOSLEEPWHOIS": waiting time between requests
    "TIMETOSLEEPSOA": waiting time between requests
    "EMAIL": email used to send email
    "PASSWORD":
    "EMAILTONOTIFY": [] emails that will receive the email
    "SHOWSTATUS": []  Domains with states listed in this list will be displayed
    "PATHRESOLVER": Path of DNS solver (resolv.conf)
    "HOURSLEEPBACKGROUNDMONITORING": Waiting time between total scans in background mode
    "CHECKRELIABILITYTIME": When scanning for typodomains fails or is unreliable this time is waited before rescanning
    "REPORTFREQUENCY" Cron regex, indicates when to email the report

```
### Domain states
The state of domains varies depending on the presence of the following records DNS

```
UNKNOWN: in case of errors in DNS requests;
INACTIVE: the SOA or NS record is present, this indicates that the typodomain is no longer available for purchase.
```
Depending on the presence of other records, the typodomain may  go into the ALIAS or ACTIVE state;
```
ALIAS: if it is the CNAME record is present;
ACTIVE: if A, AAAA or MX records are present;
```

### Typosquatting alghoritms
```
  MD	Missing Dot is created by omitting a dot from the domain.
  MDS	Missing Dashes is created by stripping all dashes from the domain.
  CO	Character Omission Omitting a character from the domain.
  CS	Character Swap Swapping two consecutive characters in a domain
  ACS	Adjacent Character Substitution replaces adjacent characters
  ACI	Adjacent Character Insertion inserts adjacent character 
  CR	Character Repeat Repeats a character of the domain name twice
  DCR	Double Character Replacement repeats a character twice.
  SD	Strip Dashes is created by omitting a dash from the domain
  SP	Singular Pluralise creates a singular domain plural and vice versa
  CM	Common Misspellings are created from a dictionary of commonly misspelled words
  VS	Vowel Swapping is created by swaps vowels
  HG	Homoglyphs replaces characters with characters that look similar
  WTLD	Wrong Top Level Domain
  W2TLD	Wrong Second Level Domain
  W3TLD	Wrong Third Level Domain
  HP	Homophones Typos are created from sets of words that sound the same
  BF	Bitsquatting relies on random bit-errors to redirect connections
  NS	Numeral Swap numbers, words and vice versa
  ALL   Apply all typosquatting algorithms
```

## TO DO Run in a docker container 
Add to monitoring or manage a domain
```
docker build -t neferpitool .

```

## DNS for testing

See [DNS](dns) 