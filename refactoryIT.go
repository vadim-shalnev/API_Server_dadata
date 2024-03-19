package main

import (
	_ "github.com/lib/pq"
)

/*
type Repository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]User, error)
}
type User struct {
	ID    int
	Name  string
	Email string
}
type UserRepository struct {
	db *sql.DB
}

type Conditions struct {
	Limit  int
	Offset int
}

func main() {

	// Подключение к базе данных
	db, err := ConnectToDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close() // Закрытие соединения с базой данных при завершении работы

	// Создание нового экземпляра репозитория пользователей
	userRepo := NewUserRepository(db)

	err = userRepo.CreateTableIfNotExists(context.Background())
	if err != nil {
		log.Fatal("Ошибка при создании таблицы пользователей:", err)
	}
	// Создание нового пользователя
	newUser := User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	err = userRepo.Create(context.Background(), newUser)
	if err != nil {
		log.Fatal("Ошибка при создании пользователя:", err)
	}

	// Получение пользователя по ID
	userByID, err := userRepo.GetByID(context.Background(), "1")
	if err != nil {
		log.Fatal("Ошибка при получении пользователя по ID:", err)
	}
	fmt.Println("Пользователь по ID:", userByID)

	// Обновление информации о пользователе
	userByID.Name = "Jane Doe"
	userByID.Email = "jane.doe@example.com"
	err = userRepo.Update(context.Background(), userByID)
	if err != nil {
		log.Fatal("Ошибка при обновлении информации о пользователе:", err)
	}

	// Удаление пользователя
	err = userRepo.Delete(context.Background(), "1")
	if err != nil {
		log.Fatal("Ошибка при удалении пользователя:", err)
	}

	// Получение списка пользователей с пагинацией
	conditions := Conditions{
		Limit:  10, // количество элементов на странице
		Offset: 0,  // смещение начала выборки данных
	}
	users, count, err := userRepo.List(context.Background(), conditions)
	if err != nil {
		log.Fatal("Ошибка при получении списка пользователей с пагинацией:", err)
	}
	fmt.Println("Список пользователей:", users)
	fmt.Println("Общее количество пользователей:", count)

}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password= dbname=mydb sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Создание нового экземпляра репозитория пользователей
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateTableIfNotExists(ctx context.Context) error {
	// Проверяем существование таблицы
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (
        SELECT 1
        FROM   information_schema.tables
        WHERE  table_schema = 'public'
        AND    table_name = 'users'
    )`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %v", err)
	}

	// Если таблица не существует, создаем ее
	if !exists {
		_, err := r.db.ExecContext(ctx, `
            CREATE TABLE users (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                email VARCHAR(255) NOT NULL
            )
        `)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	switch {
	case err == sql.ErrNoRows:
		return User{}, errors.New("user not found")
	case err != nil:
		return User{}, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, c Conditions) ([]User, int, error) {
	var users []User
	var count int

	// Сначала получаем общее количество пользователей
	countQuery := "SELECT COUNT(*) FROM users"
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %v", err)
	}

	// Затем получаем список пользователей с учетом пагинации
	query := "SELECT id, name, email FROM users"
	if c.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", c.Limit)
	}
	if c.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", c.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over user rows: %v", err)
	}

	return users, count, nil
}

*/
