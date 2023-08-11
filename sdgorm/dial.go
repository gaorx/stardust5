package sdgorm

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/driver/sqlite"

	"github.com/gaorx/stardust5/sderr"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Address struct {
	// common
	Driver string `json:"driver" toml:"driver" yaml:"driver"`
	DSN    string `json:"dsn" toml:"dsn" yaml:"dsn"`
	Logger string `json:"logger" toml:"logger"`

	// mysql
	MySqlConn                      gorm.ConnPool `json:"-" toml:"-"`
	MySqlSkipInitializeWithVersion bool          `json:"mysql_skip_initialize_with_version" toml:"mysql_skip_initialize_with_version" yaml:"mysql_skip_initialize_with_version"`
	MySqlDefaultStringSize         uint          `json:"mysql_default_string_size" toml:"mysql_default_string_size" yaml:"mysql_default_string_size"`
	MySqlDefaultDatetimePrecision  *int          `json:"mysql_default_datetime_precision" toml:"mysql_default_datetime_precision" yaml:"mysql_default_datetime_precision"`
	MySqlDisableDatetimePrecision  bool          `json:"mysql_disable_datetime_precision" toml:"mysql_disable_datetime_precision" yaml:"mysql_disable_datetime_precision"`
	MySqlDontSupportRenameIndex    bool          `json:"mysql_dont_support_rename_index" toml:"mysql_dont_support_rename_index" yaml:"mysql_dont_support_rename_index"`
	MySqlDontSupportRenameColumn   bool          `json:"mysql_dont_support_rename_column" toml:"mysql_dont_support_rename_column" yaml:"mysql_dont_support_rename_column"`
	MySqlDontSupportForShareClause bool          `json:"mysql_dont_support_for_share_clause" toml:"mysql_dont_support_for_share_clause" yaml:"mysql_dont_support_for_share_clause"`

	// postgres
	PostgresConn                 *sql.DB `json:"-" toml:"-"`
	PostgresPreferSimpleProtocol bool    `json:"postgres_prefer_simple_protocol" toml:"postgres_prefer_simple_protocol" yaml:"postgres_prefer_simple_protocol"`
	PostgresWithoutReturning     bool    `json:"postgres_without_returning" toml:"postgres_without_returning" yaml:"postgres_without_returning"`
}

var (
	ErrIllegalDriver = errors.New("illegal driver")
)

func Dial(addr Address, config *gorm.Config) (*gorm.DB, error) {
	if config == nil {
		config = &gorm.Config{}
	}
	if config.Logger == nil {
		config.Logger = LoggerOf(addr.Logger)
	}
	switch strings.ToLower(addr.Driver) {
	case "mysql":
		mysqlConfig := mysql.Config{
			DSN:                       addr.DSN,
			Conn:                      addr.MySqlConn,
			SkipInitializeWithVersion: addr.MySqlSkipInitializeWithVersion,
			DefaultStringSize:         addr.MySqlDefaultStringSize,
			DefaultDatetimePrecision:  addr.MySqlDefaultDatetimePrecision,
			DisableDatetimePrecision:  addr.MySqlDisableDatetimePrecision,
			DontSupportRenameIndex:    addr.MySqlDontSupportRenameIndex,
			DontSupportRenameColumn:   addr.MySqlDontSupportRenameColumn,
			DontSupportForShareClause: addr.MySqlDontSupportForShareClause,
		}
		db, err := gorm.Open(mysql.New(mysqlConfig), config)
		if err != nil {
			return nil, sderr.Wrap(err, "open mysql error")
		}
		return db, nil
	case "postgres":
		postgresConfig := postgres.Config{
			DSN:                  addr.DSN,
			Conn:                 addr.PostgresConn,
			PreferSimpleProtocol: addr.PostgresPreferSimpleProtocol,
			WithoutReturning:     addr.PostgresWithoutReturning,
		}
		db, err := gorm.Open(postgres.New(postgresConfig), config)
		if err != nil {
			return nil, sderr.Wrap(err, "postgres error")
		}
		return db, nil
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(addr.DSN), config)
		if err != nil {
			return nil, sderr.Wrap(err, "open sqlite error")
		}
		return db, nil
	default:
		return nil, ErrIllegalDriver
	}
}
