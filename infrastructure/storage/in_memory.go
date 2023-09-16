package storage

import "errors"

const EmptyPassword = ""

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserDoesntExist  = errors.New("user doesn't exist")
)

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{data: make(map[string]string)}
}

func (i *InMemoryStorage) AddNewUser(username, password string) error {
	_, exist := i.data[username]
	if exist {
		return ErrUserAlreadyExist
	}

	i.data[username] = password
	return nil
}

func (i *InMemoryStorage) FindUserPassword(username string) (string, error) {
	pw, exist := i.data[username]
	if !exist {
		return EmptyPassword, ErrUserDoesntExist
	}

	return pw, nil
}
