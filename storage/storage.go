package storage

import (
	"github.com/SaidovZohid/note_user_service/storage/postgres"
	"github.com/SaidovZohid/note_user_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type StoragePg struct {
	userRepo repo.UserStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo: postgres.NewUserStorage(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}
