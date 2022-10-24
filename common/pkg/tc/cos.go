package tc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/teris-io/shortid"

	"tinypro/common/lib/redis/cache"
	"github.com/chester84/libtools"
)

type wCos struct {
	pubClient *cos.Client
	priClient *cos.Client
}

var cosObj *wCos

func init() {
	uPub, _ := url.Parse(bucketPublicUrl)
	uPri, _ := url.Parse(bucketPrivateUrl)

	bPub := &cos.BaseURL{BucketURL: uPub}
	bPri := &cos.BaseURL{BucketURL: uPri}

	// 1.永久密钥
	cosObj = &wCos{}

	cosObj.pubClient = cos.NewClient(bPub, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	cosObj.priClient = cos.NewClient(bPri, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
}

func UploadFromFile2Public(s3key, filePath string) error {
	_, err := cosObj.pubClient.Object.PutFromFile(context.Background(), s3key, filePath, nil)
	if err != nil {
		logs.Error("[cos] UploadFromFile2Public err: %#v", err)
	}
	return err
}

func UploadFromStream2Public(s3key string, f io.Reader) error {
	_, err := cosObj.pubClient.Object.Put(context.Background(), s3key, f, nil)
	if err != nil {
		logs.Error("[cos] UploadFromStream2Public err: %#v", err)
	}
	return err
}

func UploadFromFile2Private(s3key, filePath string) error {
	_, err := cosObj.priClient.Object.PutFromFile(context.Background(), s3key, filePath, nil)
	if err != nil {
		logs.Error("[cos] UploadFromFile2Private err: %#v", err)
	}
	return err
}

func UploadFromStream2Private(s3key string, f io.Reader) error {
	_, err := cosObj.priClient.Object.Put(context.Background(), s3key, f, nil)
	if err != nil {
		logs.Error("[cos] UploadFromStream2Private err: %#v", err)
	}
	return err
}

func DownloadFromPublic(s3key, saveFile string) error {
	_, err := cosObj.pubClient.Object.GetToFile(context.Background(), s3key, saveFile, nil)
	if err != nil {
		logs.Error("[cos] DownloadFromPublic err: %#v", err)
	}
	return err
}

func DownloadPublic2Stream(s3key string, buf io.ReadWriter) error {
	resp, err := cosObj.pubClient.Object.Get(context.Background(), s3key, nil)
	if err != nil {
		logs.Error("[cos] DownloadPublic2Stream get exception, s3key: %s, err: %#v", s3key, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	_, err = io.Copy(buf, resp.Body)

	return err
}

func DownloadFromPrivate(s3key, saveFile string) error {
	_, err := cosObj.priClient.Object.GetToFile(context.Background(), s3key, saveFile, nil)
	if err != nil {
		logs.Error("[cos] UploadFromStream2Private err: %#v", err)
	}
	return err
}

func DownloadPrivate2Stream(s3key string, buf io.ReadWriter) error {
	resp, err := cosObj.priClient.Object.Get(context.Background(), s3key, nil)
	if err != nil {
		logs.Error("[cos] DownloadPrivate2Stream get exception, s3key: %s, err: %#v", s3key, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	_, err = io.Copy(buf, resp.Body)

	return err
}

func HeadWithPublic(s3key string) (*cos.Response, error) {
	return cosObj.pubClient.Object.Head(context.Background(), s3key, nil)
}

func HeadWithPrivate(s3key string) (*cos.Response, error) {
	return cosObj.priClient.Object.Head(context.Background(), s3key, nil)
}

// 这个方法不是常态调用,所以不用缓存
func SignUrl(s3key string, expired time.Duration) string {
	if s3key == "" {
		return ""
	}

	urlObj, err := cosObj.priClient.Object.GetPresignedURL(context.Background(), http.MethodGet, s3key, secretId, secretKey, expired, nil)
	if err != nil {
		logs.Error("[SignUrl] build get exception, s3key: %s, err: %v", s3key, err)
		return ""
	}

	return urlObj.String()
}

func PublicUrl(s3key string) string {
	return fmt.Sprintf("%s/%s", bucketPublicUrl, s3key)
}

func CDNDomain() string {
	return bucketPrivateUrl
}

func TemporaryUrl(s3key string) (s string) {
	if s3key == "" {
		logs.Warning("[TemporaryUrl] input is empty")
		return
	}

	rid, err := shortid.Generate()
	if err != nil {
		logs.Error("[TemporaryUrl] create short id exception, s3key: %s, err: %v", s3key, err)
		return
	}

	s = fmt.Sprintf(`%s/open-api/resource/%s`, libtools.InternalApiDomain(), rid)

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	ex := 900
	cKey := GenTemporaryUrlRdsKey(rid)
	_, err = cacheClient.Do("SETEX", cKey, ex, s3key)
	if err != nil {
		logs.Error("[TemporaryUrl] redis> SETEX %s %d %s, err: %v", cKey, ex, s3key, err)
	}

	return
}
