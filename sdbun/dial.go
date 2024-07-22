package sdbun

import (
	"database/sql"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"strings"
)

type Address struct {
	// common
	Driver string `json:"driver" toml:"driver" yaml:"driver"`
	DSN    string `json:"dsn" toml:"dsn" yaml:"dsn"`
	Logger string `json:"logger" toml:"logger"`

	// options
	ConnMaxLifeTimeMS int64 `json:"conn_max_lifetime" toml:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTimeMS int64 `json:"conn_max_idle_time" toml:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	MaxIdleConns      int   `json:"max_idle_conns" toml:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns      int   `json:"max_open_conns" toml:"max_open_conns" yaml:"max_open_conns"`
}

var (
	ErrIllegalDriver = sderr.Sentinel("illegal driver")
)

func Dial(addr Address, opts ...bun.DBOption) (*bun.DB, error) {
	applyOptions := func(db *bun.DB, addr *Address) *bun.DB {
		if addr.ConnMaxLifeTimeMS > 0 {
			db.SetConnMaxLifetime(sdtime.Milliseconds(addr.ConnMaxLifeTimeMS))
		}
		if addr.ConnMaxIdleTimeMS > 0 {
			db.SetConnMaxIdleTime(sdtime.Milliseconds(addr.ConnMaxIdleTimeMS))
		}
		if addr.MaxIdleConns > 0 {
			db.SetMaxIdleConns(addr.MaxIdleConns)
		}
		if addr.MaxOpenConns > 0 {
			db.SetMaxOpenConns(addr.MaxOpenConns)
		}
		db.AddQueryHook(LoggerOf(addr.Logger))
		return db
	}

	switch strings.ToLower(addr.Driver) {
	case "mysql":
		sqldb, err := sql.Open("mysql", addr.DSN)
		if err != nil {
			return nil, sderr.Wrap(err, "open mysql error")
		}
		db := bun.NewDB(sqldb, mysqldialect.New(), opts...)
		return applyOptions(db, &addr), nil
	case "postgres":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(addr.DSN)))
		db := bun.NewDB(sqldb, pgdialect.New(), opts...)
		return applyOptions(db, &addr), nil
	case "sqlite":
		sqldb, err := sql.Open(sqliteshim.ShimName, addr.DSN)
		if err != nil {
			return nil, sderr.Wrap(err, "open sqlite error")
		}
		db := bun.NewDB(sqldb, sqlitedialect.New(), opts...)
		return applyOptions(db, &addr), nil
	default:
		return nil, ErrIllegalDriver
	}
}
