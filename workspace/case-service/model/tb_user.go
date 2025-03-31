package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type DataUser struct {
	UserID     null.Int    `gorm:"column:user_id;primary_key"`
	Email      null.String `gorm:"column:email"`
	Active     null.String `gorm:"column:active"`
	TimeCreate null.Time   `gorm:"column:time_create"`
}

// TableName sets the insert table name for this struct type
func (q *DataUser) TableName() string {
	return "tb_user"
}
