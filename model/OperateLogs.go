package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type OperateLogs struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    MergeId int `gorm:"column:merge_id" json:"mergeId"`
    OperateType int `gorm:"column:operate_type" json:"operateType"`
    Changes null.String `gorm:"column:changes" json:"changes"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *OperateLogs) TableName() string {
    return "operate_logs"
}