package chatstore

import (
	"encoding/json"
	"maps"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

type chatImplementation struct {
	dataobject.DataObject
}

func NewChat() ChatInterface {
	o := &chatImplementation{}

	o.SetID(uid.HumanUid()).
		// REQUIRED:SetOwnerID("").
		SetStatus(CHAT_STATUS_ACTIVE).
		SetTitle("").
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).Format("Y-m-d H:i:s")).
		SetUpdatedAt(carbon.Now(carbon.UTC).Format("Y-m-d H:i:s")).
		SetSoftDeletedAt(sb.MAX_DATE)

	o.SetMetas(map[string]string{})

	return o
}

func NewChatFromExistingData(data map[string]string) ChatInterface {
	o := &chatImplementation{}
	o.Hydrate(data)
	return o
}

// func (o *Chat) AddMetas(metas map[string]string) error {
// 	currentMetas, err := o.Metas()
// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range metas {
// 		currentMetas[k] = v

// 	}

// 	return o.SetMetas(currentMetas)
// }

func (o *chatImplementation) IsSoftDeleted() bool {
	return o.Get(COLUMN_SOFT_DELETED_AT) != ""
}

func (o *chatImplementation) OwnerID() string {
	return o.Get(COLUMN_OWNER_ID)
}

func (o *chatImplementation) SetOwnerID(id string) ChatInterface {
	o.Set(COLUMN_OWNER_ID, id)
	return o
}

func (o *chatImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *chatImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_CREATED_AT), carbon.UTC)
}

func (o *chatImplementation) SetCreatedAt(createdAt string) ChatInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *chatImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *chatImplementation) SetID(id string) ChatInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *chatImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *chatImplementation) SetMemo(memo string) ChatInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *chatImplementation) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	return metas[key], nil
}

func (o *chatImplementation) SetMeta(key string, value string) error {
	return o.UpsertMetas(map[string]string{
		key: value,
	})
}

func (o *chatImplementation) Metas() (map[string]string, error) {
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

func (o *chatImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (o *chatImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	maps.Copy(currentMetas, metas)

	return o.SetMetas(currentMetas)
}

func (o *chatImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *chatImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_SOFT_DELETED_AT), carbon.UTC)
}

func (o *chatImplementation) SetSoftDeletedAt(softDeletedAt string) ChatInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *chatImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *chatImplementation) SetStatus(status string) ChatInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *chatImplementation) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *chatImplementation) SetTitle(title string) ChatInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *chatImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *chatImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *chatImplementation) SetUpdatedAt(updatedAt string) ChatInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
