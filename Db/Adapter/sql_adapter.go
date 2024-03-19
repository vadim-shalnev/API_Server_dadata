package Adapter

import (
	"context"
	"fmt"
	"gitlab.com/ptflp/goboilerplate/config"
	"gitlab.com/ptflp/goboilerplate/internal/infrastructure/db/scanner"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

// SQLAdapter - адаптер для работы с БД
type SQLAdapter struct {
	db         *sqlx.DB
	scanner    scanner.Scanner
	sqlBuilder sq.StatementBuilderType
}

// NewSqlAdapter - конструктор адаптера для работы с БД
func NewSqlAdapter(db *sqlx.DB, dbConf config.DB, scanner scanner.Scanner) *SQLAdapter {
	var builder sq.StatementBuilderType
	if dbConf.Driver != "mysql" {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	}

	return &SQLAdapter{db: db, scanner: scanner, sqlBuilder: builder}
}

// Create - создание записи в БД
func (s *SQLAdapter) Create(ctx context.Context, entity scanner.Tabler) error {
	createFields := s.scanner.OperationFields(entity.TableName(), scanner.Create)
	createFieldsPointers := GetFieldsPointers(entity, "create")

	queryRaw := s.sqlBuilder.Insert(entity.TableName()).Columns(createFields...).Values(createFieldsPointers...)

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)

	return err
}

// Upsert - создание записи в БД или обновление ее, если она уже существует
func (s *SQLAdapter) Upsert(ctx context.Context, entities []scanner.Tabler) error {
	if len(entities) < 1 {
		return fmt.Errorf("SQL adapter: zero entities passed")
	}
	createFields := s.scanner.OperationFields(entities[0].TableName(), scanner.Create)
	queryRaw := s.sqlBuilder.Insert(entities[0].TableName()).Columns(createFields...)

	for i := range entities {
		createFieldsPointers := GetFieldsPointers(entities[i], "create")
		queryRaw = queryRaw.Values(createFieldsPointers...)
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	conflictFields := s.scanner.OperationFields(entities[0].TableName(), scanner.Conflict)
	if len(conflictFields) > 0 {
		query = query + " ON CONFLICT (%s)"
		query = fmt.Sprintf(query, strings.Join(conflictFields, ","))
		query = query + " DO UPDATE SET"
	}
	upsertFields := s.scanner.OperationFields(entities[0].TableName(), scanner.Upsert)
	for _, field := range upsertFields {
		query += fmt.Sprintf(" %s = excluded.%s,", field, field)
	}
	if len(upsertFields) > 0 {
		query = query[0 : len(query)-1]
	}

	_, err = s.db.ExecContext(ctx, query, args...)

	return err
}

// buildSelect - сборка запроса на выборку данных из БД
func (s *SQLAdapter) buildSelect(tableName string, condition Condition, fields ...string) (string, []interface{}, error) {
	if condition.ForUpdate {
		temp := []string{"FOR UPDATE"}
		temp = append(temp, fields...)
		fields = temp
	}
	queryRaw := s.sqlBuilder.Select(fields...).From(tableName)

	if condition.Equal != nil {
		for field, args := range condition.Equal {
			queryRaw = queryRaw.Where(sq.Eq{field: args})
		}
	}

	if condition.NotEqual != nil {
		for field, args := range condition.NotEqual {
			queryRaw = queryRaw.Where(sq.NotEq{field: args})
		}
	}

	if condition.Order != nil {
		for _, order := range condition.Order {
			direction := "DESC"
			if order.Asc {
				direction = "ASC"
			}
			queryRaw = queryRaw.OrderBy(fmt.Sprintf("%s %s", order.Field, direction))
		}
	}

	if condition.LimitOffset != nil {
		if condition.LimitOffset.Limit > 0 {
			queryRaw.Limit(uint64(condition.LimitOffset.Limit))
		}
		if condition.LimitOffset.Offset > 0 {
			queryRaw.Offset(uint64(condition.LimitOffset.Offset))
		}
	}

	return queryRaw.ToSql()
}

// GetCount - получение количества записей в БД
func (s *SQLAdapter) GetCount(ctx context.Context, entity scanner.Tabler, condition Condition) (uint64, error) {
	query, args, err := s.buildSelect(entity.TableName(), condition, "COUNT(*)")
	if err != nil {
		return 0, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	var count uint64
	// iterate over each row
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	// check the error from rows
	err = rows.Err()

	return count, err
}

// List - получение списка записей из БД
func (s *SQLAdapter) List(ctx context.Context, dest interface{}, tableName string, condition Condition) error {
	fields := s.scanner.OperationFields(tableName, scanner.AllFields)
	query, args, err := s.buildSelect(tableName, condition, fields...)
	if err != nil {
		return err
	}

	if err = s.db.SelectContext(ctx, dest, query, args...); err != nil {
		return err
	}

	return nil
}

// Update - обновление записи в БД
func (s *SQLAdapter) Update(ctx context.Context, entity scanner.Tabler, condition Condition, operation string) error {
	ent := entity
	updateFields := s.scanner.OperationFields(entity.TableName(), operation)

	updateFieldsPointers := GetFieldsPointers(entity, operation)

	updateRaw := s.sqlBuilder.Update(ent.TableName())

	if condition.Equal != nil {
		for field, args := range condition.Equal {
			updateRaw = updateRaw.Where(sq.Eq{field: args})
		}
	}

	if condition.NotEqual != nil {
		for field, args := range condition.NotEqual {
			updateRaw = updateRaw.Where(sq.NotEq{field: args})
		}
	}

	for i := range updateFields {
		updateRaw = updateRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err := updateRaw.ToSql()
	if err != nil {
		return err
	}

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}
