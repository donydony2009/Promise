package authentication
import "github.com/satori/go.uuid"

type UserInfo struct{
	id uuid.UUID
	username string
	password string
	salt string
}

type AuthError struct{

}

func (f AuthError) Error() string {
	return "There was an authentication error"
}