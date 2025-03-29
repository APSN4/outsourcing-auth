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

type ClientAuth struct {
	*ClientLogin `json:"user"`
}

type ClientLogin struct {
	*ClientToken `json:"login"`
}

type ClientToken struct {
	Token string `json:"token"`
}

type ResponseClientRegister struct {
	*internal.StatusResponse
	*ResponseUser `json:"user"`
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}
