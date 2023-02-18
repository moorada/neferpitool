package whois

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/evilsocket/islazy/log"
)

// Query server const
const (
	DOMAIN_WHOIS_SERVER = "whois-servers.net"
	IANA_WHOIS_SERVER   = "whois.iana.org"
	WHOIS_PORT          = "43"
	TIMEOUT             = time.Second * 5
	PATH_LIST_TLD       = "./config/listTLD.json"
	PATH_CONFIG_FOLDER  = "./config"
)

func Whois(domain string, servers ...string) (result string, err error) {
	domain = strings.Trim(strings.TrimSpace(domain), ".")
	err = os.MkdirAll(PATH_CONFIG_FOLDER, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
	}

	tld := getTLD(domain)

	if len(servers) < 1 || servers[0] == "" {
		server, err := getWhoisServerByFile(PATH_LIST_TLD, tld)
		if err == nil {
			return query(domain, server)
		}

		r, err := whoisByDNS(domain)
		if err == nil {
			rejectionWords := []string{
				"ERROR:101",
			}
			reject := false
			for _, s := range rejectionWords {
				if strings.Contains(r, s) {
					log.Debug("Whois about %s contains rejection word :%s ", domain, s)
					reject = true
					break
				}
			}
			if !reject {
				return r, err
			}
		}

		server, err = getWhoisServerByIANA(getTLD(domain))
		if err == nil && server != "" {
			addCouple(PATH_LIST_TLD, tld, server)
			return query(domain, server)
		} else {
			return result, err
		}

	} else {
		var result string
		var err error
		for _, s := range servers {
			result, err := query(domain, s)
			if err == nil {
				return result, err
			}
		}
		return result, err
	}

}

func whoisByDNS(domain string) (result string, err error) {

	tld := getTLD(domain)
	server := tld + "." + DOMAIN_WHOIS_SERVER

	result, err = query(domain, server)
	if err == nil {
		addCouple(PATH_LIST_TLD, tld, server)
	}

	return
}

func getTLD(domain string) string {
	domains := strings.Split(domain, ".")
	if len(domains) < 2 {
		log.Fatal("Domain %s is invalid", domain)
	}
	return domains[len(domains)-1]

}

func getWhoisServerByIANA(tld string) (server string, err error) {
	result, e := query(tld, IANA_WHOIS_SERVER)
	if e == nil {
		token := "whois:"

		start := strings.Index(result, token)
		if start == -1 {
			return
		}
		start += len(token)
		end := strings.Index(result[start:], "\n")
		server = strings.TrimSpace(result[start : start+end])
	}
	return
}

func addCouple(path string, tld string, whoisAdmin string) {

	tldMap := make(map[string]interface{})

	jsonFile, err := os.Open(path)
	if err != nil {
		jsonFile, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}
		_ = json.Unmarshal(byteValue, &tldMap)

	}

	tldMap[tld] = whoisAdmin
	jsonData, err := json.Marshal(tldMap)

	jsonFile, err = os.Create(path)
	if err != nil {
		panic(err)
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		panic(err)
	}

	jsonFile.Close()

}

func query(domain string, servers ...string) (result string, err error) {

	for i := 0; i < len(servers) && result == ""; i++ {
		server := servers[i]
		server = strings.ToLower(server)
		if server == "whois.arin.net" {
			domain = "n + " + domain
		}

		conn, e := net.DialTimeout("tcp", net.JoinHostPort(server, WHOIS_PORT), TIMEOUT)

		if e != nil {
			err = e
			return
		} else {
			_, _ = conn.Write([]byte(domain + "\r\n"))
			_ = conn.SetReadDeadline(time.Now().Add(TIMEOUT))

			buffer, e := ioutil.ReadAll(conn)
			if e != nil {
				err = e
			} else {
				result = string(buffer)
				return
			}
		}

		conn.Close()
	}
	return
}

func getWhoisServerByFile(path string, tld string) (server string, err error) {

	tldMap := make(map[string]string)

	jsonFile, err := os.Open(path)
	defer jsonFile.Close()

	server = ""

	if err != nil {
		return "", err
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(byteValue, &tldMap)
		if err != nil {
			return "", err
		}
	}

	if tldMap[tld] != "" {
		server = tldMap[tld]
	} else {
		return "", errors.New("Server is a empty string")
	}

	return
}
