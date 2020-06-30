package handler

import "time"

type GenModel struct {
    Id    int    `gorm:"primary_key;auto_increment"`
    DbName       string    `gorm:"type:varchar(150);not null"`
    TableName     string `gorm:"type:varchar(150);not null"`
    ModelFieldName   string `gorm:"type:varchar(150);not null"`
    ModelFieldType    string  `gorm:"type:varchar(50);not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
} 
