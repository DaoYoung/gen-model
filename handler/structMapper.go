package handler

import "time"

type GenModel struct {
    Id    int    `gorm:"primary_key;auto_increment;"`
    DbName       string    `gorm:"type:varchar(150);not null;comment:'数据库名';"`
    TableName     string `gorm:"type:varchar(150);not null;comment:'表名';"`
    ModelFieldName   string `gorm:"type:varchar(150);not null;comment:'golang struct 里的字段名';"`
    ModelFieldType    string  `gorm:"type:varchar(50);not null;comment:'golang struct 里的字段类型';"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
} 
