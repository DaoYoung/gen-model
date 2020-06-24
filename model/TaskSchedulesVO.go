package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskSchedulesVO struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    TaskId int `gorm:"column:task_id" json:"taskId"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    DeadlineAfterTaskStart int `gorm:"column:deadline_after_task_start" json:"deadlineAfterTaskStart"`
    Type int `gorm:"column:type" json:"type"`
    RelateTaskId int `gorm:"column:relate_task_id" json:"relateTaskId"`
    Priority int `gorm:"column:priority" json:"priority"`
    UserIds string `gorm:"column:user_ids" json:"userIds"`
    IsNeedUploadDoc int `gorm:"column:is_need_upload_doc" json:"isNeedUploadDoc"`
    Subject string `gorm:"column:subject" json:"subject"`
    Status int `gorm:"column:status" json:"status"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskSchedulesVO) TableName() string {
    return "task_schedules"
}