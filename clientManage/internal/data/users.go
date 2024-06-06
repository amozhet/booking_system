package data

import (
	"clientManage/internal/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Fname     string    `json:"fname"`
	Sname     string    `json:"sname"`
	Email     string    `json:"email"`
	Password  password  `json:"password_hash"`
	Activated bool      `json:"activated"`
	UserRole  string    `json:"user_role"`
	Version   int       `json:"version"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 71 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Fname != "", "fname", "must be provided")
	v.Check(len(user.Fname) <= 500, "fname", "must not be more than 500 bytes long")

	v.Check(user.Sname != "", "sname", "must be provided")
	v.Check(len(user.Sname) <= 500, "sname", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("Missing password hash for user")
	}

}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (fname, sname, email, password_hash, user_role, activated)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
		`

	args := []any{user.Fname, user.Sname, user.Email, user.Password.hash, user.UserRole, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, created_at, fname, sname, email, password_hash, user_role, activated, version
		FROM users
		WHERE id = $1
		`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Fname,
		&user.Sname,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, fname, sname, email, password_hash, user_role, activated, version
		FROM users
		WHERE email = $1
		`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Fname,
		&user.Sname,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET fname = $1, sname=$2, email = $3, password_hash = $4, activated = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
		`

	args := []any{
		user.Fname,
		user.Sname,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
		SELECT users.id, 
		       users.created_at, 
		       users.fname, 
		       users.sname, 
		       users.email, 
		       users.password_hash, 
		       users.activated, 
		       users.version
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Fname,
		&user.Sname,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (user *UserModel) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := user.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m UserModel) GetAll(role string, sort string) ([]*User, error) {
	query := `
		SELECT id, created_at, fname, sname, email, password_hash, user_role, activated, version
		FROM users`

	args := []interface{}{}
	if role != "" {
		query += " WHERE user_role = $1"
		args = append(args, role)
	}

	// Determine the sorting order
	switch sort {
	case "fname":
		query += " ORDER BY fname"
	case "sname":
		query += " ORDER BY sname"
	default:
		query += " ORDER BY id"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Fname,
			&user.Sname,
			&user.Email,
			&user.Password.hash,
			&user.UserRole,
			&user.Activated,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
