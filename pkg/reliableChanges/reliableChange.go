package reliableChanges

import "github.com/jinzhu/gorm"

type ReliableChange struct {
	gorm.Model
	TypoDomain string
	Field      string
	Before     string
	After      string
}
