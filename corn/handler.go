package corn

import (
	"errors"
	"github.com/yuhaoyuan/RPC_server/dal"
	"log"
)

// get userinfo by name
func findUserByUserName(userName string) (dal.UserInfo, error) {
	// 先检查redis里面是否存在此用户
	userInfo, err := dal.CacherGetUserInfo(userName, dal.RedisDb)
	if err != nil {
		return userInfo, err
	}
	// 从sql中搜索
	if userInfo.Name == "" {
		userInfo, err = dal.DbGetUserInfoByName(userName, dal.SqlDb)
		if err != nil {
			return userInfo, err
		}
		if userInfo.Name != ""{
			_ = dal.CacherSetUserInfo(userInfo, dal.RedisDb)
		}
	}
	return userInfo, nil
}

func UserLogin(userName string, pwd string) (dal.UserInfo, error) {
	log.Printf("user-login!")

	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}
	// 校验身份
	if userInfo.Pwd == "" {
		return dal.UserInfo{}, errors.New("the username does not exist")
	} else {
		if pwd == userInfo.Pwd {
			return userInfo, nil
		} else {
			return dal.UserInfo{}, errors.New("wrong_password")
		}
	}
}

func UserRegister(userName string, pwd string) (dal.UserInfo, error) {
	log.Printf("user-register!")
	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}
	if userInfo.Name != ""{
		return userInfo, errors.New("the username already exist")
	}

	//
	err = dal.DbInsertUserInfo(userName, pwd, dal.SqlDb)
	if err != nil {
		return userInfo, err
	}
	rs := dal.UserInfo{
		Name: userName,
		Pwd:  pwd,
	}
	// 成功了的话，更新缓存
	_ = dal.CacherSetUserInfo(rs, dal.RedisDb)
	return rs, nil
}

func UserModifyInfo(userName, pwd, nickName, picture string) (dal.UserInfo, error) {
	log.Printf("user-modify!")
	userInfo, err := findUserByUserName(userName)
	if err != nil {
		return userInfo, err
	}
	// todo: 鉴权
	err = dal.DbModifyUserInfo(userName, pwd, nickName, picture, dal.SqlDb)
	if err != nil{
		return userInfo, err
	}
	rs := dal.UserInfo{
		Name:     userName,
		Pwd:      nickName,
		NickName: nickName,
		Picture: picture,

	}
	// 更新缓存
	_ = dal.CacherSetUserInfo(rs, dal.RedisDb)
	return rs, nil
}
