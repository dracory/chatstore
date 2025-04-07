package chatstore

import "github.com/doug-martin/goqu/v9"

// MessageQueryInterface defines the interface for querying messages
type MessageQueryInterface interface {
	// Validation method
	Validate() error

	// Basic query methods
	IsCreatedAtGteSet() bool
	GetCreatedAtGte() string
	SetCreatedAtGte(createdAt string) MessageQueryInterface

	IsCreatedAtLteSet() bool
	GetCreatedAtLte() string
	SetCreatedAtLte(createdAt string) MessageQueryInterface

	IsIDSet() bool
	GetID() string
	SetID(id string) MessageQueryInterface

	IsIDInSet() bool
	GetIDIn() []string
	SetIDIn(ids []string) MessageQueryInterface

	IsIDNotInSet() bool
	GetIDNotIn() []string
	SetIDNotIn(ids []string) MessageQueryInterface

	IsLimitSet() bool
	GetLimit() int
	SetLimit(limit int) MessageQueryInterface

	IsChatIDSet() bool
	GetChatID() string
	SetChatID(chatID string) MessageQueryInterface

	IsChatIDInSet() bool
	GetChatIDIn() []string
	SetChatIDIn(chatIDs []string) MessageQueryInterface

	IsOffsetSet() bool
	GetOffset() int
	SetOffset(offset int) MessageQueryInterface

	IsOrderBySet() bool
	GetOrderBy() string
	SetOrderBy(orderBy string) MessageQueryInterface

	IsOrderDirectionSet() bool
	GetOrderDirection() string
	SetOrderDirection(orderDirection string) MessageQueryInterface

	IsRecipientIDSet() bool
	GetRecipientID() string
	SetRecipientID(recipientID string) MessageQueryInterface

	IsSenderIDSet() bool
	GetSenderID() string
	SetSenderID(senderID string) MessageQueryInterface

	IsStatusSet() bool
	GetStatus() string
	SetStatus(status string) MessageQueryInterface
	SetStatusIn(statuses []string) MessageQueryInterface

	// Count related methods
	IsCountOnlySet() bool
	GetCountOnly() bool
	SetCountOnly(countOnly bool) MessageQueryInterface

	// Soft delete related query methods
	IsWithSoftDeletedSet() bool
	GetWithSoftDeleted() bool
	SetWithSoftDeleted(withSoftDeleted bool) MessageQueryInterface

	IsOnlySoftDeletedSet() bool
	GetOnlySoftDeleted() bool
	SetOnlySoftDeleted(onlySoftDeleted bool) MessageQueryInterface

	ToSelectDataset(store *store) (selectDataset *goqu.SelectDataset, columns []any, err error)
}
