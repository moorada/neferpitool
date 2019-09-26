package generator

import (
	"github.com/cybint/urlinsane/pkg/typo"
	"github.com/moorada/neferpitool/pkg/configuration"

	"github.com/moorada/neferpitool/pkg/domains"
)

/*return a slice of all typodomains*/
func GetUnfilledTypoDomains(mainDomain string) domains.TypoList {

	var conf typo.BasicConfig

	typos := configuration.GetConf().TYPOSALGHORITM
	conf = typo.BasicConfig{
		Domains:     []string{mainDomain},
		Keyboards:   []string{"all"},
		Typos:       typos,
		Funcs:       []string{""},
		Concurrency: 50,
		Format:      "text",
		Verbose:     false,
	}
	u := typo.New(conf.Config())

	out := u.Stream()

	var tds domains.TypoList
	for t := range out {

		td := domains.NewTypoDomain(t.Variant.Idna(), mainDomain, t.Typo.Name)
		tds = append(tds, td)
	}

	return tds

}
