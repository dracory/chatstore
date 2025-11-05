package chatstore

import (
	"errors"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

// chatQuery implements the ChatQueryInterface
type chatQuery struct {
	params map[string]any
}

var _ ChatQueryInterface = (*chatQuery)(nil)

// ChatQuery creates a new chat query
func ChatQuery() ChatQueryInterface {
	return &chatQuery{
		params: map[string]any{},
	}
}

// Validate validates the query parameters
func (q *chatQuery) Validate() error {
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

func (q *chatQuery) ToSelectDataset(st *store) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if st == nil {
		return nil, []any{}, errors.New("store cannot be nil")
	}

	if err := q.Validate(); err != nil {
		return nil, []any{}, err
	}

	sql := goqu.Dialect(st.dbDriverName).From(st.tableChat)

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

// ============================================================================
// == Getters and Setters
// ============================================================================

func (q *chatQuery) IsOwnerIDSet() bool {
	return q.hasProperty("owner_id")
}

func (q *chatQuery) GetOwnerID() string {
	if q.IsOwnerIDSet() {
		return q.params["owner_id"].(string)
	}

	return ""
}

func (q *chatQuery) SetOwnerID(ownerID string) ChatQueryInterface {
	q.params["owner_id"] = ownerID
	return q
}

func (q *chatQuery) IsCountOnlySet() bool {
	return q.hasProperty("count_only")
}

func (q *chatQuery) GetCountOnly() bool {
	if q.IsCountOnlySet() {
		return q.params["count_only"].(bool)
	}

	return false
}

func (q *chatQuery) SetCountOnly(countOnly bool) ChatQueryInterface {
	q.params["count_only"] = countOnly
	return q
}

func (q *chatQuery) IsCreatedAtGteSet() bool {
	return q.hasProperty("created_at_gte")
}

func (q *chatQuery) GetCreatedAtGte() string {
	if q.IsCreatedAtGteSet() {
		return q.params["created_at_gte"].(string)
	}

	return ""
}

func (q *chatQuery) SetCreatedAtGte(createdAtGte string) ChatQueryInterface {
	q.params["created_at_gte"] = createdAtGte
	return q
}

func (q *chatQuery) IsCreatedAtLteSet() bool {
	return q.hasProperty("created_at_lte")
}

func (q *chatQuery) GetCreatedAtLte() string {
	if q.IsCreatedAtLteSet() {
		return q.params["created_at_lte"].(string)
	}

	return ""
}

func (q *chatQuery) SetCreatedAtLte(createdAtLte string) ChatQueryInterface {
	q.params["created_at_lte"] = createdAtLte
	return q
}

func (q *chatQuery) IsIDSet() bool {
	return q.hasProperty("id")
}
func (q *chatQuery) GetID() string {
	if q.IsIDSet() {
		return q.params["id"].(string)
	}

	return ""
}

func (q *chatQuery) SetID(id string) ChatQueryInterface {
	q.params["id"] = id
	return q
}

func (q *chatQuery) IsIDInSet() bool {
	return q.hasProperty("id_in")
}

func (q *chatQuery) GetIDIn() []string {
	if q.IsIDInSet() {
		return q.params["id_in"].([]string)
	}

	return []string{}
}

func (q *chatQuery) SetIDIn(idIn []string) ChatQueryInterface {
	q.params["id_in"] = idIn
	return q
}

func (q *chatQuery) IsLimitSet() bool {
	return q.hasProperty("limit")
}

func (q *chatQuery) GetLimit() int {
	if q.IsLimitSet() {
		return q.params["limit"].(int)
	}

	return 0
}

func (q *chatQuery) SetLimit(limit int) ChatQueryInterface {
	q.params["limit"] = limit
	return q
}

func (q *chatQuery) IsOffsetSet() bool {
	return q.hasProperty("offset")
}

func (q *chatQuery) GetOffset() int {
	if q.IsOffsetSet() {
		return q.params["offset"].(int)
	}

	return 0
}

func (q *chatQuery) SetOffset(offset int) ChatQueryInterface {
	q.params["offset"] = offset
	return q
}

func (q *chatQuery) IsOnlySoftDeletedSet() bool {
	return q.hasProperty("only_soft_deleted")
}

func (q *chatQuery) GetOnlySoftDeleted() bool {
	if q.IsOnlySoftDeletedSet() {
		return q.params["only_soft_deleted"].(bool)
	}

	return false
}

func (q *chatQuery) SetOnlySoftDeleted(onlySoftDeleted bool) ChatQueryInterface {
	q.params["only_soft_deleted"] = onlySoftDeleted
	return q
}

func (q *chatQuery) IsOrderDirectionSet() bool {
	return q.hasProperty("order_direction")
}

func (q *chatQuery) GetOrderDirection() string {
	if q.IsOrderDirectionSet() {
		return q.params["order_direction"].(string)
	}

	return ""
}

func (q *chatQuery) SetOrderDirection(orderDirection string) ChatQueryInterface {
	q.params["order_direction"] = orderDirection
	return q
}

func (q *chatQuery) IsOrderBySet() bool {
	return q.hasProperty("order_by")
}

func (q *chatQuery) GetOrderBy() string {
	if q.IsOrderBySet() {
		return q.params["order_by"].(string)
	}

	return ""
}

func (q *chatQuery) SetOrderBy(orderBy string) ChatQueryInterface {
	q.params["order_by"] = orderBy
	return q
}

func (q *chatQuery) IsStatusSet() bool {
	return q.hasProperty("status")
}

func (q *chatQuery) GetStatus() string {
	if q.IsStatusSet() {
		return q.params["status"].(string)
	}

	return ""
}

func (q *chatQuery) SetStatus(status string) ChatQueryInterface {
	q.params["status"] = status
	return q
}

func (q *chatQuery) IsStatusInSet() bool {
	return q.hasProperty("status_in")
}

func (q *chatQuery) GetStatusIn() []string {
	if q.IsStatusInSet() {
		return q.params["status_in"].([]string)
	}

	return []string{}
}

func (q *chatQuery) SetStatusIn(statusIn []string) ChatQueryInterface {
	q.params["status_in"] = statusIn
	return q
}

func (q *chatQuery) IsUpdatedAtGteSet() bool {
	return q.hasProperty("updated_at_gte")
}

func (q *chatQuery) GetUpdatedAtGte() string {
	if q.IsUpdatedAtGteSet() {
		return q.params["updated_at_gte"].(string)
	}

	return ""
}

func (q *chatQuery) SetUpdatedAtGte(updatedAt string) ChatQueryInterface {
	q.params["updated_at_gte"] = updatedAt
	return q
}

func (q *chatQuery) IsUpdatedAtLteSet() bool {
	return q.hasProperty("updated_at_lte")
}

func (q *chatQuery) GetUpdatedAtLte() string {
	if q.IsUpdatedAtLteSet() {
		return q.params["updated_at_lte"].(string)
	}

	return ""
}

func (q *chatQuery) SetUpdatedAtLte(updatedAt string) ChatQueryInterface {
	q.params["updated_at_lte"] = updatedAt
	return q
}

func (q *chatQuery) IsWithSoftDeletedSet() bool {
	return q.hasProperty("with_soft_deleted")
}

func (q *chatQuery) GetWithSoftDeleted() bool {
	if q.IsWithSoftDeletedSet() {
		return q.params["with_soft_deleted"].(bool)
	}

	return false
}

func (q *chatQuery) SetWithSoftDeleted(withSoftDeleted bool) ChatQueryInterface {
	q.params["with_soft_deleted"] = withSoftDeleted
	return q
}

func (q *chatQuery) hasProperty(key string) bool {
	return q.params[key] != nil
}
