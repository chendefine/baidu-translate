package baidu_translate

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	BaiduTranslateHostOfficial = "https://fanyi-api.baidu.com"

	BaiduTranslateEndpointCommon = "/api/trans/vip/translate"
)

var (
	support         = struct{}{}
	supportLanguges = map[string]struct{}{"zh": support, "en": support}

	errEmptyTranslateText = fmt.Errorf("empty translate text")
	errEmptyTranslateTo   = fmt.Errorf("unsupport translate to")
)

type BaiduTranslateConfig struct {
	Host      string `json:"host"`
	AppId     string `json:"app_id"`
	SecretKey string `json:"secret_key"`
}

type BaiduTranslateClient struct {
	config *BaiduTranslateConfig
	client *resty.Client
}

func NewBaiduTranslateClient(config *BaiduTranslateConfig) *BaiduTranslateClient {
	if config.Host == "" {
		config.Host = BaiduTranslateHostOfficial
	}
	client := resty.New().SetBaseURL(config.Host)
	return &BaiduTranslateClient{
		config: config,
		client: client,
	}
}

type ResultItem struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type TranslateRsp struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []*ResultItem `json:"trans_result"`
	ErrCode     string        `json:"error_code"`
	ErrMsg      string        `json:"error_msg"`
}

func (c *BaiduTranslateClient) Translate(text string, from string, to string) ([]*ResultItem, error) {
	if text == "" {
		return nil, errEmptyTranslateText
	} else if _, ok := supportLanguges[to]; !ok {
		return nil, errEmptyTranslateTo
	} else if _, ok := supportLanguges[from]; !ok {
		from = "auto"
	}
	salt := strconv.Itoa(int(rand.Int31()))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(c.config.AppId+text+salt+c.config.SecretKey)))
	req := map[string]string{"q": text, "from": from, "to": to, "appid": c.config.AppId, "salt": salt, "sign": sign}
	rsp := new(TranslateRsp)
	_, err := c.client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetFormData(req).SetResult(rsp).Post(BaiduTranslateEndpointCommon)
	if err != nil {
		return nil, err
	} else if rsp.ErrCode != "" {
		return nil, fmt.Errorf("err_code: %s, err_msg: %s", rsp.ErrCode, rsp.ErrMsg)
	}
	return rsp.TransResult, nil
}
