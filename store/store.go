package store

import "github.com/jinzhu/gorm"

type Store interface {
	Transaction(t func(txStore Store) error)
	Users() UserStore
	LinkedAccounts() LinkedAccountStore
	Devices() DeviceStore
	Messages() MessageStore
	NotificationKeys() NotificationKeyStore
	UserTokens() UserTokenStore
	VerificationTokens() VerificationTokenStore
}

type store struct {
	db *gorm.DB
}

func (s *store) Transaction(t func(txStore Store) error) {
	tx := s.db.Begin()
	if err := t(&store{tx}); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}

func (s *store) Users() UserStore {
	return userStore{s.db}
}

func (s *store) LinkedAccounts() LinkedAccountStore {
	return linkedAccountStore{s.db}
}

func (s *store) Devices() DeviceStore {
	return deviceStore{s.db}
}

func (s *store) Messages() MessageStore {
	return messageStore{s.db}
}

func (s *store) NotificationKeys() NotificationKeyStore {
	return notificationKeyStore{s.db}
}

func (s *store) UserTokens() UserTokenStore {
	return userTokenStore{s.db}
}

func (s *store) VerificationTokens() VerificationTokenStore {
	return verificationTokenStore{s.db}
}

func New(db *gorm.DB) Store {
	return &store{db: db}
}
