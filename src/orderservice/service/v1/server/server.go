package serverV1

import (
	"context"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	authPBV1 "authservice/genproto/go/v1"
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
	authClient    authPBV1.AuthServiceClient
	productClient productPBV1.ProductServiceClient
}

// NewServer New service grpc server
func NewServer(
	logger log.Logger,
	conf *config.Config,
	repo *Repository,
	orderClient orderPBV1.OrderServiceClient,
	authClient authPBV1.AuthServiceClient,
	productClient productPBV1.ProductServiceClient,
) orderPBV1.OrderServiceServer {
	return &Server{
		repo:          repo,
		logger:        logger,
		conf:          conf,
		orderClient:   orderClient,
		authClient:    authClient,
		productClient: productClient,
	}
}

// Create 添加订单
func (s *Server) Create(ctx context.Context, request *orderPBV1.CreateRequest) (*emptypb.Empty, error) {

	_ = level.Error(s.logger).Log("msg", "执行了创建订单")
	resp := &emptypb.Empty{}
	err := s.repo.Create(request)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "添加订单失败，错误："+err.Error())
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	return resp, nil

}

// CreateRevert 添加订单失败补偿接口
func (s *Server) CreateRevert(ctx context.Context, request *orderPBV1.CreateRequest) (*emptypb.Empty, error) {

	_ = level.Error(s.logger).Log("msg", "执行了创建订单回滚")
	resp := &emptypb.Empty{}
	if err := s.repo.CreateRevert(ctx, request); err != nil {
		_ = level.Error(s.logger).Log("msg", "执行添加订单失败回滚失败，错误："+err.Error())
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	return resp, nil

}

// CreateSaga 添加订单事务接口
func (s *Server) CreateSaga(ctx context.Context, request *orderPBV1.CreateRequest) (*emptypb.Empty, error) {

	resp := &emptypb.Empty{}

	// 扣产品库存事务
	decreaseProductReq := &productPBV1.DecreaseStockRequest{}
	decreaseProductReq.Quantity = 1
	decreaseProductReq.Id = request.ProductId

	productDecreaseStock := "product-grpc.rgrpc-dev:50051/proto.product.v1.ProductService/DecreaseStock"
	productDecreaseStockRevert := "product-grpc.rgrpc-dev:50051/proto.product.v1.ProductService/DecreaseStockRevert"

	// 创建订单事务
	orderNo := uuid.NewString()
	request.OrderNo = orderNo
	createOrder := "order-grpc.rgrpc-dev:50051/proto.order.v1.OrderService/Create"
	createOrderRevert := "order-grpc.rgrpc-dev:50051/proto.order.v1.OrderService/CreateRevert"

	saga := dtmgrpc.NewSagaGrpc("dtm-svc.dtm-prod:36790", uuid.NewString()).
		Add(createOrder, createOrderRevert, request).
		Add(productDecreaseStock, productDecreaseStockRevert, decreaseProductReq)
	saga.WaitResult = true
	if err := saga.Submit(); err != nil {
		_ = level.Error(s.logger).Log("msg", "创建订单失败，错误："+err.Error())
		return nil, status.Error(codes.Aborted, "创建订单失败")
	}
	return resp, nil

}

// Detail 获取订单详情
func (s *Server) Detail(ctx context.Context, request *orderPBV1.DetailRequest) (*orderPBV1.Response, error) {

	resp := &orderPBV1.Response{}
	return resp, nil

}
