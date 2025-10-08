package postgres

import (
	"errors"
	"fmt"

	"github.com/K1la/event-booker/internal/config"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

var (
	ErrNoSuchImage = errors.New("there is no image with such id")
)

type Postgres struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Postgres {
	return &Postgres{db: db}
}

func NewDB(cfg *config.Config) *dbpg.DB {
	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
	)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		zlog.Logger.Fatal().Msgf("could not init db: %v", err)
	}

	return db
}
