package serverV1

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"

	orderPBV1 "orderservice/genproto/go/v1"
	"orderservice/service/model"
)

// Repository 数据仓库层
type Repository struct {
	db *gorm.DB
}

// NewRepository 实例化 Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// OrderModel order 表模型
func (r *Repository) OrderModel() *gorm.DB {
	orderModel := &model.Order{}
	return r.db.Table(orderModel.TableName())
}

// Create 添加订单
func (r *Repository) Create(request *orderPBV1.CreateRequest) error {

	jsonStr := protojson.Format(request)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return err
	}
	m["create_time"] = time.Now().Unix()
	m["update_time"] = time.Now().Unix()

	result := r.OrderModel().Create(m)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected <= 0 {
		return errors.New("添加失败")
	}
	return nil

}

// CreateRevert 添加订单失败补偿
func (r *Repository) CreateRevert(ctx context.Context, request *orderPBV1.CreateRequest) (err error) {

	m := make(map[string]any)
	m["order_status"] = 1000
	m["order_status"] = orderPBV1.OrderStatus_ORDER_STATUS_UNDEFINED
	m["update_time"] = time.Now().Unix()
	return r.OrderModel().Where("order_no = ?", request.OrderNo).Updates(m).Error

}

// Update 更新订单
func (r *Repository) Update(request *orderPBV1.UpdateRequest) error {
	return nil
}

// Delete 删除订单
func (r *Repository) Delete(request *orderPBV1.DetailRequest) error {
	return nil
}

// List 获取订单列表
func (r *Repository) List() {

}

// Detail 获取订单详情
func (r *Repository) Detail(request *orderPBV1.DetailRequest) (*model.Order, error) {

	user := &model.Order{}
	err := r.OrderModel().First(&user, request.Id).Error
	return user, err

}
