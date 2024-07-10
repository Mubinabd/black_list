package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage/postgres"
)

func TestCreateAdminAndHR(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	adminAndHRService := postgres.NewAdminAndHRStorage(db)

	req := &pb.AdminAndHRCreate{
		Username: "John Doe",
		Password: "string",
		Role:     "hr",
	}

	mock.ExpectExec("INSERT INTO adminandhr").
		WithArgs(sqlmock.AnyArg(), req.Username, req.Password, req.Role).
		WillReturnResult(sqlmock.NewResult(1, 1)) 

	resp, err := adminAndHRService.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewAdminAndHRStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role"}).
		AddRow(req.Id, "John Doe", "string", "hr")

	mock.ExpectQuery("SELECT id, username, password, role FROM adminandhr WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnRows(rows)

	resp, err := s.Get(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Id, resp.Id)
	assert.Equal(t, "John Doe", resp.Username)
	assert.Equal(t, "string", resp.Password)
	assert.Equal(t, "hr", resp.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewAdminAndHRStorage(db)

	req := &pb.GetAllAdminAndHRReq{}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role"}).
		AddRow("id1", "user1", "pass1", "hr").
		AddRow("id2", "user2", "pass2", "admin")

	mock.ExpectQuery("SELECT id, username, password, role FROM adminandhr").
		WillReturnRows(rows)

	resp, err := s.GetAll(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Admins, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprove(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewAdminAndHRStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	mock.ExpectQuery("SELECT role FROM adminandhr WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

	mock.ExpectExec("UPDATE adminandhr SET role = 'hr', status = 'approved' WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Approve(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewAdminAndHRStorage(db)

	req := &pb.AdminAndHRUpdateReq{
		Id: "eb1d7c2e-2331-4527-9a4a-995e1dc61270",
		UpdateBody: &pb.AdminAndHRCreate{
			Username: "Muhlisa",
			Password: "19166",
			Role:     "hr",
		},
	}

	mock.ExpectExec(`UPDATE adminandhr SET username = \$1, password = \$2, role = \$3, updated_at = now\(\) WHERE id = \$4`).
		WithArgs(req.UpdateBody.Username, req.UpdateBody.Password, req.UpdateBody.Role, req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Update(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewAdminAndHRStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	mock.ExpectExec("UPDATE adminandhr SET deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Delete(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
