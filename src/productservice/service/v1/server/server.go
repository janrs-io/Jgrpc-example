package serverV1

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	productPBV1 "productservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	productPBV1.UnimplementedProductServiceServer
	logger log.Logger
	repo   *Repository
}

// NewServer New service grpc server
func NewServer(
	logger log.Logger,
	repo *Repository,
) productPBV1.ProductServiceServer {
	return &Server{
		repo:   repo,
		logger: logger,
	}
}

// Create 添加产品
func (s *Server) Create(ctx context.Context, request *productPBV1.CreateRequest) (*productPBV1.Response, error) {

	if _, err := s.repo.Create(ctx, request); err != nil {
		_ = level.Info(s.logger).Log("msg", "添加产品失败，错误[1]："+err.Error())
		return nil, status.Error(codes.Aborted, "添加产品失败")
	}
	return &productPBV1.Response{}, nil

}

// Delete 删除产品
func (s *Server) Delete(ctx context.Context, request *productPBV1.DeleteRequest) (*productPBV1.Response, error) {

	if err := s.repo.Delete(ctx, request); err != nil {
		_ = level.Error(s.logger).Log("msg", "删除产品失败，错误[1]："+err.Error())
		return nil, status.Error(codes.Aborted, "删除失败")
	}
	return &productPBV1.Response{}, nil

}

// Detail 产品详情
func (s *Server) Detail(ctx context.Context, request *productPBV1.DetailRequest) (*productPBV1.Response, error) {

	resp := &productPBV1.Response{}
	detail, err := s.repo.Detail(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Msg = "数据不存在"
			return resp, nil
		}
		_ = level.Error(s.logger).Log("msg", "获取详情失败，错误[1]："+err.Error())
		return nil, status.Error(codes.Unknown, "获取详情失败")
	}

	pbDetail := &productPBV1.ProductDetail{}

	pbDetail.Id = detail.ID
	pbDetail.IsDisable = detail.IsDisable
	pbDetail.Title = detail.Title
	pbDetail.UpdateTime = detail.UpdateTime
	pbDetail.CreateTime = detail.CreateTime
	pbDetail.Price = detail.Price
	pbDetail.Stock = detail.Stock
	pbDetail.Name = detail.Name
	pbDetail.Desc = detail.Desc

	anyData, err := anypb.New(pbDetail)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取产品详情失败，错误[2]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "获取详情失败")
	}
	resp.ProtoAnyData = anyData
	return resp, nil

}

// List 获取产品列表数据
func (s *Server) List(ctx context.Context, request *productPBV1.ListRequest) (*productPBV1.Response, error) {

	resp := &productPBV1.Response{}
	list, count, err := s.repo.List(ctx, request)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取列表数据失败，错误[1]："+err.Error())
		return resp, status.Error(codes.FailedPrecondition, "获取列表失败")
	}
	var listSlice []*productPBV1.ProductDetail
	listResp := &productPBV1.ListResponse{}
	for _, v := range *list {
		detail := &productPBV1.ProductDetail{}
		detail.Name = v.Name
		detail.Id = v.ID
		detail.Desc = v.Desc
		detail.Stock = v.Stock
		detail.CreateTime = v.CreateTime
		detail.UpdateTime = v.UpdateTime
		detail.IsDisable = v.IsDisable
		detail.Title = v.Title
		detail.Price = v.Price
		listSlice = append(listSlice, detail)
	}
	listResp.Total = count
	listResp.List = listSlice
	anyData, err := anypb.New(listResp)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取产品列表失败，错误[2]："+err.Error())
		return nil, status.Error(codes.Unknown, "获取列表失败")
	}
	resp.ProtoAnyData = anyData
	return resp, nil

}

// Update 更新产品详情
func (s *Server) Update(ctx context.Context, request *productPBV1.UpdateRequest) (*productPBV1.Response, error) {

	if err := s.repo.Update(ctx, request); err != nil {
		_ = level.Error(s.logger).Log("msg", "更新产品失败，错误[1]："+err.Error())
		return nil, status.Error(codes.FailedPrecondition, "更新产品失败")
	}
	return &productPBV1.Response{}, nil

}

// DecreaseStock 减少库存操作
// 这个接口用于执行 saga 事务成功的时候调用
func (s *Server) DecreaseStock(ctx context.Context, request *productPBV1.DecreaseStockRequest) (*productPBV1.Response, error) {

	if err := s.repo.DecreaseStock(ctx, request.Id, request.Quantity); err != nil {
		_ = level.Error(s.logger).Log("msg", "减少库存失败，错误[1]："+err.Error())
		return nil, status.Error(codes.Aborted, "减少库存失败，错误[1]："+err.Error())
	}
	return &productPBV1.Response{}, nil

}

// DecreaseStockRevert 回滚库存操作
// 这个接口用于执行 saga 事务失败的时候调用
func (s *Server) DecreaseStockRevert(ctx context.Context, request *productPBV1.DecreaseStockRequest) (*productPBV1.Response, error) {

	if err := s.repo.IncreaseStock(ctx, request.Id, request.Quantity); err != nil {
		_ = level.Error(s.logger).Log("msg", "回滚库存失败，错误[1]："+err.Error())
		return nil, status.Error(codes.Aborted, "回滚库存失败，错误[1]："+err.Error())
	}
	return &productPBV1.Response{}, nil

}
