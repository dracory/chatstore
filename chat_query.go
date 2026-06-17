package chatstore

import "errors"

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

	IsStatusInSet() bool
	GetStatusIn() []string
	SetStatusIn(statuses []string) ChatQueryInterface

	IsUpdatedAtGteSet() bool
	GetUpdatedAtGte() string
	SetUpdatedAtGte(updatedAt string) ChatQueryInterface

	IsUpdatedAtLteSet() bool
	GetUpdatedAtLte() string
	SetUpdatedAtLte(updatedAt string) ChatQueryInterface
}

// ChatQuery creates a new chat query
func ChatQuery() ChatQueryInterface {
	return NewChatQuery()
}

// NewChatQuery creates a new chat query
func NewChatQuery() ChatQueryInterface {
	return &chatQueryImplementation{
		params: make(map[string]any),
	}
}

var _ ChatQueryInterface = (*chatQueryImplementation)(nil)

// chatQuery implements the ChatQueryInterface
type chatQueryImplementation struct {
	params map[string]any
}

// Validate validates the query parameters
func (q *chatQueryImplementation) Validate() error {
	if q.IsOwnerIDSet() && q.GetOwnerID() == "" {
		return errors.New("chat query: owner_id cannot be empty")
	}

	if q.IsCreatedAtGteSet() && q.GetCreatedAtGte() == "" {
		return errors.New("chat query: created_at_gte cannot be empty")
	}

	if q.IsCreatedAtLteSet() && q.GetCreatedAtLte() == "" {
		return errors.New("chat query: created_at_lte cannot be empty")
	}

	if q.IsIDSet() && q.GetID() == "" {
		return errors.New("chat query: id cannot be empty")
	}

	if q.IsIDInSet() && len(q.GetIDIn()) < 1 {
		return errors.New("chat query: id_in cannot be empty array")
	}

	if q.IsLimitSet() && q.GetLimit() < 0 {
		return errors.New("chat query: limit cannot be negative")
	}

	if q.IsOffsetSet() && q.GetOffset() < 0 {
		return errors.New("chat query: offset cannot be negative")
	}

	if q.IsStatusSet() && q.GetStatus() == "" {
		return errors.New("chat query: status cannot be empty")
	}

	if q.IsStatusInSet() && len(q.GetStatusIn()) < 1 {
		return errors.New("chat query: status_in cannot be empty array")
	}

	return nil
}

func (q *chatQueryImplementation) hasProperty(key string) bool {
	_, ok := q.params[key]
	return ok
}

// ============================================================================
// == Getters and Setters
// ============================================================================

func (q *chatQueryImplementation) IsOwnerIDSet() bool {
	return q.hasProperty("owner_id")
}

