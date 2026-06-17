package chatstore

// Column names for the chat and message tables
const (
	COLUMN_CHAT_ID         = "chat_id"
	COLUMN_CREATED_AT      = "created_at"
	COLUMN_ID              = "id"
	COLUMN_MEMO            = "memo"
	COLUMN_METAS           = "metas"
	COLUMN_RECIPIENT_ID    = "recipient_id"
	COLUMN_SENDER_ID       = "sender_id"
	COLUMN_OWNER_ID        = "owner_id"
	COLUMN_SOFT_DELETED_AT = "soft_deleted_at"
	COLUMN_STATUS          = "status"
	COLUMN_TEXT            = "text"
	COLUMN_TITLE           = "title"
	COLUMN_UPDATED_AT      = "updated_at"
)

// Status constants
const (
	CHAT_STATUS_ACTIVE   = "active"
	CHAT_STATUS_INACTIVE = "inactive"
	CHAT_STATUS_DELETED  = "deleted"

	CHAT_SYSTEM_ID = "00000000000000000000000000000001"
)

// Message status constants
const (
	MESSAGE_STATUS_ACTIVE   = "active"
	MESSAGE_STATUS_INACTIVE = "inactive"
	MESSAGE_STATUS_DELETED  = "deleted"
)

// MAX_DATETIME is a far-future datetime used as the default soft-delete sentinel.
const MAX_DATETIME = "9999-12-31 23:59:59"
