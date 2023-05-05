package sms_provider

import (
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/supabase/gotrue/internal/conf"
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

func (t *AliyunProvider) SendMessage(phone string, message string, channel string, sender string) error {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message, sender)
	default:
		return fmt.Errorf("channel type %q is not supported for Aliyun", channel)
	}
}

// Send an sms containing the OTP with Aliyun API
func (t *AliyunProvider) SendSms(phone string, message string, sender string) error {
	fmt.Println("sender: ", sender)
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

	var signName string
	var product string

	switch sender {
	case "tit":
		signName = t.Config.TitSignName
		product = t.Config.TitProduct
	case "arduino":
		signName = t.Config.ArduinoSignName
		product = t.Config.ArduinoProduct
	case "aily":
		signName = t.Config.AilySignName
		product = t.Config.AilyProduct
	case "xy":
		signName = t.Config.XiyuanSignName
		product = t.Config.XiyuanProduct
	default:
		signName = "blinker"
		product = "三塔"
	}

	fmt.Println("ConfigSignName: ", signName)
	fmt.Println("ConfigProduct: ", product)

	sendMsg := &dysmsapi20170525.SendSmsRequest{}
	sendMsg.SetPhoneNumbers(phone)
	sendMsg.SetSignName(signName)
	sendMsg.SetTemplateCode(t.Config.TemplateCode)

	params := map[string]string{
		"code":    message,
		"product": product,
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
