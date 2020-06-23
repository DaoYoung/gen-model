package model

import "time"
import "gopkg.in/guregu/null.v3"
type Merges struct {

id int `gorm:"column:id;primary_key" json:"Id"`
author_id int `gorm:"column:author_id" json:"AuthorId"`
title string `gorm:"column:title" json:"Title"`
dev_env string `gorm:"column:dev_env" json:"DevEnv"`
cache_type int `gorm:"column:cache_type" json:"CacheType"`
result_type int `gorm:"column:result_type" json:"ResultType"`
req_path string `gorm:"column:req_path" json:"ReqPath"`
req_method string `gorm:"column:req_method" json:"ReqMethod"`
status int `gorm:"column:status" json:"Status"`
created_at time.Time `gorm:"column:created_at" json:"CreatedAt"`
updated_at time.Time `gorm:"column:updated_at" json:"UpdatedAt"`
deleted_at null.Time `gorm:"column:deleted_at" json:"DeletedAt"`}

func (tc *Merges) TableName() string {
    return "merges"
}