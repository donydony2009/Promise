package promises

import (
	"github.com/donydony2009/Promise/scripts/mysql"
)

func SQLUpgradeV1(mySql *mysql.MySQL) error {
	_, err := mySql.Conn.Exec("CREATE TABLE `promises` (" +
		"`promise_id` INT NOT NULL AUTO_INCREMENT," +
		"`title` VARCHAR(50) NOT NULL," +
		"`description` VARCHAR(500) NULL," +
		"`user_id` BINARY(16) NOT NULL," +
		"`promised_to` BINARY(16)," +
		"`status` TINYINT NOT NULL DEFAULT 0," +
		"`privacy` TINYINT NOT NULL DEFAULT 0," +
		"INDEX `user_id` (`user_id`)," +
		"INDEX `promised_to` (`promised_to`)," +
		"PRIMARY KEY (`promise_id`)" +
		")" +
		"ENGINE=InnoDB;")
	return err
}

func SQLDowngradeV1(mySql *mysql.MySQL) error {
	_, err := mySql.Conn.Exec("DROP TABLE `promises`")
	return err
}
