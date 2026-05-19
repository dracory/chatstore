package chatstore

import "database/sql"

type StoreInterface interface {
	// GetChatTableName returns the chat table name
	GetChatTableName() string
	// SetChatTableName sets the chat table name
	SetChatTableName(tableName string)

	// GetMessageTableName returns the message table name
	GetMessageTableName() string
	// SetMessageTableName sets the message table name
	SetMessageTableName(tableName string)

	// MigrateDown drops the chat and message tables
	MigrateDown(tx ...*sql.Tx) error
	// MigrateUp creates the chat and message tables
	MigrateUp(tx ...*sql.Tx) error

	// EnableDebug enables or disables debug mode
	EnableDebug(enabled bool)

	// ChatCount returns the number of chats
	ChatCount(options ChatQueryInterface) (int64, error)
	ChatCreate(chat ChatInterface) error
	ChatDelete(chat ChatInterface) error
	ChatDeleteByID(id string) error
	ChatFindByID(id string) (ChatInterface, error)
	ChatList(options ChatQueryInterface) ([]ChatInterface, error)
	ChatSoftDelete(chat ChatInterface) error
	ChatSoftDeleteByID(id string) error
	ChatUpdate(chat ChatInterface) error

	MessageCount(options MessageQueryInterface) (int64, error)
	MessageCreate(message MessageInterface) error
	MessageDelete(message MessageInterface) error
	MessageDeleteByID(id string) error
	MessageFindByID(id string) (MessageInterface, error)
	MessageList(options MessageQueryInterface) ([]MessageInterface, error)
	MessageSoftDelete(message MessageInterface) error
	MessageSoftDeleteByID(id string) error
	MessageUpdate(message MessageInterface) error
}
