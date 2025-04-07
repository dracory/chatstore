package chatstore

import "github.com/dromara/carbon/v2"

type ChatInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

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
}
