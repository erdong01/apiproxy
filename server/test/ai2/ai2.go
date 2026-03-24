package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/erdong01/kit"
)

func main() {
	queryArkTask("cgt-20260323180143-wwsjr", "3c2605da-c453-4c9a-8ced-dfb7835b979d")
}

func queryArkTask(taskId, apiKey string) (*map[string]any, error) {
	url := fmt.Sprintf("https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks/%s", taskId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP 状态码错误: %d, 返回内容: %s", resp.StatusCode, string(bodyBytes))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}
	kit.DumpJson(data)
	return nil, nil
}
