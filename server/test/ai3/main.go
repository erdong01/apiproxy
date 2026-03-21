package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

const (
	baseURL         = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks"
	defaultDate     = "2026-03-18"
	defaultPageSize = 100
)

type taskListResponse struct {
	Items []task `json:"items"`
	Total int    `json:"total"`
}

type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}
type task struct {
	ID          string          `json:"id"`
	Model       string          `json:"model"`
	Status      string          `json:"status"`
	CreatedAt   json.RawMessage `json:"created_at"`
	UpdatedAt   json.RawMessage `json:"updated_at"`
	Usage       Usage           `json:"usage"`
	DraftTaskId string          `json:"draft_task_id"`
}

func main() {
	apiKey := "3c2605da-c453-4c9a-8ced-dfb7835b979d"
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ARK_API_KEY 未设置")
		os.Exit(1)
	}

	targetDate := defaultDate
	if len(os.Args) > 1 && strings.TrimSpace(os.Args[1]) != "" {
		targetDate = strings.TrimSpace(os.Args[1])
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载时区失败: %v\n", err)
		os.Exit(1)
	}

	startOfDay, err := time.ParseInLocation("2006-01-02", targetDate, loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "日期格式错误，请使用 YYYY-MM-DD: %v\n", err)
		os.Exit(1)
	}
	endOfDay := startOfDay.Add(24 * time.Hour)
	printMatches := strings.EqualFold(strings.TrimSpace(os.Getenv("PRINT_MATCHES")), "1")

	tasks, total, err := listAllTasks(apiKey, "", defaultPageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "查询任务列表失败: %v\n", err)
		os.Exit(1)
	}

	matched := make([]task, 0)
	statusStats := make(map[string]int)
	var totalAmount float64
	for _, item := range tasks {
		createdAt, err := parseUnixSeconds(item.CreatedAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "跳过无法解析 created_at 的任务 %s: %v\n", item.ID, err)
			continue
		}
		createdTime := time.Unix(createdAt, 0).In(loc)
		if !createdTime.Before(startOfDay) && createdTime.Before(endOfDay) {
			matched = append(matched, item)
			statusStats[normalizeStatus(item.Status)]++
			cType := "text"
			if item.DraftTaskId != "" {
				cType = "draft_task"
			}
			totalAmount += ai.Calculate(item.Model, cType, item.Usage.TotalTokens)
		}

	}

	fmt.Printf("统计日期：%s\n", targetDate)
	fmt.Printf("统计时区：%s（北京时间）\n", loc.String())
	fmt.Printf("接口任务总数：%d\n", total)
	fmt.Printf("实际拉取任务数：%d\n", len(tasks))
	fmt.Printf("符合日期的任务数（全部状态）：%d\n", len(matched))
	fmt.Printf("总费用：%f\n", totalAmount)

	for _, key := range []string{"succeeded", "failed", "running", "queued", "cancelling", "cancelled"} {
		if count, ok := statusStats[key]; ok {
			fmt.Printf("状态 %s：%d\n", key, count)
			delete(statusStats, key)
		}
	}
	for status, count := range statusStats {
		fmt.Printf("状态 %s：%d\n", status, count)
	}

	if len(matched) > 0 {
		firstCreatedAt, _ := parseUnixSeconds(matched[0].CreatedAt)
		lastCreatedAt, _ := parseUnixSeconds(matched[len(matched)-1].CreatedAt)
		fmt.Printf("最晚的一条匹配任务创建时间：%s\n", time.Unix(firstCreatedAt, 0).In(loc).Format(time.DateTime))
		fmt.Printf("最早的一条匹配任务创建时间：%s\n", time.Unix(lastCreatedAt, 0).In(loc).Format(time.DateTime))
	}

	if !printMatches {
		return
	}

	for _, item := range matched {
		createdAt, _ := parseUnixSeconds(item.CreatedAt)
		updatedAt, _ := parseUnixSeconds(item.UpdatedAt)
		fmt.Printf("任务ID=%s\t模型=%s\t状态=%s\t创建时间=%s\t更新时间=%s\n",
			item.ID,
			item.Model,
			item.Status,
			time.Unix(createdAt, 0).In(loc).Format(time.DateTime),
			formatUnix(updatedAt, loc),
		)
	}
}

func listAllTasks(apiKey, status string, pageSize int) ([]task, int, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	pageNum := 1
	total := 0
	allTasks := make([]task, 0)

	for {
		resp, err := queryTaskPage(client, apiKey, pageNum, pageSize, status)
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			total = resp.Total
		}
		if len(resp.Items) == 0 {
			break
		}

		allTasks = append(allTasks, resp.Items...)
		if len(allTasks) >= resp.Total || len(resp.Items) < pageSize {
			break
		}
		pageNum++
	}

	return allTasks, total, nil
}

func queryTaskPage(client *http.Client, apiKey string, pageNum, pageSize int, status string) (*taskListResponse, error) {
	query := url.Values{}
	query.Set("page_num", strconv.Itoa(pageNum))
	query.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		query.Set("filter.status", status)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result taskListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body=%s", err, strings.TrimSpace(string(body)))
	}

	return &result, nil
}

func parseUnixSeconds(raw json.RawMessage) (int64, error) {
	s := strings.TrimSpace(string(raw))
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" {
		return 0, fmt.Errorf("empty value")
	}
	return strconv.ParseInt(s, 10, 64)
}

func formatUnix(ts int64, loc *time.Location) string {
	if ts <= 0 {
		return ""
	}
	return time.Unix(ts, 0).In(loc).Format(time.DateTime)
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "unknown"
	}
	return status
}
