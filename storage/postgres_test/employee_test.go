package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage/postgres"
)

func TestGetEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewEmployeeStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "position", "department", "created_by"}).
		AddRow("test-id", "John", "Doe", "Developer", "IT", "admin")

	mock.ExpectQuery("SELECT id, first_name, last_name, position, department, created_by FROM employee").
		WithArgs(req.Id).
		WillReturnRows(rows)

	resp, err := s.Get(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-id", resp.Id)
	assert.Equal(t, "John", resp.FirstName)
	assert.Equal(t, "Doe", resp.LastName)
	assert.Equal(t, "Developer", resp.Position)
	assert.Equal(t, "IT", resp.Department)
	assert.Equal(t, "admin", resp.CreateBy)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllEmployees(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewEmployeeStorage(db)

	req := &pb.GetAllEmployeeReq{}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "position", "department", "comment", "created_by"}).
		AddRow("emp-id1", "John", "Doe", "Developer", "IT", "Good Employee", "admin").
		AddRow("emp-id2", "Jane", "Smith", "Manager", "HR", "Experienced Manager", "admin")

	mock.ExpectQuery("SELECT id, first_name, last_name, position, department, comment, created_by FROM employee").
		WillReturnRows(rows)

	resp, err := s.GetAll(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Employees, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewEmployeeStorage(db)

	req := &pb.EmployeeUpdateReq{
		Id: "test-id",
		UpdateBody: &pb.EmployeeCreate{
			FirstName: "test-first-name",
		},
	}

	mock.ExpectExec("UPDATE employee").
		WithArgs(req.UpdateBody.FirstName, req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Update(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewEmployeeStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	mock.ExpectExec("UPDATE employee SET deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Delete(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
