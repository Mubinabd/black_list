package service

import (
	"context"

	pb "gitlab.com/black_list/black_list/genproto/blacklist"
	"gitlab.com/black_list/black_list/storage"
)

type BlackListService struct {
	storage storage.StorageI
	pb.UnimplementedBlackListServiceServer
}

func NewBlackListStorage(storage storage.StorageI) *BlackListService {
	return &BlackListService{
		storage: storage,
	}
}

func (s *BlackListService) Create(c context.Context, req *pb.BlackListCreate) (*pb.Void, error) {
	return s.storage.BlackList().Create(req)
}
func (s *BlackListService) Delete(c context.Context, id *pb.GetByIdReq) (*pb.Void, error) {
	return s.storage.BlackList().Delete(id)
}

func (s *BlackListService) Get(c context.Context, id *pb.GetByIdReq) (*pb.BlackListRes, error) {
	return s.storage.BlackList().Get(id)
}

func (s *BlackListService) GetAll(c context.Context, req *pb.GetAllBlackListReq) (*pb.GetAllBlackListRes, error) {
	return s.storage.BlackList().GetAll(req)
}
