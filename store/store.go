package store

import (
	"github.com/jinzhu/gorm"
	"log"
)

type Store interface {
	Transaction(t func(txStore Store) error)
	Users() UserStore
	LinkedAccounts() LinkedAccountStore
	Devices() DeviceStore
	EncryptionKeys() EncryptionKeyStore
	Messages() MessageStore
	NotificationKeys() NotificationKeyStore
	UserTokens() UserTokenStore
	VerificationTokens() VerificationTokenStore
}

type store struct {
	db                 *gorm.DB
	users              userStore
	linkedAccounts     linkedAccountStore
	devices            deviceStore
	encryptionKeys     encryptionKeyStore
	messages           messageStore
	notificationKeys   notificationKeyStore
	userTokens         userTokenStore
	verificationTokens verificationTokenStore
}

func (s *store) Transaction(t func(txStore Store) error) {
	tx := s.db.Begin()
	txStore := New(tx)
	if err := t(txStore); err != nil {
		tx.Rollback()
		log.Println("Database rollback caused by %v", err)
		return
	}
	tx.Commit()
}

func (s *store) Users() UserStore                           { return s.users }
func (s *store) LinkedAccounts() LinkedAccountStore         { return s.linkedAccounts }
func (s *store) Devices() DeviceStore                       { return s.devices }
func (s *store) EncryptionKeys() EncryptionKeyStore         { return s.encryptionKeys }
func (s *store) Messages() MessageStore                     { return s.messages }
func (s *store) NotificationKeys() NotificationKeyStore     { return s.notificationKeys }
func (s *store) UserTokens() UserTokenStore                 { return s.userTokens }
func (s *store) VerificationTokens() VerificationTokenStore { return s.verificationTokens }

func New(db *gorm.DB) Store {
	return &store{
		db:                 db,
		users:              userStore{db},
		linkedAccounts:     linkedAccountStore{db},
		devices:            deviceStore{db},
		encryptionKeys:     encryptionKeyStore{db},
		messages:           messageStore{db},
		notificationKeys:   notificationKeyStore{db},
		userTokens:         userTokenStore{db},
		verificationTokens: verificationTokenStore{db},
	}
}
