package serverV1

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	productPBV1 "productservice/genproto/go/v1"
	"productservice/service/model"
	"strconv"
	"time"
)

// Repository 数据仓库层
type Repository struct {
	db *gorm.DB
}

// NewRepository 实例化 Repository
func NewRepository(
	db *gorm.DB,
) *Repository {
	return &Repository{
		db: db,
	}
}

// ProductModel 获取 product 模型
func (r *Repository) ProductModel() *gorm.DB {
	productModel := &model.Product{}
	return r.db.Table(productModel.TableName())
}

// Create 添加产品
func (r *Repository) Create(request *productPBV1.CreateRequest) (bool, error) {

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
		return false, result.Error
	}
	return true, nil

}

// Update 更新产品
func (r *Repository) Update(request *productPBV1.UpdateRequest) error {

	jsonStr := protojson.Format(request)
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return err
	}
	// 只有 ID 字段，则不更新
	if len(m) <= 1 {
		return nil
	}
	m["update_time"] = time.Now().Unix()
	return r.ProductModel().Where("id = ?", request.Id).Updates(m).Error

}

// Detail 获取产品详情
func (r *Repository) Detail(request *productPBV1.DetailRequest) (*model.Product, error) {

	product := &model.Product{}
	if err := r.ProductModel().First(&product, request.Id).Error; err != nil {
		return nil, err
	}
	return product, nil

}

// Delete 删除产品
func (r *Repository) Delete(request *productPBV1.DeleteRequest) error {
	return r.ProductModel().Delete(&model.Product{}, request.Id).Error
}

// List 获取产品列表
func (r *Repository) List(request *productPBV1.ListRequest) (*[]model.Product, int64, error) {

	var products []model.Product
	var count int64

	model := r.ProductModel().Where("name LIKE ?", ""+request.Name+"%").Count(&count)
	err := model.Limit(100).Offset(0).Order("create_time DESC").Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return &products, count, nil

}

// DecreaseStock 减少库存
func (r *Repository) DecreaseStock(productId int64, quantity int64) error {

	return r.ProductModel().
		Where("id = ?", productId).
		Update("stock", gorm.Expr("stock - "+strconv.FormatInt(quantity, 10))).
		Error

}

// IncreaseStock 增加库存
func (r *Repository) IncreaseStock(productId int64, quantity int64) error {

	return r.ProductModel().
		Where("id = ?", productId).
		Update("stock", gorm.Expr("stock + "+strconv.FormatInt(quantity, 10))).
		Error

}
