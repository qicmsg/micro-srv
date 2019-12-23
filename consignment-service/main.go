package main

import (
	pb "consignment-service/proto/consignment"
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
	vesselProto "vessel-service/proto/vessel"
)

//
// 仓库接口
//
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
	GetAll() []*pb.Consignment                                   // 获取仓库中所有的货物
}

//
// 我们存放多批货物的仓库，实现了 IRepository 接口
//
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

//
// 定义微服务
//
type service struct {
	repo IRepository // 这里应该是IRepository接口类型更合理
	// consignment-service 作为客户端调用 vessel-service 的函数
	vesselClient vesselProto.VesselService
}

//
// 实现 consignment.pb.go 中的 ShippingServiceHandler 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {

	// 检查是否有适合的货轮
	vReq := &vesselProto.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	log.Printf("found vessel: %s\n", vResp.Vessel.Name)

	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	//*resp = pb.Response{Created: true, Consignment: consignment}
	resp.Created = true
	resp.Consignment = consignment
	return nil
}

// 获取目前所有托运的货物
// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	//*resp = pb.Response{Consignments: allConsignments}
	resp.Consignments = allConsignments
	return nil
}

func main() {
	// 会启动默认的Transport[http]来做服务间的调用
	/*server := micro.NewService(
		// 必须和 consignment.proto 中的 package 一致
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)*/

	// 用grpc协议做服务间调用
	server := grpc.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// 解析命令行参数
	server.Init()
	repo := &Repository{} //
	// 作为 vessel-service 的客户端
	vClient := vesselProto.NewVesselService("go.micro.srv.vessel", server.Client())

	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
