package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	models "github.com/vadim-shalnev/API_Server_dadata/Models"
	"gitlab.com/ptflp/goboilerplate/internal/db/adapter"
	"gitlab.com/ptflp/goboilerplate/internal/infrastructure/cache"
	"gitlab.com/ptflp/goboilerplate/internal/infrastructure/db/scanner"
	"time"
)

// UserStorage - хранилище пользователей
type UserStorage struct {
	adapter *adapter.SQLAdapter
	cache   cache.Cache
}

const (
	userCacheKey     = "user:%d"
	userCacheTTL     = 15
	userCacheTimeout = 50
)

// NewUserStorage - конструктор хранилища пользователей
func NewUserStorage(sqlAdapter *adapter.SQLAdapter, cache cache.Cache) *UserStorage {
	return &UserStorage{adapter: sqlAdapter, cache: cache}
}

// Create - создание пользователя в БД
func (s *UserStorage) Create(ctx context.Context, u models.UserDTO) (int, error) {
	err := s.adapter.Create(ctx, &u)

	return 0, err
}

// Update - обновление пользователя в БД
func (s *UserStorage) Update(ctx context.Context, u models.UserDTO) error {
	err := s.adapter.Update(ctx, &u, adapter.Condition{
		Equal: sq.Eq{
			"id": u.GetID(),
		},
	}, scanner.Update)
	if err != nil {
		return err
	}
	go func() {
		_ = s.cache.Expire(ctx, fmt.Sprintf(userCacheKey, u.GetID()), 0)
	}()

	return nil
}

// GetByID - получение пользователя по ID из БД
func (s *UserStorage) GetByID(ctx context.Context, userID int) (models.UserDTO, error) {
	var dto models.UserDTO
	var err error

	// создаем контекст с таймаутом, чтобы не ждать ответа от кеша
	timeout, cancel := context.WithTimeout(context.Background(), userCacheTimeout*time.Millisecond)
	defer cancel()
	// получаем данные из кеша по ключу user:{id}
	err = s.cache.Get(timeout, fmt.Sprintf(userCacheKey, userID), &dto)
	// если данные есть в кеше, то возвращаем их
	if err == nil {
		return dto, nil
	}

	// если данных нет в кеше, то получаем их из БД
	var list []models.UserDTO
	err = s.adapter.List(ctx, &list, dto.TableName(), adapter.Condition{
		Equal: map[string]interface{}{
			"id": userID,
		},
	})
	if err != nil {
		return models.UserDTO{}, err
	}
	if len(list) < 1 {
		return models.UserDTO{}, fmt.Errorf("user storage: GetByID not found")
	}

	// запускаем горутину для записи данных в кеш
	go func() {
		// создаем контекст с таймаутом, чтобы не ждать ответа от кеша
		timeout, cancel = context.WithTimeout(context.Background(), userCacheTimeout*time.Millisecond)
		defer cancel()
		s.cache.Set(timeout, fmt.Sprintf(userCacheKey, userID), list[0], userCacheTTL*time.Minute)
	}()

	return list[0], nil
}

// GetByIDs - получение пользователей по ID из БД
func (s *UserStorage) GetByIDs(ctx context.Context, ids []int) ([]models.UserDTO, error) {
	var users []models.UserDTO
	err := s.adapter.List(ctx, &users, "users", adapter.Condition{
		Equal: map[string]interface{}{
			"id": ids,
		},
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) GetByFilter(ctx context.Context, condition adapter.Condition) ([]models.UserDTO, error) {
	var users []models.UserDTO
	err := s.adapter.List(ctx, &users, "users", condition)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail - получение пользователя по email из БД
func (s *UserStorage) GetByEmail(ctx context.Context, email string) (models.UserDTO, error) {
	var users []models.UserDTO
	err := s.adapter.List(ctx, &users, "users", adapter.Condition{
		Equal: map[string]interface{}{
			"email": email,
		},
	})
	if err != nil {
		return models.UserDTO{}, err
	}
	if len(users) < 1 {
		return models.UserDTO{}, fmt.Errorf("user with email %s not found", email)
	}

	return users[0], nil
}
