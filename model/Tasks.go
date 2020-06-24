package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type Tasks struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    TemplateId int `gorm:"column:template_id" json:"templateId"`
    Title string `gorm:"column:title" json:"title"`
    Description string `gorm:"column:description" json:"description"`
    ParentId int `gorm:"column:parent_id" json:"parentId"`
    ModuleId int `gorm:"column:module_id" json:"moduleId"`
    IterateId int `gorm:"column:iterate_id" json:"iterateId"`
    FeatureId int `gorm:"column:feature_id" json:"featureId"`
    AcceptId int `gorm:"column:accept_id" json:"acceptId"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    ScheduleId int `gorm:"column:schedule_id" json:"scheduleId"`
    Type int `gorm:"column:type" json:"type"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *Tasks) TableName() string {
    return "tasks"
}