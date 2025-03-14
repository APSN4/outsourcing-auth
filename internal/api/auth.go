package api

import "core/internal"

type Register struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password"`
	Photo        string `json:"photo"`
	Type         string `json:"type"`
}

type Client struct {
	*Register `json:"register"`
}

type ClientRegister struct {
	*Client `json:"user"`
}

type ResponseClientRegister struct {
	*internal.StatusResponse `json:"status"`
	*ResponseUser            `json:"user"`
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}
