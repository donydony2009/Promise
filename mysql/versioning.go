/*
Package mysql implements a simple way of obtaining a SQL Connection
*/
package mysql

import (
	"database/sql"
	"strconv"

)
const versioningTableName = "service_versions"
type DatabaseFunc func(mySql *MySQL) error
type DBVersion struct{
	Upgrade DatabaseFunc
	Downgrade DatabaseFunc
}

type VersionManager struct{
	mySql *MySQL
	serviceName string
	versions []DBVersion
}

func CreateVersionManager(mySql *MySQL, serviceName string) *VersionManager{
	versionManager := new(VersionManager)
	versionManager.mySql = mySql
	versionManager.serviceName = serviceName
	if(!versionManager.mySql.DoesTableExist("service_versions")){
		versionManager.mySql.Conn.Exec(
			"CREATE TABLE `" + versioningTableName + "` (" +
				"`name` VARCHAR(100) NOT NULL," +
				"`version` INT NOT NULL," +
				"PRIMARY KEY (`name`)" +
			")" +
			"ENGINE=InnoDB;")
	}
	return versionManager
}

func (f *VersionManager) VersionCount() int{
	return len(f.versions)
}

func (f *VersionManager) AddVersion(version DBVersion){
	f.versions = append(f.versions, version)
}

func (f *VersionManager) AddVersions(versions []DBVersion){
	f.versions = append(f.versions, versions...)
}

func (f *VersionManager) UpgradeToLatest(){
	version := 0

	row := f.mySql.Conn.QueryRow("SELECT version FROM " + versioningTableName + " WHERE name = \"" + f.serviceName + "\"")
	err := row.Scan(&version)
	if err == sql.ErrNoRows {
		f.mySql.Conn.Exec("INSERT INTO " + versioningTableName + "(name, version) VALUES(\"" + 
			f.serviceName + "\",0)")
	}
	for i := version; i < len(f.versions); i++ {
		err = f.versions[i].Upgrade(f.mySql)
		if(err != nil){
			panic(err)
		}
	}

	f.mySql.Conn.Exec("UPDATE " + versioningTableName + " SET version=" + strconv.Itoa(len(f.versions)) + " WHERE name=\"" + f.serviceName + "\"")

}