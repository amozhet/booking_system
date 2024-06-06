package model

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
