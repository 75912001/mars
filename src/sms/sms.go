package sms

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
//手机短信
var GphoneSms phoneSms_t

type phoneSms_t struct {
	UrlPattern      string
	AppKey          string
	AppSecret       string
	Method          string
	SignMethod      string
	SmsFreeSignName string
	SmsTemplateCode string
	SmsType         string
	Versions        string
	SmsParamProduct string

	SmsFreeSignNameChangePwd string
	SmsTemplateCodeChangePwd string
}

//初始化
func (p *phoneSms_t) Init() (err error) {
	const benchFileSection string = "ict_phone_sms"

	p.UrlPattern = ict_cfg.Gbench.FileIni.Get(benchFileSection, "UrlPattern", " ")
	p.AppKey = ict_cfg.Gbench.FileIni.Get(benchFileSection, "AppKey", " ")
	p.AppSecret = ict_cfg.Gbench.FileIni.Get(benchFileSection, "AppSecret", " ")
	p.Method = ict_cfg.Gbench.FileIni.Get(benchFileSection, "Method", " ")
	p.SignMethod = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SignMethod", " ")
	p.SmsFreeSignName = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsFreeSignName", " ")
	p.SmsTemplateCode = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsTemplateCode", " ")
	p.SmsType = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsType", " ")
	p.Versions = ict_cfg.Gbench.FileIni.Get(benchFileSection, "Versions", " ")
	p.SmsParamProduct = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsParamProduct", " ")

	p.SmsFreeSignNameChangePwd = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsFreeSignNameChangePwd", " ")
	p.SmsTemplateCodeChangePwd = ict_cfg.Gbench.FileIni.Get(benchFileSection, "SmsTemplateCodeChangePwd", " ")
	return err
}

//生成短信请求url
func (p *phoneSms_t) GenReqUrl(recNum string, smsParam string, SmsFreeSignName string, SmsTemplateCode string) (value string, err error) {
	//时间戳格式"2015-11-26 20:32:42"
	var timeStamp = zutility.StringSubstr(time.Now().String(), 19)

	var strMd5 string = p.genSign(recNum, smsParam, timeStamp, SmsFreeSignName, SmsTemplateCode)

	var reqUrl = p.UrlPattern +
		"?sign=" + strMd5 +
		"&app_key=" + p.AppKey +
		"&method=" + p.Method +
		"&rec_num=" + recNum +
		"&sign_method=" + p.SignMethod +
		"&sms_free_sign_name=" + SmsFreeSignName +
		"&sms_param=" + smsParam +
		"&sms_template_code=" + SmsTemplateCode +
		"&sms_type=" + p.SmsType +
		"&timestamp=" + timeStamp +
		"&v=" + p.Versions

	url, err := url.Parse(reqUrl)
	if nil != err {
		fmt.Println("######PhoneRegister.genReqUrl err:", reqUrl, err)
		return reqUrl, err
	}
	reqUrl = p.UrlPattern + "?" + url.Query().Encode()
	return reqUrl, err
}

