package service

import (
	"context"

	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage"
)

type EmployeeService struct {
	storage storage.StorageI
	pb.UnimplementedEmployeesServiceServer
}

func NewEmployeeStorage(storage storage.StorageI) *EmployeeService {
	return &EmployeeService{
		storage: storage,
	}
}

func (s *EmployeeService) Create(c context.Context, req *pb.EmployeeCreate) (*pb.Void, error) {
	return s.storage.Employee().Create(req)
}
func (s *EmployeeService) Update(c context.Context, req *pb.EmployeeUpdateReq) (*pb.Void, error) {
	return s.storage.Employee().Update(req)
}
func (s *EmployeeService) Delete(c context.Context, id *pb.GetByIdReq) (*pb.Void, error) {
	return s.storage.Employee().Delete(id)
}

func (s *EmployeeService) Get(c context.Context, id *pb.GetByIdReq) (*pb.Employee, error) {
	return s.storage.Employee().Get(id)
}

func (s *EmployeeService) GetAll(c context.Context, req *pb.GetAllEmployeeReq) (*pb.GetAllEmployeeRes, error) {
	return s.storage.Employee().GetAll(req)
}
