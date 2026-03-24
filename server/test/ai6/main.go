package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/video"
	"github.com/shopspring/decimal"
)

type Data struct {
	Row         int         `json:"row"`
	InputPrompt InputPrompt `json:"input_prompt"`
	TaskId      string      `json:"task_id"`
	Task        task        `json:"task"`

	TotalTokensAmount float64 // 火山
	CostPerSecond     float64 // 火山每秒成本
	PanQuAmount       float64 // 盼趣成本
	DifAmount         float64 // 差价
	IsDraftVideo      string  // 是否包含视频
}

type InputPrompt struct {
	Content []Content `json:"content"`
}

type Content struct {
	Text     string   `json:"text,omitempty"`
	Type     string   `json:"type"` // video_url image_url text
	Role     string   `json:"role,omitempty"`
	VideoURL *FileURL `json:"video_url,omitempty"`
	ImageURL *FileURL `json:"image_url,omitempty"`
}

type FileURL struct {
	URL string `json:"url"`
}

var videoURLPattern = regexp.MustCompile(`"video_url"\s*:\s*\{\s*"url"\s*:\s*"([^"]+)"`)

func main() {
	var data []Data
	filePath, err := resolveCSVPath("pq_volcengine_ai_task_202603241533任务23-24 0点.csv")
	if err != nil {
		panic(err)
	}
	if err := ReadDataFromCSV(filePath, &data); err != nil {
		panic(err)
	}

	apiKey := "3c2605da-c453-4c9a-8ced-dfb7835b979d"
	if apiKey == "" {
		panic("ARK_API_KEY 未设置")
	}

	taskIDs := make([]string, 0, len(data))
	for _, item := range data {
		taskIDs = append(taskIDs, item.TaskId)
	}
	tasks, err := QueryTasksByIDs(apiKey, taskIDs)
	if err != nil {
		panic(err)
	}
	BindTasksToData(data, tasks)
	for i := range data {
		taskData := &data[i]
		contentType := "text"
		var videoMetadata *video.VideoMetadata
		for _, v := range taskData.InputPrompt.Content {
			if v.Type == "video_url" && v.VideoURL != nil && v.VideoURL.URL != "" {
				contentType = "draft_task"

				videoMetadata, _ = video.GetVideoMetadata(context.Background(), v.VideoURL.URL)
				break
			}
		}

		taskData.TotalTokensAmount = ai.Calculate(taskData.Task.Model, contentType, taskData.Task.Usage.TotalTokens)
		if taskData.Task.Duration > 0 {
			dDuration := decimal.NewFromInt(taskData.Task.Duration)
			taskData.CostPerSecond, _ = decimal.NewFromFloat(taskData.TotalTokensAmount).Div(dDuration).Truncate(2).Float64()
		} else {
			taskData.CostPerSecond = 0
		}
		if videoMetadata != nil {

			taskData.PanQuAmount = ai.PanQuModelPriceCalculate(taskData.Task.Model, taskData.Task.Resolution,
				contentType,
				taskData.Task.Duration,
				int64(videoMetadata.DurationSeconds))
		} else {
			taskData.PanQuAmount = ai.PanQuModelPriceCalculate(taskData.Task.Model, taskData.Task.Resolution, "", taskData.Task.Duration, 0)
		}

		taskData.DifAmount = taskData.PanQuAmount - taskData.TotalTokensAmount
		if contentType == "draft_task" {
			taskData.IsDraftVideo = "是"
		} else {
			taskData.IsDraftVideo = "否"
		}
	}

	updatedRows, err := WriteCSVCell(filePath, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("csv 写入完成，更新 %d 行: %s\n", updatedRows, filePath)

}

type taskListResponse struct {
	Items []task `json:"items"`
	Total int    `json:"total"`
}

type task struct {
	ID          string          `json:"id"`
	Model       string          `json:"model"`
	Status      string          `json:"status"`
	CreatedAt   json.RawMessage `json:"created_at"`
	UpdatedAt   json.RawMessage `json:"updated_at"`
	Usage       Usage           `json:"usage"`
	DraftTaskId string          `json:"draft_task_id"`
	Resolution  string          `json:"resolution"`
	Duration    int64           `json:"duration"`
}

type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

const (
	maxTaskIDsPerRequest = 100
	taskQueryConcurrency = 8
)

var queryTaskListURL = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks"
var newTaskQueryHTTPClient = func() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
var errTaskNotFound = errors.New("ark task not found")

// QueryTasksByIDs 使用 filter.task_ids 查询任务列表并返回 []task。
func QueryTasksByIDs(apiKey string, taskIDs []string) ([]task, error) {
	taskIDs = normalizeTaskIDs(taskIDs)
	if len(taskIDs) == 0 {
		return nil, nil
	}

	filteredTaskIDs := make([]string, 0, len(taskIDs))
	skippedTaskIDs := make([]string, 0)
	for _, taskID := range taskIDs {
		if isLikelyArkTaskID(taskID) {
			filteredTaskIDs = append(filteredTaskIDs, taskID)
			continue
		}
		skippedTaskIDs = append(skippedTaskIDs, taskID)
	}
	if len(skippedTaskIDs) > 0 {
		fmt.Fprintf(os.Stderr, "skip %d non-Ark task ids, examples: %s\n", len(skippedTaskIDs), previewTaskIDs(skippedTaskIDs, 5))
	}
	if len(filteredTaskIDs) == 0 {
		return nil, nil
	}

	client := newTaskQueryHTTPClient()
	results := make([]task, 0, len(filteredTaskIDs))

	for start := 0; start < len(filteredTaskIDs); start += maxTaskIDsPerRequest {
		end := start + maxTaskIDsPerRequest
		if end > len(filteredTaskIDs) {
			end = len(filteredTaskIDs)
		}

		items, err := queryTasksByIDBatch(client, apiKey, filteredTaskIDs[start:end])
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	return results, nil
}

func BindTasksToData(data []Data, tasks []task) {
	taskMap := make(map[string]task, len(tasks))
	for _, item := range tasks {
		taskMap[item.ID] = item
	}

	for i := range data {
		if item, ok := taskMap[data[i].TaskId]; ok {
			data[i].Task = item
		}
	}
}

func queryTasksByIDBatch(client *http.Client, apiKey string, taskIDs []string) ([]task, error) {
	requested := make(map[string]struct{}, len(taskIDs))
	for _, taskID := range taskIDs {
		requested[taskID] = struct{}{}
	}

	taskMap := make(map[string]task, len(taskIDs))
	if items, err := queryTasksByListFilter(client, apiKey, taskIDs); err == nil {
		for _, item := range items {
			if _, ok := requested[item.ID]; ok {
				taskMap[item.ID] = item
			}
		}
	}

	missingIDs := make([]string, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		if _, ok := taskMap[taskID]; !ok {
			missingIDs = append(missingIDs, taskID)
		}
	}

	if len(missingIDs) > 0 {
		items, err := queryTasksByDetail(client, apiKey, missingIDs)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			taskMap[item.ID] = item
		}
	}

	results := make([]task, 0, len(taskMap))
	for _, taskID := range taskIDs {
		if item, ok := taskMap[taskID]; ok {
			results = append(results, item)
		}
	}

	return results, nil
}

func queryTasksByListFilter(client *http.Client, apiKey string, taskIDs []string) ([]task, error) {
	query := url.Values{}
	query.Set("page_num", "1")
	query.Set("page_size", fmt.Sprintf("%d", len(taskIDs)))
	query.Set("filter.task_ids", strings.Join(taskIDs, ","))

	req, err := http.NewRequest(http.MethodGet, queryTaskListURL+"?"+query.Encode(), nil)
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
		return nil, fmt.Errorf("parse response failed: %w, body=%s", err, strings.TrimSpace(string(body)))
	}

	return result.Items, nil
}

