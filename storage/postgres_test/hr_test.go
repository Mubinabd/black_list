package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/black_list/black_list/genproto/hr_service"
	"gitlab.com/black_list/black_list/storage/postgres"
)


func TestCommentEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewHRStorage(db)

	req := &pb.CommentReq{
		HrId:       "HR",
		EmployeeId: "emp-id",
		Comment:    "Good employee, hard worker",
	}

	mock.ExpectExec(`UPDATE employee SET comment = \$1 WHERE id = \$2`).
		WithArgs(req.Comment, req.EmployeeId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err = s.Comment(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
