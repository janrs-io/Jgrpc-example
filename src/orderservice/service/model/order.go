package model

import (
	"gorm.io/gorm"
)

// MigrateOrderTable 迁移 order 表
func MigrateOrderTable(db *gorm.DB) {

	m := db.Migrator()
	if !m.HasTable(&Order{}) {
		if err := m.CreateTable(&Order{}); err != nil {
			panic("migrate Failed.[ERROR]=>create order table failed.")
		}
		db.Exec("ALTER TABLE `order` COMMENT 'order table'")
	}

}

// Order 订单表
type Order struct {
	// 主键 ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:int(10);unique;autoIncrement;comment:主键id"`
	// 订单编号
	OrderNo string `json:"order_no" gorm:"column:order_no;type:varchar(255);default:'';not null;comment:订单编号"`
	// 支付方式
	PaymentType int64 `json:"payment_type" gorm:"column:payment_type;tinyint(2);default:0;not null;comment:支付方式"`
	// 支付状态
	PayStatus int64 `json:"pay_status" gorm:"column:pay_status;tinyint(2);default:0;not null;comment:支付状态"`
	// 支付时间
	PayTime int64 `json:"pay_time" gorm:"column:pay_time;type:int(10);default:0;not null;comment:支付时间"`
	// 用户 ID
	UserID int64 `json:"user_id" gorm:"column:user_id;type:int(10);default:0;not null;comment:用户id"`
	// 产品 ID
	ProductID int64 `json:"product_id" gorm:"column:product_id;type:int(10);default:0;not null;comment:用户id"`
	// 订单流程状态
	OrderStatus int64 `json:"order_status" gorm:"column:order_status;type:int(10);default:0;not null;comment:订单状态"`
	// 金额
	Amount float32 `json:"amount" gorm:"column:amount;type:decimal(10,4);default:0;not null;comment:订单金额"`
	// 添加时间 / 更新时间
	CreateTime int64 `json:"create_time" gorm:"column:create_time;type:int(10);default:0;comment:create time'"`
	UpdateTime int64 `json:"update_time" gorm:"column:update_time;type:int(10);default:0;comment:update time"`
}

// TableName 标名称
func (*Order) TableName() string {
	return "order"
}