func queryTasksByDetail(client *http.Client, apiKey string, taskIDs []string) ([]task, error) {
	type result struct {
		task task
		err  error
		skip bool
	}

	workers := taskQueryConcurrency
	if len(taskIDs) < workers {
		workers = len(taskIDs)
	}
	if workers <= 0 {
		return nil, nil
	}

	taskCh := make(chan string)
	resultCh := make(chan result, len(taskIDs))

	for i := 0; i < workers; i++ {
		go func() {
			for taskID := range taskCh {
				item, err := queryTaskByID(client, apiKey, taskID)
				if err != nil {
					if errors.Is(err, errTaskNotFound) {
						fmt.Fprintf(os.Stderr, "skip ark task not found: %s\n", taskID)
						resultCh <- result{skip: true}
						continue
					}
					resultCh <- result{err: err}
					continue
				}
				resultCh <- result{task: *item}
			}
		}()
	}

	for _, taskID := range taskIDs {
		taskCh <- taskID
	}
	close(taskCh)

	items := make([]task, 0, len(taskIDs))
	var errs []error
	for i := 0; i < len(taskIDs); i++ {
		res := <-resultCh
		if res.skip {
			continue
		}
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		items = append(items, res.task)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return items, nil
}

func queryTaskByID(client *http.Client, apiKey, taskID string) (*task, error) {
	detailURL := strings.TrimRight(queryTaskListURL, "/") + "/" + url.PathEscape(taskID)
	req, err := http.NewRequest(http.MethodGet, detailURL, nil)
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
		if resp.StatusCode == http.StatusNotFound && strings.Contains(string(body), "ResourceNotFound") {
			return nil, fmt.Errorf("%w: %s", errTaskNotFound, taskID)
		}
		return nil, fmt.Errorf("query task %s http %d: %s", taskID, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var item task
	if err := json.Unmarshal(body, &item); err == nil && item.ID != "" {
		return &item, nil
	}

	var wrapped struct {
		Data task `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, fmt.Errorf("parse task %s failed: %w, body=%s", taskID, err, strings.TrimSpace(string(body)))
	}
	if wrapped.Data.ID == "" {
		return nil, fmt.Errorf("task %s not found in response: %s", taskID, strings.TrimSpace(string(body)))
	}

	return &wrapped.Data, nil
}

func normalizeTaskIDs(taskIDs []string) []string {
	seen := make(map[string]struct{}, len(taskIDs))
	result := make([]string, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		taskID = strings.TrimSpace(taskID)
		if taskID == "" {
			continue
		}
		if _, ok := seen[taskID]; ok {
			continue
		}
		seen[taskID] = struct{}{}
		result = append(result, taskID)
	}
	return result
}

func isLikelyArkTaskID(taskID string) bool {
	return strings.HasPrefix(taskID, "cgt-")
}

func previewTaskIDs(taskIDs []string, limit int) string {
	if len(taskIDs) == 0 {
		return ""
	}
	if limit <= 0 || len(taskIDs) <= limit {
		return strings.Join(taskIDs, ", ")
	}
	return strings.Join(taskIDs[:limit], ", ") + fmt.Sprintf(" ... (+%d more)", len(taskIDs)-limit)
}

func ReadDataFromCSV(filePath string, data *[]Data) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open csv failed: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv failed: %w", err)
	}
	if len(records) == 0 {
		*data = nil
		return nil
	}

	header := make(map[string]int, len(records[0]))
	for i, name := range records[0] {
		header[normalizeHeaderName(name)] = i
	}

	taskIDIndex, ok := header["task_id"]
	if !ok {
		return fmt.Errorf("csv missing task_id column")
	}

	inputPromptIndex, ok := header["input_prompt"]
	if !ok {
		return fmt.Errorf("csv missing input_prompt column")
	}

	result := make([]Data, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		record := records[i]
		if isEmptyRecord(record) {
			continue
		}

		item := Data{
			Row: i,
		}

		if taskIDIndex < len(record) {
			item.TaskId = normalizeCSVValue(record[taskIDIndex])
		}

		if inputPromptIndex < len(record) {
			rawInputPrompt := normalizeCSVValue(record[inputPromptIndex])
			if rawInputPrompt == "" {
				result = append(result, item)
				continue
			}

			inputPrompt, err := parseInputPrompt(rawInputPrompt)
			if err != nil {
				fmt.Printf("warn: parse input_prompt failed, row=%d task_id=%s: %v\n", i, item.TaskId, err)
			} else {
				item.InputPrompt = inputPrompt
			}
		}

		result = append(result, item)
	}

	*data = result
	return nil
}

func parseInputPrompt(raw string) (InputPrompt, error) {
	var inputPrompt InputPrompt
	if err := json.Unmarshal([]byte(raw), &inputPrompt); err == nil {
		return inputPrompt, nil
	}

	matches := videoURLPattern.FindAllStringSubmatch(raw, -1)
	if len(matches) == 0 {
		return InputPrompt{}, fmt.Errorf("invalid json and no video_url fallback matched")
	}

	content := make([]Content, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 || strings.TrimSpace(match[1]) == "" {
			continue
		}
		content = append(content, Content{
			Type:     "video_url",
			VideoURL: &FileURL{URL: match[1]},
		})
	}
	if len(content) == 0 {
		return InputPrompt{}, fmt.Errorf("invalid json and extracted video_url is empty")
	}

	return InputPrompt{Content: content}, nil
}

func WriteCSVCell(filePath string, data []Data) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open csv failed: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("read csv failed: %w", err)
	}
	if len(records) == 0 {
		return 0, fmt.Errorf("csv is empty")
	}

	header := records[0]
	columnIndex := make(map[string]int, len(header))
	for i, name := range header {
		columnIndex[normalizeHeaderName(name)] = i
	}

	rowIndexByTaskID := make(map[string]int, len(records)-1)
	taskIDIdx, hasTaskID := columnIndex["task_id"]
	if hasTaskID {
		for rowIdx := 1; rowIdx < len(records); rowIdx++ {
			row := records[rowIdx]
			if taskIDIdx >= len(row) {
				continue
			}
			taskID := strings.TrimSpace(row[taskIDIdx])
			if taskID == "" {
				continue
			}
			rowIndexByTaskID[taskID] = rowIdx
		}
	}

	ensureColumn := func(name string) int {
		if idx, ok := columnIndex[name]; ok {
			return idx
		}
		header = append(header, name)
		idx := len(header) - 1
		columnIndex[name] = idx
		return idx
	}

	resolutionIdx := ensureColumn("分辨率")
	durationIdx := ensureColumn("生成视频秒数")
	totalTokensIdx := ensureColumn("火山销售总tokens")
	totalTokensAmountIdx := ensureColumn("火山费用")
	costPerSecondIdx := ensureColumn("火山每秒价格")
	panQuAmountIdx := ensureColumn("盼趣成本")
	diffAmountIdx := ensureColumn("盼趣成本-火山费用=差额")
	isDraftVideoIdx := ensureColumn("是否包含视频")
	records[0] = header

	maxIdx := len(header) - 1
	updatedRows := 0
	for _, item := range data {
		rowIdx := item.Row
		if matchedRowIdx, ok := rowIndexByTaskID[strings.TrimSpace(item.TaskId)]; ok {
			rowIdx = matchedRowIdx
		}
		if rowIdx <= 0 || rowIdx >= len(records) {
			continue
		}

		row := records[rowIdx]
		for len(row) <= maxIdx {
			row = append(row, "")
		}

		row[resolutionIdx] = item.Task.Resolution
		row[durationIdx] = strconv.FormatInt(item.Task.Duration, 10)
		row[totalTokensIdx] = strconv.FormatInt(item.Task.Usage.TotalTokens, 10)
		row[totalTokensAmountIdx] = strconv.FormatFloat(item.TotalTokensAmount, 'f', 6, 64)
		row[costPerSecondIdx] = strconv.FormatFloat(item.CostPerSecond, 'f', 6, 64)
		row[panQuAmountIdx] = strconv.FormatFloat(item.PanQuAmount, 'f', 6, 64)
		row[diffAmountIdx] = strconv.FormatFloat(item.DifAmount, 'f', 6, 64)
		row[isDraftVideoIdx] = item.IsDraftVideo

		records[rowIdx] = row
		updatedRows++
	}

	tempFilePath := filePath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return 0, fmt.Errorf("create temp csv failed: %w", err)
	}
	defer func() {
		_ = tempFile.Close()
	}()

	if _, err := tempFile.WriteString("\xef\xbb\xbf"); err != nil {
		return 0, fmt.Errorf("write bom failed: %w", err)
	}

	writer := csv.NewWriter(tempFile)
	if err := writer.WriteAll(records); err != nil {
		return 0, fmt.Errorf("write csv failed: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return 0, fmt.Errorf("flush csv failed: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		return 0, fmt.Errorf("sync temp csv failed: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return 0, fmt.Errorf("close temp csv failed: %w", err)
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		return 0, fmt.Errorf("replace csv failed: %w", err)
	}

	return updatedRows, nil
}

func normalizeHeaderName(name string) string {
	name = strings.TrimPrefix(name, "\ufeff")
	name = strings.TrimPrefix(name, "\xef\xbb\xbf")
	return strings.TrimSpace(name)
}

func normalizeCSVValue(value string) string {
	value = strings.TrimPrefix(value, "\ufeff")
	value = strings.TrimPrefix(value, "\xef\xbb\xbf")
	return strings.TrimSpace(value)
}

func isEmptyRecord(record []string) bool {
	for _, value := range record {
		if normalizeCSVValue(value) != "" {
			return false
		}
	}
	return true
}

func resolveCSVPath(fileName string) (string, error) {
	if _, err := os.Stat(fileName); err == nil {
		return fileName, nil
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve csv path failed: runtime caller unavailable")
	}

	resolvedPath := filepath.Join(filepath.Dir(currentFile), fileName)
	if _, err := os.Stat(resolvedPath); err != nil {
		return "", fmt.Errorf("resolve csv path failed: %w", err)
	}
	return resolvedPath, nil
}
