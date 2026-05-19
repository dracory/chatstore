package chatstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

// ============================================================================
// == TYPE
// ============================================================================

type storeImplementation struct {
	tableChat          string
	tableMessage       string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
	logger             *slog.Logger
}

// ============================================================================
// == INTERFACE
// ============================================================================

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// ============================================================================
// == PUBLIC METHODS
// ============================================================================

// MigrateUp creates the chat and message tables
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	if store.db == nil {
		return errors.New("chatstore: database is nil")
	}

	sqlChat, err := store.sqlChatTableCreate()
	if err != nil {
		return fmt.Errorf("failed to generate chat table SQL: %w", err)
	}

	sqlMessage, err := store.sqlMessageTableCreate()
	if err != nil {
		return fmt.Errorf("failed to generate message table SQL: %w", err)
	}

	sqls := []string{
		sqlChat,
		sqlMessage,
	}

	for _, sql := range sqls {
		if sql == "" {
			return errors.New("table create sql is empty")
		}

		var errExec error
		if txToUse != nil {
			_, errExec = txToUse.ExecContext(ctx, sql)
		} else {
			_, errExec = store.db.ExecContext(ctx, sql)
		}

		if errExec != nil {
			return errExec
		}
	}

	return nil
}

// MigrateDown drops the chat and message tables
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	if store.db == nil {
		return errors.New("chatstore: database is nil")
	}

	sqlChat, err := store.sqlChatTableDrop()
	if err != nil {
		return fmt.Errorf("failed to generate chat table drop SQL: %w", err)
	}

	sqlMessage, err := store.sqlMessageTableDrop()
	if err != nil {
		return fmt.Errorf("failed to generate message table drop SQL: %w", err)
	}

	// Drop in reverse order of creation
	sqls := []string{
		sqlMessage,
		sqlChat,
	}

	for _, sql := range sqls {
		if sql == "" {
			return errors.New("table drop sql is empty")
		}

		var errExec error
		if txToUse != nil {
			_, errExec = txToUse.ExecContext(ctx, sql)
		} else {
			_, errExec = store.db.ExecContext(ctx, sql)
		}

		if errExec != nil {
			return errExec
		}
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// GetChatTableName returns the chat table name
func (st *storeImplementation) GetChatTableName() string {
	return st.tableChat
}

// SetChatTableName sets the chat table name
func (st *storeImplementation) SetChatTableName(tableName string) {
	st.tableChat = tableName
}

// GetMessageTableName returns the message table name
func (st *storeImplementation) GetMessageTableName() string {
	return st.tableMessage
}

// SetMessageTableName sets the message table name
func (st *storeImplementation) SetMessageTableName(tableName string) {
	st.tableMessage = tableName
}

// GetLogger - returns the logger
func (st *storeImplementation) GetLogger() *slog.Logger {
	if st.logger == nil {
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return st.logger
}

// SetLogger - sets the logger
func (st *storeImplementation) SetLogger(logger *slog.Logger) {
	st.logger = logger
}
