package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskTemplateSchedules struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    TemplateId int `gorm:"column:template_id" json:"templateId"`
    DeadlineAfterTaskStart int `gorm:"column:deadline_after_task_start" json:"deadlineAfterTaskStart"`
    Type int `gorm:"column:type" json:"type"`
    Priority int `gorm:"column:priority" json:"priority"`
    UserIds string `gorm:"column:user_ids" json:"userIds"`
    IsNeedUploadDoc int `gorm:"column:is_need_upload_doc" json:"isNeedUploadDoc"`
    Subject string `gorm:"column:subject" json:"subject"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskTemplateSchedules) TableName() string {
    return "task_template_schedules"
}