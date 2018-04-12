package authentication
import (
	"mysql"
	"fmt"
)

func SQLUpgradeV1(mySql *mysql.MySQL) error{
	fmt.Println("Creating Auth DB")
	_, err := mySql.Conn.Exec(
		"CREATE TABLE `user_info` (" +
			"`id` BINARY(16) NOT NULL," +
			"`username` VARCHAR(50) NOT NULL DEFAULT '0'," +
			"`password` VARCHAR(256) NOT NULL DEFAULT '0'," +
			"`salt` VARCHAR(50) NOT NULL DEFAULT '0'," +
			"`email` VARCHAR(100) NOT NULL DEFAULT '0'," +
			"PRIMARY KEY (`id`)" +
		")" +
		"ENGINE=InnoDB;")
	return err
}

func SQLDowngradeV1(mySql *mysql.MySQL) error{
	_,err := mySql.Conn.Exec("DROP TABLE `user_info`")
	return err
}

func SQLUpgradeV2(mySql *mysql.MySQL) error{
	_,err := mySql.Conn.Exec("ALTER TABLE user_info ADD CONSTRAINT unique_username UNIQUE KEY (username)")
	if err != nil{
		return err
	}
	_, err = mySql.Conn.Exec("ALTER TABLE user_info ADD CONSTRAINT unique_email UNIQUE KEY (email)")
	return err
}

func SQLDowngradeV2(mySql *mysql.MySQL) error{
	_,err := mySql.Conn.Exec("ALTER TABLE user_info DROP CONSTRAINT unique_username")
	if err != nil{
		return err
	}
	_,err = mySql.Conn.Exec("ALTER TABLE user_info DROP CONSTRAINT unique_email")
	return err
}

func SQLUpgradeV3(mySql *mysql.MySQL) error{
	_,err := mySql.Conn.Exec(
		"CREATE TABLE `user_tickets` (" +
			"`user` BINARY(16) NOT NULL," +
			"`ticket` VARCHAR(64) NOT NULL," +
			"`creation_date` TIMESTAMP NOT NULL" +
		")" +
		"ENGINE=InnoDB;")
	return err
}

func SQLDowngradeV3(mySql *mysql.MySQL) error{
	_,err := mySql.Conn.Exec("DROP TABLE `user_tickets`")
	return err
}