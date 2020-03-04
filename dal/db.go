package dal

import (
	"database/sql"
	"fmt"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
)

// SQLDB sql连接池.
var SQLDB *sql.DB

// DbInit init.
func DbInit() {
	SQLDB, _ = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", config.BaseConf.DbUser, config.BaseConf.DbPwd, config.BaseConf.DbAddr, config.BaseConf.DbDatabase))
	err := SQLDB.Ping()
	if err != nil {
		log.Println("sql-db ping error, err = ", err)
	}
}

// checkSQLDB sql检查
func checkSQLDB(Db *sql.DB) error {
	err := Db.Ping()
	if err != nil {
		log.Println("Error!!! Sql DB Ping Error! info=")
		log.Println(fmt.Sprintf("%s:%s@tcp(%s)/%s", config.BaseConf.DbUser, config.BaseConf.DbPwd, config.BaseConf.DbAddr, config.BaseConf.DbDatabase))
		return err
	}
	return nil
}

// DbGetUserInfoByName get user info by name in sql.
func DbGetUserInfoByName(userName string, Db *sql.DB) (UserInfo, error) {
	if err := checkSQLDB(Db); err != nil {
		return UserInfo{}, err
	}
	var rs UserInfo
	row := Db.QueryRow("select user_name, pwd, nick_name, picture from user_info where user_name = ? and is_enabled = 1", userName)
	err := row.Scan(&rs.Name, &rs.Pwd, &rs.NickName, &rs.Picture)
	if err != nil && err != sql.ErrNoRows {
		log.Println("DbGetUserInfoByName, row-scan-error, err = ", err)
		return rs, err
	}
	return rs, nil
}

// DbInsertUserInfo insert user info to sql.
func DbInsertUserInfo(userName, pwd string, Db *sql.DB) error {
	if err := checkSQLDB(Db); err != nil {
		return err
	}
	inserSQL := fmt.Sprintf("insert into user_info set user_name='%s', pwd='%s', nick_name=''", userName, pwd)
	_, err := Db.Exec(inserSQL)
	if err != nil {
		log.Println("DbInsertUserInfo, inserSQL-error, err = ", err)
		return err
	}
	return nil
}

// DbModifyUserInfo modify user info by name in sql.
func DbModifyUserInfo(userName, pwd, nickName, picture string, Db *sql.DB) error {
	if err := checkSQLDB(Db); err != nil {
		return err
	}
	modifySQL := ""
	if picture != "" {
		modifySQL = fmt.Sprintf("update user_info set nick_name = '%s', picture = '%s' where user_name = '%s'", nickName, picture, userName)
	} else {
		modifySQL = fmt.Sprintf("update user_info set nick_name = '%s' where user_name = '%s'", nickName, userName)
	}
	_, err := Db.Exec(modifySQL)
	if err != nil {
		log.Println("DbModifyUserInfo, modifySql-error, err = ", err)
		return err
	}
	return nil
}
