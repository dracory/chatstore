package chatstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// == INTERFACE ===============================================================

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
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	// MigrateUp creates the chat and message tables
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

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

// == TYPE ====================================================================

var _ StoreInterface = (*storeImplementation)(nil)

// storeImplementation implements StoreInterface for chat operations.
type storeImplementation struct {
	tableChat          string
	tableMessage       string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	logger             *slog.Logger
}

// == MIGRATE =================================================================

// MigrateUp creates the chat and message tables if they do not already exist.
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.tableChat) && st.db.Schema().HasTable(st.tableMessage) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: tables already exist", "chat_table", st.tableChat, "message_table", st.tableMessage)
		}
		return nil
	}

	err := st.db.Schema().Create(st.tableChat, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 21)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 40)
		table.String(COLUMN_OWNER_ID, 40)
		table.String(COLUMN_TITLE, 255)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateUp chat table failed", "error", err)
		}
		return err
	}

	err = st.db.Schema().Create(st.tableMessage, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 21)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_CHAT_ID, 21)
		table.String(COLUMN_STATUS, 40)
		table.String(COLUMN_SENDER_ID, 40)
		table.String(COLUMN_RECIPIENT_ID, 40)
		table.Text(COLUMN_TEXT)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateUp message table failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the chat and message tables.
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.tableMessage) {
		err := st.db.Schema().Drop(st.tableMessage)
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown message table failed", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.tableChat) {
		err := st.db.Schema().Drop(st.tableChat)
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown chat table failed", "error", err)
			}
			return err
		}
	}

	return nil
}

// == DEBUG ===================================================================

// EnableDebug enables or disables debug mode.
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// == TABLE NAME ==============================================================

// GetChatTableName returns the chat table name.
func (st *storeImplementation) GetChatTableName() string {
	return st.tableChat
}

// SetChatTableName sets the chat table name.
func (st *storeImplementation) SetChatTableName(tableName string) {
	st.tableChat = tableName
}

// GetMessageTableName returns the message table name.
func (st *storeImplementation) GetMessageTableName() string {
	return st.tableMessage
}

// SetMessageTableName sets the message table name.
func (st *storeImplementation) SetMessageTableName(tableName string) {
	st.tableMessage = tableName
}

// == CHAT METHODS ============================================================

// ChatCount counts the number of chats that match the query.
func (st *storeImplementation) ChatCount(options ChatQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("query is nil")
	}

	q := st.buildChatQuery(options)

	var count int64
	err := q.Table(st.tableChat).Count(&count)
	return count, err
}

// ChatCreate creates a new chat.
func (st *storeImplementation) ChatCreate(chat ChatInterface) error {
	if chat == nil {
		return errors.New("chat is nil")
	}

	if chat.ID() == "" {
		return errors.New("chat ID is required")
	}

	chat.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	chat.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_ID:              chat.ID(),
		COLUMN_STATUS:          chat.Status(),
		COLUMN_OWNER_ID:        chat.OwnerID(),
		COLUMN_TITLE:           chat.Title(),
		COLUMN_MEMO:            chat.Memo(),
		COLUMN_METAS:           chat.(*chatImplementation).MetasField,
		COLUMN_CREATED_AT:      chat.CreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      chat.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: chat.SoftDeletedAtCarbon().StdTime(),
	}

	if st.debugEnabled {
		st.logger.Debug("Chat create", "id", chat.ID())
	}

	return st.db.Query().Table(st.tableChat).Create(row)
}

// ChatDelete permanently deletes a chat.
func (st *storeImplementation) ChatDelete(chat ChatInterface) error {
	if chat == nil {
		return errors.New("chat is nil")
	}
	return st.ChatDeleteByID(chat.ID())
}

