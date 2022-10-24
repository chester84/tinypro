package tc

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

var (
	credential *common.Credential

	ocrCpf    *profile.ClientProfile
	ocrClient *ocr.Client

	smsCpf    *profile.ClientProfile
	smsClient *sms.Client
)

func init() {
	var err error

	credential = common.NewCredential(
		secretId,
		secretKey,
	)
	if credential == nil {
		panic("tencentcloud credential init fail")
	}

	ocrCpf = profile.NewClientProfile()
	if ocrCpf == nil {
		panic("tencentcloud ocr cpf init fail")
	}

	ocrCpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	ocrClient, err = ocr.NewClient(credential, "ap-guangzhou", ocrCpf)
	if err != nil {
		panic(fmt.Sprintf(`ocr client init fail, err: %v`, err))
	} else {
		fmt.Println("tencentcloud ocr client init success")
	}

	smsCpf = profile.NewClientProfile()
	if smsCpf == nil {
		panic("tencentcloud sms cpf init fail")
	}

	smsCpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	smsClient, err = sms.NewClient(credential, "ap-guangzhou", smsCpf)
	if err != nil {
		panic(fmt.Sprintf(`sms client init fail, err: %v`, err))
	} else {
		fmt.Println("tencentcloud sms client init success")
	}

}

func OcrClient() *ocr.Client {
	return ocrClient
}

func SmsClient() *sms.Client {
	return smsClient
}
