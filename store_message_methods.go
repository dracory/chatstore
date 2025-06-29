package chatstore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// MessageCount counts the number of messages that match the query
func (st *store) MessageCount(options MessageQueryInterface) (int64, error) {
	if st.db == nil {
		return 0, errors.New("database is not initialized")
	}

	if options == nil {
		return 0, errors.New("query is nil")
	}

	options.SetCountOnly(true)

	q, _, err := options.ToSelectDataset(st)

	if err != nil {
		return -1, err
	}

	sqlStr, sqlParams, err := q.
		//Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if err != nil {
		return -1, err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	mapped, err := database.SelectToMapString(database.Context(context.Background(), st.db), sqlStr, sqlParams...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	count, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err
	}

	return count, nil
}

// MessageCreate creates a new message
func (st *store) MessageCreate(message MessageInterface) error {
	if st.db == nil {
		return errors.New("database is not initialized")
	}

	if message.ID() == "" {
		return errors.New("message ID is required")
	}

	message.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	message.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := message.Data()

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		Insert(st.tableMessage).
		Prepared(true).
		Rows(data).
		ToSQL()

	if err != nil {
		return err
	}

	if st.debugEnabled {
		st.logger.Debug("Chat create query", "query", sqlStr, "params", sqlParams)
	}

	_, err = database.Execute(database.Context(context.Background(), st.db), sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	message.MarkAsNotDirty()

	return nil
}

// MessageDelete permanently deletes an message
func (st *store) MessageDelete(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}

	return st.MessageDeleteByID(message.ID())
}

// MessageDeleteByID permanently deletes an message by ID
func (st *store) MessageDeleteByID(id string) error {
	if st.db == nil {
		return errors.New("database is not initialized")
	}

	if id == "" {
		return errors.New("message ID is required")
	}

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		Delete(st.tableMessage).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if err != nil {
		return err
	}

	if st.debugEnabled {
		st.logger.Debug("Chat delete query", "query", sqlStr, "params", sqlParams)
	}

	_, err = database.Execute(database.Context(context.Background(), st.db), sqlStr, sqlParams...)
	if err != nil {
		return err
	}

	return nil
}

