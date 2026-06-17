package chatstore

import "errors"

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

	IsStatusInSet() bool
	GetStatusIn() []string
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
}

// MessageQuery creates a new message query
func MessageQuery() MessageQueryInterface {
	return NewMessageQuery()
}

// NewMessageQuery creates a new message query
func NewMessageQuery() MessageQueryInterface {
	return &messageQueryImplementation{
		params: make(map[string]any),
	}
}

var _ MessageQueryInterface = (*messageQueryImplementation)(nil)

// messageQuery implements the MessageQueryInterface
type messageQueryImplementation struct {
	params map[string]any
}

// Validate validates the query parameters
func (q *messageQueryImplementation) Validate() error {
	if q.IsChatIDSet() && q.GetChatID() == "" {
		return errors.New("message query: chat_id cannot be empty")
	}

	if q.IsCreatedAtGteSet() && q.GetCreatedAtGte() == "" {
		return errors.New("message query: created_at_gte cannot be empty")
	}

	if q.IsCreatedAtLteSet() && q.GetCreatedAtLte() == "" {
		return errors.New("message query: created_at_lte cannot be empty")
	}

	if q.IsIDSet() && q.GetID() == "" {
		return errors.New("message query: id cannot be empty")
	}

	if q.IsIDInSet() && len(q.GetIDIn()) < 1 {
		return errors.New("message query: id_in cannot be empty array")
	}

	if q.IsIDNotInSet() && len(q.GetIDNotIn()) < 1 {
		return errors.New("message query: id_not_in cannot be empty array")
	}

	if q.IsLimitSet() && q.GetLimit() < 0 {
		return errors.New("message query: limit cannot be negative")
	}

	if q.IsOffsetSet() && q.GetOffset() < 0 {
		return errors.New("message query: offset cannot be negative")
	}

	if q.IsOrderBySet() && q.GetOrderBy() == "" {
		return errors.New("message query: order_by cannot be empty")
	}

	if q.IsOrderDirectionSet() && q.GetOrderDirection() == "" {
		return errors.New("message query: order_direction cannot be empty")
	}

	if q.IsRecipientIDSet() && q.GetRecipientID() == "" {
		return errors.New("message query: recipient_id cannot be empty")
	}

	if q.IsSenderIDSet() && q.GetSenderID() == "" {
		return errors.New("message query: sender_id cannot be empty")
	}

	if q.IsStatusSet() && q.GetStatus() == "" {
		return errors.New("message query: status cannot be empty")
	}

	if q.IsStatusInSet() && len(q.GetStatusIn()) < 1 {
		return errors.New("message query: status_in cannot be empty array")
	}

	return nil
}

func (q *messageQueryImplementation) hasProperty(key string) bool {
	_, ok := q.params[key]
	return ok
}

// ============================================================================
// == Getters and Setters
// ============================================================================

func (q *messageQueryImplementation) IsCountOnlySet() bool {
	return q.hasProperty("count_only")
}

func (q *messageQueryImplementation) GetCountOnly() bool {
	if q.IsCountOnlySet() {
		return q.params["count_only"].(bool)
	}
	return false
}

func (q *messageQueryImplementation) SetCountOnly(countOnly bool) MessageQueryInterface {
	q.params["count_only"] = countOnly
	return q
}

func (q *messageQueryImplementation) IsCreatedAtGteSet() bool {
	return q.hasProperty("created_at_gte")
}

