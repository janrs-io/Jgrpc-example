package model

import (
	"gorm.io/gorm"
)

// MigrateProductTable Migrate Product table
func MigrateProductTable(db *gorm.DB) {

	m := db.Migrator()
	if !m.HasTable(&Product{}) {
		if err := m.CreateTable(&Product{}); err != nil {
			panic("migrate Failed.[ERROR]=>create product table failed.")
		}
		db.Exec("ALTER TABLE `product` COMMENT 'product table'")
	}

}

// Product Product Table
type Product struct {
	// 主键 ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:int(10);unique;autoIncrement;comment:primary id"`
	// 产品名称
	Name string `json:"name" gorm:"column:name;type:varchar(255);default:'';not null;comment:产品名称"`
	// 产品价格
	Price float32 `json:"price" gorm:"column:price;type:decimal(10,4);default:0;not null;comment:产品价格"`
	// 产品简介
	Desc string `json:"desc" gorm:"column:desc;type:varchar(255);default:'';not null;comment:产品简介"`
	// 产品标题
	Title string `json:"title" gorm:"column:title;type:varchar(100);default:'';not null;comment:产品标题"`
	// 产品库存
	Stock int64 `json:"stock" gorm:"column:stock;type:int(10);default:0;not null;comment:产品库存"`
	// 是否禁用
	IsDisable int64 `json:"is_disable" gorm:"column:is_disable;type:tinyint(2);default:2;not null;comment:是否禁用[1=是2=否]"`
	// 添加时间 / 更新时间
	CreateTime int64 `json:"create_time" gorm:"column:create_time;type:int(10);default:0;comment:create time'"`
	UpdateTime int64 `json:"update_time" gorm:"column:update_time;type:int(10);default:0;comment:update time"`
}

// TableName Table Name
func (*Product) TableName() string {
	return "product"
}
