package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	User        UserModel
	Token       TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:        UserModel{DB: db},
		Token:       TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}
