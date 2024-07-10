package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	pb "gitlab.com/black_list/black_list/genproto/hr_service"
)

type AdminAndHRService struct {
	db *sql.DB
}

func NewAdminAndHRStorage(db *sql.DB) *AdminAndHRService {
	return &AdminAndHRService{
		db: db,
	}
}

func (s *AdminAndHRService) Create(req *pb.AdminAndHRCreate) (*pb.Void, error) {
	id := uuid.NewString()
	query := `
	INSERT INTO adminandhr 
		(id, username, password, role) 
	VALUES 
		($1, $2, $3, $4)`

	_, err := s.db.Exec(query, id, req.Username, req.Password, req.Role)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil

}
func (s *AdminAndHRService) Get(req *pb.GetByIdReq) (*pb.AdminAndHR, error) {
	query := `
	SELECT 
		id, username, password, role 
	FROM 
		adminandhr 
	WHERE 
		id = $1
	AND 
		deleted_at = 0`

	var resp pb.AdminAndHR
	err := s.db.QueryRow(query, req.Id).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.Role)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *AdminAndHRService) GetAll(req *pb.GetAllAdminAndHRReq) (*pb.GetAllAdminAndHRRes, error) {
	query := "SELECT id, username, password, role FROM adminandhr"
	var args []interface{}
	argCount := 1

	filter := []string{}

	if req.Role != "" {
		filter = append(filter, "role = $"+fmt.Sprint(argCount))
		args = append(args, req.Role)
		argCount++
	}

	if len(filter) > 0 {
		query += " WHERE " + strings.Join(filter, " AND ")
	}
	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += " LIMIT $" + fmt.Sprint(argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += " OFFSET $" + fmt.Sprint(argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp pb.GetAllAdminAndHRRes
	for rows.Next() {
		var admin pb.AdminAndHR
		if err := rows.Scan(&admin.Id, &admin.Username, &admin.Password, &admin.Role); err != nil {
			return nil, err
		}
		resp.Admins = append(resp.Admins, &admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *AdminAndHRService) Approve(req *pb.GetByIdReq) (*pb.Void, error) {
	var role string
	err := s.db.QueryRow("SELECT role FROM adminandhr WHERE id = $1", req.Id).Scan(&role)
	if err != nil {
		return nil, err
	}

	if role == "admin" {
		return nil, fmt.Errorf("only admin users can approve roles")
	}

	query := "UPDATE adminandhr SET status = 'approved' WHERE id = $1"
	_, err = s.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *AdminAndHRService) Update(req *pb.AdminAndHRUpdateReq) (*pb.Void, error) {
	query := "UPDATE adminandhr SET"
	var args []interface{}
	argIndex := 1

	if req.UpdateBody.Username != "" {
		query += " username = $" + strconv.Itoa(argIndex) + ","
		args = append(args, req.UpdateBody.Username)
		argIndex++
	}
	if req.UpdateBody.Password != "" {
		query += " password = $" + strconv.Itoa(argIndex) + ","
		args = append(args, req.UpdateBody.Password)
		argIndex++
	}
	if req.UpdateBody.Role != "" {
		query += " role = $" + strconv.Itoa(argIndex) + ","
		args = append(args, req.UpdateBody.Role)
		argIndex++
	}

	query += " updated_at = now() WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, req.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *AdminAndHRService) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	query := `
	UPDATE 
		adminandhr 
	SET 	
		deleted_at = EXTRACT(EPOCH FROM NOW()) 
	WHERE 
		id = $1`
	_, err := s.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
