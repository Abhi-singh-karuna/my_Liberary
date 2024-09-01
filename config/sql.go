package config

type SQL struct {
	Host     string `validate:"required"`
	Database string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Port     string `validate:"required"`
}

func (sql *SQL) GetHost() string {
	return sql.Host
}

func (sql *SQL) GetDatabase() string {
	return sql.Database
}

func (sql *SQL) GetUser() string {
	return sql.User
}

func (sql *SQL) GetPassword() string {
	return sql.Password
}

func (sql *SQL) GetPort() string {
	return sql.Port
}
