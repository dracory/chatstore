package chatstore

import "github.com/doug-martin/goqu/v9"

// ChatQueryInterface defines the interface for querying chats
type ChatQueryInterface interface {
	// Validation method
	Validate() error

	// Count related methods
	IsCountOnlySet() bool
	GetCountOnly() bool
	SetCountOnly(countOnly bool) ChatQueryInterface

	// Soft delete related query methods
	IsWithSoftDeletedSet() bool
	GetWithSoftDeleted() bool
	SetWithSoftDeleted(withSoftDeleted bool) ChatQueryInterface

	IsOnlySoftDeletedSet() bool
	GetOnlySoftDeleted() bool
	SetOnlySoftDeleted(onlySoftDeleted bool) ChatQueryInterface

	// Dataset conversion methods
	ToSelectDataset(store *store) (selectDataset *goqu.SelectDataset, columns []any, err error)

	// Field query methods

	IsOwnerIDSet() bool
	GetOwnerID() string
	SetOwnerID(ownerID string) ChatQueryInterface

	IsCreatedAtGteSet() bool
	GetCreatedAtGte() string
	SetCreatedAtGte(createdAt string) ChatQueryInterface

	IsCreatedAtLteSet() bool
	GetCreatedAtLte() string
	SetCreatedAtLte(createdAt string) ChatQueryInterface

	IsIDSet() bool
	GetID() string
	SetID(id string) ChatQueryInterface

	IsIDInSet() bool
	GetIDIn() []string
	SetIDIn(ids []string) ChatQueryInterface

	IsLimitSet() bool
	GetLimit() int
	SetLimit(limit int) ChatQueryInterface

	IsOffsetSet() bool
	GetOffset() int
	SetOffset(offset int) ChatQueryInterface

	IsOrderBySet() bool
	GetOrderBy() string
	SetOrderBy(orderBy string) ChatQueryInterface

	IsOrderDirectionSet() bool
	GetOrderDirection() string
	SetOrderDirection(orderDirection string) ChatQueryInterface

	IsStatusSet() bool
	GetStatus() string
	SetStatus(status string) ChatQueryInterface
	SetStatusIn(statuses []string) ChatQueryInterface

	IsUpdatedAtGteSet() bool
	GetUpdatedAtGte() string
	SetUpdatedAtGte(updatedAt string) ChatQueryInterface

	IsUpdatedAtLteSet() bool
	GetUpdatedAtLte() string
	SetUpdatedAtLte(updatedAt string) ChatQueryInterface
}
