package datastore

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/cuvva-sample-code/datastore/migrations"
)

// Connection helps setup database
type Connection struct {
	Logger     *logrus.Logger
	DB         *sql.DB
	SQLBuilder squirrel.StatementBuilderType
}

// NewConnection create a new connection to the database
func NewConnection(logger *logrus.Logger, srcs ...interface{}) (*Connection, error) {
	var (
		err  error
		conn *pgx.ConnConfig
	)
	if len(srcs) > 0 {
		conn, err = pgx.ParseConfig(srcs[0].(string))
	} else {
		conn, err = pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/cuvva_merchant?sslmode=disable")
	}
	if err != nil {
		return nil, err
	}
	conn.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*conn)
	err = validateSchema(db)
	if err != nil {
		return nil, err
	}
	fmt.Println("how working")

	return &Connection{
		Logger:     logger,
		DB:         db,
		SQLBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}, err
}

//Close connection
func (c *Connection) Close() error {
	return c.DB.Close()
}

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const version = 1

func validateSchema(db *sql.DB) error {
	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return err
	}

	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}
