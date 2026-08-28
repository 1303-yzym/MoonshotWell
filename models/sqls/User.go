package sqls

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

func (u UserModel) TableName() string {
	return "user"
}
