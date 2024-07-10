package postgres_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/black_list/black_list/genproto/blacklist"
	"gitlab.com/black_list/black_list/storage/postgres"
)

func TestCreateBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blackListStorage := postgres.NewBlackListStorage(db)

	req := &pb.BlackListCreate{
		EmployeId: "employee-123",
		Reason:    "Violation of policy",
		AddedBy:   "admin",
	}

	mock.ExpectExec("INSERT INTO blacklist").
		WithArgs(sqlmock.AnyArg(), req.EmployeId, req.Reason, req.AddedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := blackListStorage.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.NoError(t, mock.ExpectationsWereMet())

	reqInvalid := &pb.BlackListCreate{
		EmployeId: "employee-456",
		Reason:    "Policy breach",
		AddedBy:   "HR",
	}

	resp, err = blackListStorage.Create(reqInvalid)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "only HR can create blacklist entries", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := postgres.NewBlackListStorage(db)

	req := &pb.GetByIdReq{Id: "test-id"}

	mock.ExpectExec("UPDATE blacklist SET deleted_at = EXTRACT\\(EPOCH FROM NOW\\(\\)\\) WHERE id = \\$1").
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = s.Delete(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blackListStorage := postgres.NewBlackListStorage(db)

	t.Run("Successful retrieval without filters", func(t *testing.T) {
		req := &pb.GetAllBlackListReq{}

		rows := sqlmock.NewRows([]string{"b.id", "e.id", "e.first_name", "e.last_name", "e.position", "e.department", "e.created_by", "b.reason", "b.added_by"}).
			AddRow("1", "101", "John", "Doe", "Engineer", "IT", "admin", "Policy Violation", "admin").
			AddRow("2", "102", "Jane", "Smith", "Manager", "HR", "hr", "Performance Issue", "admin")

		mock.ExpectQuery(`^SELECT b\.id, e\.id, e\.first_name, e\.last_name, e\.position, e\.department, e\.created_by, b\.reason, b\.added_by FROM blacklist AS b JOIN employee AS e ON b\.employee_id = e\.id$`).
			WillReturnRows(rows)

		resp, err := blackListStorage.GetAll(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.BlackLists, 2)

		assert.Equal(t, "1", resp.BlackLists[0].Id)
		assert.Equal(t, "101", resp.BlackLists[0].Employe.Id)
		assert.Equal(t, "John", resp.BlackLists[0].Employe.FirstName)
		assert.Equal(t, "Doe", resp.BlackLists[0].Employe.LastName)
		assert.Equal(t, "Engineer", resp.BlackLists[0].Employe.Position)
		assert.Equal(t, "IT", resp.BlackLists[0].Employe.Department)
		assert.Equal(t, "admin", resp.BlackLists[0].Employe.CreateBy)
		assert.Equal(t, "Policy Violation", resp.BlackLists[0].Reason)
		assert.Equal(t, "admin", resp.BlackLists[0].AddedBy)

		assert.Equal(t, "2", resp.BlackLists[1].Id)
		assert.Equal(t, "102", resp.BlackLists[1].Employe.Id)
		assert.Equal(t, "Jane", resp.BlackLists[1].Employe.FirstName)
		assert.Equal(t, "Smith", resp.BlackLists[1].Employe.LastName)
		assert.Equal(t, "Manager", resp.BlackLists[1].Employe.Position)
		assert.Equal(t, "HR", resp.BlackLists[1].Employe.Department)
		assert.Equal(t, "hr", resp.BlackLists[1].Employe.CreateBy)
		assert.Equal(t, "Performance Issue", resp.BlackLists[1].Reason)
		assert.Equal(t, "admin", resp.BlackLists[1].AddedBy)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Successful retrieval with AddedBy filter", func(t *testing.T) {
		req := &pb.GetAllBlackListReq{AddedBy: "admin"}

		rows := sqlmock.NewRows([]string{"b.id", "e.id", "e.first_name", "e.last_name", "e.position", "e.department", "e.created_by", "b.reason", "b.added_by"}).
			AddRow("1", "101", "John", "Doe", "Engineer", "IT", "admin", "Policy Violation", "admin")

		mock.ExpectQuery(`^SELECT b\.id, e\.id, e\.first_name, e\.last_name, e\.position, e\.department, e\.created_by, b\.reason, b\.added_by FROM blacklist AS b JOIN employee AS e ON b\.employee_id = e\.id WHERE b\.added_by = \$1$`).
			WithArgs("admin").
			WillReturnRows(rows)

		resp, err := blackListStorage.GetAll(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.BlackLists, 1)

		assert.Equal(t, "1", resp.BlackLists[0].Id)
		assert.Equal(t, "101", resp.BlackLists[0].Employe.Id)
		assert.Equal(t, "John", resp.BlackLists[0].Employe.FirstName)
		assert.Equal(t, "Doe", resp.BlackLists[0].Employe.LastName)
		assert.Equal(t, "Engineer", resp.BlackLists[0].Employe.Position)
		assert.Equal(t, "IT", resp.BlackLists[0].Employe.Department)
		assert.Equal(t, "admin", resp.BlackLists[0].Employe.CreateBy)
		assert.Equal(t, "Policy Violation", resp.BlackLists[0].Reason)
		assert.Equal(t, "admin", resp.BlackLists[0].AddedBy)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Successful retrieval with Limit and Offset filters", func(t *testing.T) {
		req := &pb.GetAllBlackListReq{
			Filter: &pb.Filter{
				Limit:  1,
				Offset: 1,
			},
		}

		rows := sqlmock.NewRows([]string{"b.id", "e.id", "e.first_name", "e.last_name", "e.position", "e.department", "e.created_by", "b.reason", "b.added_by"}).
			AddRow("2", "102", "Jane", "Smith", "Manager", "HR", "hr", "Performance Issue", "admin")

		mock.ExpectQuery(`^SELECT b\.id, e\.id, e\.first_name, e\.last_name, e\.position, e\.department, e\.created_by, b\.reason, b\.added_by FROM blacklist AS b JOIN employee AS e ON b\.employee_id = e\.id LIMIT \$1 OFFSET \$2$`).
			WithArgs(1, 1).
			WillReturnRows(rows)

		resp, err := blackListStorage.GetAll(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.BlackLists, 1)

		assert.Equal(t, "2", resp.BlackLists[0].Id)
		assert.Equal(t, "102", resp.BlackLists[0].Employe.Id)
		assert.Equal(t, "Jane", resp.BlackLists[0].Employe.FirstName)
		assert.Equal(t, "Smith", resp.BlackLists[0].Employe.LastName)
		assert.Equal(t, "Manager", resp.BlackLists[0].Employe.Position)
		assert.Equal(t, "HR", resp.BlackLists[0].Employe.Department)
		assert.Equal(t, "hr", resp.BlackLists[0].Employe.CreateBy)
		assert.Equal(t, "Performance Issue", resp.BlackLists[0].Reason)
		assert.Equal(t, "admin", resp.BlackLists[0].AddedBy)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error on query execution", func(t *testing.T) {
		req := &pb.GetAllBlackListReq{}

		mock.ExpectQuery(`^SELECT b\.id, e\.id, e\.first_name, e\.last_name, e\.position, e\.department, e\.created_by, b\.reason, b\.added_by FROM blacklist AS b JOIN employee AS e ON b\.employee_id = e\.id$`).
			WillReturnError(fmt.Errorf("query error"))

		resp, err := blackListStorage.GetAll(req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "query error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
