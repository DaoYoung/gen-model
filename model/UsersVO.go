package model

import (
    "gopkg.in/guregu/null.v3"
    "time"
)

type UsersVO struct {
    Id int `gorm:"column:id;primary_key" json:"id"`
    GitId int `gorm:"column:git_id" json:"gitId"`
    Role int `gorm:"column:role" json:"role"`
    GitUsername string `gorm:"column:git_username" json:"gitUsername"`
    GitPassword string `gorm:"column:git_password" json:"gitPassword"`
    GitToken string `gorm:"column:git_token" json:"gitToken"`
    GitTokenName string `gorm:"column:git_token_name" json:"gitTokenName"`
    Email null.String `gorm:"column:email" json:"email"`
    Phone null.String `gorm:"column:phone" json:"phone"`
    TeamId int `gorm:"column:team_id" json:"teamId"`
    IsLeader int `gorm:"column:is_leader" json:"isLeader"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt null.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func (model *UsersVO) TableName() string {
    return "users"
}