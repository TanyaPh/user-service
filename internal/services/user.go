package services

import "user-service/internal/models"

type UserRepo interface {
	Add(user models.User) (int64, error)
	Get(id int64) (models.User, error)
	GetList() ([]models.User, error)
	Update(models.User) error
	Remove(id int64) error
}

type User struct {
	repo UserRepo
}

func NewUserService(userRepo UserRepo) *User {
	return &User{
		repo: userRepo,
	}
}

func (us *User) CreateUser(user models.User) (int64, error) {
	return us.repo.Add(user)
}

func (us *User) GetUserById(userID int64) (models.User, error) {
	return us.repo.Get(userID)
}

func (us *User) GetUsers() ([]models.User, error) {
	return us.repo.GetList()
}

func (us *User) EditUser(userID int64, name, adderss, phone *string) error {
	user, err := us.repo.Get(userID)
	if err != nil {
		return err
	}

	if name != nil {
		user.Name = *name
	}
	if adderss != nil {
		user.Address = *adderss
	}
	if phone != nil {
		user.Phone = *phone
	}

	return us.repo.Update(user)
}

func (us *User) DeleteUser(userID int64) error {
	return us.repo.Remove(userID)
}
