package sms

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/75912001/goz/zredis"
	"github.com/75912001/goz/zutility"
	"github.com/garyburd/redigo/redis"
)

//初始化
func (p *smsRegister) Init(bench *zutility.Ini) (err error) {
	const benchFileSection string = "sms"
	p.Pattern = bench.GetString(benchFileSection, "register_http_handler_pattern", "")
	//redis
	p.redisKeyPerfix = bench.GetString(benchFileSection, "redis_key_perfix_register", "")
	return err
}

func (p *smsRegister) IsExist(recNum string) (value bool) {
	//检查是否有记录
	const commandName string = "get"
	key := p.genRedisKey(recNum)
	reply, err := redis.server.Conn.Do(commandName, key)
	if nil != err {
		gLog.Error("redis get err:", err)
		return false
	}
	if nil == reply {
		return false
	}

	return true
}

/*
	http://gw.api.taobao.com/router/rest
	?sign=5A3BF0982B182890A900852CA6076CA9
	&app_key=23273583
	&method=alibaba.aliqin.fc.sms.num.send
	&rec_num=17721027200
	&sign_method=md5
	&sms_free_sign_name=%E6%B3%A8%E5%86%8C%E9%AA%8C%E8%AF%81
	&sms_param=%7B%27code%27%3A%27123%27%2C%27product%27%3A%27%E7%88%B1%E8%BD%A6%E6%97%85%27%7D
	&sms_template_code=SMS_2515091
	&sms_type=normal
	&timestamp=2015-11-26+19%3A29%3A56
	&v=2.0
*/
//用户请求获取sms验证码
func SmsRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	const paraNumber string = "number"

	var recNum string
	{ //解析手机号码
		err := req.ParseForm()
		if nil != err {
			gLog.Error("SmsRegisterHttpHandler err:", err)
			w.Write([]byte(strconv.Itoa(zutility.ECParam)))
			return
		}
		if len(req.Form[paraNumber]) > 0 {
			recNum = req.Form[paraNumber][0]
		}

		if phoneNumberLen != len(recNum) {
			gLog.Error("SmsRegisterHttpHandler number:", recNum)
			w.Write([]byte(strconv.Itoa(zutility.ECParam)))
			return
		}
	}

	{ //检查是否有记录 来自redis
		isExist := GSmsRegister.IsExist(recNum)
		if isExist {
			//有记录就返回，短信已发出，请收到后重试
			w.Write([]byte(strconv.Itoa(zutility.ECSMSSending)))
			return
		}
	}

	{ //检查手机号是否绑定
		bind, err := GphoneRegister.IsSmsBind(recNum)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zutility.ECSYS)))
			return
		}
		if bind {
			w.Write([]byte(strconv.Itoa(zutility.ECSMSBind)))
			return
		}
	}
	//生成短信内容参数
	smsCode := genSmsCode()

	{ //设置到redis中
		err := GSmsRegister.InsertSmsCode(recNum, smsCode)
		if nil != err {
			w.Write([]byte(strconv.Itoa(zutility.ECRedisSYS)))
			return
		}
	}
	/*
		参数名称 			参数类型 		必填与否 	样例取值 							 参数说明
		PhoneNumbers 		String 		必须 		 15000000000 						短信接收号码,支持以逗号分隔的形式进行批量调用，批量上限为1000个手机号码,批量调用相对于单条调用及时性稍有延迟,验证码类型的短信推荐使用单条调用的方式；发送国际/港澳台消息时，接收号码格式为：国际区号+号码，如“85200000000”
		SignName 			String 		必须 		 云通信 							 短信签名
		TemplateCode 		String 		必须 		 SMS_0000 							短信模板ID，发送国际/港澳台消息时，请使用国际/港澳台短信模版
		TemplateParam 		String 		可选 		 {“code”:”1234”,”product”:”ytx”} 	短信模板变量替换JSON串,友情提示:如果JSON中需要带换行符,请参照标准的JSON协议。
		SmsUpExtendCode 	String 		可选 		 90999 								上行短信扩展码,无特殊需要此字段的用户请忽略此字段
		OutId 				String 		可选 		 abcdefgh 							外部流水扩展字段
	*/

	var smsParam = "{'code':'" + smsCode + "','product':'" + GphoneSms.SmsParamProduct + "'}"
	reqUrl, err := GphoneSms.GenReqUrl(recNum, smsParam, GphoneSms.SmsFreeSignName, GphoneSms.SmsTemplateCode)
	if nil != err {
		gLog.Error("GenReqUrl err:", err)
		w.Write([]byte(strconv.Itoa(zutility.ECSYS)))
		return
	}
	//fmt.Println(reqUrl)

	{ //发送消息到短信服务器
		resp, err := http.Get(reqUrl)
		if nil != err {
			gLog.Error("err:", err, reqUrl)
			w.Write([]byte(strconv.Itoa(zutility.ECSYS)))
			return
		}
		defer resp.Body.Close()
		fmt.Println(resp)
		//fmt.Println(resp.Body)
	}
	w.Write([]byte(strconv.Itoa(zutility.ECSucc)))
}

type smsRegister struct {
	Pattern string
	//redis
	redisKeyPerfix string
	redis          *zredis.Server
}

//生成redis的键值
func (p *smsRegister) genRedisKey(key string) (value string) {
	return p.redisKeyPerfix + key
}

func (p *smsRegister) InsertSmsCode(recNum string, smsCode string) (err error) {
	//设置到redis中
	commandName := "setex"
	key := p.genRedisKey(recNum)
	timeout := "300" //5分钟
	_, err = p.redis.Conn.Do(commandName, key, timeout, smsCode)
	if nil != err {
		gLog.Error("redis setex err:", err)
	}

	return err
}

func (p *smsRegister) IsExistSmsCode(recNum string, smsCode string) (exist bool, err error) {
	//检查是否有短信验证码记录
	commandName := "get"
	key := p.genRedisKey(recNum)
	reply, err := redis.server.Conn.Do(commandName, key)
	if nil != err {
		return false, err
	}
	if nil == reply {
		return false, err
	}
	getRecNum, _ := redis.String(reply, err)
	if smsCode != getRecNum {
		fmt.Println("IsExistSmsCode,", recNum, smsCode, getRecNum)
		return false, err
	}

	return true, err
}

func (p *smsRegister) Del(recNum string) {
	//删除有短信验证码记录
	commandName := "del"
	key := p.genRedisKey(recNum)
	redis.server.Conn.Do(commandName, key)
}
