package postgres

import (
	"database/sql"
	"fmt"
	"user-service/internal/models"
    sq "github.com/Masterminds/squirrel"
)

type UserRepo interface {
	Add(user models.User) (int64, error)
	Get(id int64) (models.User, error)
	GetList() ([]models.User, error)
	Update(user models.User) error
	Remove(id int64) error
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (us *UserStorage) Add(user models.User) (int64, error) {
	query := sq.Insert("users").Columns("name", "address", "phone", "created_at").
        Values(user.Name, user.Address, user.Phone, user.CreatedAt).
        Suffix("RETURNING \"id\"").
        RunWith(us.db).
        PlaceholderFormat(sq.Dollar)
    
    err := query.QueryRow().Scan(&user.Id)
    if err != nil {
        return 0, fmt.Errorf("couldn't add user. Error: %w", err)
    }

    return user.Id, nil
}

func (us *UserStorage) Get(id int64) (models.User, error) {
    var user models.User
    err := sq.Select("id", "name", "address", "phone", "created_at").From("users").
        Where(sq.Eq{"id": id}).
        RunWith(us.db).
        PlaceholderFormat(sq.Dollar).
        QueryRow().Scan(&user.Id, &user.Name, &user.Address, &user.Phone, &user.CreatedAt)

    if err == sql.ErrNoRows {
        return user, fmt.Errorf("user not found. Error: %w", err)
    } else if err != nil {
        return user, fmt.Errorf("couldn't get user. Error: %w", err)
    }
    return user, nil
}

func (us *UserStorage) GetList() ([]models.User, error) {
	rows, err := sq.Select("id", "name", "address", "phone", "created_at").From("users").
        RunWith(us.db).
        Query()

    if err != nil {
        return nil, fmt.Errorf("couldn't get users. Error: %w", err)
    }
    defer rows.Close()

    users := make([]models.User, 0)
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.Id, &user.Name, &user.Address, &user.Phone, &user.CreatedAt)
        if err != nil {
            continue
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("couldn't get users. Error: %w", err)
    }
    return users, nil
}

func (us *UserStorage) Update(user models.User) error {
	result, err := sq.Update("users").
        Set("name", user.Name).
        Set("address", user.Address).
        Set("phone", user.Phone).
        Set("created_at", user.CreatedAt).
        Where(sq.Eq{"id": user.Id}).
        RunWith(us.db).
        PlaceholderFormat(sq.Dollar).
        Exec()
    if err != nil {
        return fmt.Errorf("couldn't update users. Error: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("nothing to update. Error: %w", err)
    }
    return nil
}

func (us *UserStorage) Remove(id int64) error {
	result, err := sq.Delete("users").
        Where(sq.Eq{"id": id}).
        RunWith(us.db).
        PlaceholderFormat(sq.Dollar).
        Exec()
    if err != nil {
        return fmt.Errorf("couldn't delete users. Error: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("nothing to delete. Error: %w", err)
    }
    return nil
}
