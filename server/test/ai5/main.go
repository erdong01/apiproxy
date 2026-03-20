package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

const (
	baseURL         = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks"
	defaultStatus   = "succeeded"
	defaultPageNum  = 1
	defaultPageSize = 500
)

type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"` //模型输出视频花费的 token 数量。
	TotalTokens      int64 `json:"total_tokens"`      //本次请求消耗的总 token 数量。视频生成模型不统计输入 token，输入 token 为 0，故 total_tokens=completion_tokens。
}

type Tasks struct {
	ID              string      `json:"id"`         // 视频生成任务 ID。
	Model           string      `json:"model"`      // 任务使用的模型名称和版本。
	Status          string      `json:"status"`     //任务状态
	Resolution      string      `json:"resolution"` //生成视频的分辨率。
	Duration        interface{} `json:"duration"`   //生成视频的时长，单位：秒。
	Frames          interface{} `json:"frames"`     //生成视频的帧数。
	FramesPerSecond interface{} `json:"framespersecond"`
	Usage           Usage       `json:"usage"`         //本次请求的 token 用量。
	DraftTaskId     string      `json:"draft_task_id"` //Draft 参考视频任务 ID。基于 Draft 视频生成正式视频时，会返回该参数。
}

type listTasksResponse struct {
	Total int     `json:"total"`
	Items []Tasks `json:"items"`
}

type ModelResolutionStatisticalResults struct {
	Model                       string             `json:"model"`
	Resolution                  string             `json:"resolution"`
	StatisticalResults          StatisticalResults `json:"statisticalResults"`          // 无参考视频id 平均成本
	DraftTaskStatisticalResults StatisticalResults `json:"draftTaskStatisticalResults"` // 有参考视频id 平均成本
}

type StatisticalResults struct {
	AverageCostPerSecond float64           `json:"averageCostPerSecond"` // 总的每条平均成本
	EachHighCost         EachHighCost      `json:"eachHighCost"`         // 最高的20%数据的 平均成本
	MedianAverageCost    MedianAverageCost `json:"medianAverageCost"`    // 中间的60%数据的 平均成本
	EveryLowCost         EveryLowCost      `json:"everyLowCost"`         // 最低的20%数据的 平均成本
	Count                int               `json:"count"`                // 条数统计
}

type EachHighCost struct {
	AverageCostPerSecond float64 `json:"averageCostPerSecond"` // 每条平均成本
	NumberItems          int     `json:"numberItems"`          // 条数
	HighestPrice         float64 `json:"highestPrice"`         // 区间内最高金额
	LowestPrice          float64 `json:"lowestPrice"`          // 区间内最低金额
}

type MedianAverageCost struct {
	AverageCostPerSecond float64 `json:"averageCostPerSecond"` // 每条平均成本
	NumberItems          int     `json:"numberItems"`          // 条数
	HighestPrice         float64 `json:"highestPrice"`         // 区间内最高金额
	LowestPrice          float64 `json:"lowestPrice"`          // 区间内最低金额
}

type EveryLowCost struct {
	AverageCostPerSecond float64 `json:"averageCostPerSecond"` // 每条平均成本
	NumberItems          int     `json:"numberItems"`          // 条数
	HighestPrice         float64 `json:"highestPrice"`         // 区间内最高金额
	LowestPrice          float64 `json:"lowestPrice"`          // 区间内最低金额
}

type taskMetric struct {
	ID                string
	Model             string
	Resolution        string
	DurationSeconds   float64
	CostPerSecond     float64
	HasDraftReference bool
}

type groupMetrics struct {
	withoutDraft []taskMetric
	withDraft    []taskMetric
}

