package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TeamsVO struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    ParentId int `gorm:"column:parent_id" json:"parentId"`
    Title null.String `gorm:"column:title" json:"title"`
    IsExecute int `gorm:"column:is_execute" json:"isExecute"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TeamsVO) TableName() string {
    return "teams"
}