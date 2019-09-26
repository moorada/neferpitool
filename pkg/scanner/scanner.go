package scanner

import (
	"sync"
	"time"

	"github.com/sheerun/queue"

	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
)

func UpdateTypoDomains(tds domains.TypoList, c chan int) map[string]error {

	errs := make(map[string]error)

	start := time.Now()
	log.Info("Scanning of typodomains in progress...")
	var wg sync.WaitGroup
	wg.Add(len(tds))
	var wgWhois sync.WaitGroup

	mw := configuration.GetConf().TIMETOSLEEPWHOIS

	dnsDone := false
	q := queue.New()
	go func(c chan int) {
		for {
			time.Sleep(time.Duration(mw) * time.Millisecond)
			if q.Length() == 0 {
				if dnsDone {
					return
				}
			} else {
				item := q.Pop()
				pos := item.(int)
				go func(i int, c chan int) {
					log.Debug("Checking Whois of %s", tds[i].Name)
					err := tds[i].UpdateWhois()
					if err != nil {
						log.Debug("error updatewhois about %s: %s", tds[i].Name, err)
						errs[tds[i].Name] = err
					}
					c <- i
					wgWhois.Done()
				}(pos, c)
			}
		}
	}(c)

	ms := configuration.GetConf().TIMETOSLEEPSOA
	for i := range tds {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		go func(i int, c chan int) {
			log.Debug("Checking Status of %s", tds[i].Name)
			err := tds[i].UpdateStatus()
			if err != nil {
				log.Debug("error updatestatus about %s: %s", tds[i].Name, err)
				errs[tds[i].Name] = err
			}

			if tds[i].Status == constants.INACTIVE || tds[i].Status == constants.ACTIVE {
				wgWhois.Add(1)
				q.Append(i)
			} else {
				c <- i
			}
			wg.Done()
		}(i, c)
	}
	wg.Wait()
	dnsDone = true
	wgWhois.Wait()

	elapsed := time.Since(start)
	log.Info("Scanned %v typodomains in time: %s", len(tds), elapsed.String())
	log.Debug("When crash %v", 6)

	return errs

}
