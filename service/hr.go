package service

import (
	"context"

	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage"
)

type HRService struct {
	storage storage.StorageI
	pb.UnimplementedHRServiceServer
}

func NewHRStorage(storage storage.StorageI) *HRService {
	return &HRService{
		storage: storage,
	}
}


func (s *HRService) Comment(c context.Context, req *pb.CommentReq) (*pb.Void, error) {
	return s.storage.HR().Comment(req)
}

