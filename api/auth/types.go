package auth

type UserCredsHandler interface {
	GetPassFromDB(username string) (string, error)
	GetHashedPass(password string) (string, error)
}

type MongoBasic struct {
	ClientHandler UserCredsHandler
}