func main() {
	apiKey := "3c2605da-c453-4c9a-8ced-dfb7835b979d"
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ARK_API_KEY environment variable is required")
		os.Exit(1)
	}

	tasks, err := fetchAllTasks(apiKey, defaultStatus, defaultPageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch tasks failed: %v\n", err)
		os.Exit(1)
	}

	results, skipped := buildStatistics(tasks)
	if len(skipped) > 0 {
		fmt.Fprintf(os.Stderr, "skipped %d tasks due to missing duration/resolution or unsupported pricing\n", len(skipped))
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal results failed: %v\n", err)
		os.Exit(1)
	}

	if err := writeResultsCSV(results); err != nil {
		fmt.Fprintf(os.Stderr, "write csv failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func fetchAllTasks(apiKey, status string, pageSize int) ([]Tasks, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	pageNum := defaultPageNum
	var all []Tasks
	total := -1

	for {
		resp, err := fetchTaskPage(client, apiKey, status, pageNum, pageSize)
		if err != nil {
			return nil, err
		}

		if total < 0 {
			total = resp.Total
		}

		all = append(all, resp.Items...)
		if len(resp.Items) == 0 || len(all) >= resp.Total {
			break
		}

		pageNum++
	}

	return all, nil
}

func fetchTaskPage(client *http.Client, apiKey, status string, pageNum, pageSize int) (*listTasksResponse, error) {
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
		return nil, fmt.Errorf("list tasks http status %d: %s", resp.StatusCode, string(body))
	}

	var result listTasksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func buildStatistics(tasks []Tasks) ([]ModelResolutionStatisticalResults, []string) {
	grouped := make(map[string]*groupMetrics)
	var skipped []string

	for _, task := range tasks {
		metric, ok := buildTaskMetric(task)
		if !ok {
			skipped = append(skipped, task.ID)
			continue
		}

		groupKey := metric.Model + "|" + metric.Resolution
		group, exists := grouped[groupKey]
		if !exists {
			group = &groupMetrics{}
			grouped[groupKey] = group
		}

		if metric.HasDraftReference {
			group.withDraft = append(group.withDraft, metric)
			continue
		}

		group.withoutDraft = append(group.withoutDraft, metric)
	}

	results := make([]ModelResolutionStatisticalResults, 0, len(grouped))
	for key, group := range grouped {
		model, resolution, _ := strings.Cut(key, "|")
		results = append(results, ModelResolutionStatisticalResults{
			Model:                       model,
			Resolution:                  resolution,
			StatisticalResults:          computeStatisticalResults(group.withoutDraft),
			DraftTaskStatisticalResults: computeStatisticalResults(group.withDraft),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Model == results[j].Model {
			return results[i].Resolution < results[j].Resolution
		}
		return results[i].Model < results[j].Model
	})

	return results, skipped
}

func buildTaskMetric(task Tasks) (taskMetric, bool) {
	duration, ok := extractDurationSeconds(task)
	if !ok || duration <= 0 {
		return taskMetric{}, false
	}

	resolution := strings.TrimSpace(task.Resolution)
	if resolution == "" {
		return taskMetric{}, false
	}

	contentType := "text"
	hasDraftReference := strings.TrimSpace(task.DraftTaskId) != ""
	if hasDraftReference {
		contentType = "draft_task"
	}

	totalPrice := ai.Calculate(task.Model, contentType, task.Usage.TotalTokens)
	if totalPrice <= 0 {
		return taskMetric{}, false
	}

	return taskMetric{
		ID:                task.ID,
		Model:             task.Model,
		Resolution:        resolution,
		DurationSeconds:   duration,
		CostPerSecond:     roundTo(totalPrice / duration),
		HasDraftReference: hasDraftReference,
	}, true
}

func extractDurationSeconds(task Tasks) (float64, bool) {
	if duration, ok := parseFloat(task.Duration); ok && duration > 0 {
		return duration, true
	}

	frames, okFrames := parseFloat(task.Frames)
	fps, okFPS := parseFloat(task.FramesPerSecond)
	if okFrames && okFPS && fps > 0 {
		return frames / fps, true
	}

	return 0, false
}

func parseFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case nil:
		return 0, false
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func computeStatisticalResults(metrics []taskMetric) StatisticalResults {
	if len(metrics) == 0 {
		return StatisticalResults{}
	}

	sorted := append([]taskMetric(nil), metrics...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CostPerSecond < sorted[j].CostPerSecond
	})

	lowCount := int(math.Floor(float64(len(sorted)) * 0.2))
	highCount := int(math.Floor(float64(len(sorted)) * 0.2))
	middleStart := lowCount
	middleEnd := len(sorted) - highCount
	if middleEnd < middleStart {
		middleEnd = middleStart
	}

	lowItems := sorted[:lowCount]
	middleItems := sorted[middleStart:middleEnd]
	highItems := sorted[middleEnd:]

	return StatisticalResults{
		Count:                len(metrics),
		AverageCostPerSecond: averageCost(metrics),
		EachHighCost:         toEachHighCost(highItems),
		MedianAverageCost:    toMedianAverageCost(middleItems),
		EveryLowCost:         toEveryLowCost(lowItems),
	}
}

func averageCost(metrics []taskMetric) float64 {
	if len(metrics) == 0 {
		return 0
	}

	var sum float64
	for _, metric := range metrics {
		sum += metric.CostPerSecond
	}

	return roundTo(sum / float64(len(metrics)))
}

func toEachHighCost(metrics []taskMetric) EachHighCost {
	stats := buildIntervalStats(metrics)
	return EachHighCost{
		AverageCostPerSecond: stats.averageCostPerSecond,
		NumberItems:          stats.numberItems,
		HighestPrice:         stats.highestPrice,
		LowestPrice:          stats.lowestPrice,
	}
}

func toMedianAverageCost(metrics []taskMetric) MedianAverageCost {
	stats := buildIntervalStats(metrics)
	return MedianAverageCost{
		AverageCostPerSecond: stats.averageCostPerSecond,
		NumberItems:          stats.numberItems,
		HighestPrice:         stats.highestPrice,
		LowestPrice:          stats.lowestPrice,
	}
}

func toEveryLowCost(metrics []taskMetric) EveryLowCost {
	stats := buildIntervalStats(metrics)
	return EveryLowCost{
		AverageCostPerSecond: stats.averageCostPerSecond,
		NumberItems:          stats.numberItems,
		HighestPrice:         stats.highestPrice,
		LowestPrice:          stats.lowestPrice,
	}
}

type intervalStats struct {
	averageCostPerSecond float64
	numberItems          int
	highestPrice         float64
	lowestPrice          float64
}

func buildIntervalStats(metrics []taskMetric) intervalStats {
	if len(metrics) == 0 {
		return intervalStats{}
	}

	var sum float64
	lowest := metrics[0].CostPerSecond
	highest := metrics[0].CostPerSecond

	for _, metric := range metrics {
		sum += metric.CostPerSecond
		if metric.CostPerSecond < lowest {
			lowest = metric.CostPerSecond
		}
		if metric.CostPerSecond > highest {
			highest = metric.CostPerSecond
		}
	}

	return intervalStats{
		averageCostPerSecond: roundTo(sum / float64(len(metrics))),
		numberItems:          len(metrics),
		highestPrice:         roundTo(highest),
		lowestPrice:          roundTo(lowest),
	}
}

func roundTo(value float64) float64 {
	return math.Round(value*1_000_000) / 1_000_000
}

func writeResultsCSV(results []ModelResolutionStatisticalResults) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(wd, "model_resolution_statistical_results.csv")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString("\xef\xbb\xbf"); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"模型",
		"分辨率",
		"无参考视频总条数",
		"无参考视频总平均成本/秒",
		"无参考视频最高20%条数",
		"无参考视频最高20%平均成本/秒",
		"无参考视频最高20%最高金额",
		"无参考视频最高20%最低金额",
		"无参考视频中间60%条数",
		"无参考视频中间60%平均成本/秒",
		"无参考视频中间60%最高金额",
		"无参考视频中间60%最低金额",
		"无参考视频最低20%条数",
		"无参考视频最低20%平均成本/秒",
		"无参考视频最低20%最高金额",
		"无参考视频最低20%最低金额",
		"有参考视频总条数",
		"有参考视频总平均成本/秒",
		"有参考视频最高20%条数",
		"有参考视频最高20%平均成本/秒",
		"有参考视频最高20%最高金额",
		"有参考视频最高20%最低金额",
		"有参考视频中间60%条数",
		"有参考视频中间60%平均成本/秒",
		"有参考视频中间60%最高金额",
		"有参考视频中间60%最低金额",
		"有参考视频最低20%条数",
		"有参考视频最低20%平均成本/秒",
		"有参考视频最低20%最高金额",
		"有参考视频最低20%最低金额",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		row := []string{
			result.Model,
			result.Resolution,
			strconv.Itoa(result.StatisticalResults.Count),
			formatFloat(result.StatisticalResults.AverageCostPerSecond),
			strconv.Itoa(result.StatisticalResults.EachHighCost.NumberItems),
			formatFloat(result.StatisticalResults.EachHighCost.AverageCostPerSecond),
			formatFloat(result.StatisticalResults.EachHighCost.HighestPrice),
			formatFloat(result.StatisticalResults.EachHighCost.LowestPrice),
			strconv.Itoa(result.StatisticalResults.MedianAverageCost.NumberItems),
			formatFloat(result.StatisticalResults.MedianAverageCost.AverageCostPerSecond),
			formatFloat(result.StatisticalResults.MedianAverageCost.HighestPrice),
			formatFloat(result.StatisticalResults.MedianAverageCost.LowestPrice),
			strconv.Itoa(result.StatisticalResults.EveryLowCost.NumberItems),
			formatFloat(result.StatisticalResults.EveryLowCost.AverageCostPerSecond),
			formatFloat(result.StatisticalResults.EveryLowCost.HighestPrice),
			formatFloat(result.StatisticalResults.EveryLowCost.LowestPrice),
			strconv.Itoa(result.DraftTaskStatisticalResults.Count),
			formatFloat(result.DraftTaskStatisticalResults.AverageCostPerSecond),
			strconv.Itoa(result.DraftTaskStatisticalResults.EachHighCost.NumberItems),
			formatFloat(result.DraftTaskStatisticalResults.EachHighCost.AverageCostPerSecond),
			formatFloat(result.DraftTaskStatisticalResults.EachHighCost.HighestPrice),
			formatFloat(result.DraftTaskStatisticalResults.EachHighCost.LowestPrice),
			strconv.Itoa(result.DraftTaskStatisticalResults.MedianAverageCost.NumberItems),
			formatFloat(result.DraftTaskStatisticalResults.MedianAverageCost.AverageCostPerSecond),
			formatFloat(result.DraftTaskStatisticalResults.MedianAverageCost.HighestPrice),
			formatFloat(result.DraftTaskStatisticalResults.MedianAverageCost.LowestPrice),
			strconv.Itoa(result.DraftTaskStatisticalResults.EveryLowCost.NumberItems),
			formatFloat(result.DraftTaskStatisticalResults.EveryLowCost.AverageCostPerSecond),
			formatFloat(result.DraftTaskStatisticalResults.EveryLowCost.HighestPrice),
			formatFloat(result.DraftTaskStatisticalResults.EveryLowCost.LowestPrice),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return writer.Error()
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}
