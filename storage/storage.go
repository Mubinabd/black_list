package storage

import pb "gitlab.com/black_list/black_list/genproto/hr_service"
import bl "gitlab.com/black_list/black_list/genproto/blacklist"


type StorageI interface {
	Admin() AdminI
	Employee() EmployeeI
	HR() HRI
	BlackList() BlackListI
}

type AdminI interface {
	Create(req *pb.AdminAndHRCreate) (*pb.Void, error)
	Update(req *pb.AdminAndHRUpdateReq) (*pb.Void, error)
	Delete(id *pb.GetByIdReq) (*pb.Void, error)
	Get(req *pb.GetByIdReq) (*pb.AdminAndHR, error)
	GetAll(req *pb.GetAllAdminAndHRReq) (*pb.GetAllAdminAndHRRes, error)
	Approve(id *pb.GetByIdReq) (*pb.Void, error)
}

type EmployeeI interface {
	Create(req *pb.EmployeeCreate) (*pb.Void, error)
	Update(req *pb.EmployeeUpdateReq) (*pb.Void, error)
	Delete(id *pb.GetByIdReq) (*pb.Void, error)
	Get(id *pb.GetByIdReq) (*pb.Employee, error)
	GetAll(req *pb.GetAllEmployeeReq) (*pb.GetAllEmployeeRes, error)
}


type HRI interface {
	Comment(req *pb.CommentReq) (*pb.Void, error)
}

type BlackListI interface {
	Create(req *bl.BlackListCreate) (*bl.Void, error)
	Delete(id *bl.GetByIdReq) (*bl.Void, error)
	Get(id *bl.GetByIdReq) (*bl.BlackListRes, error)
	GetAll(req *bl.GetAllBlackListReq) (*bl.GetAllBlackListRes, error)
}