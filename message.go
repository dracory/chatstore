package chatstore

import (
	"encoding/json"
	"maps"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// MessageInterface defines the interface for a message record.
type MessageInterface interface {
	IsSoftDeleted() bool

	ID() string
	SetID(id string) MessageInterface

	ChatID() string
	SetChatID(chatID string) MessageInterface

	SenderID() string
	SetSenderID(id string) MessageInterface

	RecipientID() string
	SetRecipientID(id string) MessageInterface

	Status() string
	SetStatus(status string) MessageInterface

	Memo() string
	SetMemo(memo string) MessageInterface

	Meta(key string) (string, error)
	SetMeta(key string, value string) error

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Text() string
	SetText(text string) MessageInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) MessageInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) MessageInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) MessageInterface

	MarkAsNotDirty()
}

var _ MessageInterface = (*messageImplementation)(nil)

// == TYPE ===================================================================

// messageImplementation is the private implementation of MessageInterface.
type messageImplementation struct {
	orm.ShortID

	ChatIDField      string `db:"chat_id"`
	StatusField      string `db:"status"`
	SenderIDField    string `db:"sender_id"`
	RecipientIDField string `db:"recipient_id"`
	TextField        string `db:"text"`
	MemoField        string `db:"memo"`
	MetasField       string `db:"metas"`
	CreatedAtField   orm.CreatedAt
	UpdatedAtField   orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS ============================================================

// NewMessage creates a new message.
func NewMessage() MessageInterface {
	o := &messageImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetStatus(MESSAGE_STATUS_ACTIVE)
	o.SetMemo("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	o.SetMetas(map[string]string{})
	return o
}

// NewMessageFromExistingData creates a new message from a raw column map (e.g. query results).
func NewMessageFromExistingData(data map[string]string) MessageInterface {
	o := &messageImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetChatID(data[COLUMN_CHAT_ID])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetSenderID(data[COLUMN_SENDER_ID])
	o.SetRecipientID(data[COLUMN_RECIPIENT_ID])
	o.SetText(data[COLUMN_TEXT])
	o.SetMemo(data[COLUMN_MEMO])
	o.MetasField = data[COLUMN_METAS]
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS =================================================================

// == SETTERS AND GETTERS =====================================================

// ID returns the id of the message.
func (o *messageImplementation) ID() string {
	return o.ShortID.ID
}

// SetID sets the id of the message.
func (o *messageImplementation) SetID(id string) MessageInterface {
	o.ShortID.ID = id
	return o
}

// ChatID returns the chat id of the message.
func (o *messageImplementation) ChatID() string {
	return o.ChatIDField
}

// SetChatID sets the chat id of the message.
func (o *messageImplementation) SetChatID(chatID string) MessageInterface {
	o.ChatIDField = chatID
	return o
}

// SenderID returns the sender id of the message.
func (o *messageImplementation) SenderID() string {
	return o.SenderIDField
}

// SetSenderID sets the sender id of the message.
func (o *messageImplementation) SetSenderID(id string) MessageInterface {
	o.SenderIDField = id
	return o
}

// RecipientID returns the recipient id of the message.
func (o *messageImplementation) RecipientID() string {
	return o.RecipientIDField
}

// SetRecipientID sets the recipient id of the message.
func (o *messageImplementation) SetRecipientID(id string) MessageInterface {
	o.RecipientIDField = id
	return o
}

// Status returns the status of the message.
func (o *messageImplementation) Status() string {
	return o.StatusField
}

// SetStatus sets the status of the message.
func (o *messageImplementation) SetStatus(status string) MessageInterface {
	o.StatusField = status
	return o
}

// Memo returns the memo of the message.
func (o *messageImplementation) Memo() string {
	return o.MemoField
}

// SetMemo sets the memo of the message.
func (o *messageImplementation) SetMemo(memo string) MessageInterface {
	o.MemoField = memo
	return o
}

// Meta returns a single meta value by key.
func (o *messageImplementation) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	return metas[key], nil
}

// SetMeta sets a single meta key-value pair.
func (o *messageImplementation) SetMeta(key string, value string) error {
	return o.UpsertMetas(map[string]string{
		key: value,
	})
}

// Metas returns the metas map of the message.
func (o *messageImplementation) Metas() (map[string]string, error) {
	metasStr := o.MetasField
	if metasStr == "" {
		metasStr = "{}"
	}
	var metasJson map[string]string
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}
	return metasJson, nil
}

// SetMetas sets the metas map of the message.
func (o *messageImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.MetasField = string(mapString)
	return nil
}

// UpsertMetas merges the given metas into the existing metas.
func (o *messageImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()
	if err != nil {
		return err
	}
	maps.Copy(currentMetas, metas)
	return o.SetMetas(currentMetas)
}

// Text returns the text of the message.
func (o *messageImplementation) Text() string {
	return o.TextField
}

// SetText sets the text of the message.
func (o *messageImplementation) SetText(text string) MessageInterface {
	o.TextField = text
	return o
}

// CreatedAt returns the created at time of the message.
func (o *messageImplementation) CreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// CreatedAtCarbon returns the created at time of the message as a carbon object.
func (o *messageImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the created at time of the message.
func (o *messageImplementation) SetCreatedAt(createdAt string) MessageInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// SoftDeletedAt returns the soft deleted at time of the message as a string.
func (o *messageImplementation) SoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// SoftDeletedAtCarbon returns the soft deleted at time of the message as a carbon object.
func (o *messageImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the message.
func (o *messageImplementation) SetSoftDeletedAt(softDeletedAt string) MessageInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

// UpdatedAt returns the updated at time of the message.
func (o *messageImplementation) UpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// UpdatedAtCarbon returns the updated at time of the message as a carbon object.
func (o *messageImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the updated at time of the message.
func (o *messageImplementation) SetUpdatedAt(updatedAt string) MessageInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// MarkAsNotDirty is a no-op for backward compatibility.
func (o *messageImplementation) MarkAsNotDirty() {
}
