package serverV1

import (
	"context"
	"encoding/json"
	"time"

	Jgrpc_otelspan "github.com/janrs-io/Jgrpc-otel-span"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"

	"orderservice/config"
	orderPBV1 "orderservice/genproto/go/v1"
	"orderservice/service/model"
)

// Repository 数据仓库层
type Repository struct {
	mysqlDB *gorm.DB
	conf    *config.Config
	span    *Jgrpc_otelspan.OtelSpan
}

// NewRepository 实例化 Repository
func NewRepository(
	mysqlDB *gorm.DB,
	conf *config.Config,
	span *Jgrpc_otelspan.OtelSpan,
) *Repository {
	return &Repository{
		mysqlDB: mysqlDB,
		conf:    conf,
		span:    span,
	}
}

// OrderModel order 表模型
func (r *Repository) OrderModel() *gorm.DB {
	orderModel := &model.Order{}
	return r.mysqlDB.Table(orderModel.TableName())
}

// Create 添加订单
func (r *Repository) Create(ctx context.Context, request *orderPBV1.CreateRequest) error {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	jsonStr := protojson.Format(request)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return err
	}
	m["create_time"] = time.Now().Unix()
	m["update_time"] = time.Now().Unix()

	result := r.OrderModel().Create(m)
	if result.Error != nil {
		return r.span.Error(span, result.Error.Error())
	}

	if result.RowsAffected <= 0 {
		return r.span.Error(span, "添加失败")
	}
	return nil

}

// CreateRevert 添加订单失败补偿
func (r *Repository) CreateRevert(ctx context.Context, request *orderPBV1.CreateRequest) (err error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	m := make(map[string]any)
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
func (r *Repository) Detail(ctx context.Context, request *orderPBV1.DetailRequest) (*model.Order, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	user := &model.Order{}
	err := r.OrderModel().First(&user, request.Id).Error
	if err != nil {
		return user, r.span.Error(span, err.Error())
	}
	return user, err

}
