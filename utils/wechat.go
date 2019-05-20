package utils

import (
	"fmt"
	"github.com/yigger/JZ-back/conf"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

const (
	code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

// ResCode2Session 登录凭证校验的返回结果
type ResCode2Session struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func Code2Session(jsCode string) (result ResCode2Session, err error) {
	urlStr := fmt.Sprintf(code2SessionURL, conf.Conf.AppId, conf.Conf.AppSecret, jsCode)
	var response []byte
	response, err = HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("Code2Session error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
