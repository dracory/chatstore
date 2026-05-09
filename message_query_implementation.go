package chatstore

import (
	"errors"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

// messageQuery implements the MessageQueryInterface
type messageQueryImplementation struct {
	params map[string]interface{}
}

var _ MessageQueryInterface = (*messageQueryImplementation)(nil)

// MessageQuery creates a new message query
func MessageQuery() MessageQueryInterface {
	return &messageQueryImplementation{
		params: map[string]interface{}{},
	}
}

func (q *messageQueryImplementation) ToSelectDataset(st *store) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if st == nil {
		return nil, []any{}, errors.New("store cannot be nil")
	}

	if err := q.Validate(); err != nil {
		return nil, []any{}, err
	}

	sql := goqu.Dialect(st.dbDriverName).From(st.tableMessage)

	// Chat ID filter
	if q.IsChatIDSet() {
		sql = sql.Where(goqu.C(COLUMN_CHAT_ID).Eq(q.GetChatID()))
	}

	// Chat ID IN filter
	if q.IsChatIDInSet() {
		sql = sql.Where(goqu.C(COLUMN_CHAT_ID).In(q.GetChatIDIn()))
	}

	// Created At filter
	if q.IsCreatedAtGteSet() {
		sql = sql.Where(goqu.C(COLUMN_CREATED_AT).Gte(q.GetCreatedAtGte()))
	}

	if q.IsCreatedAtLteSet() {
		sql = sql.Where(goqu.C(COLUMN_CREATED_AT).Lte(q.GetCreatedAtLte()))
	}

	// ID filter
	if q.IsIDSet() {
		sql = sql.Where(goqu.C(COLUMN_ID).Eq(q.GetID()))
	}

	// ID IN filter
	if q.IsIDInSet() {
		sql = sql.Where(goqu.C(COLUMN_ID).In(q.GetIDIn()))
	}

	// Recipient ID filter
	if q.IsRecipientIDSet() {
		sql = sql.Where(goqu.C(COLUMN_RECIPIENT_ID).Eq(q.GetRecipientID()))
	}

	// Sender ID filter
	if q.IsSenderIDSet() {
		sql = sql.Where(goqu.C(COLUMN_SENDER_ID).Eq(q.GetSenderID()))
	}

	// Status filter
	if q.IsStatusSet() {
		sql = sql.Where(goqu.C(COLUMN_STATUS).Eq(q.GetStatus()))
	}

	// Status IN filter
	if q.IsStatusInSet() {
		sql = sql.Where(goqu.C(COLUMN_STATUS).In(q.GetStatusIn()))
	}

	// Updated At filter
	if q.IsUpdatedAtGteSet() {
		sql = sql.Where(goqu.C(COLUMN_UPDATED_AT).Gte(q.GetUpdatedAtGte()))
	}

	if q.IsUpdatedAtLteSet() {
		sql = sql.Where(goqu.C(COLUMN_UPDATED_AT).Lte(q.GetUpdatedAtLte()))
	}

	if !q.IsCountOnlySet() {
		if q.IsLimitSet() {
			sql = sql.Limit(uint(q.GetLimit()))
		}

		if q.IsOffsetSet() {
			sql = sql.Offset(uint(q.GetOffset()))
		}
	}

	sortOrder := sb.DESC
	if q.IsOrderDirectionSet() {
		sortOrder = q.GetOrderDirection()
	}

	if q.IsOrderBySet() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Asc())
		} else {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Desc())
		}
	}

	// Limit (if count only is not set)
	if !q.IsCountOnlySet() || !q.GetCountOnly() {
		if q.IsLimitSet() {
			sql = sql.Limit(uint(q.GetLimit()))
		}

		if q.IsOffsetSet() {
			sql = sql.Offset(uint(q.GetOffset()))
		}
	}

	// Sort order
	if q.IsOrderBySet() {
		sortOrder := q.GetOrderDirection()

		if strings.EqualFold(sortOrder, sb.ASC) {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Asc())
		} else {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Desc())
		}
	}

	// Soft delete filters

	// Only soft deleted
	if q.IsOnlySoftDeletedSet() && q.GetOnlySoftDeleted() {
		sql = sql.Where(goqu.C(COLUMN_SOFT_DELETED_AT).Lte(carbon.Now(carbon.UTC).ToDateTimeString()))
		return sql, []any{}, nil
	}

	// Include soft deleted
	if q.IsWithSoftDeletedSet() && q.GetWithSoftDeleted() {
		return sql, []any{}, nil
	}

	// Exclude soft deleted, not in the past (default)
	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	sql = sql.Where(softDeleted)

	return sql, []any{}, nil
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

func (q *messageQueryImplementation) IsLimitSet() bool {
	return q.hasProperty("limit")
}

func (q *messageQueryImplementation) GetLimit() int {
	if q.IsLimitSet() {
		return q.params["limit"].(int)
	}

	return 0
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

func (q *messageQueryImplementation) SetStatusIn(statusIn []string) MessageQueryInterface {
	q.params["status_in"] = statusIn
	return q
}

func (q *messageQueryImplementation) IsUpdatedAtGteSet() bool {
	return q.hasProperty("updated_at_gte")
}

func (q *messageQueryImplementation) GetUpdatedAtGte() string {
	if q.IsUpdatedAtGteSet() {
		return q.params["updated_at_gte"].(string)
	}

	return ""
}

func (q *messageQueryImplementation) SetUpdatedAtGte(updatedAt string) MessageQueryInterface {
	q.params["updated_at_gte"] = updatedAt
	return q
}

func (q *messageQueryImplementation) IsUpdatedAtLteSet() bool {
	return q.hasProperty("updated_at_lte")
}

func (q *messageQueryImplementation) GetUpdatedAtLte() string {
	if q.IsUpdatedAtLteSet() {
		return q.params["updated_at_lte"].(string)
	}

	return ""
}

func (q *messageQueryImplementation) SetUpdatedAtLte(updatedAt string) MessageQueryInterface {
	q.params["updated_at_lte"] = updatedAt
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

func (q *messageQueryImplementation) hasProperty(key string) bool {
	return q.params[key] != nil
}
