package models

import "time"

type Media struct {
	Id uint64 `json:"id" gorm:"column:id;bigint unsigned;PRI;AUTO_INCREMENT;not null"` 
	DiskType string `json:"diskType" gorm:"column:disk_type;varchar(20);not null"` //存储类型
	SrcType string `json:"srcType" gorm:"column:src_type;tinyint(1);not null"` //链接类型 1相对路径 2外链
	Src string `json:"src" gorm:"column:src;varchar(191);not null"` //资源链接
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;datetime"` 
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;datetime"` 
}
func (entity *Media) TableName() string {
	return "media"
}