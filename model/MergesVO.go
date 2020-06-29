package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type MergesVO struct {
    Id int `gorm:"column:id" json:"id"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    Title string `gorm:"column:title" json:"title"`
    DevEnv string `gorm:"column:dev_env" json:"devEnv"`
    CacheType int `gorm:"column:cache_type" json:"cacheType"`
    ResultType int `gorm:"column:result_type" json:"resultType"`
    ReqPath string `gorm:"column:req_path" json:"reqPath"`
    ReqMethod string `gorm:"column:req_method" json:"reqMethod"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *MergesVO) TableName() string {
    return "merges"
}