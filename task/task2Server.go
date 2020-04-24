package task

import (
	"database/sql"
	"fmt"
	pb "gotask/provinceCityService"
	"io"

	u "gotask/util"
	"log"
	"net"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type Server struct {
	db   u.MysqlOps
	pool *redis.Pool
}

func (s *Server) GetCityByProvince(stream pb.ProvinceCityService_GetCityByProvinceServer) error {

	conn := s.pool.Get()
	defer conn.Close()

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result, err := redis.Strings(conn.Do("zrangebyscore", request.ProvinceName, 0, 0))
		if err != nil {
			return err
		}
		//change to utf-8 encode
		for i, v := range result {
			result[i], _ = u.DecodeGBK(v)
		}
		note := &pb.SelectByKeyReponse{CityName: result}
		if err := stream.Send(note); err != nil {
			return err
		}
	}
}

func (s *Server) SetCityAndProvince(stream pb.ProvinceCityService_SetCityAndProvinceServer) error {

	conn := s.pool.Get()
	defer conn.Close()

	if err := s.db.Ping(); err != nil {
		return err
	}

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SetKeyValueReponse{})
		}
		if err != nil {
			return err
		}
		//sync mysql
		sqlstring := fmt.Sprintf("select city_name from city where city_name = %s", request.CityName)
		row := s.db.QueryRow(sqlstring)
		var tempName string
		err = row.Scan(&tempName)
		if err == sql.ErrNoRows {
			//select province id
			sqlstring = fmt.Sprintf("select province_name from province where province_name = %s", request.ProvinceName)
			row = s.db.QueryRow(sqlstring)
			var pid int
			err = row.Scan(&pid)

			if err == nil {
				sqlstring = fmt.Sprintf("insert into city values (0,'%s',%v)", request.CityName, pid)
				_, err := s.db.Exec(sqlstring)
				if err != nil {
					return err
				}
			} else {
				if err == sql.ErrNoRows {
					sqlstring = fmt.Sprintf("insert into province values (0,'%s')", request.ProvinceName)
					_, err := s.db.Exec(sqlstring)
					if err != nil {
						log.Fatalf("%s,%v\n", sqlstring, err)
					}

					//
					sqlstring = fmt.Sprintf("select province_id from province where province_name = %s", request.ProvinceName)
					row = s.db.QueryRow(sqlstring)
					err = row.Scan(&pid)
					sqlstring = fmt.Sprintf("insert into city values (0,'%s',%v)", request.CityName, pid)
					_, err = s.db.Exec(sqlstring)
					if err != nil {
						log.Fatalf("%s,%v\n", sqlstring, err)
					}
				} else {
					// query error!
					return err
				}
			}
		}

		//sync redis
		gbkValue, _ := u.DecodeUTF(request.CityName)
		reply, err := conn.Do("zadd", request.ProvinceName, 0, string(gbkValue))
		log.Printf("reply=%#v,err=%v\n", reply, err)
	}
}

func (s *Server) DeleteCity(stream pb.ProvinceCityService_DeleteCityServer) error {
	conn := s.pool.Get()
	defer conn.Close()

	if err := s.db.Ping(); err != nil {
		return err
	}

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.DeleteKeyReponse{})
		}
		if err != nil {
			return err
		}
		//sync mysql
		sqlstring := fmt.Sprintf("select city_name,province_id from city where city_name = %s", request.KeyName)
		row := s.db.QueryRow(sqlstring)
		var name string
		var pid int
		err = row.Scan(&name, &pid)
		if err == nil {
			sqlstring = fmt.Sprintf("delete from city where city_name = %s ", request.KeyName)
			_, err := s.db.Exec(sqlstring)
			if err != nil {
				log.Fatalf("%s,%v\n", sqlstring, err)
			}
			//select province name
			sqlstring = fmt.Sprintf("select province_name from province where province_id = %v ", pid)
			row = s.db.QueryRow(sqlstring)
			err = row.Scan(&name)
			if err != nil {
				log.Fatalf("%s,%v\n", sqlstring, err)
			}
			//sync redis
			gbkValue, _ := u.DecodeUTF(request.KeyName)
			reply, err := conn.Do("zrem", name, string(gbkValue))
			log.Printf("reply=%#v,err=%v\n", reply, err)
		}
	}
}

func RunGrpcServer() {
	listener, err := net.Listen("tcp", ":60051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	d := u.MysqlOps{UserName: "root", PassWord: "123456", IP: "127.0.0.1", Port: "3306", DataBase: "provincecity"}
	p := redis.NewPool(func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") }, 100)
	pb.RegisterProvinceCityServiceServer(grpcServer, &Server{db: d, pool: p})
	// determine whether to use TLS
	grpcServer.Serve(listener)
}
