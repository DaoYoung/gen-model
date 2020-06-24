package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type TaskIterateFeatureResources struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    FeatureId int `gorm:"column:feature_id" json:"featureId"`
    Title null.String `gorm:"column:title" json:"title"`
    LinkUrl null.String `gorm:"column:link_url" json:"linkUrl"`
    AuthorId int `gorm:"column:author_id" json:"authorId"`
    DocId int `gorm:"column:doc_id" json:"docId"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *TaskIterateFeatureResources) TableName() string {
    return "task_iterate_feature_resources"
}