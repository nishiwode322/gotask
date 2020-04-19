package task

import (
	pb "gotask/provinceCityService"
	"io"

	u "gotask/util"

	"github.com/garyburd/redigo/redis"
)

const (
	port = ":60051"
)

type Server struct {
	db u.MysqlOps
	
}

func (s *Server) GetCityByProvince(stream pb.ProvinceCityService_GetCityByProvinceServer) error {
	return nil
}

func (s *Server) SetCityAndProvince(stream pb.ProvinceCityService_SetCityAndProvinceServer) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SetKeyValueReponse{
				Reponse: true,
			})
		}
		if err != nil {
			return err
		}
		//sync mysql

		//sync redis

	}
	return nil
}

func (s *Server) DeleteCity(stream pb.ProvinceCityService_DeleteCityServer) error {
	return nil
}

func (s *Server) DeleteProvince(stream pb.ProvinceCityService_DeleteProvinceServer) error {
	return nil
}
