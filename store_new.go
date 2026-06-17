package chatstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/neat"
)

// NewStoreOptions defines the options for creating a new chat store.
type NewStoreOptions struct {
	TableChatName      string
	TableMessageName   string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
	Logger             *slog.Logger
}

// NewStore creates a new chat store.
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.TableChatName == "" {
		return nil, errors.New("chat store: TableChatName is required")
	}

	if opts.TableMessageName == "" {
		return nil, errors.New("chat store: TableMessageName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("chat store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	if opts.Logger == nil {
		opts.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	store := &storeImplementation{
		tableChat:          opts.TableChatName,
		tableMessage:       opts.TableMessageName,
		db:                 neatDB,
		automigrateEnabled: opts.AutomigrateEnabled,
		debugEnabled:       opts.DebugEnabled,
		logger:             opts.Logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
