package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	pb "gitlab.com/black_list/black_list/genproto/blacklist"
)

type BlackListStorage struct {
	db *sql.DB
}

func NewBlackListStorage(db *sql.DB) *BlackListStorage {
	return &BlackListStorage{
		db: db,
	}
}

func (s *BlackListStorage) Create(req *pb.BlackListCreate) (*pb.Void, error) {

	if req.AddedBy != "HR" && req.AddedBy != "hr" {
		return nil, fmt.Errorf("only HR can create blacklist entries")
	}

	id := uuid.NewString()
	query := "INSERT INTO blacklist (id, employee_id, reason, added_by) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, id, req.EmployeId, req.Reason, req.AddedBy)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *BlackListStorage) Get(req *pb.GetByIdReq) (*pb.BlackListRes, error) {
	query := `
		SELECT 
			b.id,
			b.reason,
			b.added_by,
			e.id,
			e.first_name,
			e.last_name,
			e.position,
			e.department,
			e.created_by
		FROM 
			blacklist  as b
		JOIN 
			employee as e
		ON
			b.employee_id = e.id
		WHERE
			b.id = $1
		AND
			b.deleted_at = 0`

	var resp pb.BlackListRes
	resp.Employe = &pb.Employee{}
	err := s.db.QueryRow(query, req.Id).
		Scan(
			&resp.Id,
			&resp.Employe.Id,
			&resp.Reason,
			&resp.AddedBy,
			&resp.Employe.FirstName,
			&resp.Employe.LastName,
			&resp.Employe.Position,
			&resp.Employe.Department,
			&resp.Employe.CreateBy,
		)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BlackListStorage) GetAll(req *pb.GetAllBlackListReq) (*pb.GetAllBlackListRes, error) {
	query := "SELECT b.id, e.id, e.first_name, e.last_name, e.position, e.department, e.created_by, b.reason, b.added_by FROM blacklist AS b JOIN employee AS e ON b.employee_id = e.id"
	var args []interface{}
	argCount := 1

	res := []string{}
	if req.AddedBy != "" {
		res = append(res, "b.added_by = $"+fmt.Sprint(argCount))
		args = append(args, req.AddedBy)
		argCount++
	}

	if len(res) > 0 {
		query += " WHERE " + strings.Join(res, " AND ")
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

	var resp pb.GetAllBlackListRes
	for rows.Next() {
		var blackList pb.BlackListRes
		blackList.Employe = &pb.Employee{}
		if err := rows.Scan(
			&blackList.Id,
			&blackList.Employe.Id,
			&blackList.Employe.FirstName,
			&blackList.Employe.LastName,
			&blackList.Employe.Position,
			&blackList.Employe.Department,
			&blackList.Employe.CreateBy,
			&blackList.Reason,
			&blackList.AddedBy,
		); err != nil {
			return nil, err
		}
		resp.BlackLists = append(resp.BlackLists, &blackList)
	}

	return &resp, nil
}

func (s *BlackListStorage) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	query := `
		UPDATE 
			blacklist 
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