func (q *chatQueryImplementation) GetOwnerID() string {
	if q.IsOwnerIDSet() {
		return q.params["owner_id"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetOwnerID(ownerID string) ChatQueryInterface {
	q.params["owner_id"] = ownerID
	return q
}

func (q *chatQueryImplementation) IsCountOnlySet() bool {
	return q.hasProperty("count_only")
}

func (q *chatQueryImplementation) GetCountOnly() bool {
	if q.IsCountOnlySet() {
		return q.params["count_only"].(bool)
	}
	return false
}

func (q *chatQueryImplementation) SetCountOnly(countOnly bool) ChatQueryInterface {
	q.params["count_only"] = countOnly
	return q
}

func (q *chatQueryImplementation) IsCreatedAtGteSet() bool {
	return q.hasProperty("created_at_gte")
}

func (q *chatQueryImplementation) GetCreatedAtGte() string {
	if q.IsCreatedAtGteSet() {
		return q.params["created_at_gte"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetCreatedAtGte(createdAtGte string) ChatQueryInterface {
	q.params["created_at_gte"] = createdAtGte
	return q
}

func (q *chatQueryImplementation) IsCreatedAtLteSet() bool {
	return q.hasProperty("created_at_lte")
}

func (q *chatQueryImplementation) GetCreatedAtLte() string {
	if q.IsCreatedAtLteSet() {
		return q.params["created_at_lte"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetCreatedAtLte(createdAtLte string) ChatQueryInterface {
	q.params["created_at_lte"] = createdAtLte
	return q
}

func (q *chatQueryImplementation) IsIDSet() bool {
	return q.hasProperty("id")
}

func (q *chatQueryImplementation) GetID() string {
	if q.IsIDSet() {
		return q.params["id"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetID(id string) ChatQueryInterface {
	q.params["id"] = id
	return q
}

func (q *chatQueryImplementation) IsIDInSet() bool {
	return q.hasProperty("id_in")
}

func (q *chatQueryImplementation) GetIDIn() []string {
	if q.IsIDInSet() {
		return q.params["id_in"].([]string)
	}
	return []string{}
}

func (q *chatQueryImplementation) SetIDIn(idIn []string) ChatQueryInterface {
	q.params["id_in"] = idIn
	return q
}

func (q *chatQueryImplementation) IsLimitSet() bool {
	return q.hasProperty("limit")
}

func (q *chatQueryImplementation) GetLimit() int {
	if q.IsLimitSet() {
		return q.params["limit"].(int)
	}
	return 0
}

func (q *chatQueryImplementation) SetLimit(limit int) ChatQueryInterface {
	q.params["limit"] = limit
	return q
}

func (q *chatQueryImplementation) IsOffsetSet() bool {
	return q.hasProperty("offset")
}

func (q *chatQueryImplementation) GetOffset() int {
	if q.IsOffsetSet() {
		return q.params["offset"].(int)
	}
	return 0
}

func (q *chatQueryImplementation) SetOffset(offset int) ChatQueryInterface {
	q.params["offset"] = offset
	return q
}

func (q *chatQueryImplementation) IsOnlySoftDeletedSet() bool {
	return q.hasProperty("only_soft_deleted")
}

func (q *chatQueryImplementation) GetOnlySoftDeleted() bool {
	if q.IsOnlySoftDeletedSet() {
		return q.params["only_soft_deleted"].(bool)
	}
	return false
}

func (q *chatQueryImplementation) SetOnlySoftDeleted(onlySoftDeleted bool) ChatQueryInterface {
	q.params["only_soft_deleted"] = onlySoftDeleted
	return q
}

func (q *chatQueryImplementation) IsWithSoftDeletedSet() bool {
	return q.hasProperty("with_soft_deleted")
}

func (q *chatQueryImplementation) GetWithSoftDeleted() bool {
	if q.IsWithSoftDeletedSet() {
		return q.params["with_soft_deleted"].(bool)
	}
	return false
}

func (q *chatQueryImplementation) SetWithSoftDeleted(withSoftDeleted bool) ChatQueryInterface {
	q.params["with_soft_deleted"] = withSoftDeleted
	return q
}

func (q *chatQueryImplementation) IsOrderBySet() bool {
	return q.hasProperty("order_by")
}

func (q *chatQueryImplementation) GetOrderBy() string {
	if q.IsOrderBySet() {
		return q.params["order_by"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetOrderBy(orderBy string) ChatQueryInterface {
	q.params["order_by"] = orderBy
	return q
}

func (q *chatQueryImplementation) IsOrderDirectionSet() bool {
	return q.hasProperty("order_direction")
}

func (q *chatQueryImplementation) GetOrderDirection() string {
	if q.IsOrderDirectionSet() {
		return q.params["order_direction"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetOrderDirection(orderDirection string) ChatQueryInterface {
	q.params["order_direction"] = orderDirection
	return q
}

func (q *chatQueryImplementation) IsStatusSet() bool {
	return q.hasProperty("status")
}

func (q *chatQueryImplementation) GetStatus() string {
	if q.IsStatusSet() {
		return q.params["status"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetStatus(status string) ChatQueryInterface {
	q.params["status"] = status
	return q
}

func (q *chatQueryImplementation) IsStatusInSet() bool {
	return q.hasProperty("status_in")
}

func (q *chatQueryImplementation) GetStatusIn() []string {
	if q.IsStatusInSet() {
		return q.params["status_in"].([]string)
	}
	return []string{}
}

func (q *chatQueryImplementation) SetStatusIn(statuses []string) ChatQueryInterface {
	q.params["status_in"] = statuses
	return q
}

func (q *chatQueryImplementation) IsUpdatedAtGteSet() bool {
	return q.hasProperty("updated_at_gte")
}

func (q *chatQueryImplementation) GetUpdatedAtGte() string {
	if q.IsUpdatedAtGteSet() {
		return q.params["updated_at_gte"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetUpdatedAtGte(updatedAt string) ChatQueryInterface {
	q.params["updated_at_gte"] = updatedAt
	return q
}

func (q *chatQueryImplementation) IsUpdatedAtLteSet() bool {
	return q.hasProperty("updated_at_lte")
}

func (q *chatQueryImplementation) GetUpdatedAtLte() string {
	if q.IsUpdatedAtLteSet() {
		return q.params["updated_at_lte"].(string)
	}
	return ""
}

func (q *chatQueryImplementation) SetUpdatedAtLte(updatedAt string) ChatQueryInterface {
	q.params["updated_at_lte"] = updatedAt
	return q
}
