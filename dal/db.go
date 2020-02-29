package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
)

var SqlDb *sql.DB

func DbInit() {
	SqlDb, _ = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", config.BaseConf.DbUser, config.BaseConf.DbPwd, config.BaseConf.DbAddr, config.BaseConf.DbDatabase))
	err := SqlDb.Ping()
	if err != nil {
		log.Println("sql-db ping error, err = ", err)
	}
}

func checkSqlDB(Db *sql.DB) error {
	err := Db.Ping()
	if err != nil {
		log.Println("Error!!! Sql DB Ping Error! info=")
		log.Println(fmt.Sprintf("%s:%s@tcp(%s)/%s", config.BaseConf.DbUser, config.BaseConf.DbPwd, config.BaseConf.DbAddr, config.BaseConf.DbDatabase))
		return err
	}
	return nil
}

func DbGetUserInfoByName(userName string, Db *sql.DB) (UserInfo, error) {
	if err := checkSqlDB(Db); err != nil {
		return UserInfo{}, err
	}
	var rs UserInfo
	row := Db.QueryRow("select user_name, pwd, nick_name, picture from user_info where user_name = ? and is_enabled = 1", userName)
	//defer row.Close() todo: 调试的时候注释掉了
	//if  row. !=nil{
	//	log.Println("DbGetUserInfoByName-Query-failed, err=", err)
	//	return UserInfo{}, err
	//}
	err := row.Scan(&rs.Name, &rs.Pwd, &rs.NickName, &rs.Picture)
	if err != nil && err != sql.ErrNoRows {
		log.Println("DbGetUserInfoByName, row-scan-error, err = ", err)
		return rs, err
	}
	return rs, nil
}

func DbInsertUserInfo(userName, pwd string, Db *sql.DB) error {
	if err := checkSqlDB(Db); err != nil {
		return err
	}
	inserSql := fmt.Sprintf("insert into user_info set user_name='%s', pwd='%s', nick_name=''", userName, pwd)
	_, err := Db.Exec(inserSql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DbModifyUserInfo(userName, pwd, nickName, picture string, Db *sql.DB) error {
	if err := checkSqlDB(Db); err != nil {
		return err
	}
	modifySql := fmt.Sprintf("update user_info set nick_name = '%s', picture = '%s' where user_name = '%s'",nickName, picture,userName)
	_, err := Db.Exec(modifySql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}