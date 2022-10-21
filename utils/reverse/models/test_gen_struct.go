package models

import "time"

import "github.com/go-sql-driver/mysql"

type TestGenStruct struct {
	Id uint64 `json:"id" gorm:"column:id;bigint unsigned;PRI;AUTO_INCREMENT;not null"` 
	Name string `json:"name" gorm:"column:name;varchar(30);not null"` //用户名称
	Mobile string `json:"mobile" gorm:"column:mobile;varchar(24);not null"` //用户手机号
	Password string `json:"password" gorm:"column:password;varchar(191);not null"` //用户密码
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;datetime"` 
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;datetime"` 
	DeletedAt mysql.NullTime `json:"deletedAt" gorm:"column:deleted_at;datetime"` 
	TestedAt mysql.NullTime `json:"testedAt" gorm:"column:tested_at;datetime"` 
}
func (entity *TestGenStruct) TableName() string {
	return "test_gen_struct"
}