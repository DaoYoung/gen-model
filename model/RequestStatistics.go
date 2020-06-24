package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type RequestStatistics struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    MergeId int `gorm:"column:merge_id" json:"mergeId"`
    DevEnv string `gorm:"column:dev_env" json:"devEnv"`
    SpendTime int `gorm:"column:spend_time" json:"spendTime"`
    SavingTime int `gorm:"column:saving_time" json:"savingTime"`
    MaxTime int `gorm:"column:max_time" json:"maxTime"`
    MaxTimeRequestId int `gorm:"column:max_time_request_id" json:"maxTimeRequestId"`
    RequestTimes int `gorm:"column:request_times" json:"requestTimes"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *RequestStatistics) TableName() string {
    return "request_statistics"
}