//生成sign(MD5)
func (p *phoneSms_t) genSign(recNum string, smsParam string, timeStamp string, SmsFreeSignName string, SmsTemplateCode string) (value string) {
	var signSource = p.AppSecret +
		"app_key" + p.AppKey +
		"method" + p.Method +
		"rec_num" + recNum +
		"sign_method" + p.SignMethod +
		"sms_free_sign_name" + SmsFreeSignName +
		"sms_param" + smsParam +
		"sms_template_code" + SmsTemplateCode +
		"sms_type" + p.SmsType +
		"timestamp" + timeStamp +
		"v" + p.Versions +
		p.AppSecret
	strMd5 := zutility.GenMd5(signSource)
	strMd5 = strings.ToUpper(strMd5)
	fmt.Println(signSource)
	fmt.Println(strMd5)
	return strMd5
}



    /**
     * 发送短信
     * @return stdClass
     */
	  func sendSms() {

        // 初始化SendSmsRequest实例用于设置发送短信的参数
        $request = new SendSmsRequest();

        // 必填，设置短信接收号码
        $request->setPhoneNumbers("12345678901");

        // 必填，设置签名名称，应严格按"签名名称"填写，请参考: https://dysms.console.aliyun.com/dysms.htm#/develop/sign
        $request->setSignName("短信签名");

        // 必填，设置模板CODE，应严格按"模板CODE"填写, 请参考: https://dysms.console.aliyun.com/dysms.htm#/develop/template
        $request->setTemplateCode("SMS_0000001");

        // 可选，设置模板参数, 假如模板中存在变量需要替换则为必填项
        $request->setTemplateParam(json_encode(array(  // 短信模板中字段的值
            "code"=>"12345",
            "product"=>"dsd"
        ), JSON_UNESCAPED_UNICODE));

        // 可选，设置流水号
        $request->setOutId("yourOutId");

        // 选填，上行短信扩展码（扩展码字段控制在7位或以下，无特殊需求用户请忽略此字段）
        $request->setSmsUpExtendCode("1234567");

        // 发起访问请求
        $acsResponse = static::getAcsClient()->getAcsResponse($request);

        return $acsResponse;
	}
	
type  SendSmsRequest struct
{
	public function  __construct()
	{
		parent::__construct("Dysmsapi", "2017-05-25", "SendSms");
		$this->setMethod("POST");
	}

	private  $templateCode;

	private  $phoneNumbers;

	private  $signName;

	private  $resourceOwnerAccount;

	private  $templateParam;

	private  $resourceOwnerId;

	private  $ownerId;

	private  $outId;

    private  $smsUpExtendCode;

	public function getTemplateCode() {
		return $this->templateCode;
	}

	public function setTemplateCode($templateCode) {
		$this->templateCode = $templateCode;
		$this->queryParameters["TemplateCode"]=$templateCode;
	}

	public function getPhoneNumbers() {
		return $this->phoneNumbers;
	}

	public function setPhoneNumbers($phoneNumbers) {
		$this->phoneNumbers = $phoneNumbers;
		$this->queryParameters["PhoneNumbers"]=$phoneNumbers;
	}

	public function getSignName() {
		return $this->signName;
	}

	public function setSignName($signName) {
		$this->signName = $signName;
		$this->queryParameters["SignName"]=$signName;
	}

	public function getResourceOwnerAccount() {
		return $this->resourceOwnerAccount;
	}

	public function setResourceOwnerAccount($resourceOwnerAccount) {
		$this->resourceOwnerAccount = $resourceOwnerAccount;
		$this->queryParameters["ResourceOwnerAccount"]=$resourceOwnerAccount;
	}

	public function getTemplateParam() {
		return $this->templateParam;
	}

	public function setTemplateParam($templateParam) {
		$this->templateParam = $templateParam;
		$this->queryParameters["TemplateParam"]=$templateParam;
	}

	public function getResourceOwnerId() {
		return $this->resourceOwnerId;
	}

	public function setResourceOwnerId($resourceOwnerId) {
		$this->resourceOwnerId = $resourceOwnerId;
		$this->queryParameters["ResourceOwnerId"]=$resourceOwnerId;
	}

	public function getOwnerId() {
		return $this->ownerId;
	}

	public function setOwnerId($ownerId) {
		$this->ownerId = $ownerId;
		$this->queryParameters["OwnerId"]=$ownerId;
	}

	public function getOutId() {
		return $this->outId;
	}

	public function setOutId($outId) {
		$this->outId = $outId;
		$this->queryParameters["OutId"]=$outId;
	}

    public function getSmsUpExtendCode() {
        return $this->smsUpExtendCode;
    }

    public function setSmsUpExtendCode($smsUpExtendCode) {
        $this->smsUpExtendCode = $smsUpExtendCode;
        $this->queryParameters["SmsUpExtendCode"]=$smsUpExtendCode;
    }
}