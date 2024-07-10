package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	pb "gitlab.com/black_list/black_list/genproto/hr_service"
)

type EmployeeService struct {
	db *sql.DB
}

func NewEmployeeStorage(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		db: db,
	}
}

func (s *EmployeeService) Create(req *pb.EmployeeCreate) (*pb.Void, error) {
	comment := " "
	id := uuid.NewString()
	query := `INSERT INTO employee (id, first_name, last_name, position, department, created_by,comment) 
			VALUES ($1, $2, $3, $4, $5, $6,$7)`

	_, err := s.db.Exec(query, id, req.FirstName, req.LastName, req.Position, req.Department, req.CreateBy,comment)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *EmployeeService) Get(req *pb.GetByIdReq) (*pb.Employee, error) {
	query := "SELECT id, first_name, last_name, position, department, created_by FROM employee WHERE id = $1 AND deleted_at  = 0"

	var resp pb.Employee
	err := s.db.QueryRow(query, req.Id).
		Scan(
			&resp.Id,
			&resp.FirstName,
			&resp.LastName,
			&resp.Position,
			&resp.Department,
			&resp.CreateBy)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (s *EmployeeService) GetAll(req *pb.GetAllEmployeeReq) (*pb.GetAllEmployeeRes, error) {
	query := "SELECT id, first_name, last_name, position, department, comment, created_by FROM employee"
	var args []interface{}
	argCount := 1

	filters := []string{}
	if req.Position != "" {
		filters = append(filters, "position = $"+fmt.Sprint(argCount))
		args = append(args, req.Position)
		argCount++
	}

	if req.Department != "" {
		filters = append(filters, "department = $"+fmt.Sprint(argCount))
		args = append(args, req.Department)
		argCount++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
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

	var resp pb.GetAllEmployeeRes
	for rows.Next() {
		var employee pb.Employee
		if err := rows.Scan(
			&employee.Id,
			&employee.FirstName,
			&employee.LastName,
			&employee.Position,
			&employee.Department,
			&employee.Comment,
			&employee.CreateBy,
		); err != nil {
			return nil, err
		}
		resp.Employees = append(resp.Employees, &employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}


func (s *EmployeeService) Update(req *pb.EmployeeUpdateReq) (*pb.Void, error) {
	query := "UPDATE employee"

	var updates []string
	var args []interface{}

	if req.UpdateBody.FirstName != "" {
		updates = append(updates, "first_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.FirstName)
	}
	if req.UpdateBody.LastName != "" {
		updates = append(updates, "last_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.LastName)
	}
	if req.UpdateBody.Position != "" {
		updates = append(updates, "position = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.Position)
	}
	if req.UpdateBody.Department != "" {
		updates = append(updates, "department = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.Department)
	}
	if req.UpdateBody.Comment != "" {
		updates = append(updates, "comment = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.Comment)
	}
	if req.UpdateBody.CreateBy != "" {
		updates = append(updates, "created_by = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UpdateBody.CreateBy)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	query += " SET " + strings.Join(updates, ", ") + ", updated_at = now() WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, req.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}



func (s *EmployeeService) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	query := `
	UPDATE 
		employee 
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
