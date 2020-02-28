package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
)


func DbInit() {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", config.BaseConf.DbUser, config.BaseConf.DbPwd, config.BaseConf.DbAddr, config.BaseConf.DbDatabase))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
