package db

import (
	"github.com/moorada/neferpitool/pkg/reliableChanges"
)

/*Add a typo-DomainName in the DataBase*/
func AddReliableChangeToDB(rc reliableChanges.ReliableChange) {
	db.Create(&rc)
}
