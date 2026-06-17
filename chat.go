package chatstore

import (
	"encoding/json"
	"maps"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// ChatInterface defines the interface for a chat record.
type ChatInterface interface {
	IsSoftDeleted() bool

	ID() string
	SetID(id string) ChatInterface

	OwnerID() string
	SetOwnerID(id string) ChatInterface

	Status() string
	SetStatus(status string) ChatInterface

	Title() string
	SetTitle(title string) ChatInterface

	Memo() string
	SetMemo(memo string) ChatInterface

	Meta(key string) (string, error)
	SetMeta(key string, value string) error

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) ChatInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) ChatInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) ChatInterface

	MarkAsNotDirty()
}

var _ ChatInterface = (*chatImplementation)(nil)

// == TYPE ===================================================================

// chatImplementation is the private implementation of ChatInterface.
type chatImplementation struct {
	orm.ShortID

	StatusField    string `db:"status"`
	OwnerIDField   string `db:"owner_id"`
	TitleField     string `db:"title"`
	MemoField      string `db:"memo"`
	MetasField     string `db:"metas"`
	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS ============================================================

// NewChat creates a new chat.
func NewChat() ChatInterface {
	o := &chatImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetStatus(CHAT_STATUS_ACTIVE)
	o.SetTitle("")
	o.SetMemo("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	o.SetMetas(map[string]string{})
	return o
}

// NewChatFromExistingData creates a new chat from a raw column map (e.g. query results).
func NewChatFromExistingData(data map[string]string) ChatInterface {
	o := &chatImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetOwnerID(data[COLUMN_OWNER_ID])
	o.SetTitle(data[COLUMN_TITLE])
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

// ID returns the id of the chat.
func (o *chatImplementation) ID() string {
	return o.ShortID.ID
}

// SetID sets the id of the chat.
func (o *chatImplementation) SetID(id string) ChatInterface {
	o.ShortID.ID = id
	return o
}

// OwnerID returns the owner id of the chat.
func (o *chatImplementation) OwnerID() string {
	return o.OwnerIDField
}

// SetOwnerID sets the owner id of the chat.
func (o *chatImplementation) SetOwnerID(id string) ChatInterface {
	o.OwnerIDField = id
	return o
}

// Status returns the status of the chat.
func (o *chatImplementation) Status() string {
	return o.StatusField
}

// SetStatus sets the status of the chat.
func (o *chatImplementation) SetStatus(status string) ChatInterface {
	o.StatusField = status
	return o
}

// Title returns the title of the chat.
func (o *chatImplementation) Title() string {
	return o.TitleField
}

// SetTitle sets the title of the chat.
func (o *chatImplementation) SetTitle(title string) ChatInterface {
	o.TitleField = title
	return o
}

// Memo returns the memo of the chat.
func (o *chatImplementation) Memo() string {
	return o.MemoField
}

// SetMemo sets the memo of the chat.
func (o *chatImplementation) SetMemo(memo string) ChatInterface {
	o.MemoField = memo
	return o
}

// Meta returns a single meta value by key.
func (o *chatImplementation) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	return metas[key], nil
}

// SetMeta sets a single meta key-value pair.
func (o *chatImplementation) SetMeta(key string, value string) error {
	return o.UpsertMetas(map[string]string{
		key: value,
	})
}

// Metas returns the metas map of the chat.
func (o *chatImplementation) Metas() (map[string]string, error) {
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

// SetMetas sets the metas map of the chat.
func (o *chatImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.MetasField = string(mapString)
	return nil
}

// UpsertMetas merges the given metas into the existing metas.
func (o *chatImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()
	if err != nil {
		return err
	}
	maps.Copy(currentMetas, metas)
	return o.SetMetas(currentMetas)
}

// CreatedAt returns the created at time of the chat.
func (o *chatImplementation) CreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// CreatedAtCarbon returns the created at time of the chat as a carbon object.
func (o *chatImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the created at time of the chat.
func (o *chatImplementation) SetCreatedAt(createdAt string) ChatInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// SoftDeletedAt returns the soft deleted at time of the chat as a string.
func (o *chatImplementation) SoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// SoftDeletedAtCarbon returns the soft deleted at time of the chat as a carbon object.
func (o *chatImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the chat.
func (o *chatImplementation) SetSoftDeletedAt(softDeletedAt string) ChatInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}

// UpdatedAt returns the updated at time of the chat.
func (o *chatImplementation) UpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// UpdatedAtCarbon returns the updated at time of the chat as a carbon object.
func (o *chatImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the updated at time of the chat.
func (o *chatImplementation) SetUpdatedAt(updatedAt string) ChatInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// MarkAsNotDirty is a no-op for backward compatibility.
func (o *chatImplementation) MarkAsNotDirty() {
}
