package sqlhandler

import (
	"database/sql"
	"strings"

	pvtconfig "github.com/Abhi-singh-karuna/my_Liberary/config"
	"github.com/Abhi-singh-karuna/my_Liberary/errs"
	"github.com/Abhi-singh-karuna/my_Liberary/gateway"
	"github.com/Abhi-singh-karuna/my_Liberary/logger"

	_ "github.com/go-sql-driver/mysql"
)

const (
	optionSingleStatement = "?parseTime=true&loc=UTC&multiStatements=false"
	optionMultiStatements = "?parseTime=true&loc=UTC&multiStatements=true"
)

type SqlHandler struct {
	log     logger.Logger
	DB      *sql.DB
	connect string
}

func newDB(connect, option string) (*sql.DB, error) {
	dbms := "mysql"
	connect = strings.Join([]string{connect, option}, "")
	return sql.Open(dbms, connect)
}

func newConnect(host, database, user, password string) string {
	return strings.Join([]string{user, ":", password, "@", "tcp(", host, ":3306)/", database}, "")
}

func NewSqlHandler(log logger.Logger, config pvtconfig.SQL) gateway.SqlHandler {
	host := config.GetHost()
	database := config.GetDatabase()
	user := config.GetUser()
	password := config.GetPassword()
	log.Debug("SqlHandler created variables from Config")

	connect := newConnect(host, database, user, password)
	db, err := newDB(connect, optionSingleStatement)
	if err != nil {
		log.Panic(err)
	}
	log.Debug("SqlHandler prepared connection to database in single statement mode")

	db.SetMaxIdleConns(300)
	db.SetMaxOpenConns(300)

	return &SqlHandler{
		log:     log,
		DB:      db,
		connect: connect,
	}
}

func NewMapSqlHandler(log logger.Logger, sqlconfigs map[string]pvtconfig.SQL) map[string]gateway.SqlHandler {
	if sqlconfigs == nil {
		return nil
	}

	sNumber := len(sqlconfigs)
	if sNumber <= 0 {
		return nil
	}

	mapSqlHandlers := make(map[string]gateway.SqlHandler, sNumber)

	for i, config := range sqlconfigs {
		host := config.GetHost()
		database := config.GetDatabase()
		user := config.GetUser()
		password := config.GetPassword()
		log.Debugf("SqlHandler created variables from Config for host:%s and dbtype:%s", host, i)

		connect := newConnect(host, database, user, password)
		db, err := newDB(connect, optionSingleStatement)
		if err != nil {
			log.Panic(err)
		}
		log.Debugf("SqlHandler prepared connection to database in single statement mode for host:%s", host)

		db.SetMaxIdleConns(300)
		db.SetMaxOpenConns(300)

		mapSqlHandlers[i] = &SqlHandler{
			log:     log,
			DB:      db,
			connect: connect,
		}
	}
	return mapSqlHandlers
}

func (handler *SqlHandler) MultiExec(multiStatements string) error {
	handler.log.Debug("Connect to MySQL Database in multi statement mode")
	db, err := newDB(handler.connect, optionMultiStatements)
	if err != nil {
		handler.log.Error(err)
		return err
	}
	defer db.Close()

	handler.log.Debug("Exec multi statements SQL")
	_, err = db.Exec(multiStatements)
	if err != nil {
		handler.log.Error(err)
	}
	return err
}

func (handler *SqlHandler) Exec(statement string, args ...interface{}) (gateway.Result, error) {
	handler.log.Debug("Prepare SQL statement for execution")
	stmt, err := handler.DB.Prepare(statement)
	if err != nil {
		handler.log.Error(err)
		return nil, errs.Failed.Wrap(err, err.Error())
	}
	defer stmt.Close()

	handler.log.Debug("Execute prepared SQL statement")
	res, err := stmt.Exec(args...)
	if err != nil {
		handler.log.Error(err)
		return nil, errs.Failed.Wrap(err, err.Error())
	}

	return &SqlResult{Result: res}, nil
}

func (handler *SqlHandler) Query(statement string, args ...interface{}) (gateway.Row, error) {
	handler.log.Debug("Prepare SQL statement for query")
	stmt, err := handler.DB.Prepare(statement)
	if err != nil {
		handler.log.Error(err)
		return nil, errs.Failed.Wrap(err, err.Error())
	}
	defer stmt.Close()

	handler.log.Debug("Query prepared SQL statement")
	rows, err := stmt.Query(args...)
	if err != nil {
		handler.log.Error(err)
		return nil, errs.Failed.Wrap(err, err.Error())
	}

	return &SqlRow{Rows: rows}, nil
}

func (handler *SqlHandler) Transaction(f func() (interface{}, error)) (interface{}, error) {
	handler.log.Debug("Begin SQL transaction")
	tx, err := handler.DB.Begin()
	if err != nil {
		handler.log.Error(err)
		return nil, errs.Failed.Wrap(err, err.Error())
	}

	v, err := f()
	if err != nil {
		handler.log.Error(err)
		handler.log.Warn("Rollback transaction")
		eRollback := tx.Rollback()
		if eRollback != nil {
			err = errs.Failed.New(err.Error())
			err = errs.Wrap(err, eRollback.Error())
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		handler.log.Error(err)
		handler.log.Warn("Rollback transaction")
		tx.Rollback()
		return nil, errs.Failed.Wrap(err, err.Error())
	}

	return v, nil
}

type SqlResult struct {
	Result sql.Result
}

func (r *SqlResult) LastInsertId() (int64, error) {
	res, err := r.Result.LastInsertId()
	if err != nil {
		return res, errs.Failed.Wrap(err, err.Error())
	}
	return res, nil
}

func (r *SqlResult) RowsAffected() (int64, error) {
	res, err := r.Result.RowsAffected()
	if err != nil {
		return res, errs.Failed.Wrap(err, err.Error())
	}
	return res, nil
}

type SqlRow struct {
	Rows *sql.Rows
}

func (r *SqlRow) Scan(dest ...interface{}) error {
	if err := r.Rows.Scan(dest...); err != nil {
		return errs.Failed.Wrap(err, err.Error())
	}
	return nil
}

func (r *SqlRow) Next() bool {
	return r.Rows.Next()
}

func (r SqlRow) Close() error {
	if err := r.Rows.Close(); err != nil {
		return errs.Failed.Wrap(err, err.Error())
	}
	return nil
}
