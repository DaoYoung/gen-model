package handler

import "time"

type structMapper struct {
    Id    int    `gorm:"primary_key;auto_increment;"`
    DbName       string    `gorm:"type:varchar(150);not null;comment:'database name';"`
    TableName     string `gorm:"type:varchar(150);not null;comment:'table name';"`
    ModelFieldName   string `gorm:"type:varchar(150);not null;comment:'golang struct field name';"`
    ModelFieldType    string  `gorm:"type:varchar(50);not null;comment:'golang struct field type';"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
} 
