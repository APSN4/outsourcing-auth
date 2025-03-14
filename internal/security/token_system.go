package security

type Token interface {
	CreateToken() string
	CheckToken() string
}
