package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.com/black_list/black_list/config"
	"gitlab.com/black_list/black_list/storage"
)

type Storage struct {
	db        *sql.DB
	AdminS    storage.AdminI
	HRS       storage.HRI
	EmployeeS storage.EmployeeI
	BlackListS storage.BlackListI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	adminS := NewAdminAndHRStorage(db)
	employeeS := NewEmployeeStorage(db)
	hrS := NewHRStorage(db)
	blacklistS := NewBlackListStorage(db)
	return &Storage{
		db:        db,
		AdminS:   	adminS,
		EmployeeS:     employeeS,
		HRS: hrS,
		BlackListS: blacklistS,
	}, nil
}
func (s *Storage) Admin() storage.AdminI {
	if s.AdminS == nil {
		s.AdminS = NewAdminAndHRStorage(s.db)
	}
	return s.AdminS
}

func (s *Storage) Employee() storage.EmployeeI {
	if s.EmployeeS == nil {
		s.EmployeeS = NewEmployeeStorage(s.db)
	}
	return s.EmployeeS
}

func (s *Storage) HR() storage.HRI {
	if s.HRS == nil {
		s.HRS = NewHRStorage(s.db)
	}
	return s.HRS
}

func (s *Storage) BlackList() storage.BlackListI {
	if s.BlackListS == nil {
		s.BlackListS = NewBlackListStorage(s.db)
	}
	return s.BlackListS
}
