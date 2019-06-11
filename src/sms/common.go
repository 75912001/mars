package sms

import (
	"math/rand"
	"strconv"

	"github.com/75912001/goz/zutility"
)

var GSmsRegister smsRegister

//SetLog 设置log
func SetLog(v *zutility.Log) {
	gLog = v
}

////////////////////////////////////////////////////////////////////////////////
//手机短信注册(发送手机号,接收验证码)

//手机验证码个数 5位,[10000-100000)
//手机上5位数字 会有下划线，可以长按复制，方便用户使用
const smsCodeBegin = 10000
const smsCodeEnd = 99999 + 1

//手机号码长度
const phoneNumberLen int = 11

var gLog *zutility.Log

func genSmsCode() (value string) {
	{ //生成短信内容参数
		index := rand.Intn(smsCodeEnd)
		if index < smsCodeBegin {
			index += smsCodeBegin
		}
		value = strconv.Itoa(index)
	}
	return value
}
