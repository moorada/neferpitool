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

func GetRelaibleChangesFromDBWithExpression(expression string) (err error, changes reliableChanges.ReliableChangeList) {

	var changesTemp reliableChanges.ReliableChangeList
	err = db.Model(&reliableChanges.ReliableChange{}).Preload("Crons").Find(&changesTemp).Error

	for _, c := range changesTemp {
		if reliableChanges.Contains(c.Crons, expression) >= 0 {
			changes = append(changes, c)
		}
	}

	return
}

func DeleteExprToChDB(rc reliableChanges.ReliableChange, cron reliableChanges.CronExpression) {
	db.Model(&rc).Association("Crons").Delete(cron)
}

func GetCronExpressionFromDB(expression string) reliableChanges.CronExpression {
	var cron reliableChanges.CronExpression
	db.Where("exrpression = ?", expression).Last(&cron)
	return cron
}
