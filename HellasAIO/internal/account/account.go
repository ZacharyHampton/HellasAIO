package account

import (
	"errors"
	"sync"
)

var (
	accountMutex = sync.RWMutex{}

	AccountDoesNotExistErr = errors.New("accounut does not exist")
	accounts               = make(map[int]map[string]*Account)
)

// DoesAccountExist checks if an account exists
func DoesAccountExist(siteId int, email string) bool {
	accountMutex.RLock()
	defer accountMutex.RUnlock()

	_, ok := accounts[siteId][email]
	return ok
}

// CreateAccount creates a new profile
func CreateAccount(account *Account) (int, string) {
	accountMutex.RLock()
	defer accountMutex.RUnlock()

	if len(accounts) == 0 || accounts == nil {
		accounts = map[int]map[string]*Account{account.SiteId: {account.Email: account}}
	}

	if accounts[account.SiteId] == nil {
		accounts[account.SiteId] = map[string]*Account{account.Email: account}
	} else {
		accounts[account.SiteId][account.Email] = account
	}

	return account.SiteId, account.Email
}

// RemoveAccount removes an account
func RemoveAccount(siteId int, email string) error {
	if !DoesAccountExist(siteId, email) {
		return AccountDoesNotExistErr
	}

	accountMutex.RLock()
	defer accountMutex.RUnlock()

	delete(accounts[siteId], email)

	return nil
}

// GetAccount gets an account
func GetAccount(siteId int, email string) (*Account, error) {
	if !DoesAccountExist(siteId, email) {
		return &Account{}, AccountDoesNotExistErr
	}

	accountMutex.RLock()
	defer accountMutex.RUnlock()

	return accounts[siteId][email], nil
}

// GetAllAccountEmailsInSite gets all profile ids in a specific site
func GetAllAccountEmailsInSite(siteId int) []string {
	var emails []string

	for email := range accounts[siteId] {
		emails = append(emails, email)
	}

	return emails
}
