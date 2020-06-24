package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type MergeParams struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    MergeId int `gorm:"column:merge_id" json:"mergeId"`
    MergeField null.String `gorm:"column:merge_field" json:"mergeField"`
    MergeFieldSource int `gorm:"column:merge_field_source" json:"mergeFieldSource"`
    RequestId int `gorm:"column:request_id" json:"requestId"`
    RequestField string `gorm:"column:request_field" json:"requestField"`
    RequestFieldDefault null.String `gorm:"column:request_field_default" json:"requestFieldDefault"`
    RequestFieldSource int `gorm:"column:request_field_source" json:"requestFieldSource"`
    IsOptional int `gorm:"column:is_optional" json:"isOptional"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *MergeParams) TableName() string {
    return "merge_params"
}