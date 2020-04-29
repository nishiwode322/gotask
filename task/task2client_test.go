package task

import (
	pb "gotask/provinceCityService"
	"log"
	"testing"

	"google.golang.org/grpc"
)

func TestGetCityByProvince(t *testing.T) {
	type args struct {
		client       pb.ProvinceCityServiceClient
		provinceList []*pb.SelectByKeyRequest
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect server err,", err)
	}
	defer conn.Close()
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "1", args: args{client: pb.NewProvinceCityServiceClient(conn), provinceList: []*pb.SelectByKeyRequest{&pb.SelectByKeyRequest{ProvinceName: "上海"}, &pb.SelectByKeyRequest{ProvinceName: "江苏"}}}},
		{name: "2", args: args{client: pb.NewProvinceCityServiceClient(conn), provinceList: []*pb.SelectByKeyRequest{&pb.SelectByKeyRequest{ProvinceName: "北京"}, &pb.SelectByKeyRequest{ProvinceName: "安徽"}}}},
		{name: "3", args: args{client: pb.NewProvinceCityServiceClient(conn), provinceList: []*pb.SelectByKeyRequest{&pb.SelectByKeyRequest{ProvinceName: "江苏"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCityByProvince(tt.args.client, tt.args.provinceList)
		})
	}
}

func TestSetCityAndProvince(t *testing.T) {
	type args struct {
		client              pb.ProvinceCityServiceClient
		cityAndProvinceList []*pb.SetKeyValueRequest
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect server err,", err)
	}
	defer conn.Close()
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "1", args: args{client: pb.NewProvinceCityServiceClient(conn), cityAndProvinceList: []*pb.SetKeyValueRequest{&pb.SetKeyValueRequest{CityName: "上海", ProvinceName: "上海"}, &pb.SetKeyValueRequest{CityName: "南京", ProvinceName: "江苏"}}}},
		{name: "2", args: args{client: pb.NewProvinceCityServiceClient(conn), cityAndProvinceList: []*pb.SetKeyValueRequest{&pb.SetKeyValueRequest{CityName: "天津", ProvinceName: "天津"}}}},
		{name: "3", args: args{client: pb.NewProvinceCityServiceClient(conn), cityAndProvinceList: []*pb.SetKeyValueRequest{&pb.SetKeyValueRequest{CityName: "六安", ProvinceName: "安徽"}, &pb.SetKeyValueRequest{CityName: "苏州", ProvinceName: "江苏"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCityAndProvince(tt.args.client, tt.args.cityAndProvinceList)
		})
	}
}

func TestDeleteCity(t *testing.T) {
	type args struct {
		client   pb.ProvinceCityServiceClient
		cityList []*pb.DeleteKeyRequest
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect server err,", err)
	}
	defer conn.Close()
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "1", args: args{client: pb.NewProvinceCityServiceClient(conn), cityList: []*pb.DeleteKeyRequest{&pb.DeleteKeyRequest{KeyName: "上海"}, &pb.DeleteKeyRequest{KeyName: "南京"}}}},
		{name: "1", args: args{client: pb.NewProvinceCityServiceClient(conn), cityList: []*pb.DeleteKeyRequest{&pb.DeleteKeyRequest{KeyName: "上海"}}}},
		{name: "3", args: args{client: pb.NewProvinceCityServiceClient(conn), cityList: []*pb.DeleteKeyRequest{&pb.DeleteKeyRequest{KeyName: "六安"}, &pb.DeleteKeyRequest{KeyName: "北京"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteCity(tt.args.client, tt.args.cityList)
		})
	}
}