// ChatDeleteByID permanently deletes a chat by ID.
func (st *storeImplementation) ChatDeleteByID(id string) error {
	if id == "" {
		return errors.New("chat ID is required")
	}

	_, err := st.db.Query().
		Table(st.tableChat).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// ChatFindByID finds a chat by ID.
func (st *storeImplementation) ChatFindByID(chatID string) (ChatInterface, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	list, err := st.ChatList(ChatQuery().
		SetID(chatID).
		SetLimit(1))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// ChatList lists chats based on the query.
func (st *storeImplementation) ChatList(query ChatQueryInterface) ([]ChatInterface, error) {
	if query == nil {
		return nil, errors.New("query is nil")
	}

	type chatRow struct {
		ID            string    `db:"id"`
		Status        string    `db:"status"`
		OwnerID       string    `db:"owner_id"`
		Title         string    `db:"title"`
		Memo          string    `db:"memo"`
		Metas         string    `db:"metas"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := st.buildChatQuery(query)

	var rows []chatRow
	if err := q.Table(st.tableChat).Get(&rows); err != nil {
		return []ChatInterface{}, err
	}

	list := make([]ChatInterface, 0, len(rows))
	for _, r := range rows {
		chat := &chatImplementation{}
		chat.SetID(r.ID)
		chat.StatusField = r.Status
		chat.OwnerIDField = r.OwnerID
		chat.TitleField = r.Title
		chat.MemoField = r.Memo
		chat.MetasField = r.Metas
		chat.CreatedAtField.CreatedAt = r.CreatedAt
		chat.UpdatedAtField.UpdatedAt = r.UpdatedAt
		chat.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, chat)
	}

	return list, nil
}

// ChatSoftDelete soft deletes a chat.
func (st *storeImplementation) ChatSoftDelete(chat ChatInterface) error {
	if chat == nil {
		return errors.New("chat is nil")
	}

	chat.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: chat.SoftDeletedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).StdTime(),
	}

	_, err := st.db.Query().Table(st.tableChat).Where(COLUMN_ID+" = ?", chat.ID()).Update(row)
	return err
}

// ChatSoftDeleteByID soft deletes a chat by ID.
func (st *storeImplementation) ChatSoftDeleteByID(id string) error {
	chat, err := st.ChatFindByID(id)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}
	return st.ChatSoftDelete(chat)
}

// ChatUpdate updates a chat.
func (st *storeImplementation) ChatUpdate(chat ChatInterface) error {
	if chat == nil {
		return errors.New("chat is nil")
	}

	if chat.ID() == "" {
		return errors.New("chat ID is required")
	}

	chat.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_STATUS:          chat.Status(),
		COLUMN_OWNER_ID:        chat.OwnerID(),
		COLUMN_TITLE:           chat.Title(),
		COLUMN_MEMO:            chat.Memo(),
		COLUMN_METAS:           chat.(*chatImplementation).MetasField,
		COLUMN_UPDATED_AT:      chat.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: chat.SoftDeletedAtCarbon().StdTime(),
	}

	_, err := st.db.Query().Table(st.tableChat).Where(COLUMN_ID+" = ?", chat.ID()).Update(row)
	return err
}

// == MESSAGE METHODS =========================================================

// MessageCount counts the number of messages that match the query.
func (st *storeImplementation) MessageCount(options MessageQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("query is nil")
	}

	q := st.buildMessageQuery(options)

	var count int64
	err := q.Table(st.tableMessage).Count(&count)
	return count, err
}

// MessageCreate creates a new message.
func (st *storeImplementation) MessageCreate(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}

	if message.ID() == "" {
		return errors.New("message ID is required")
	}

	message.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	message.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_ID:              message.ID(),
		COLUMN_CHAT_ID:         message.ChatID(),
		COLUMN_STATUS:          message.Status(),
		COLUMN_SENDER_ID:       message.SenderID(),
		COLUMN_RECIPIENT_ID:    message.RecipientID(),
		COLUMN_TEXT:            message.Text(),
		COLUMN_MEMO:            message.Memo(),
		COLUMN_METAS:           message.(*messageImplementation).MetasField,
		COLUMN_CREATED_AT:      message.CreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      message.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: message.SoftDeletedAtCarbon().StdTime(),
	}

	if st.debugEnabled {
		st.logger.Debug("Message create", "id", message.ID())
	}

	return st.db.Query().Table(st.tableMessage).Create(row)
}

// MessageDelete permanently deletes a message.
func (st *storeImplementation) MessageDelete(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}
	return st.MessageDeleteByID(message.ID())
}

// MessageDeleteByID permanently deletes a message by ID.
func (st *storeImplementation) MessageDeleteByID(id string) error {
	if id == "" {
		return errors.New("message ID is required")
	}

	_, err := st.db.Query().
		Table(st.tableMessage).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// MessageFindByID finds a message by ID.
func (st *storeImplementation) MessageFindByID(messageID string) (MessageInterface, error) {
	if messageID == "" {
		return nil, errors.New("message ID is required")
	}

	list, err := st.MessageList(MessageQuery().
		SetID(messageID).
		SetLimit(1))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// MessageList lists messages based on the query.
func (st *storeImplementation) MessageList(query MessageQueryInterface) ([]MessageInterface, error) {
	if query == nil {
		return nil, errors.New("query is nil")
	}

	type messageRow struct {
		ID            string    `db:"id"`
		ChatID        string    `db:"chat_id"`
		Status        string    `db:"status"`
		SenderID      string    `db:"sender_id"`
		RecipientID   string    `db:"recipient_id"`
		Text          string    `db:"text"`
		Memo          string    `db:"memo"`
		Metas         string    `db:"metas"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := st.buildMessageQuery(query)

	var rows []messageRow
	if err := q.Table(st.tableMessage).Get(&rows); err != nil {
		return []MessageInterface{}, err
	}

	list := make([]MessageInterface, 0, len(rows))
	for _, r := range rows {
		msg := &messageImplementation{}
		msg.SetID(r.ID)
		msg.ChatIDField = r.ChatID
		msg.StatusField = r.Status
		msg.SenderIDField = r.SenderID
		msg.RecipientIDField = r.RecipientID
		msg.TextField = r.Text
		msg.MemoField = r.Memo
		msg.MetasField = r.Metas
		msg.CreatedAtField.CreatedAt = r.CreatedAt
		msg.UpdatedAtField.UpdatedAt = r.UpdatedAt
		msg.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, msg)
	}

	return list, nil
}

// MessageSoftDelete soft deletes a message.
func (st *storeImplementation) MessageSoftDelete(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}

	message.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: message.SoftDeletedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).StdTime(),
	}

	_, err := st.db.Query().Table(st.tableMessage).Where(COLUMN_ID+" = ?", message.ID()).Update(row)
	return err
}

