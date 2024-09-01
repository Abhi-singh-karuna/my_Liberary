package psqlhandler

type Config interface {
	GetHost() string
	GetDatabase() string
	GetUser() string
	GetPassword() string
	GetPort() string
}