func (q *messageQueryImplementation) GetCreatedAtGte() string {
	if q.IsCreatedAtGteSet() {
		return q.params["created_at_gte"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetCreatedAtGte(createdAtGte string) MessageQueryInterface {
	q.params["created_at_gte"] = createdAtGte
	return q
}

func (q *messageQueryImplementation) IsCreatedAtLteSet() bool {
	return q.hasProperty("created_at_lte")
}

func (q *messageQueryImplementation) GetCreatedAtLte() string {
	if q.IsCreatedAtLteSet() {
		return q.params["created_at_lte"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetCreatedAtLte(createdAtLte string) MessageQueryInterface {
	q.params["created_at_lte"] = createdAtLte
	return q
}

func (q *messageQueryImplementation) IsChatIDSet() bool {
	return q.hasProperty("chat_id")
}

func (q *messageQueryImplementation) GetChatID() string {
	if q.IsChatIDSet() {
		return q.params["chat_id"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetChatID(chatID string) MessageQueryInterface {
	q.params["chat_id"] = chatID
	return q
}

func (q *messageQueryImplementation) IsChatIDInSet() bool {
	return q.hasProperty("chat_id_in")
}

func (q *messageQueryImplementation) GetChatIDIn() []string {
	if q.IsChatIDInSet() {
		return q.params["chat_id_in"].([]string)
	}
	return []string{}
}

func (q *messageQueryImplementation) SetChatIDIn(chatIDIn []string) MessageQueryInterface {
	q.params["chat_id_in"] = chatIDIn
	return q
}

func (q *messageQueryImplementation) IsIDSet() bool {
	return q.hasProperty("id")
}

func (q *messageQueryImplementation) GetID() string {
	if q.IsIDSet() {
		return q.params["id"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetID(id string) MessageQueryInterface {
	q.params["id"] = id
	return q
}

func (q *messageQueryImplementation) IsIDInSet() bool {
	return q.hasProperty("id_in")
}

func (q *messageQueryImplementation) GetIDIn() []string {
	if q.IsIDInSet() {
		return q.params["id_in"].([]string)
	}
	return []string{}
}

func (q *messageQueryImplementation) SetIDIn(idIn []string) MessageQueryInterface {
	q.params["id_in"] = idIn
	return q
}

func (q *messageQueryImplementation) IsIDNotInSet() bool {
	return q.hasProperty("id_not_in")
}

func (q *messageQueryImplementation) GetIDNotIn() []string {
	if q.IsIDNotInSet() {
		return q.params["id_not_in"].([]string)
	}
	return []string{}
}

func (q *messageQueryImplementation) SetIDNotIn(idNotIn []string) MessageQueryInterface {
	q.params["id_not_in"] = idNotIn
	return q
}

func (q *messageQueryImplementation) IsLimitSet() bool {
	return q.hasProperty("limit")
}

func (q *messageQueryImplementation) GetLimit() int {
	if q.IsLimitSet() {
		return q.params["limit"].(int)
	}
	return 0
}

func (q *messageQueryImplementation) SetLimit(limit int) MessageQueryInterface {
	q.params["limit"] = limit
	return q
}

func (q *messageQueryImplementation) IsOffsetSet() bool {
	return q.hasProperty("offset")
}

func (q *messageQueryImplementation) GetOffset() int {
	if q.IsOffsetSet() {
		return q.params["offset"].(int)
	}
	return 0
}

func (q *messageQueryImplementation) SetOffset(offset int) MessageQueryInterface {
	q.params["offset"] = offset
	return q
}

func (q *messageQueryImplementation) IsOnlySoftDeletedSet() bool {
	return q.hasProperty("only_soft_deleted")
}

func (q *messageQueryImplementation) GetOnlySoftDeleted() bool {
	if q.IsOnlySoftDeletedSet() {
		return q.params["only_soft_deleted"].(bool)
	}
	return false
}

func (q *messageQueryImplementation) SetOnlySoftDeleted(onlySoftDeleted bool) MessageQueryInterface {
	q.params["only_soft_deleted"] = onlySoftDeleted
	return q
}

func (q *messageQueryImplementation) IsWithSoftDeletedSet() bool {
	return q.hasProperty("with_soft_deleted")
}

func (q *messageQueryImplementation) GetWithSoftDeleted() bool {
	if q.IsWithSoftDeletedSet() {
		return q.params["with_soft_deleted"].(bool)
	}
	return false
}

func (q *messageQueryImplementation) SetWithSoftDeleted(withSoftDeleted bool) MessageQueryInterface {
	q.params["with_soft_deleted"] = withSoftDeleted
	return q
}

func (q *messageQueryImplementation) IsOrderBySet() bool {
	return q.hasProperty("order_by")
}

func (q *messageQueryImplementation) GetOrderBy() string {
	if q.IsOrderBySet() {
		return q.params["order_by"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetOrderBy(orderBy string) MessageQueryInterface {
	q.params["order_by"] = orderBy
	return q
}

func (q *messageQueryImplementation) IsOrderDirectionSet() bool {
	return q.hasProperty("order_direction")
}

func (q *messageQueryImplementation) GetOrderDirection() string {
	if q.IsOrderDirectionSet() {
		return q.params["order_direction"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetOrderDirection(orderDirection string) MessageQueryInterface {
	q.params["order_direction"] = orderDirection
	return q
}

func (q *messageQueryImplementation) IsRecipientIDSet() bool {
	return q.hasProperty("recipient_id")
}

func (q *messageQueryImplementation) GetRecipientID() string {
	if q.IsRecipientIDSet() {
		return q.params["recipient_id"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetRecipientID(recipientID string) MessageQueryInterface {
	q.params["recipient_id"] = recipientID
	return q
}

func (q *messageQueryImplementation) IsSenderIDSet() bool {
	return q.hasProperty("sender_id")
}

func (q *messageQueryImplementation) GetSenderID() string {
	if q.IsSenderIDSet() {
		return q.params["sender_id"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetSenderID(senderID string) MessageQueryInterface {
	q.params["sender_id"] = senderID
	return q
}

func (q *messageQueryImplementation) IsStatusSet() bool {
	return q.hasProperty("status")
}

func (q *messageQueryImplementation) GetStatus() string {
	if q.IsStatusSet() {
		return q.params["status"].(string)
	}
	return ""
}

func (q *messageQueryImplementation) SetStatus(status string) MessageQueryInterface {
	q.params["status"] = status
	return q
}

func (q *messageQueryImplementation) IsStatusInSet() bool {
	return q.hasProperty("status_in")
}

func (q *messageQueryImplementation) GetStatusIn() []string {
	if q.IsStatusInSet() {
		return q.params["status_in"].([]string)
	}
	return []string{}
}

func (q *messageQueryImplementation) SetStatusIn(statuses []string) MessageQueryInterface {
	q.params["status_in"] = statuses
	return q
}
