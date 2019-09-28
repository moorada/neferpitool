package configuration

import (
	"encoding/json"
	"os"

	"github.com/tkanos/gonfig"

	"github.com/moorada/neferpitool/pkg/log"
)

type configuration struct {
	TYPOSALGHORITM                []string
	EXPIRATIONTIME                int
	MAXATTEMPTSWHOIS              int
	MAXATTEMPTSSOA                int
	TIMETOSLEEPWHOIS              int
	TIMETOSLEEPSOA                int
	EMAIL                         string
	PASSWORD                      string
	EMAILTONOTIFY                 []string
	SHOWSTATUS                    []string
	PATHRESOLVER                  string
	HOURSLEEPBACKGROUNDMONITORING int
}

const (
	pathConfig = "./config/config.json"
)

var (
	initVar bool

	standardConf = configuration{
		TYPOSALGHORITM:                []string{"all"},
		EXPIRATIONTIME:                7,
		MAXATTEMPTSWHOIS:              2,
		MAXATTEMPTSSOA:                2,
		TIMETOSLEEPWHOIS:              2000,
		TIMETOSLEEPSOA:                100,
		EMAIL:                         "",
		PASSWORD:                      "",
		EMAILTONOTIFY:                 []string{},
		SHOWSTATUS:                    []string{"Inactive", "Active", "Available", "Unknown", "Alias"},
		PATHRESOLVER:                  "./config/resolv.conf",
		HOURSLEEPBACKGROUNDMONITORING: 1,
	}
)

func initConf() {

	configuration := configuration{}
	err := gonfig.GetConf(pathConfig, &configuration)
	if err != nil {
		log.Debug("Basic configuration")
	} else {
		log.Debug("configuration by file")
		standardConf = configuration
	}
}

func GetConf() configuration {
	if !initVar {
		initConf()
		initVar = true
	}
	return standardConf
}

func MakeConfigFile() error {

	config := GetConf()

	jsonFile, err := os.Create(pathConfig)
	if err != nil {
		panic(err)
	} else {
		jsonData, err := json.Marshal(config)

		_, err = jsonFile.WriteAt(jsonData, 0)
		if err != nil {
			panic(err)
		}

	}

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("File \"config.json\" created")
	}

	jsonFile.Close()
	return nil
}
