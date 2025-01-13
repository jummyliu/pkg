package dingding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jummyliu/pkg/db/types"
	"github.com/jummyliu/pkg/request"
)

const BotURL = "https://oapi.dingtalk.com/robot/send?access_token=%s"

func PushToBot(apiKey string, body map[string]any, proxyURL string) (err error) {
	url := fmt.Sprintf(BotURL, apiKey)
	reqData, _ := json.Marshal(body)
	headers := map[string]string{
		// "Accept-Encoding": "gzip, deflate",
		"Accept":          "application/json",
		"Content-Type":    "application/json;charset=UTF-8",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Connection":      "close",
	}
	_, resp, _, err := request.DoRequest(
		url,
		request.WithMethod(http.MethodPost),
		request.WithHeader(headers),
		request.WithData(reqData),
		request.WithProxy(proxyURL),
	)
	if err != nil {
		return err
	}
	respData := map[string]any{}
	err = json.Unmarshal(resp, &respData)
	if err != nil {
		return err
	}
	errcode := types.MustGetValue[float64](respData, "errcode")
	if errcode == 0 {
		return nil
	}
	return errors.New(types.MustGetValue[string](respData, "errmsg"))
}
