package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

type userInfo struct {
	Name     string `json:"name"`
	Pwd      string `json:"pwd"`
	NickName string `json:"nick_name"`
	Picture  string `json:"picture"`
}

func TestDb(t *testing.T) {
	db, err := sql.Open("mysql", "root:12345678@(127.0.0.1:3306)/yhy")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	}

	// insert
	inserSql := "insert into user_info set name='test001', nick_name='哈苏滴孩奴'"
	_, err = db.Exec(inserSql) // OK
	if err != nil {
		fmt.Println(err)
	}

	inserSql2 := "insert into user_info set name='test002', nick_name='🌶️🔟🤨🐂🍺'"
	_, err = db.Exec(inserSql2) // OK

	if err != nil {
		fmt.Println(err)
	}

	// get
	rows, err := db.Query("select user_name, pwd, nick_name, picture from user_info")
	defer rows.Close()
	for rows.Next() {
		info := userInfo{}
		err := rows.Scan(&info.Name, &info.Pwd, &info.NickName, &info.Picture)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(info)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// update

}
