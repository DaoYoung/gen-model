package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type MergeRequestsVO struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    MergeId int `gorm:"column:merge_id" json:"mergeId"`
    ReqMethod string `gorm:"column:req_method" json:"reqMethod"`
    ReqUri string `gorm:"column:req_uri" json:"reqUri"`
    Weight int `gorm:"column:weight" json:"weight"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *MergeRequestsVO) TableName() string {
    return "merge_requests"
}