package main

import (
	"encoding/json"
	"github.com/micro/go-micro/service/grpc"
	"io/ioutil"
	"log"
	"os"

	pb "consignment-service/proto/consignment"
	"golang.org/x/net/context"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	// 会启动默认的[http]初始化
	//cmd.Init()

	server := grpc.NewService()
	server.Init()

	// Create new greeter client
	// http的注册方式
	//client := pb.NewShippingService("go.micro.srv.consignment", microclient.DefaultClient)
	// Create new greeter client
	// grpc的方式
	client := pb.NewShippingService("go.micro.srv.consignment", server.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.TODO(), consignment)
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
