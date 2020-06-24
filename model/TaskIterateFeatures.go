package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskIterateFeatures struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    Title string `gorm:"column:title" json:"title"`
    Description string `gorm:"column:description" json:"description"`
    IterateId int `gorm:"column:iterate_id" json:"iterateId"`
    ModuleId int `gorm:"column:module_id" json:"moduleId"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    Supervisor int `gorm:"column:supervisor" json:"supervisor"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskIterateFeatures) TableName() string {
    return "task_iterate_features"
}