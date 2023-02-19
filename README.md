<p align="center">
  <img alt="Neferpitool" src="https://raw.githubusercontent.com/moorada/neferpitool/master/logo.png" width="40%" />

</p>


  <p align="center"><a href="https://github.com/moorada/neferpitool/blob/master/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a></p>


Neferpitool is a powerful tool that seamlessly integrates DNS and WHOIS protocols to automatically monitor domain name changes. By utilizing the faster and more unrestricted DNS protocol, Neferpitool is able to analyze the DNS records of all monitored domains and verify their status in real time. This enables the tool to quickly filter out domains that are not registered or are otherwise unavailable for WHOIS monitoring.

For the remaining registered domains, Neferpitool utilizes the more limited but highly accurate WHOIS protocol to extract critical information about the domain. This includes details such as the domain owner, registration date, and any recent changes to the domain's registration information.

By combining these two protocols, Neferpitool is able to provide comprehensive and highly accurate monitoring of domain names, giving users peace of mind that they will be alerted to any changes that could potentially impact their online presence.

In addition to its comprehensive domain name monitoring capabilities, Neferpitool can also be a valuable tool in the fight against phishing. Monitoring the registration information and status of all domains, including those that are recently registered or have recently changed ownership can be a critical early warning sign of a potential phishing attempt, allowing users to take proactive steps to protect themselves and their information.

Furthermore, Neferpitool is built with a user-friendly CLI interface that makes it easy to set up and manage domain monitoring. It also provides detailed reports and alerts via email to ensure that users are always up-to-date on the status of their monitored domains.

Overall, Neferpitool is a powerful and indispensable tool for anyone who wants to stay on top of their online presence and protect their valuable domain names.


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