package db

import (
	"github.com/moorada/neferpitool/pkg/reliableChanges"
)

/*Add a typo-DomainName in the DataBase*/
func SaveReliableChangeToDB(rc reliableChanges.ReliableChange) {
	db.Save(&rc)
}

func AddReliableChangeToDB(rc reliableChanges.ReliableChange) {
	db.Create(&rc)
}

func GetRelaibleChangesFromDB() (err error, changes reliableChanges.ReliableChangeList) {
	err = db.Model(&reliableChanges.ReliableChange{}).Preload("Crons").Find(&changes).Error
	return
}

func GetRelaibleChangesFromDBWithoutExpression(expression string) (err error, changes reliableChanges.ReliableChangeList) {

	var changesTemp reliableChanges.ReliableChangeList
	err = db.Model(&reliableChanges.ReliableChange{}).Preload("Crons").Find(&changesTemp).Error

	for _, c := range changesTemp {
		if !reliableChanges.Contains(c.Crons, expression) {
			changes = append(changes, c)
		}
	}

	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
