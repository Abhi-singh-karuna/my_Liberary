package sqlhandler

type Config interface {
	GetHost() string
	GetDatabase() string
	GetUser() string
	GetPassword() string
}