// MessageSoftDeleteByID soft deletes a message by ID.
func (st *storeImplementation) MessageSoftDeleteByID(id string) error {
	message, err := st.MessageFindByID(id)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("message not found")
	}
	return st.MessageSoftDelete(message)
}

// MessageUpdate updates a message.
func (st *storeImplementation) MessageUpdate(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}

	if message.ID() == "" {
		return errors.New("message ID is required")
	}

	message.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_CHAT_ID:         message.ChatID(),
		COLUMN_STATUS:          message.Status(),
		COLUMN_SENDER_ID:       message.SenderID(),
		COLUMN_RECIPIENT_ID:    message.RecipientID(),
		COLUMN_TEXT:            message.Text(),
		COLUMN_MEMO:            message.Memo(),
		COLUMN_METAS:           message.(*messageImplementation).MetasField,
		COLUMN_UPDATED_AT:      message.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: message.SoftDeletedAtCarbon().StdTime(),
	}

	_, err := st.db.Query().Table(st.tableMessage).Where(COLUMN_ID+" = ?", message.ID()).Update(row)
	return err
}

// == QUERY BUILDERS ==========================================================

// buildChatQuery builds a neat query from the chat query interface.
func (st *storeImplementation) buildChatQuery(query ChatQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := st.db.Query().Model(&chatImplementation{})

	if query == nil {
		return q
	}

	if query.IsOwnerIDSet() && query.GetOwnerID() != "" {
		q = q.Where(COLUMN_OWNER_ID+" = ?", query.GetOwnerID())
	}

	if query.IsStatusSet() && query.GetStatus() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", query.GetStatus())
	}

	if query.IsStatusInSet() && len(query.GetStatusIn()) > 0 {
		q = q.Where(COLUMN_STATUS+" IN ?", query.GetStatusIn())
	}

	if query.IsIDSet() && query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if query.IsIDInSet() && len(query.GetIDIn()) > 0 {
		q = q.Where(COLUMN_ID+" IN ?", query.GetIDIn())
	}

	if query.IsCreatedAtGteSet() && query.GetCreatedAtGte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", query.GetCreatedAtGte())
	}

	if query.IsCreatedAtLteSet() && query.GetCreatedAtLte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", query.GetCreatedAtLte())
	}

	if query.IsUpdatedAtGteSet() && query.GetUpdatedAtGte() != "" {
		q = q.Where(COLUMN_UPDATED_AT+" >= ?", query.GetUpdatedAtGte())
	}

	if query.IsUpdatedAtLteSet() && query.GetUpdatedAtLte() != "" {
		q = q.Where(COLUMN_UPDATED_AT+" <= ?", query.GetUpdatedAtLte())
	}

	if query.IsLimitSet() && query.GetLimit() > 0 {
		q = q.Limit(query.GetLimit())
	}

	if query.IsOffsetSet() && query.GetOffset() > 0 {
		q = q.Offset(query.GetOffset())
	}

	if query.IsOrderBySet() && query.GetOrderBy() != "" {
		direction := lo.CoalesceOrEmpty(query.GetOrderDirection(), "DESC")
		q = q.OrderBy(query.GetOrderBy() + " " + direction)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if query.IsWithSoftDeletedSet() && query.GetWithSoftDeleted() {
		q = q.WithSoftDeleted()
	} else if query.IsOnlySoftDeletedSet() && query.GetOnlySoftDeleted() {
		q = q.OnlySoftDeleted()
	}

	return q
}

