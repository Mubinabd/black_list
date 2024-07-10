package service

import (
	"context"

	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage"
)

type AdminAndHRService struct {
	storage storage.StorageI
	pb.UnimplementedAdminAndHRServiceServer
}

func NewAdminAndHRStorage(storage storage.StorageI) *AdminAndHRService {
	return &AdminAndHRService{
		storage: storage,
	}
}

func (s *AdminAndHRService) Create(c context.Context, req *pb.AdminAndHRCreate) (*pb.Void, error) {
	return s.storage.Admin().Create(req)
}
func (s *AdminAndHRService) Update(c context.Context, req *pb.AdminAndHRUpdateReq) (*pb.Void, error) {
	return s.storage.Admin().Update(req)
}
func (s *AdminAndHRService) Delete(c context.Context, id *pb.GetByIdReq) (*pb.Void, error) {
	return s.storage.Admin().Delete(id)
}

func (s *AdminAndHRService) Get(c context.Context, id *pb.GetByIdReq) (*pb.AdminAndHR, error) {
	return s.storage.Admin().Get(id)
}

func (s *AdminAndHRService) GetAll(c context.Context, req *pb.GetAllAdminAndHRReq) (*pb.GetAllAdminAndHRRes, error) {
	return s.storage.Admin().GetAll(req)
}

func (s *AdminAndHRService) Approve(c context.Context, req *pb.GetByIdReq) (*pb.Void, error) {
	return s.storage.Admin().Approve(req)
}