package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskModulesVO struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    Title string `gorm:"column:title" json:"title"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    ParentId int `gorm:"column:parent_id" json:"parentId"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskModulesVO) TableName() string {
    return "task_modules"
}