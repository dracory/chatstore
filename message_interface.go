package chatstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
)

type MessageInterface interface {
	dataobject.DataObjectInterface

	IsSoftDeleted() bool

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
}
