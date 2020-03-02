package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
		fmt.Println(pingErr)
	}

	// insert

	//file, err := os.Open("/Users/yuhaoyuan/Downloads/kelason.jpg")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer file.Close()
	//
	//stats, statsErr := file.Stat()
	//if statsErr != nil || stats ==nil{
	//	fmt.Println(statsErr)
	//}
	//
	//var size int64 = stats.Size()
	//imgBytes := make([]byte, size)
	//
	//bufr := bufio.NewReader(file)
	//_,err = bufr.Read(imgBytes)
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//}


	inserSql := "insert into user_info set user_name='test001', nick_name='哈苏滴孩奴', picture=''"
	_, err = db.Exec(inserSql) // OK
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	inserSql2 := "insert into user_info set user_name='test002', nick_name='🌶️🔟🤨🐂🍺'"
	_, err = db.Exec(inserSql2) // OK

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// get
	rows, err := db.Query("select user_name, pwd, nick_name, picture from user_info")
	defer rows.Close()
	for rows.Next() {
		info := UserInfo{}
		err := rows.Scan(&info.Name, &info.Pwd, &info.NickName, &info.Picture)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		log.Println(info)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
