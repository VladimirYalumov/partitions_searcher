package models

import (
	"partitions_searcher/settings"
	"time"
)

type Task struct {
	Id          int64  `gorm:"primary_key" json:"id"`
	Title       string `gorm:"type:varchar(255);default:''" json:"title"`
	Description string `gorm:"type:text;default:''" json:"description"`
	CreatedAt   time.Time
}

func (model *Task) Create() error {
	err := settings.Db.Create(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
