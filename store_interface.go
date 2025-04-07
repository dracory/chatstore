package chatstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(enabled bool)

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
