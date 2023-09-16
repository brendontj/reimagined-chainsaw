package gateway

type Db interface {
	AddNewUser(username, password string) error
	FindUserPassword(username string) (string, error)
}
