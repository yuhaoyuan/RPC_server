package dal

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"testing"
)

func TestDb(t *testing.T) {
	db, err := sql.Open("mysql", "root:12345678@(127.0.0.1:3306)/yhy")

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer db.Close()
	pingErr := db.Ping()
	if pingErr != nil {
		t.Errorf("err = %v want nil", err)
	}

	// insert

	inserSQL := "insert into user_info set user_name='test001', nick_name='哈苏滴孩奴', picture=''"
	_, err = db.Exec(inserSQL) // OK
	if err != nil && !strings.Contains(err.Error(), "Duplicate") {
		t.Errorf("inser-sql = %v want nil", err)
	}

	inserSQL2 := "insert into user_info set user_name='test002', nick_name='🌶️🔟🤨🐂🍺'"
	_, err = db.Exec(inserSQL2) // OK

	if err != nil && !strings.Contains(err.Error(), "Duplicate") {
		t.Errorf("inser-sql = %v want nil", err)
	}

	// get
	rows, err := db.Query("select user_name, pwd, nick_name, picture from user_info where user_name = 'test001'")
	defer rows.Close()
	for rows.Next() {
		info := UserInfo{}
		err := rows.Scan(&info.Name, &info.Pwd, &info.NickName, &info.Picture)
		if err != nil {
			t.Errorf("rows Scan error = %v want nil", err)
		}
		if info.Name != "test001" {
			t.Errorf("rows-get info.Name != 'test001'")
		}
	}
	err = rows.Err()
	if err != nil {
		t.Errorf("rows.ERR err = %v want nil", err)
	}
}
