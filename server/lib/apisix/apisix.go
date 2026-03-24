package apisix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/erdong01/kit/httpClient"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func DELETEConsumers(username string) error {
	cfg := global.GVA_CONFIG.Apisix
	deleteURL := consumersURL(cfg.Url, username)

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

func POSTConsumers(username, apiKey, userKey string, rate int) error {
	if rate == 0 {
		rate = 5
	}
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
				"rate":          rate,
				"burst":         5,
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
	cfg := global.GVA_CONFIG.Apisix
	putURL := consumersURL(cfg.Url, username)

	client := httpClient.New(putURL)
	h := make(http.Header)
	h.Set("X-API-KEY", cfg.ApiKey)
	client.Header = h
	res := client.SetMethod(http.MethodPut).Do(b)
	if res.Err != nil {
		return res.Err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("apisix request failed: %s - %s", res.Status, string(res.ResponseBody))
	}

	return nil
}

func consumersURL(base string, paths ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}

	adminPath := "/apisix/admin/consumers"
	if !hasConsumersPath(u.Path) {
		u.Path = joinURLPath(u.Path, adminPath)
	}
	for _, p := range paths {
		u.Path = joinURLPath(u.Path, url.PathEscape(p))
	}
	return u.String()
}

func hasConsumersPath(path string) bool {
	return path == "/apisix/admin/consumers" || path == "/apisix/admin/consumers/"
}

func joinURLPath(base, elem string) string {
	base = trimTrailingSlash(base)
	elem = trimLeadingSlash(elem)
	if base == "" {
		return "/" + elem
	}
	if elem == "" {
		return base
	}
	return base + "/" + elem
}

func trimLeadingSlash(s string) string {
	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
