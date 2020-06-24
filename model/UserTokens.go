package model

import (
    "gopkg.in/guregu/null.v3"
)

type UserTokens struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    CreatedAt null.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt null.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
    UserId int `gorm:"column:user_id" json:"userId"`
    Token null.String `gorm:"column:token" json:"token"`
}

func (model *UserTokens) TableName() string {
    return "user_tokens"
}