package flag_user

import (
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/models/enum"
	pwdUtils "TPC-EDM-Server/utils/pwd"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type FlagUser struct{}

func (FlagUser) Create() {
	var role enum.RoleType
	fmt.Println("选择角色 1管理员 2普通用户")
	_, err := fmt.Scan(&role)
	if err != nil {
		logrus.Errorf("输入错误: %s", err)
		return
	}
	if !(role == 1 || role == 2) {
		logrus.Errorf("输入角色错误,请输入1~2内的值")
		return
	}
	var username string
	fmt.Println("输入用户名")
	_, err = fmt.Scan(&username)
	if err != nil {
		logrus.Errorf("输入错误: %s", err)
		return
	}
	var model models.UserModel
	err = global.DB.Take(&model, "username = ?", username).Error
	if err == nil {
		logrus.Errorf("用户名已经存在")
		return
	}
	fmt.Println("输入密码")
	pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		logrus.Errorf("读取密码出错:%s", err)
		return
	}
	fmt.Println("重复密码")
	rpwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		logrus.Errorf("读取密码出错:%s", err)
		return
	}
	if string(pwd) != string(rpwd) {
		logrus.Errorf("密码输入不一致")
		return
	}
	hashPwd, _ := pwdUtils.GenerateFromPassword(string(pwd))
	err = global.DB.Create(&models.UserModel{
		Username: username,
		Password: hashPwd,
		Role:     role,
	}).Error
	if err != nil {
		logrus.Errorf("用户创建失败")
		return
	}
	logrus.Infoln("用户创建成功")
}
