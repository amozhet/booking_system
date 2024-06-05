package model

import (
	"golang.org/x/crypto/bcrypt"
)

type Client struct {
	ID           int64  `json:"id"`
	Name         string `json:"fname"`
	Surname      string `json:"sname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"user_role"`
	Activated    bool   `json:"activated"`
	Version      int    `json:"-"`
}

func (c *Client) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.PasswordHash = string(hash)
	return nil
}

func (c *Client) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password))
}
