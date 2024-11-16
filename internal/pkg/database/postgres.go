package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go-starter-kit/internal/server/config"
	"go.uber.org/zap"
	"net/url"
)

type Postgres struct {
	writeDB *sqlx.DB
	readDB  *sqlx.DB
}

type Conn interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type connectionInfo struct {
	Database string
	Host     string
	User     string
	Password string
	MaxOpen  int
	MaxIdle  int
}

func NewPostgres(conf *config.Config, logger *zap.Logger) (*Postgres, error) {
	masterInfo := connectionInfo{
		Database: conf.Connection.Postgresql.Master.DB,
		Host:     conf.Connection.Postgresql.Master.Host,
		User:     conf.Connection.Postgresql.Master.User,
		Password: conf.Connection.Postgresql.Master.Password,
		MaxOpen:  conf.Connection.Postgresql.Master.MaxOpen,
		MaxIdle:  conf.Connection.Postgresql.Master.MaxIdle,
	}

	slaveInfo := connectionInfo{
		Database: conf.Connection.Postgresql.Slave.DB,
		Host:     conf.Connection.Postgresql.Slave.Host,
		User:     conf.Connection.Postgresql.Slave.User,
		Password: conf.Connection.Postgresql.Slave.Password,
		MaxOpen:  conf.Connection.Postgresql.Slave.MaxOpen,
		MaxIdle:  conf.Connection.Postgresql.Slave.MaxIdle,
	}

	var readInfo connectionInfo
	switch conf.Connection.Postgresql.FixedReadInstance {
	case "master":
		readInfo = masterInfo
	case "slave":
		readInfo = slaveInfo
	default:

	}

	writeDB, err := connectPostgres(masterInfo)
	if err != nil {
		return nil, fmt.Errorf("can't not open write database connection: %w", err)
	}
	readDB, err := connectPostgres(readInfo)
	if err != nil {
		return nil, fmt.Errorf("can't not open read database connection: %w", err)
	}

	return &Postgres{
		writeDB: writeDB,
		readDB:  readDB,
	}, nil
}

func connectPostgres(inf connectionInfo) (*sqlx.DB, error) {
	source := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", url.QueryEscape(inf.User),
		url.QueryEscape(inf.Password), inf.Host, inf.Database)
	fmt.Println(source)
	conf, err := pgx.ParseConfig(source)
	if err != nil {
		return nil, fmt.Errorf("pgx parse config failed: %w", err)
	}

	db := stdlib.OpenDB(*conf)

	DB := sqlx.NewDb(db, "pgx")
	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("pig to database failed: %w", err)
	}

	DB.SetMaxOpenConns(inf.MaxOpen)
	DB.SetMaxIdleConns(inf.MaxIdle)

	return DB, err
}

func (p *Postgres) GetReadConnection(ctx context.Context) (Conn, error) {
	if transactionCtx, ok := ctx.Value(TransactionCtxKey).(*TransactionCtx); ok {
		conn := transactionCtx.Conn
		if conn != nil {
			return conn, nil
		}
	}

	if customSettingCtx, ok := ctx.Value(CustomSettingCtxKey).(*CustomSettingCtx); ok {
		if customSettingCtx.IsJobAfterTxCommit {
			return p.writeDB, nil
		}
	}
	return p.readDB, nil
}

func (p *Postgres) GetWriteConnection(ctx context.Context) (Conn, error) {
	if transactionCtx, ok := ctx.Value(TransactionCtxKey).(*TransactionCtx); ok {
		if transactionCtx == nil {
			transactionCtx.Mu.Lock()
			defer transactionCtx.Mu.Unlock()

			conn, err := p.writeDB.Beginx()
			if err != nil {
				return nil, fmt.Errorf("can't get database write connection: %w", err)
			}
			transactionCtx.Conn = conn
		}
		return transactionCtx.Conn, nil
	}
	return p.writeDB, nil
}

func (p *Postgres) Ping() error {
	if p.writeDB != nil {
		if err := p.writeDB.Ping(); err != nil {
			return fmt.Errorf("postgres write pig failed: %w", err)
		}
	}

	if p.readDB != nil {
		if err := p.readDB.Ping(); err != nil {
			return fmt.Errorf("postgres read pig failed: %w", err)
		}
	}
	return nil
}

func (p *Postgres) Shutdown() {
	if p.writeDB != nil {
		_ = p.writeDB.Close()
	}
	if p.readDB != nil {
		_ = p.readDB.Close()
	}
}

func (p *Postgres) InitCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, TransactionCtxKey, &TransactionCtx{})
	return ctx, cancel
}

func (p *Postgres) EndCtx(ctx context.Context, err error) error {
	if transactionCtx, ok := ctx.Value(TransactionCtxKey).(*TransactionCtx); ok {
		if tx := transactionCtx.Conn; tx != nil {
			if p := recover(); p != nil {
				return tx.Rollback()
			} else if err != nil {
				return tx.Rollback()
			} else {
				return tx.Commit()
			}
		}
	}
	return errors.New("TransactionCtx Not found")
}
