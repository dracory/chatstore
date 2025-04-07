package chatstore

import (
	"database/sql"
	"errors"
	"os"

	"github.com/gouniverse/utils"
)

const testChat_O1 = "00000000000000000000000000000010"
const testChat_O2 = "00000000000000000000000000000020"
const testMessage_O1 = "00000000000000000000000000000050"
const testMessage_O2 = "00000000000000000000000000000060"
const testUser_O1 = "00000000000000000000000000000030"
const testUser_O2 = "00000000000000000000000000000040"

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil {
			panic(err)
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func initStore(filepath string) (StoreInterface, error) {
	db := initDB(filepath)

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableChatName:      "chat_table",
		TableMessageName:   "message_table",
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}
