package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type Notices struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    UserId null.Int `gorm:"column:user_id" json:"userId"`
    TriggerUserId null.Int `gorm:"column:trigger_user_id" json:"triggerUserId"`
    Title null.Int `gorm:"column:title" json:"title"`
    Content null.Int `gorm:"column:content" json:"content"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *Notices) TableName() string {
    return "notices"
}