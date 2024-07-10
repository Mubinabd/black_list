package main

import (
	"log"
	"net"

	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	bl "gitlab.com/black_list/black_list/genproto/blacklist"
	"gitlab.com/black_list/black_list/service"
	"gitlab.com/black_list/black_list/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	liss, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterAdminAndHRServiceServer(s, service.NewAdminAndHRStorage(db))
	pb.RegisterEmployeesServiceServer(s, service.NewEmployeeStorage(db))
	pb.RegisterHRServiceServer(s, service.NewHRStorage(db))
	bl.RegisterBlackListServiceServer(s, service.NewBlackListStorage(db))

	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
