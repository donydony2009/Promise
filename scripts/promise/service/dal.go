package promises

import (
	"encoding/hex"
	"strconv"

	"github.com/donydony2009/Promise/scripts/mysql"
	"github.com/satori/go.uuid"
)

type Dal struct {
	mySQL          *mysql.MySQL
	versionManager *mysql.VersionManager
}

func CreateDal(serviceName string) *Dal {
	dal := new(Dal)
	dal.mySQL = mysql.CreateConnection()
	dal.versionManager = mysql.CreateVersionManager(dal.mySQL, serviceName)
	dal.addDBVersions()
	dal.versionManager.UpgradeToLatest()

	return dal
}

func (f *Dal) addDBVersions() {
	var v1 mysql.DBVersion
	v1.Upgrade = SQLUpgradeV1
	v1.Downgrade = SQLDowngradeV1
	f.versionManager.AddVersion(v1)
}

func (f *Dal) fromUUIDToHex(id uuid.UUID) string {
	return hex.EncodeToString(id.Bytes())
}

func (f *Dal) AddPromise(user_id uuid.UUID, title string, description string) {
	hexUserId := f.fromUUIDToHex(user_id)
	f.mySQL.Conn.Exec("INSERT INTO promises(user_id, title, description)" +
		" VALUES(X'" + hexUserId + "','" + title + "','" + description + "')")
}

func (f *Dal) EditPromise(promiseId int, title string, description string, promisedTo uuid.UUID, status Status, privacy Privacy) {
	hexPromisedTo := f.fromUUIDToHex(promisedTo)
	f.mySQL.Conn.Exec("UPDATE promises SET " +
		"title = '" + title + "'," +
		"description = '" + description + "'," +
		"promised_to = X'" + hexPromisedTo + "'," +
		"status = " + strconv.Itoa(int(status)) + "," +
		"privacy = " + strconv.Itoa(int(privacy)) +
		" WHERE promise_id = " + strconv.Itoa(promiseId))
}

func (f *Dal) DeletePromise(promiseId int) {
	f.mySQL.Conn.Exec("DELETE FROM promises WHERE promise_id = " + strconv.Itoa(promiseId))
}

func (f *Dal) GetPromises(user_id uuid.UUID) []LEPromise {
	var results []LEPromise
	hexId := f.fromUUIDToHex(user_id)
	rows, _ := f.mySQL.Conn.Query("SELECT * FROM promises WHERE user_id = X'" + hexId + "'")
	defer rows.Close()
	for rows.Next() {
		var result LEPromise
		var userId []byte
		var promisedTo []byte
		rows.Scan(&result.PromiseId, &result.Title, &result.Description, &userId, &promisedTo, &result.Status, &result.Privacy)
		result.UserId, _ = uuid.FromBytes(userId)
		result.PromisedTo = uuid.FromBytesOrNil(promisedTo)
		results = append(results, result)
	}

	return results
}

func (f *Dal) GetPromise(promiseId int) LEPromise {
	row := f.mySQL.Conn.QueryRow("SELECT * FROM promises WHERE promise_id = " + strconv.Itoa(promiseId))
	var result LEPromise
	var userId []byte
	var promisedTo []byte
	row.Scan(&result.PromiseId, &result.Title, &result.Description, &userId, &promisedTo, &result.Status, &result.Privacy)
	result.UserId, _ = uuid.FromBytes(userId)
	result.PromisedTo = uuid.FromBytesOrNil(promisedTo)
	return result
}

func (f *Dal) CloseConnection() {
	f.mySQL.Close()
}
