package sms_provider

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/netlify/gotrue/conf"
)

const (
	defaultApi = "dysmsapi.aliyuncs.com"
)

type AliyunProvider struct {
	Config  *conf.AliyunProviderConfiguration
	APIPath string
}

// Creates a SmsProvider with the Aliyun Config
func NewAliyunProvider(config conf.AliyunProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &AliyunProvider{
		Config:  &config,
		APIPath: defaultApi,
	}, nil
}

// Send an sms containing the OTP with Aliyun API
func (t *AliyunProvider) SendSms(phone string, message string) error {
	config := &openapi.Config{}
	config.SetAccessKeyId(t.Config.ApiKey)
	config.SetAccessKeySecret(t.Config.ApiSecret)
	config.Endpoint = tea.String(t.APIPath)

	client := &dysmsapi20170525.Client{}
	var _err error
	client, _err = dysmsapi20170525.NewClient(config)
	if _err != nil {
		fmt.Println("SMS Client init error")
	}

	sendMsg := &dysmsapi20170525.SendSmsRequest{}
	sendMsg.SetPhoneNumbers(phone)
	sendMsg.SetSignName(t.Config.SignName)
	sendMsg.SetTemplateCode(t.Config.TemplateCode)

	params := map[string]string{
		"code":    message,
		"product": t.Config.Product,
	}
	jsonStr, _ := json.Marshal(params)
	sendMsg.SetTemplateParam(string(jsonStr))

	res, sErr := client.SendSms(sendMsg)
	if sErr != nil {
		fmt.Println(sErr)
	}

	fmt.Println(res)

	return nil
}
