package serverV1

import (
	"context"
	"encoding/json"
	Jgrpc_otelspan "github.com/janrs-io/Jgrpc-otel-span"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"productservice/config"
	productPBV1 "productservice/genproto/go/v1"
	"productservice/service/model"
	"strconv"
	"time"
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

// ProductModel 获取 product 模型
func (r *Repository) ProductModel() *gorm.DB {
	productModel := &model.Product{}
	return r.mysqlDB.Table(productModel.TableName())
}

// Create 添加产品
func (r *Repository) Create(ctx context.Context, request *productPBV1.CreateRequest) (bool, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	product := &model.Product{}
	product.Name = request.Name
	product.Desc = request.Desc
	product.Price = request.Price
	product.Stock = request.Stock
	product.Title = request.Title
	product.IsDisable = request.IsDisable
	product.CreateTime = time.Now().Unix()
	product.UpdateTime = time.Now().Unix()
	result := r.ProductModel().Create(&product)
	if result.Error != nil {
		return false, r.span.Error(span, result.Error.Error())
	}
	return true, nil

}

// Update 更新产品
func (r *Repository) Update(ctx context.Context, request *productPBV1.UpdateRequest) error {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	jsonStr := protojson.Format(request)
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return r.span.Error(span, err.Error())
	}
	// 只有 ID 字段，则不更新
	if len(m) <= 1 {
		return nil
	}
	m["update_time"] = time.Now().Unix()
	if err := r.ProductModel().Where("id = ?", request.Id).Updates(m); err != nil {
		return r.span.Error(span, err.Error.Error())
	}
	return nil

}

// Detail 获取产品详情
func (r *Repository) Detail(ctx context.Context, request *productPBV1.DetailRequest) (*model.Product, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	product := &model.Product{}
	if err := r.ProductModel().First(&product, request.Id).Error; err != nil {
		return nil, r.span.Error(span, err.Error())
	}
	return product, nil

}

// Delete 删除产品
func (r *Repository) Delete(ctx context.Context, request *productPBV1.DeleteRequest) error {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()
	if err := r.ProductModel().Delete(&model.Product{}, request.Id).Error; err != nil {
		return r.span.Error(span, err.Error())
	}
	return nil
}

// List 获取产品列表
func (r *Repository) List(ctx context.Context, request *productPBV1.ListRequest) (*[]model.Product, int64, error) {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	var products []model.Product
	var count int64

	result := r.ProductModel().Where("name LIKE ?", ""+request.Name+"%").Count(&count)
	err := result.Limit(100).Offset(0).Order("create_time DESC").Find(&products).Error

	if err != nil {
		return nil, 0, r.span.Error(span, err.Error())
	}

	return &products, count, nil

}

// DecreaseStock 减少库存
func (r *Repository) DecreaseStock(ctx context.Context, productId int64, quantity int64) error {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	if err := r.ProductModel().
		Where("id = ?", productId).
		Update("stock", gorm.Expr("stock -"+strconv.FormatInt(quantity, 10))).
		Error; err != nil {
		return r.span.Error(span, err.Error())
	}
	return nil

}

// IncreaseStock 增加库存
func (r *Repository) IncreaseStock(ctx context.Context, productId int64, quantity int64) error {

	_, span := r.span.Record(ctx, r.conf.Trace.TracerName)
	defer span.End()

	if err := r.ProductModel().
		Where("id = ?", productId).
		Update("stock", gorm.Expr("stock + "+strconv.FormatInt(quantity, 10))).
		Error; err != nil {
		return r.span.Error(span, err.Error())
	}
	return nil

}
