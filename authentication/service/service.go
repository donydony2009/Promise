package authentication

import "github.com/satori/go.uuid"
import "crypto/sha256"
import "encoding/hex"
import "hash"
import "sync"
import "fmt"

const saltSize = 16

type Authentication struct{
	dal *Dal
	name string
	sha256 hash.Hash
}

var instance *Authentication
var once sync.Once

func GetServiceInstance() *Authentication{
	once.Do(func() {
        instance = createAuthService()
    })
    return instance
}

func createAuthService() *Authentication{
	auth := new(Authentication)
	auth.name = "Authentication"
	auth.dal = CreateDal(auth.name)
	auth.sha256 = sha256.New()
	
	return auth
}

func (f *Authentication) GenerateSalt() string{
	return RandStringBytes(saltSize)
}

func (f *Authentication) CreateAccount(username string, password string, email string){
	var userUUID uuid.UUID
	salt := f.GenerateSalt()
	hashedPassword := hex.EncodeToString(f.sha256.Sum([]byte(password + salt)))
	fmt.Println(username, password, salt)
	userUUID, _ = uuid.NewV4()
	f.dal.CreateAccount(userUUID, username, hashedPassword, salt, email)
}

func (f *Authentication) Login(username string, password string) string{
	
	userInfo := f.dal.GetUserInfo(username, password)
	hashedPassword := hex.EncodeToString(f.sha256.Sum([]byte(password + userInfo.salt)))
	fmt.Println(username, password, hashedPassword, userInfo.password, userInfo.salt, password + userInfo.salt)
	if hashedPassword == userInfo.password{
		return f.dal.RefreshTicket(userInfo.id)
	}
	return "error"
}

func (f *Authentication) CheckTicket(ticket string) uuid.UUID{
	return f.dal.CheckTicket(ticket)
}

func (f *Authentication) Close(){
	f.dal.CloseConnection()
}

