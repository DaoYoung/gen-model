package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)
type Merges struct {
    Id int `gorm:"column:id;primary_key" json:"Id"`
    AuthorId int `gorm:"column:author_id" json:"AuthorId"`
    Title string `gorm:"column:title" json:"Title"`
    DevEnv string `gorm:"column:dev_env" json:"DevEnv"`
    CacheType int `gorm:"column:cache_type" json:"CacheType"`
    ResultType int `gorm:"column:result_type" json:"ResultType"`
    ReqPath string `gorm:"column:req_path" json:"ReqPath"`
    ReqMethod string `gorm:"column:req_method" json:"ReqMethod"`
    Status int `gorm:"column:status" json:"Status"`
    CreatedAt time.Time `gorm:"column:created_at" json:"CreatedAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"UpdatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"DeletedAt"`
}

func (model *Merges) TableName() string {
    return "merges"
}