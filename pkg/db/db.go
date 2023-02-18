package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/reliableChanges"
)

var db *gorm.DB

/*Initialize the DataBase*/
func InitDB(nameDB string) {

	var err error
	db, err = gorm.Open("sqlite3", "./"+nameDB+".db")
	if err != nil {
		log.Fatal("Failed to connect database %s", err.Error())
	}

	// Migrate the schema
	db.AutoMigrate(&domains.LegitDomain{})
	db.AutoMigrate(&domains.TypoDomain{})
	db.AutoMigrate(&reliableChanges.ReliableChange{})
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Error(err.Error())
	}
}

func uniqueTypoDomains(ds domains.TypoList) domains.TypoList {

	domainsMap := make(map[string]domains.TypoDomain)

	for _, d := range ds {

		_, bol := domainsMap[d.Name]
		if bol {
			if domainsMap[d.Name].ID < d.ID {
				domainsMap[d.Name] = d
			}
		} else {
			domainsMap[d.Name] = d
		}
	}

	var dsUnique []domains.TypoDomain
	for _, dm := range domainsMap {
		dsUnique = append(dsUnique, dm)
	}
	return dsUnique

}

/*Return a list of all the typo-domains about mainDomain*/
func GetTypoDomainListFromDB(mainDomain string) domains.TypoList {
	var tds []domains.TypoDomain
	db.Where("legit_domain = ?", mainDomain).Find(&tds)
	/*senza uniqueTypoDomains:
	Select * , MAX(id) from typo_domains group by name
	*/
	//db.Exec("SELECT	 * , MAX(updated_at) from typo_domains group by name").Scan(&tds)
	return uniqueTypoDomains(tds)
	//return tds
}

func GetTypoDomainListWithStatusFromDB(mainDomain string, status []int) domains.TypoList {
	var tds []domains.TypoDomain
	db.Where("legit_domain = ? AND status IN (?)", mainDomain, status).Find(&tds)
	/*senza uniqueTypoDomains:
	Select * , MAX(id) from typo_domains group by name
	*/
	//db.Exec("SELECT	 * , MAX(updated_at) from typo_domains group by name").Scan(&tds)
	return uniqueTypoDomains(tds)
}

/*Return a list of all the typo-domains about mainDomain*/
func GetAllTypoDomainListFromDB() domains.TypoList {
	var tds domains.TypoList
	db.Find(&tds)
	/*senza uniqueTypoDomains:
	Select * , MAX(id) from typo_domains group by name
	*/
	//db.Exec("SELECT	 * , MAX(updated_at) from typo_domains group by name").Scan(&tds)
	return tds
}

func AddTypoListToDB(tds domains.TypoList) {
	for _, td := range tds {
		AddTypoDomainToDB(td)
	}
}

func AddReliableChangeListToDB(tds []reliableChanges.ReliableChange) {
	for _, td := range tds {
		AddReliableChangeToDB(td)
	}
}

/*Return a list of all the main-domains*/
func GetMainDomainListFromDB() domains.LegitList {
	var ds []domains.LegitDomain
	db.Find(&ds)
	return ds
}
