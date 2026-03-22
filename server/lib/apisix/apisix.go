package apisix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/erdong01/kit/httpClient"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

var client *httpClient.HttpClient

func Init() {
	cfg := global.GVA_CONFIG.Apisix
	client = httpClient.New(cfg.Url)
	h := make(http.Header)
	h.Set("X-API-KEY", cfg.ApiKey)
	client.Header = h
}

func DELETEConsumers(username string) error {
	cfg := global.GVA_CONFIG.Apisix
	deleteURL := strings.TrimRight(cfg.Url, "/") + "/" + url.PathEscape(username)

	deleteClient := httpClient.New(deleteURL)
	h := make(http.Header)
	h.Set("X-API-KEY", cfg.ApiKey)
	deleteClient.Header = h

	res := deleteClient.SetMethod(http.MethodDelete).Do()
	if res.Err != nil {
		return res.Err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("apisix delete consumer failed: %s - %s", res.Status, string(res.ResponseBody))
	}

	return nil
}

func POSTConsumers(username, apiKey, userKey string) error {
	reqData := map[string]interface{}{
		"username": username,
		"plugins": map[string]interface{}{
			"key-auth": map[string]interface{}{
				"key": userKey,
			},
			"serverless-pre-function": map[string]interface{}{
				"functions": []string{
					"return function(conf, ctx) ngx.req.set_header(\"Authorization\", \"Bearer 960de025-b69f-417e-a105-ae7955bd51b2\") end",
				},
			},
			"proxy-rewrite": map[string]interface{}{
				"headers": map[string]interface{}{
					"set": map[string]interface{}{
						"Authorization": "Bearer " + apiKey,
					},
				},
			},
			"limit-req": map[string]interface{}{
				"rate":          5,
				"burst":         2,
				"rejected_code": 429,
				"key_type":      "var",
				"key":           "consumer_name",
			},
		},
	}

	b, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	res := client.POST(b)
	if res.Err != nil {
		return res.Err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("apisix request failed: %s - %s", res.Status, string(res.ResponseBody))
	}

	return nil
}
