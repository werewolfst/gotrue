package sms_provider

import (
	"encoding/json"
	"fmt"
	"strings"

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
	var signName string
	var product string
	var templateCode string
	var params map[string]string
	var apiKey string
	var apiSecret string

	switch sender {
	case "tit":
		apiKey = t.Config.TitApiKey
		apiSecret = t.Config.TitApiSecret
		signName = t.Config.TitSignName
		product = t.Config.TitProduct
		templateCode = t.Config.TitTemplateCode
		params = map[string]string{
			"code": message,
		}
	case "arduino":
		apiKey = t.Config.ApiKey
		apiSecret = t.Config.ApiSecret
		signName = t.Config.ArduinoSignName
		product = t.Config.ArduinoProduct
		templateCode = t.Config.ArduinoTemplateCode
		params = map[string]string{
			"code":    message,
			"product": product,
		}
	case "aily":
		apiKey = t.Config.ApiKey
		apiSecret = t.Config.ApiSecret
		signName = t.Config.AilySignName
		product = t.Config.AilyProduct
		templateCode = t.Config.AilyTemplateCode
		params = map[string]string{
			"code":    message,
			"product": product,
		}
	case "xy":
		apiKey = t.Config.ApiKey
		apiSecret = t.Config.ApiSecret
		signName = t.Config.XiyuanSignName
		product = t.Config.XiyuanProduct
		templateCode = t.Config.XiyuanTemplateCode
		params = map[string]string{
			"code":    message,
			"product": product,
		}
	default:
		apiKey = t.Config.ApiKey
		apiSecret = t.Config.ApiSecret
		signName = "blinker"
		product = "三塔"
		params = map[string]string{
			"code":    message,
			"product": product,
		}
	}

	fmt.Println("ConfigSignName: ", signName)
	fmt.Println("ConfigProduct: ", product)
	fmt.Println("TemplateCode: ", templateCode)
	fmt.Println("ApiKey: ", apiKey)
	fmt.Println("ApiSecret: ", apiSecret)

	config := &openapi.Config{}
	config.SetAccessKeyId(apiKey)
	config.SetAccessKeySecret(apiSecret)
	config.Endpoint = tea.String(t.APIPath)

	client := &dysmsapi20170525.Client{}
	var _err error
	client, _err = dysmsapi20170525.NewClient(config)
	if _err != nil {
		fmt.Println("SMS Client init error")
	}

	sendMsg := &dysmsapi20170525.SendSmsRequest{}
	sendMsg.SetPhoneNumbers(phone)
	sendMsg.SetSignName(signName)
	sendMsg.SetTemplateCode(templateCode)

	jsonStr, _ := json.Marshal(params)
	sendMsg.SetTemplateParam(string(jsonStr))

	res, sErr := client.SendSms(sendMsg)
	if sErr != nil {
		fmt.Println(sErr)
	}

	fmt.Println(res)

	return nil
}

func (t *AliyunProvider) GetSender(request_refer string) (string, error) {
	if strings.Contains(request_refer, "titlab") {
		return "tit", nil
	} else if strings.Contains(request_refer, "arduino") {
		return "arduino", nil
	} else if strings.Contains(request_refer, "aily") {
		return "aily", nil
	} else if strings.Contains(request_refer, "xy") {
		return "xy", nil
	} else {
		return "", fmt.Errorf("not support refer: ", request_refer)
	}
}
