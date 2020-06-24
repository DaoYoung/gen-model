package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskTemplates struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    Title string `gorm:"column:title" json:"title"`
    AcceptId int `gorm:"column:accept_id" json:"acceptId"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskTemplates) TableName() string {
    return "task_templates"
}