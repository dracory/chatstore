package chatstore

import (
	"encoding/json"
	"maps"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

const CHAT_MESSAGE_STATUS_ACTIVE = "active"
const CHAT_MESSAGE_STATUS_INACTIVE = "inactive"
const CHAT_MESSAGE_STATUS_DELETED = "deleted"

type messageImplementation struct {
	dataobject.DataObject
}

func NewMessage() MessageInterface {
	o := &messageImplementation{}
	o.SetID(uid.HumanUid())
	o.SetStatus(CHAT_MESSAGE_STATUS_ACTIVE)
	// REQUIRED: o.SetChatID("")
	// REQUIRED: o.SetSenderID("")
	// REQUIRED: o.SetRecipientID("")
	// REQUIRED: o.SetText("")
	o.SetMemo("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).Format("Y-m-d H:i:s"))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).Format("Y-m-d H:i:s"))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)

	o.SetMetas(map[string]string{})

	return o
}

func NewMessageFromExistingData(data map[string]string) MessageInterface {
	o := &messageImplementation{}
	o.Hydrate(data)
	return o
}

// IsSoftDeleted checks if the message is soft deleted
func (o *messageImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAt() != sb.MAX_DATETIME
}

func (o *messageImplementation) ChatID() string {
	return o.Get(COLUMN_CHAT_ID)
}

func (o *messageImplementation) SetChatID(id string) MessageInterface {
	o.Set(COLUMN_CHAT_ID, id)
	return o
}

func (o *messageImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *messageImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_CREATED_AT), carbon.UTC)
}
func (o *messageImplementation) SetCreatedAt(createdAt string) MessageInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *messageImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *messageImplementation) SetMemo(memo string) MessageInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *messageImplementation) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	return metas[key], nil
}

func (o *messageImplementation) SetMeta(key string, value string) error {
	return o.UpsertMetas(map[string]string{
		key: value,
	})
}

func (o *messageImplementation) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

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

func (o *messageImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (o *messageImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	maps.Copy(currentMetas, metas)

	return o.SetMetas(currentMetas)
}

func (o *messageImplementation) RecipientID() string {
	return o.Get(COLUMN_RECIPIENT_ID)
}

func (o *messageImplementation) SetRecipientID(id string) MessageInterface {
	o.Set(COLUMN_RECIPIENT_ID, id)
	return o
}

func (o *messageImplementation) SenderID() string {
	return o.Get(COLUMN_SENDER_ID)
}

func (o *messageImplementation) SetSenderID(id string) MessageInterface {
	o.Set(COLUMN_SENDER_ID, id)
	return o
}

func (o *messageImplementation) SetSoftDeletedAt(softDeletedAt string) MessageInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *messageImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *messageImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *messageImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *messageImplementation) SetStatus(status string) MessageInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *messageImplementation) Text() string {
	return o.Get(COLUMN_TEXT)
}

func (o *messageImplementation) SetText(text string) MessageInterface {
	o.Set(COLUMN_TEXT, text)
	return o
}

func (o *messageImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *messageImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *messageImplementation) SetUpdatedAt(updatedAt string) MessageInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
