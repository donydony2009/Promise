package authentication
import (
	"github.com/satori/go.uuid"
	"github.com/donydony2009/Promise/scripts/mysql"
	"encoding/hex"
	"time"
	"database/sql"
	"strconv"
	"fmt"
)

const ticketSize = 32
const ticketLifetime = 3 * 60 * 60
type Dal struct {
	mySQL *mysql.MySQL
	versionManager *mysql.VersionManager
}

func CreateDal(serviceName string) *Dal{
	dal := new(Dal)
	dal.mySQL = mysql.CreateConnection()
	dal.versionManager = mysql.CreateVersionManager(dal.mySQL, serviceName)
	dal.addDBVersions()
	dal.versionManager.UpgradeToLatest()
	
	return dal
}

func (f *Dal) addDBVersions(){
	var v1 mysql.DBVersion
	v1.Upgrade = SQLUpgradeV1
	v1.Downgrade = SQLDowngradeV1
	f.versionManager.AddVersion(v1)

	var v2 mysql.DBVersion
	v2.Upgrade = SQLUpgradeV2
	v2.Downgrade = SQLDowngradeV2
	f.versionManager.AddVersion(v2)

	var v3 mysql.DBVersion
	v3.Upgrade = SQLUpgradeV3
	v3.Downgrade = SQLDowngradeV3
	f.versionManager.AddVersion(v3)
}

func (f *Dal) CloseConnection(){
	f.mySQL.Close()
}

func (f *Dal) fromUUIDToHex(id uuid.UUID) string{
	return hex.EncodeToString(id.Bytes())
}

func (f *Dal) CreateAccount(id uuid.UUID, username string, password string, salt string, email string){
	hexId := f.fromUUIDToHex(id)
	fmt.Println("INSERT INTO user_info(id, username, password, salt, email)" +
		" VALUES(X'" + hexId + "','" + username + "','" +
		password + "','" + salt + "','" + email + "')")
	f.mySQL.Conn.Exec("INSERT INTO user_info(id, username, password, salt, email)" +
		" VALUES(X'" + hexId + "','" + username + "','" +
		password + "','" + salt + "','" + email + "')")
}

func (f *Dal) GetUserInfo(username string, password string) UserInfo{
	var userInfo UserInfo
	var idBytes []byte
	result := f.mySQL.Conn.QueryRow("SELECT id, username, password, salt FROM user_info WHERE username='" + username + "'")
	result.Scan(&idBytes, &userInfo.username, &userInfo.password, &userInfo.salt)
	userInfo.id, _ = uuid.FromBytes(idBytes)
	return userInfo
}

func (f *Dal) generateTicket() string{
	return RandStringBytes(ticketSize)
}

func (f *Dal) RefreshTicket(userId uuid.UUID) string{
	var ticket string
	var creationDate time.Time
	_, timezoneOffset := time.Now().Zone()
	lifetime,_ := time.ParseDuration(strconv.Itoa(ticketLifetime - timezoneOffset) + "s")
	
	hexId := f.fromUUIDToHex(userId)
	result := f.mySQL.Conn.QueryRow("SELECT ticket, creation_date FROM user_tickets WHERE user=X'" + hexId + "'")
	
	err := result.Scan(&ticket, &creationDate)
	if (err == sql.ErrNoRows){
		ticket = f.generateTicket()
		f.mySQL.Conn.Exec("INSERT INTO user_tickets(user, ticket)" +
			" VALUES(X'" + hexId + "','" + ticket + "')")
	}else{
		creationDate = creationDate.Add(lifetime)
		if creationDate.Before(time.Now()){
			ticket = f.generateTicket()
		}
		f.mySQL.Conn.Exec("UPDATE user_tickets SET " +
			" ticket = '" + ticket + "'," +
			" creation_date = null " + 
			"WHERE user = X'" + hexId + "'")
	}
	return ticket
}

func (f *Dal) CheckTicket(ticket string) uuid.UUID{
	var creationDate time.Time
	var userID []byte
	_, timezoneOffset := time.Now().Zone()
	lifetime,_ := time.ParseDuration(strconv.Itoa(ticketLifetime - timezoneOffset) + "s")

	result := f.mySQL.Conn.QueryRow("SELECT user, creation_date FROM user_tickets WHERE ticket='" + ticket + "'")
	err := result.Scan(&userID, &creationDate)
	if (err != sql.ErrNoRows){
		creationDate = creationDate.Add(lifetime)
		if creationDate.After(time.Now()){
			user, _ := uuid.FromBytes(userID)
			return user
		}
	}

	return uuid.Nil
}
	


