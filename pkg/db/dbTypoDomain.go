package db

import "github.com/moorada/neferpitool/pkg/domains"

/*Add a typo-DomainName in the DataBase*/
func AddTypoDomainToDB(td domains.TypoDomain) {
	db.Create(&td)
}

/*Get a typo-DomainName from the DataBase*/
func GetTypoDomainFromDB(nameTypoDomain string) domains.TypoDomain {

	var typoDomain domains.TypoDomain
	db.Where("name = ?", nameTypoDomain).Last(&typoDomain)
	return typoDomain
}

func RemoveTypoDomainFromDB(nameTypoDomain string) {
	db.Where("name = ?", nameTypoDomain).Delete(&domains.TypoDomain{})
}

/*Return a list of all the typo-domains about mainDomain*/
func GetTypoDomainHistoryFromDB(typoDomain string) domains.TypoList {
	var tds []domains.TypoDomain
	db.Where("name = ?", typoDomain).Find(&tds)
	return tds
}
