package user

import (
	"goddns/bolts/crypto"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB, autoMigrate bool) (*Service, error) {
	svc := &Service{
		db: db,
	}
	if autoMigrate {
		return svc, db.AutoMigrate(&User{})
	}
	return svc, nil
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	var user User
	tx := s.db.First(&user, "username = ?", username)
	return &user, tx.Error
}

func (s *Service) MatchUserWithPlainPassword(user *User) (bool, error) {
	matUser, err := s.GetUserByUsername(user.Username)
	if err != nil {
		return false, err
	}
	if err := crypto.CompareHashAndData(matUser.Password, user.Password); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) CreateUser(user *User) error {
	var err error
	user.Password, err = crypto.Hash(user.Password)
	if err != nil {
		return err
	}
	tx := s.db.Create(user)
	return tx.Error
}

func (s *Service) UpdateUser(old, new *User) error {
	user, err := s.GetUserByUsername(old.Username)
	if err != nil {
		return err
	}
	if new.Password != "" {
		new.Password, err = crypto.Hash(new.Password)
	}
	tx := s.db.Model(user).Updates(new)
	return tx.Error
}

func (s *Service) DeleteUser(user *User) error {
	tx := s.db.Delete(user)
	return tx.Error
}
