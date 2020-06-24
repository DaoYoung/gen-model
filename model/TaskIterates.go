package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskIterates struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    Title null.String `gorm:"column:title" json:"title"`
    StartLine time.Time `gorm:"column:start_line" json:"startLine"`
    Deadline time.Time `gorm:"column:deadline" json:"deadline"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskIterates) TableName() string {
    return "task_iterates"
}