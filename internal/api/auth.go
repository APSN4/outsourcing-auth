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
	*ClientLoginAttributes `json:"login"`
}

type ClientLoginAttributes struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type ClientToken struct {
	Token string `json:"token"`
}

type TokenAccess struct {
	User struct {
		Login struct {
			Token string `json:"token"`
		} `json:"login"`
	} `json:"user"`
}

type AccountInfo struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	Token    string `json:"token"`
	Type     string `json:"type"`
}

type ResponseSuccessAccess struct {
	*internal.StatusResponse
	*ResponseUser `json:"user"`
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type ResponseAccount struct {
	*internal.StatusResponse
	User struct {
		Account AccountInfo `json:"account"`
	} `json:"user"`
}
