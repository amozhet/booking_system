package repository

import (
	"clientManage/internal/domain/model"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrEditConflict = errors.New("edit conflict")
)

type ClientRepository interface {
	CreateClient(client *model.Client) error
	GetClientByID(id int64) (*model.Client, error)
	GetClientByEmail(email string) (*model.Client, error)

	UpdateClient(client *model.Client) error
	DeleteClient(id int64) error
	ListClients(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Client, error)
}

type ClientRepositoryImpl struct {
	DB *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepositoryImpl {
	return &ClientRepositoryImpl{DB: db}
}

func (r *ClientRepositoryImpl) CreateClient(client *model.Client) error {
	query := `
INSERT INTO clientdb (fname, sname, email, password_hash, user_role, activated)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, version`
	return r.DB.QueryRow(query, client.Name, client.Surname, client.Email, client.PasswordHash, client.Role, client.Activated).Scan(&client.ID, &client.Version)
}

func (r *ClientRepositoryImpl) GetClientByID(id int64) (*model.Client, error) {
	query := `SELECT id, fname, sname, email, password_hash, user_role, activated, version
FROM clientdb WHERE id = $1`
	var client model.Client
	err := r.DB.QueryRow(query, id).Scan(&client.ID, &client.Name, &client.Surname, &client.Email, &client.PasswordHash, &client.Role, &client.Activated, &client.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepositoryImpl) GetClientByEmail(email string) (*model.Client, error) {
	query := "SELECT id, fname, sname, email, password_hash, user_role, activated, version FROM clientdb WHERE email = $1"

	var client model.Client
	err := r.DB.QueryRow(query, email).Scan(&client.ID, &client.Name, &client.Surname, &client.Email, &client.PasswordHash, &client.Role, &client.Activated, &client.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepositoryImpl) UpdateClient(client *model.Client) error {
	query := `UPDATE clientdb
SET fname = $1, sname = $2, email = $3, password_hash = $4, user_role=$5, activated = $6, version = version + 1
WHERE id = $7 AND version = $8
RETURNING version`
	args := []any{client.Name, client.Surname, client.Email, client.PasswordHash, client.Role, client.Activated, client.ID, client.Version}
	err := r.DB.QueryRow(query, args...).Scan(&client.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (r *ClientRepositoryImpl) DeleteClient(id int64) error {
	query := `DELETE FROM clientdb WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *ClientRepositoryImpl) ListClients(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Client, error) {
	query := `SELECT id, fname, sname, email, password_hash, user_role, activated, version 
FROM clientdb`
	var whereClauses []string
	var args []interface{}
	i := 1

	for key, value := range filters {
		whereClauses = append(whereClauses, key+" = $"+strconv.Itoa(i))
		args = append(args, value)
		i++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY " + sortBy + " " + sortOrder
	query += " LIMIT $" + strconv.Itoa(i) + " OFFSET $" + strconv.Itoa(i+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*model.Client
	for rows.Next() {
		var client model.Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Surname, &client.Email, &client.Role, &client.Activated, &client.Version); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}
