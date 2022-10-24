package tc

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

func generateHmacSHA1(secretToken, payloadBody string) []byte {
	mac := hmac.New(sha1.New, []byte(secretToken))
	sha1.New()
	mac.Write([]byte(payloadBody))
	return mac.Sum(nil)
}

func GetSignV1() string {
	rand.Seed(time.Now().Unix())
	// timestamp := time.Now().Unix()
	timestamp := time.Now().Unix()
	expireTime := timestamp + 86400*365*10
	timestampStr := strconv.FormatInt(timestamp, 10)
	expireTimeStr := strconv.FormatInt(expireTime, 10)
	random := 220625
	randomStr := strconv.Itoa(random)
	original := "secretId=" + secretId + "&current1TimeStamp=" + timestampStr + "&expireTime=" + expireTimeStr + "&random=" + randomStr
	signature := generateHmacSHA1(secretKey, original)
	signature = append(signature, []byte(original)...)
	signatureB64 := base64.StdEncoding.EncodeToString(signature)
	return signatureB64
}
