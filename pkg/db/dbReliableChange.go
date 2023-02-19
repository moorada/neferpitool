package db

import (
	"time"

	"github.com/moorada/neferpitool/pkg/reliableChanges"
)

/*Add a typo-DomainName in the DataBase*/
func AddReliableChangeToDB(rc reliableChanges.ReliableChange) {
	db.Create(&rc)
}

func GetRelaibleChangesFromDB(from time.Time) reliableChanges.ReliableChangeList {
	var changes reliableChanges.ReliableChangeList
	db.Where("created_at > ?", from).Find(&changes)
	return changes
}
