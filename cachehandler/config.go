package cachehandler

type Redis struct {
	Host     string `validate:"required"`
	Database int `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Port     string `validate:"required"`
}

func (sql *Redis) GetHost() string {
	return sql.Host
}

func (sql *Redis) GetDatabase() int {
	return sql.Database
}

func (sql *Redis) GetUser() string {
	return sql.User
}

func (sql *Redis) GetPassword() string {
	return sql.Password
}

func (sql *Redis) GetPort() string {
	return sql.Port
}


// 	log.Infof("Host :-  %v   -- port  %v  ---  Pass %v  --- DB  %v", cfg.Redis.Write.Host, cfg.Redis.Write.Port, cfg.Redis.Write.Password, cfg.Redis.Write.Db)
