package postgres

import (
	"database/sql"
	"fmt"

	pb "gitlab.com/black_list/black_list/genproto/hr_service"
)

type HRService struct {
	db *sql.DB
}

func NewHRStorage(db *sql.DB) *HRService {
	return &HRService{
		db: db,
	}
}

func (s *HRService) Comment(req *pb.CommentReq) (*pb.Void, error) {
	if req.HrId != "HR" && req.HrId != "hr" {
		return nil, fmt.Errorf("only HR can write comment entries")
	}

	query := "UPDATE employee SET comment = $1 WHERE id = $2"
	_, err := s.db.Exec(query, req.Comment, req.EmployeeId)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}
