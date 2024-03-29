package configuration

import (
	"encoding/json"
	"os"

	"github.com/tkanos/gonfig"

	"github.com/moorada/neferpitool/pkg/log"
)

type configuration struct {
	TYPOSALGHORITM                  []string
	EXPIRATIONTIME                  int
	MAXATTEMPTSWHOIS                int
	MAXATTEMPTSSOA                  int
	TIMETOSLEEPWHOIS                int
	TIMETOSLEEPSOA                  int
	EMAIL                           string
	PASSWORD                        string
	EMAILTONOTIFY                   []string
	SHOWSTATUS                      []string
	PATHRESOLVER                    string
	MINUTESLEEPBACKGROUNDMONITORING int
	CHECKRELIABILITYTIME            int
	REPORTFREQUENCY                 []string
}

const (
	pathConfig = "./config/config.json"
)

var (
	initVar bool

	standardConf = configuration{
		TYPOSALGHORITM:                  []string{"all"},
		EXPIRATIONTIME:                  7,
		MAXATTEMPTSWHOIS:                2,
		MAXATTEMPTSSOA:                  2,
		TIMETOSLEEPWHOIS:                2000,
		TIMETOSLEEPSOA:                  100,
		EMAIL:                           "",
		PASSWORD:                        "",
		EMAILTONOTIFY:                   []string{},
		SHOWSTATUS:                      []string{"Inactive", "Active", "Available", "Unknown", "Alias"},
		PATHRESOLVER:                    "./config/resolv.conf",
		MINUTESLEEPBACKGROUNDMONITORING: 1,
		CHECKRELIABILITYTIME:            2000,
		REPORTFREQUENCY:                 []string{"0 0/4 * * * "},
	}
)

func initConf() {

	configuration := configuration{}
	err := gonfig.GetConf(pathConfig, &configuration)
	if err != nil {
		log.Debug("Basic configuration, err: %s", err)
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
