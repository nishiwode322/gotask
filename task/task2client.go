package task

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "gotask/provinceCityService"
)

const (
	address = "localhost:60051"
)

func GetCityByProvince(client pb.ProvinceCityServiceClient, provinceList []*pb.SelectByKeyRequest) {
	stream, err := client.GetCityByProvince(context.Background())
	if err != nil {
		log.Fatalf("client call GetCityByProvince function err:%v", err)
	}
	// for sync
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive : %v", err)
			}
			//
			log.Printf("received response :%v", response.CityName)
		}
	}()
	for _, province := range provinceList {
		if err := stream.Send(province); err != nil {
			log.Fatalf("failed to send a select province name request: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func SetCityAndProvince(client pb.ProvinceCityServiceClient, cityAndProvinceList []*pb.SetKeyValueRequest) {
	stream, err := client.SetCityAndProvince(context.Background())
	if err != nil {
		log.Fatalf("client call SetCityAndProvince function err:%v", err)
	}
	for _, pair := range cityAndProvinceList {
		if err := stream.Send(pair); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("send %v err:%v", pair, err)
		}
	}
	//receive response
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv response err:%v", err)
	}
	log.Printf("response summary: %v", response)
}

func DeleteCity(client pb.ProvinceCityServiceClient, cityList []*pb.DeleteKeyRequest) {
	stream, err := client.DeleteCity(context.Background())
	if err != nil {
		log.Fatalf("client call DeleteCity function err:%v", err)
	}
	for _, city := range cityList {
		if err := stream.Send(city); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("send %v response err:%v", city, err)
		}
	}
	//receive response
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv response err:%v", err)
	}
	log.Printf("response summary: %v", response)
}

func RunGrpcClient() {
	// set up a connection to the RPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect server err , %v", err)
	}
	defer conn.Close()
}
