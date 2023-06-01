package serverV1

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"orderservice/config"
	orderPBV1 "orderservice/genproto/go/v1"
	productPBV1 "productservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	orderPBV1.UnimplementedOrderServiceServer
	logger        log.Logger
	conf          *config.Config
	repo          *Repository
	orderClient   orderPBV1.OrderServiceClient
	productClient productPBV1.ProductServiceClient
}

// NewServer New service grpc server
func NewServer(
	logger log.Logger,
	conf *config.Config,
	repo *Repository,
	orderClient orderPBV1.OrderServiceClient,
	productClient productPBV1.ProductServiceClient,
) orderPBV1.OrderServiceServer {
	return &Server{
		repo:          repo,
		logger:        logger,
		conf:          conf,
		orderClient:   orderClient,
		productClient: productClient,
	}
}

// Create 添加订单
func (s *Server) Create(ctx context.Context, request *orderPBV1.CreateRequest) (*orderPBV1.Response, error) {

	err := s.repo.Create(ctx, request)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "执行了创建订单失败")
		return nil, status.Error(codes.Aborted, "创建订单失败，错误[1]："+err.Error())
	}
	return &orderPBV1.Response{}, nil

}

// CreateRevert 添加订单失败补偿接口
func (s *Server) CreateRevert(ctx context.Context, request *orderPBV1.CreateRequest) (*orderPBV1.Response, error) {

	if err := s.repo.CreateRevert(ctx, request); err != nil {
		_ = level.Error(s.logger).Log("msg", "执行添加订单失败回滚失败，错误："+err.Error())
		return nil, status.Error(codes.Aborted, "执行创建订单失败补偿回滚失败，错误[1]："+err.Error())
	}
	return &orderPBV1.Response{}, nil

}

// CreateSaga 添加订单事务接口
func (s *Server) CreateSaga(ctx context.Context, request *orderPBV1.CreateRequest) (*orderPBV1.Response, error) {

	// 扣产品库存事务
	decreaseProductReq := &productPBV1.DecreaseStockRequest{}
	decreaseProductReq.Quantity = 1
	decreaseProductReq.Id = request.ProductId

	productDecreaseStock := "product.rgrpc-dev:50051/proto.product.v1.ProductService/DecreaseStock"
	productDecreaseStockRevert := "product.rgrpc-dev:50051/proto.product.v1.ProductService/DecreaseStockRevert"

	// 创建订单事务
	orderNo := uuid.NewString()
	request.OrderNo = orderNo
	createOrder := "order.rgrpc-dev:50051/proto.order.v1.OrderService/Create"
	createOrderRevert := "order.rgrpc-dev:50051/proto.order.v1.OrderService/CreateRevert"

	saga := dtmgrpc.NewSagaGrpc("dtm-svc.dtm-prod:36790", uuid.NewString())
	saga.Add(createOrder, createOrderRevert, request)
	saga.Add(productDecreaseStock, productDecreaseStockRevert, decreaseProductReq)
	saga.WaitResult = true
	if err := saga.Submit(); err != nil {
		_ = level.Error(s.logger).Log("msg", "创建订单失败，错误："+err.Error())
		return nil, status.Error(codes.Aborted, "创建订单失败")
	}
	return &orderPBV1.Response{}, nil

}

// Detail 获取订单详情
func (s *Server) Detail(ctx context.Context, request *orderPBV1.DetailRequest) (*orderPBV1.Response, error) {

	resp := &orderPBV1.Response{}

	result, err := s.repo.Detail(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, nil
		}
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误："+err.Error())
		return nil, status.Errorf(codes.FailedPrecondition, "获取订单详情失败")
	}
	orderDetail := &orderPBV1.OrderDetail{}
	orderDetail.Id = result.ID
	orderDetail.OrderNo = result.OrderNo
	orderDetail.ProductId = result.ProductID
	orderDetail.PaymentType = result.PaymentType
	orderDetail.OrderStatus = result.OrderStatus
	orderDetail.UserId = result.UserID
	orderDetail.UpdateTime = result.UpdateTime
	orderDetail.CreateTime = result.CreateTime
	orderDetail.PayTime = result.PayTime
	orderDetail.Amount = result.Amount

	anyData, err := anypb.New(orderDetail)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "获取订单详情失败，错误："+err.Error())
		return nil, status.Errorf(codes.FailedPrecondition, "获取订单想详情失败")
	}
	resp.ProtoAnyData = anyData
	return resp, nil

}