// buildMessageQuery builds a neat query from the message query interface.
func (st *storeImplementation) buildMessageQuery(query MessageQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := st.db.Query().Model(&messageImplementation{})

	if query == nil {
		return q
	}

	if query.IsChatIDSet() && query.GetChatID() != "" {
		q = q.Where(COLUMN_CHAT_ID+" = ?", query.GetChatID())
	}

	if query.IsChatIDInSet() && len(query.GetChatIDIn()) > 0 {
		q = q.Where(COLUMN_CHAT_ID+" IN ?", query.GetChatIDIn())
	}

	if query.IsStatusSet() && query.GetStatus() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", query.GetStatus())
	}

	if query.IsStatusInSet() && len(query.GetStatusIn()) > 0 {
		q = q.Where(COLUMN_STATUS+" IN ?", query.GetStatusIn())
	}

	if query.IsIDSet() && query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if query.IsIDInSet() && len(query.GetIDIn()) > 0 {
		q = q.Where(COLUMN_ID+" IN ?", query.GetIDIn())
	}

	if query.IsIDNotInSet() && len(query.GetIDNotIn()) > 0 {
		q = q.Where(COLUMN_ID+" NOT IN ?", query.GetIDNotIn())
	}

	if query.IsSenderIDSet() && query.GetSenderID() != "" {
		q = q.Where(COLUMN_SENDER_ID+" = ?", query.GetSenderID())
	}

	if query.IsRecipientIDSet() && query.GetRecipientID() != "" {
		q = q.Where(COLUMN_RECIPIENT_ID+" = ?", query.GetRecipientID())
	}

	if query.IsCreatedAtGteSet() && query.GetCreatedAtGte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", query.GetCreatedAtGte())
	}

	if query.IsCreatedAtLteSet() && query.GetCreatedAtLte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", query.GetCreatedAtLte())
	}

	if query.IsLimitSet() && query.GetLimit() > 0 {
		q = q.Limit(query.GetLimit())
	}

	if query.IsOffsetSet() && query.GetOffset() > 0 {
		q = q.Offset(query.GetOffset())
	}

	if query.IsOrderBySet() && query.GetOrderBy() != "" {
		direction := lo.CoalesceOrEmpty(query.GetOrderDirection(), "DESC")
		q = q.OrderBy(query.GetOrderBy() + " " + direction)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if query.IsWithSoftDeletedSet() && query.GetWithSoftDeleted() {
		q = q.WithSoftDeleted()
	} else if query.IsOnlySoftDeletedSet() && query.GetOnlySoftDeleted() {
		q = q.OnlySoftDeleted()
	}

	return q
}
