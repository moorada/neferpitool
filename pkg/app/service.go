package app

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/generator"
	"github.com/moorada/neferpitool/pkg/reliableChanges"
	"github.com/moorada/neferpitool/pkg/scanner"
)

type ProgressFn func(done, total int)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) DomainPresence(domainNames []string) map[string]bool {
	domainSet := map[string]bool{}
	for _, d := range db.GetMainDomainListFromDB() {
		domainSet[d.Name] = true
	}

	presence := map[string]bool{}
	for _, name := range domainNames {
		presence[name] = domainSet[name]
	}
	return presence
}

func (s *Service) ScanTypoDomains(tds domains.TypoList, progress ProgressFn) map[string]error {
	c := make(chan int, len(tds))
	errsCh := make(chan map[string]error, 1)

	go func() {
		errsCh <- scanner.UpdateTypoDomains(tds, c)
	}()

	done := 0
	total := len(tds)
	for done < total {
		<-c
		done++
		if progress != nil {
			progress(done, total)
		}
	}

	return <-errsCh
}

func (s *Service) AddDomainAndTypos(domain string, progress ProgressFn) (domains.TypoList, map[string]error, error) {
	tds := generator.GetUnfilledTypoDomains(domain)
	errs := s.ScanTypoDomains(tds, progress)
	db.AddTypoListToDB(tds)

	md := domains.NewLegitDomain(domain)
	if err := md.Update(); err != nil {
		return tds, errs, err
	}

	db.AddLegitDomainToDB(md)
	return tds, errs, nil
}

func (s *Service) ImportTypos(domain, path string, progress ProgressFn) (domains.TypoList, map[string]error, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	var tds domains.TypoList
	for sc.Scan() {
		name := strings.TrimSpace(sc.Text())
		if name == "" {
			continue
		}
		tds = append(tds, domains.NewTypoDomain(name, domain, "imported"))
	}

	if err := sc.Err(); err != nil {
		return nil, nil, err
	}

	if len(tds) == 0 {
		return nil, nil, errors.New("empty file")
	}

	errs := s.ScanTypoDomains(tds, progress)
	db.AddTypoListToDB(tds)
	return tds, errs, nil
}

func (s *Service) GetTypoDomainsInExpiration() domains.TypoList {
	var total domains.TypoList
	expirationDays := configuration.GetConf().EXPIRATIONTIME

	for _, d := range db.GetMainDomainListFromDB() {
		tds := db.GetTypoDomainListWithStatusFromDB(d.Name, []int{constants.INACTIVE, constants.ACTIVE, constants.ALIAS})
		total = append(total, tds.FilterInExpiration(expirationDays)...)
	}

	return total
}

func (s *Service) IterateCheckGetChanges(tds domains.TypoList, progress ProgressFn) (tdsReliable []domains.TypoDomain, changesReliable []changes.Change, scanErrs map[string]error) {
	tdsNew := tds.GetUnfilledCopy()
	scanErrs = s.ScanTypoDomains(tdsNew, progress)

	tdsOldCh, tdsNewCh, chs := changes.MakeChangeList(tds, tdsNew)
	for i := 0; i < 2 && len(tdsOldCh) > 0; i++ {
		time.Sleep(time.Duration(configuration.GetConf().CHECKRELIABILITYTIME) * time.Millisecond)
		scanErrs = mergeErrs(scanErrs, s.ScanTypoDomains(tdsNewCh, progress))

		tdsOldChNext, tdsNewChNext, chsNext := changes.MakeChangeList(tdsOldCh, tdsNewCh)
		tdsReliable, changesReliable = chsNext.FilterReliableWithPrev(chs, tdsNewCh, tdsNewChNext)
		tdsOldCh, tdsNewCh, chs = tdsOldChNext, tdsReliable, changesReliable
	}

	return tdsReliable, changesReliable, scanErrs
}

func (s *Service) SaveReliableChanges(changesList changes.ChangeList) {
	var crons []*reliableChanges.CronExpression
	for _, expression := range configuration.GetConf().REPORTFREQUENCY {
		cronExpr := db.GetCronExpressionFromDB(expression)
		if cronExpr.ID != 0 {
			crons = append(crons, &cronExpr)
		} else {
			crons = append(crons, &reliableChanges.CronExpression{Exrpression: expression})
		}
	}

	var rChanges []reliableChanges.ReliableChange
	for _, c := range changesList {
		rChanges = append(rChanges, reliableChanges.ReliableChange{
			TypoDomain: c.TypoDomain,
			Field:      c.Field,
			Before:     c.Before,
			After:      c.After,
			Crons:      crons,
		})
	}
	db.AddReliableChangeListToDB(rChanges)
}

func mergeErrs(first, second map[string]error) map[string]error {
	if first == nil && second == nil {
		return map[string]error{}
	}
	if first == nil {
		first = map[string]error{}
	}
	for k, v := range second {
		first[k] = v
	}
	return first
}
