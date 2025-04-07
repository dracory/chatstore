package chatstore

import (
	"maps"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

const CHAT_MESSAGE_STATUS_ACTIVE = "active"
const CHAT_MESSAGE_STATUS_INACTIVE = "inactive"
const CHAT_MESSAGE_STATUS_DELETED = "deleted"

type Message struct {
	dataobject.DataObject
}

func NewMessage() MessageInterface {
	o := &Message{}
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
	o := &Message{}
	o.Hydrate(data)
	return o
}

// IsSoftDeleted checks if the message is soft deleted
func (o *Message) IsSoftDeleted() bool {
	return o.SoftDeletedAt() != sb.MAX_DATETIME
}

func (o *Message) ChatID() string {
	return o.Get(COLUMN_CHAT_ID)
}

func (o *Message) SetChatID(id string) MessageInterface {
	o.Set(COLUMN_CHAT_ID, id)
	return o
}

func (o *Message) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *Message) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_CREATED_AT), carbon.UTC)
}
func (o *Message) SetCreatedAt(createdAt string) MessageInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *Message) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *Message) SetMemo(memo string) MessageInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *Message) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	return metas[key], nil
}

func (o *Message) SetMeta(key string, value string) error {
	return o.UpsertMetas(map[string]string{
		key: value,
	})
}

func (o *Message) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (o *Message) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)
	return nil
}

func (o *Message) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	maps.Copy(currentMetas, metas)

	return o.SetMetas(currentMetas)
}

func (o *Message) RecipientID() string {
	return o.Get(COLUMN_RECIPIENT_ID)
}

func (o *Message) SetRecipientID(id string) MessageInterface {
	o.Set(COLUMN_RECIPIENT_ID, id)
	return o
}

func (o *Message) SenderID() string {
	return o.Get(COLUMN_SENDER_ID)
}

func (o *Message) SetSenderID(id string) MessageInterface {
	o.Set(COLUMN_SENDER_ID, id)
	return o
}

func (o *Message) SetSoftDeletedAt(softDeletedAt string) MessageInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *Message) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *Message) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *Message) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *Message) SetStatus(status string) MessageInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *Message) Text() string {
	return o.Get(COLUMN_TEXT)
}

func (o *Message) SetText(text string) MessageInterface {
	o.Set(COLUMN_TEXT, text)
	return o
}

func (o *Message) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *Message) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *Message) SetUpdatedAt(updatedAt string) MessageInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
