package db

import (
	"github.com/moorada/neferpitool/pkg/domains"
)

/*Add a main-DomainName in the DataBase*/
func AddLegitDomainToDB(d domains.LegitDomain) {
	legitDomain := domains.NewLegitDomain(d.Name)
	db.Where(legitDomain).Assign(d).FirstOrCreate(&d)

}

/*Get a main-DomainName from the DataBase*/
func GetLegitDomainFromDB(nameMainDomain string) (domains.LegitDomain, error) {
	var mainDomain domains.LegitDomain
	err := db.Where("name = ?", nameMainDomain).Last(&mainDomain).Error
	return mainDomain, err
}

/*remove a main-DomainName from the DataBase*/
func RemoveLegitDomainFromDB(nameLegitDomain string) {
	var mainDomain domains.LegitDomain
	db.Where("name = ?", nameLegitDomain).Last(&mainDomain)
	db.Delete(mainDomain)
	db.Where("legit_domain = ?", nameLegitDomain).Delete(&domains.TypoDomain{})
}