// MessageExists checks if an message exists
func (st *store) MessageExists(id string) (bool, error) {
	if st.db == nil {
		return false, errors.New("database is not initialized")
	}

	if id == "" {
		return false, errors.New("message ID is required")
	}

	count, err := st.MessageCount(MessageQuery().SetID(id))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// MessageFindByID finds an message by ID
func (st *store) MessageFindByID(messageID string) (MessageInterface, error) {
	if st.db == nil {
		return nil, errors.New("database is not initialized")
	}

	if messageID == "" {
		return nil, errors.New("message ID is required")
	}

	list, err := st.MessageList(MessageQuery().
		SetID(messageID).
		SetLimit(1))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// MessageList lists messages based on the query
func (st *store) MessageList(query MessageQueryInterface) ([]MessageInterface, error) {
	if st.db == nil {
		return nil, errors.New("database is not initialized")
	}

	if query == nil {
		return nil, errors.New("query is nil")
	}

	q, columns, err := query.ToSelectDataset(st)

	if err != nil {
		return []MessageInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Select(columns...).Prepared(true).ToSQL()

	if errSql != nil {
		return []MessageInterface{}, nil
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.Context(context.Background(), st.db), sqlStr, sqlParams...)

	if err != nil {
		return []MessageInterface{}, err
	}

	list := []MessageInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewMessageFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// MessageSoftDelete soft deletes an message
func (st *store) MessageSoftDelete(message MessageInterface) error {
	if message == nil {
		return errors.New("message is nil")
	}

	message.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return st.MessageUpdate(message)
}

// MessageSoftDeleteByID soft deletes an message by ID
func (st *store) MessageSoftDeleteByID(id string) error {
	message, err := st.MessageFindByID(id)

	if err != nil {
		return err
	}

	return st.MessageSoftDelete(message)
}

// MessageUpdate updates an message
func (st *store) MessageUpdate(message MessageInterface) error {
	if st.db == nil {
		return errors.New("database is not initialized")
	}

	if message == nil {
		return errors.New("message is nil")
	}

	if message.ID() == "" {
		return errors.New("message ID is required")
	}

	message.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := message.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(st.dbDriverName).
		Update(st.tableMessage).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(message.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr, params...)

	message.MarkAsNotDirty()

	return err
}

// func (store *store) incidentSelectQuery(options IncidentQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
// 	if options == nil {
// 		return nil, []any{}, errors.New("site options cannot be nil")
// 	}

// 	if err := options.Validate(); err != nil {
// 		return nil, []any{}, err
// 	}

// 	q := goqu.Dialect(store.dbDriverName).From(store.tableIncidents)

// 	// ID filter
// 	if options.IsIDInSet() {
// 		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
// 	}

// // ID IN filter
// if ids, exists := optionsMap["id_in"].([]string); exists && len(ids) > 0 {
// 	sql = sql.Where(goqu.C(COLUMN_ID).In(ids))
// }

// // Monitor ID filter
// if monitorID, exists := optionsMap["monitor_id"].(string); exists && monitorID != "" {
// 	sql = sql.Where(goqu.C(COLUMN_MONITOR_ID).Eq(monitorID))
// }

// // Monitor ID IN filter
// if monitorIDs, exists := optionsMap["monitor_id_in"].([]string); exists && len(monitorIDs) > 0 {
// 	sql = sql.Where(goqu.C(COLUMN_MONITOR_ID).In(monitorIDs))
// }

// // Status filter
// if status, exists := optionsMap["status"].(string); exists && status != "" {
// 	sql = sql.Where(goqu.C(COLUMN_STATUS).Eq(status))
// }

// 	return q, []any{}, nil
// }

// // applyIncidentFilters applies filters to the incident query
// func (st *store) applyIncidentFilters(sql *goqu.SelectDataset, optionsMap map[string]interface{}) *goqu.SelectDataset {
// 	// ID filter
// 	if id, exists := optionsMap["id"].(string); exists && id != "" {
// 		sql = sql.Where(goqu.C(COLUMN_ID).Eq(id))
// 	}

// 	// ID IN filter
// 	if ids, exists := optionsMap["id_in"].([]string); exists && len(ids) > 0 {
// 		sql = sql.Where(goqu.C(COLUMN_ID).In(ids))
// 	}

// 	// Monitor ID filter
// 	if monitorID, exists := optionsMap["monitor_id"].(string); exists && monitorID != "" {
// 		sql = sql.Where(goqu.C(COLUMN_MONITOR_ID).Eq(monitorID))
// 	}

// 	// Monitor ID IN filter
// 	if monitorIDs, exists := optionsMap["monitor_id_in"].([]string); exists && len(monitorIDs) > 0 {
// 		sql = sql.Where(goqu.C(COLUMN_MONITOR_ID).In(monitorIDs))
// 	}

// 	// Status filter
// 	if status, exists := optionsMap["status"].(string); exists && status != "" {
// 		sql = sql.Where(goqu.C(COLUMN_STATUS).Eq(status))
// 	}

// 	// Status IN filter
// 	if statuses, exists := optionsMap["status_in"].([]string); exists && len(statuses) > 0 {
// 		sql = sql.Where(goqu.C(COLUMN_STATUS).In(statuses))
// 	}

// 	// Created at GTE filter
// 	if createdAtGte, exists := optionsMap["created_at_gte"].(string); exists && createdAtGte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_CREATED_AT).Gte(createdAtGte))
// 	}

// 	// Created at LTE filter
// 	if createdAtLte, exists := optionsMap["created_at_lte"].(string); exists && createdAtLte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_CREATED_AT).Lte(createdAtLte))
// 	}

// 	// Start time GTE filter
// 	if startTimeGte, exists := optionsMap["start_time_gte"].(string); exists && startTimeGte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_START_TIME).Gte(startTimeGte))
// 	}

// 	// Start time LTE filter
// 	if startTimeLte, exists := optionsMap["start_time_lte"].(string); exists && startTimeLte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_START_TIME).Lte(startTimeLte))
// 	}

// 	// End time GTE filter
// 	if endTimeGte, exists := optionsMap["end_time_gte"].(string); exists && endTimeGte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_END_TIME).Gte(endTimeGte))
// 	}

// 	// End time LTE filter
// 	if endTimeLte, exists := optionsMap["end_time_lte"].(string); exists && endTimeLte != "" {
// 		sql = sql.Where(goqu.C(COLUMN_END_TIME).Lte(endTimeLte))
// 	}

// 	// Duration GTE filter
// 	if durationGte, exists := optionsMap["duration_gte"].(int); exists && durationGte > 0 {
// 		sql = sql.Where(goqu.C(COLUMN_DURATION).Gte(durationGte))
// 	}

// 	// Duration LTE filter
// 	if durationLte, exists := optionsMap["duration_lte"].(int); exists && durationLte > 0 {
// 		sql = sql.Where(goqu.C(COLUMN_DURATION).Lte(durationLte))
// 	}

// 	// Root cause filter
// 	if rootCause, exists := optionsMap["root_cause"].(string); exists && rootCause != "" {
// 		sql = sql.Where(goqu.C(COLUMN_ROOT_CAUSE).Eq(rootCause))
// 	}

// 	// Root cause LIKE filter
// 	if rootCauseLike, exists := optionsMap["root_cause_like"].(string); exists && rootCauseLike != "" {
// 		sql = sql.Where(goqu.C(COLUMN_ROOT_CAUSE).Like("%" + rootCauseLike + "%"))
// 	}

// 	// Soft delete filters
// 	withSoftDeleted, withSoftDeletedExists := optionsMap["with_soft_deleted"].(bool)
// 	onlySoftDeleted, onlySoftDeletedExists := optionsMap["only_soft_deleted"].(bool)

// 	// Only soft deleted
// 	if onlySoftDeletedExists && onlySoftDeleted {
// 		return sql.Where(goqu.C(COLUMN_SOFT_DELETED_AT).Lte(carbon.Now(carbon.UTC).ToDateTimeString()))
// 	}

// 	// Include soft deleted
// 	if withSoftDeletedExists && withSoftDeleted {
// 		return sql
// 	}

// 	// Exclude soft deleted (default)
// 	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
// 		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

// 	return sql.Where(softDeleted)
// }
