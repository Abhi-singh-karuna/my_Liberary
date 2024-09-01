package gateway

type SqlHandler interface {
	Exec(string, ...interface{}) (Result, error)
	Query(string, ...interface{}) (Row, error)
	Transaction(func() (interface{}, error)) (interface{}, error)
	MultiExec(string) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Row interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}
