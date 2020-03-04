package corn

import (
	"errors"
	"github.com/yuhaoyuan/RPC_server/auth"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/dal"
	"log"
)

// findUserByUserName get userinfo by name
func findUserByUserName(userName string) (dal.UserInfo, error) {
	// 先检查redis里面是否存在此用户
	log.Println("findUserByUserName---in")
	userInfo, err := dal.CacherGetUserInfo(userName, dal.RedisDb)
	if err != nil {
		log.Println("CacherGetUserInfo - error , err = ", err)
		return userInfo, err
	}
	log.Println("findUserByUserName---cacher-get-user=", userInfo)
	// 从sql中搜索
	if userInfo.Name == "" {
		userInfo, err = dal.DbGetUserInfoByName(userName, dal.SQLDB)
		if err != nil {
			log.Println("findUserByUserName-DbGetUserInfoByName error = ", err)
			return userInfo, err
		}
		if userInfo.Name != "" {
			_ = dal.CacherSetUserInfo(userInfo, dal.RedisDb)
		}
		log.Println("findUserByUserName--sql-get-user=", userInfo)
	}
	return userInfo, nil
}

// UserLogin login
func UserLogin(userName string, pwd string) (dal.UserInfo, error) {
	log.Println("user-login-----")

	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}
	// 校验身份
	if userInfo.Pwd == "" {
		return dal.UserInfo{}, errors.New("the username does not exist")
	}
	if pwd == userInfo.Pwd {
		orginData := auth.CacherGetToken(userName, dal.RedisDb)
		cryptedToken, _ := auth.AesEncrypt([]byte(orginData), []byte(config.BaseConf.AesTokenKey))
		userInfo.Token = string(cryptedToken)
		log.Println("UserLogin----return-data=", userInfo.Name)
		return userInfo, nil
	}
	return dal.UserInfo{}, errors.New("wrong_password")
}

// UserRegister register
func UserRegister(userName string, pwd string) (dal.UserInfo, error) {
	log.Printf("user-register!")
	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}
	if userInfo.Name != "" {
		return userInfo, errors.New("the username already exist")
	}

	//
	err = dal.DbInsertUserInfo(userName, pwd, dal.SQLDB)
	if err != nil {
		return userInfo, err
	}
	rs := dal.UserInfo{
		Name: userName,
		Pwd:  pwd,
	}
	// 成功了的话，更新缓存
	_ = dal.CacherSetUserInfo(rs, dal.RedisDb)

	orginData := auth.CacherGetToken(userName, dal.RedisDb)
	cryptedToken, _ := auth.AesEncrypt([]byte(orginData), []byte(config.BaseConf.AesTokenKey))
	rs.Token = string(cryptedToken)
	return rs, nil
}

// UserModifyInfo  修改资料
func UserModifyInfo(userName, pwd, nickName, picture string) (dal.UserInfo, error) {
	log.Printf("user-modify!")
	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}

	err = dal.DbModifyUserInfo(userName, pwd, nickName, picture, dal.SQLDB)
	if err != nil {
		return userInfo, err
	}
	rs := dal.UserInfo{
		Name:     userName,
		Pwd:      nickName,
		NickName: nickName,
		Picture:  picture,
	}
	// 更新缓存
	_ = dal.CacherDelUserInfo(userName, dal.RedisDb)
	return rs, nil
}

// GetUserInfoByToken 校验token
func GetUserInfoByToken(userName, token string) (dal.UserInfo, error) {
	log.Printf("GetUserInfoByToken in, user_name=%s", userName)
	// 插入一个管理员token逻辑
	tokenUserName := ""
	if token == "ckQSpDXWcJVTWfFidRkh" {
		tokenUserName = userName
	} else {
		encodeToken, _ := auth.AesDecrypt([]byte(token), []byte(config.BaseConf.AesTokenKey))
		tokenUserName = auth.CacherGetUserNameByToken(string(encodeToken), dal.RedisDb)
	}
	if tokenUserName != "" && tokenUserName == userName {
		userInfo, _ := findUserByUserName(userName)
		userInfo.Token = token
		log.Printf("GetUserInfoByToken out, return=%s", userInfo.Name)
		return userInfo, nil
	}
	log.Println("can not find token or wrong user_name")
	return dal.UserInfo{}, errors.New("can not find token or wrong user_name")
}
