package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type ArkTaskResult struct {
	CompletionTokens  interface{}
	TotalTokens       interface{}
	Status            string
	ErrorCode         string
	ErrorMessage      string
	TotalTokensAmount float64 //tokens 火山云费用金额
	Duration          float64 // 时长
	TotalDuration     float64 // 时长
	Resolution        string
	CostPerSecond     string // 每条成本  TotalTokensAmount / Duration
	DraftTaskId       string
	PanQuAmount       float64 //盼趣成本
	AmountDiff        float64 //盼趣成本 和 火山云费用差异 AmountDiff -  TotalTokensAmount

}

func main() {
	apiKey := "3c2605da-c453-4c9a-8ced-dfb7835b979d"
	if apiKey == "" {
		fmt.Println("Warning: ARK_API_KEY 环境变量未设置。如果有配置问题，请在运行前通过 export ARK_API_KEY=xxx 设置。")
	}

	fileName := "火山2.0视频任务列表(20260318).csv" // 替换为真实的CSV文件名
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("无法打开CSV文件 %s: %v\n", fileName, err)
		return
	}

	reader := csv.NewReader(f)
	// 防止因为 CSV 内包含不标准的引号格式而报错
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // 允许每行的列数不一致
	rows, err := reader.ReadAll()
	f.Close() // 提前关闭，方便稍后重新写入

	if err != nil {
		fmt.Printf("读取CSV失败: %v\n", err)
		return
	}

	if len(rows) == 0 {
		fmt.Println("表格中没有数据")
		return
	}

	normalizeCSVText := func(s string) string {
		s = strings.TrimPrefix(s, "\xef\xbb\xbf")
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "\"")
		return s
	}

	headerRow := rows[0]
	taskIdColIdx := -1
	completionColIdx := -1
	totalTokensColIdx := -1
	statusColIdx := -1
	codeColIdx := -1
	messageColIdx := -1
	totalTokensAmountIdx := -1
	durationColIdx := -1
	totalDurationColIdx := -1
	resolutionColIdx := -1
	costPerSecondColIdx := -1
	panQuAmountColIdx := -1
	amountDiffColIdx := -1
	hasReferenceVideoColIdx := -1

	// 查找列索引
	for i, colCell := range headerRow {
		hdr := strings.ToLower(normalizeCSVText(colCell))
		if strings.Contains(hdr, "task_id") {
			taskIdColIdx = i
		}
		if hdr == "output_result" && taskIdColIdx == -1 {
			// 在截图中，cgt-任务id似乎在 output_result 列
			taskIdColIdx = i
		}
		if hdr == "completion_tokens" || hdr == "completion" {
			completionColIdx = i
		}
		if hdr == "total_tokensq" || hdr == "total_tokens" || hdr == "total_token" {
			totalTokensColIdx = i
		}
		if hdr == "status" {
			statusColIdx = i
		}
		if hdr == "code" {
			codeColIdx = i
		}

		if hdr == "message" || hdr == "error_message" {
			messageColIdx = i
		}
		if hdr == "total_tokens_amount" {
			totalTokensAmountIdx = i
		}
		if hdr == "duration" || strings.Contains(hdr, "duration") {
			durationColIdx = i
		}
		if hdr == "total_duration" || hdr == "totalduration" {
			totalDurationColIdx = i
		}
		if hdr == "resolution" || strings.Contains(hdr, "resolution") {
			resolutionColIdx = i
		}
		if hdr == "cost_per_second" || strings.Contains(hdr, "cost_per_second") {
			costPerSecondColIdx = i
		}
		if hdr == "pan_qu_amount" {
			panQuAmountColIdx = i
		}
		if hdr == "amount_diff" {
			amountDiffColIdx = i
		}
		if hdr == "has_reference_video" {
			hasReferenceVideoColIdx = i
		}
	}

	// 如果没有找到写入的列，则在表头末尾追加
	if completionColIdx == -1 {
		completionColIdx = len(headerRow)
		headerRow = append(headerRow, "completion")
	}
	if totalTokensColIdx == -1 {
		totalTokensColIdx = len(headerRow)
		headerRow = append(headerRow, "total_token")
	}
	if statusColIdx == -1 {
		statusColIdx = len(headerRow)
		headerRow = append(headerRow, "status")
	}
	if codeColIdx == -1 {
		codeColIdx = len(headerRow)
		headerRow = append(headerRow, "code")
	}
	if messageColIdx == -1 {
		messageColIdx = len(headerRow)
		headerRow = append(headerRow, "message")
	}
	if totalTokensAmountIdx == -1 {
		totalTokensAmountIdx = len(headerRow)
		headerRow = append(headerRow, "total_tokens_amount")
	}
	if durationColIdx == -1 {
		durationColIdx = len(headerRow)
		headerRow = append(headerRow, "duration")
	}
	if totalDurationColIdx == -1 {
		totalDurationColIdx = len(headerRow)
		headerRow = append(headerRow, "total_duration")
	}
	if resolutionColIdx == -1 {
		resolutionColIdx = len(headerRow)
		headerRow = append(headerRow, "resolution")
	}
	if costPerSecondColIdx == -1 {
		costPerSecondColIdx = len(headerRow)
		headerRow = append(headerRow, "cost_per_second")
	}
	if panQuAmountColIdx == -1 {
		panQuAmountColIdx = len(headerRow)
		headerRow = append(headerRow, "pan_qu_amount")
	}
	if amountDiffColIdx == -1 {
		amountDiffColIdx = len(headerRow)
		headerRow = append(headerRow, "amount_diff")
	}
	if hasReferenceVideoColIdx == -1 {
		hasReferenceVideoColIdx = len(headerRow)
		headerRow = append(headerRow, "has_reference_video")
	}
	rows[0] = headerRow

	updates := 0

	for rIdx := 1; rIdx < len(rows); rIdx++ {
		row := rows[rIdx]

		var taskId string
		// 优先从 task_id/output_result 对应列找
		if taskIdColIdx != -1 && taskIdColIdx < len(row) {
			taskId = normalizeCSVText(row[taskIdColIdx])
		}

		// 容错：如果那一列不是cgt-开头，遍历当前行的每一列寻找真正的task_id
		if !strings.HasPrefix(taskId, "cgt-") {
			for _, cell := range row {
				nCell := normalizeCSVText(cell)
				if strings.HasPrefix(nCell, "cgt-") {
					taskId = nCell
					break
				}
			}
		}

		if taskId == "" || !strings.HasPrefix(taskId, "cgt-") {
			continue // 该行没有有效的任务ID跳过
		}

		// 检查是否已经存在完整的费用、时长、分辨率等数据
		hasAmount := totalTokensAmountIdx != -1 && totalTokensAmountIdx < len(row) && row[totalTokensAmountIdx] != ""
		hasDuration := durationColIdx != -1 && durationColIdx < len(row) && row[durationColIdx] != ""
		hasResolution := resolutionColIdx != -1 && resolutionColIdx < len(row) && row[resolutionColIdx] != ""
		hasCost := costPerSecondColIdx != -1 && costPerSecondColIdx < len(row) && row[costPerSecondColIdx] != ""
		isFailed := statusColIdx != -1 && statusColIdx < len(row) && row[statusColIdx] == "failed"

		// 满足条件直接跳过查询接口
		if (hasAmount && hasDuration && hasResolution && hasCost) || isFailed {
			// fmt.Printf("Task ID: %s (Row %d) 已有数据或已失败，跳过查询，直接同步...\n", taskId, rIdx+1)
			// continue
		}

		fmt.Printf("正在查询 Task ID: %s (Row %d)...\n", taskId, rIdx+1)
		res, err := queryArkTask(taskId, apiKey)
		if err != nil {
			fmt.Printf("  -> 查询Http层出错: %v\n", err)
			continue
		}

		if res.Status == "failed" {
			fmt.Printf("  -> 任务失败: code=%s, message=%s\n", res.ErrorCode, res.ErrorMessage)
		} else {
			fmt.Printf("  -> 查询成功: status=%s, completion: %v, total: %v\n", res.Status, res.CompletionTokens, res.TotalTokens)
		}

		// 确保当前行的长度足够容纳所有新列（以防有空缺）
		maxIdx := completionColIdx
		if totalTokensColIdx > maxIdx {
			maxIdx = totalTokensColIdx
		}
		if statusColIdx > maxIdx {
			maxIdx = statusColIdx
		}
		if codeColIdx > maxIdx {
			maxIdx = codeColIdx
		}
		if messageColIdx > maxIdx {
			maxIdx = messageColIdx
		}
		if totalTokensAmountIdx > maxIdx {
			maxIdx = totalTokensAmountIdx
		}
		if durationColIdx > maxIdx {
			maxIdx = durationColIdx
		}
		if totalDurationColIdx > maxIdx {
			maxIdx = totalDurationColIdx
		}
		if resolutionColIdx > maxIdx {
			maxIdx = resolutionColIdx
		}
		if costPerSecondColIdx > maxIdx {
			maxIdx = costPerSecondColIdx
		}
		if panQuAmountColIdx > maxIdx {
			maxIdx = panQuAmountColIdx
		}
		if amountDiffColIdx > maxIdx {
			maxIdx = amountDiffColIdx
		}
		if hasReferenceVideoColIdx > maxIdx {
			maxIdx = hasReferenceVideoColIdx
		}

		for len(row) <= maxIdx {
			row = append(row, "")
		}

		// 回写数据
		if res.CompletionTokens != nil {
			row[completionColIdx] = fmt.Sprintf("%v", res.CompletionTokens)
		}
		if res.TotalTokens != nil {
			row[totalTokensColIdx] = fmt.Sprintf("%v", res.TotalTokens)
		}
		row[statusColIdx] = res.Status
		row[codeColIdx] = res.ErrorCode
		row[messageColIdx] = res.ErrorMessage
		if res.TotalTokensAmount > 0 {
			row[totalTokensAmountIdx] = strconv.FormatFloat(res.TotalTokensAmount, 'f', 6, 64)
		}
		if res.Duration > 0 {
			row[durationColIdx] = strconv.FormatFloat(res.Duration, 'f', 2, 64)
		}
		if res.TotalDuration > 0 {
			row[totalDurationColIdx] = strconv.FormatFloat(res.TotalDuration, 'f', 2, 64)
		}
		if res.Resolution != "" {
			row[resolutionColIdx] = res.Resolution
		}

		// 计算 CostPerSecond (总花费 / 时长 或者 分辨率)
		var costPerSec float64
		if res.Duration > 0 {
			costPerSec = res.TotalTokensAmount / res.Duration
		} else if res.Resolution != "" {
			resFloat, err := strconv.ParseFloat(strings.TrimRight(res.Resolution, "p"), 64)
			if err == nil && resFloat > 0 {
				costPerSec = res.TotalTokensAmount / resFloat
			}
		}
		if costPerSec > 0 {
			res.CostPerSecond = strconv.FormatFloat(costPerSec, 'f', 6, 64)
			row[costPerSecondColIdx] = res.CostPerSecond
		}
		if res.PanQuAmount > 0 {
			row[panQuAmountColIdx] = strconv.FormatFloat(res.PanQuAmount, 'f', 6, 64)
		}
		row[amountDiffColIdx] = strconv.FormatFloat(res.AmountDiff, 'f', 6, 64)
		if strings.TrimSpace(res.DraftTaskId) != "" {
			row[hasReferenceVideoColIdx] = "是"
		} else {
			row[hasReferenceVideoColIdx] = "否"
		}

		rows[rIdx] = row
		updates++
	}

	if updates > 0 {
		// 覆盖写回 CSV 文件
		outFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println("创建CSV文件以供保存时失败:", err)
			return
		}
		defer outFile.Close()

		// 写入 BOM 头，这对于在 Excel 中打开包含中文的 CSV 文件很有帮助，防止乱码
		outFile.WriteString("\xef\xbb\xbf")

		writer := csv.NewWriter(outFile)
		if err := writer.WriteAll(rows); err != nil {
			fmt.Println("保存CSV文件失败:", err)
		} else {
			// 一定要 Flush 并且检查是否有错
			writer.Flush()
			if err := writer.Error(); err != nil {
				fmt.Println("写入CSV文件失败:", err)
			} else {
				fmt.Printf("\n执行完毕！成功更新了 %d 行数据，已保存至 %s\n", updates, fileName)
			}
		}
	} else {
		fmt.Println("\n执行完毕！没有需要更新的行数据。")
	}

	// === 开始同步到 pq_score_log_202603182351.csv ===
	fmt.Println("\n开始同步数据到 pq_score_log_202603182351.csv ...")
	normalizeKey := func(s string) string {
		s = strings.TrimPrefix(s, "\xef\xbb\xbf")
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "\"")
		return s
	}

	hsIdColIdx := -1
	for i, colCell := range rows[0] {
		h := normalizeKey(colCell)
		if hsIdColIdx == -1 && h == "id" {
			hsIdColIdx = i
		}
	}
	if hsIdColIdx == -1 {
		fmt.Println("源CSV未找到 id 列，无法按 id -> task_id 同步。")
		return
	}

	hsRowByID := make(map[string][]string)
	for rIdx := 1; rIdx < len(rows); rIdx++ {
		r := rows[rIdx]
		var idKey string
		if hsIdColIdx != -1 && hsIdColIdx < len(r) {
			idKey = normalizeKey(r[hsIdColIdx])
		}
		if idKey != "" {
			hsRowByID[idKey] = r
		}
	}

	pqFileName := "pq_score_log_202603182351.csv"
	pqFile, err := os.Open(pqFileName)
	if err != nil {
		fmt.Printf("无法打开CSV文件 %s: %v\n", pqFileName, err)
		return
	}
	pqReader := csv.NewReader(pqFile)
	pqReader.LazyQuotes = true
	pqReader.FieldsPerRecord = -1 // 允许每行的列数不一致
	pqRows, err := pqReader.ReadAll()
	pqFile.Close()

	if err != nil {
		fmt.Printf("读取 pq_score_log CSV 失败: %v\n", err)
		return
	}

	if len(pqRows) > 0 {
		pqHeader := pqRows[0]
		pqTaskIdIdx := -1
		for i, h := range pqHeader {
			if normalizeKey(h) == "task_id" {
				pqTaskIdIdx = i
				break
			}
		}

		appendCols := []string{"completion_tokens", "total_tokens", "status", "code", "message", "total_tokens_amount", "duration", "total_duration", "resolution", "cost_per_second", "pan_qu_amount", "amount_diff", "has_reference_video"}
		pqColMap := make(map[string]int)
		for _, colName := range appendCols {
			idx := -1
			for i, h := range pqHeader {
				if normalizeKey(h) == colName {
					idx = i
					break
				}
			}
			if idx == -1 {
				idx = len(pqHeader)
				pqHeader = append(pqHeader, colName)
			}
			pqColMap[colName] = idx
		}
		pqRows[0] = pqHeader

		pqUpdates := 0
		for i := 1; i < len(pqRows); i++ {
			row := pqRows[i]
			if pqTaskIdIdx != -1 && pqTaskIdIdx < len(row) {
				taskId := normalizeKey(row[pqTaskIdIdx])
				hsRow, ok := hsRowByID[taskId]
				if ok {
					maxIdx := 0
					for _, idx := range pqColMap {
						if idx > maxIdx {
							maxIdx = idx
						}
					}
					for len(row) <= maxIdx {
						row = append(row, "")
					}

					safeAssign := func(pqIdx, hsIdx int) {
						if hsIdx != -1 && hsIdx < len(hsRow) {
							row[pqIdx] = hsRow[hsIdx]
						}
					}

					safeAssign(pqColMap["completion_tokens"], completionColIdx)
					safeAssign(pqColMap["total_tokens"], totalTokensColIdx)
					safeAssign(pqColMap["status"], statusColIdx)
					safeAssign(pqColMap["code"], codeColIdx)
					safeAssign(pqColMap["message"], messageColIdx)
					safeAssign(pqColMap["total_tokens_amount"], totalTokensAmountIdx)
					safeAssign(pqColMap["duration"], durationColIdx)
					safeAssign(pqColMap["total_duration"], totalDurationColIdx)
					safeAssign(pqColMap["resolution"], resolutionColIdx)
					safeAssign(pqColMap["cost_per_second"], costPerSecondColIdx)
					safeAssign(pqColMap["pan_qu_amount"], panQuAmountColIdx)
					safeAssign(pqColMap["amount_diff"], amountDiffColIdx)
					safeAssign(pqColMap["has_reference_video"], hasReferenceVideoColIdx)

					pqRows[i] = row
					pqUpdates++
				}
			}
		}

		if pqUpdates > 0 {
			pqOutFile, err := os.Create(pqFileName)
			if err != nil {
				fmt.Println("创建 pq_score_log CSV 失败:", err)
			} else {
				defer pqOutFile.Close()
				pqOutFile.WriteString("\xef\xbb\xbf")
				pqWriter := csv.NewWriter(pqOutFile)
				if err := pqWriter.WriteAll(pqRows); err != nil {
					fmt.Println("保存 pq_score_log CSV 失败:", err)
				} else {
					pqWriter.Flush()
					if err := pqWriter.Error(); err != nil {
						fmt.Println("写入 pq_score_log CSV 失败:", err)
					} else {
						fmt.Printf("成功同步了 %d 行数据至 %s\n", pqUpdates, pqFileName)
					}
				}
			}
		} else {
			fmt.Println("没有行同步到 pq_score_log_202603182351.csv。")
		}
	}
}

func parseDurationValue(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func extractDurationFromTaskData(data map[string]interface{}) float64 {
	if durRaw, ok := data["duration"]; ok && durRaw != nil {
		if dur, ok := parseDurationValue(durRaw); ok {
			return dur
		}
	}
	if nested, ok := data["data"].(map[string]interface{}); ok {
		if durRaw, ok := nested["duration"]; ok && durRaw != nil {
			if dur, ok := parseDurationValue(durRaw); ok {
				return dur
			}
		}
	}
	return 0
}

func queryTaskDuration(taskId, apiKey string) (float64, error) {
	url := fmt.Sprintf("https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks/%s", taskId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP 状态码错误: %d, 返回内容: %s", resp.StatusCode, string(bodyBytes))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return 0, err
	}

	return extractDurationFromTaskData(data), nil
}

// queryArkTask 调用 Volcengine Ark 获取任务信息
func queryArkTask(taskId, apiKey string) (*ArkTaskResult, error) {
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

	res := &ArkTaskResult{}

	if s, ok := data["status"].(string); ok {
		res.Status = s
	}

	if draftTaskId, ok := data["draft_task_id"].(string); ok {
		res.DraftTaskId = draftTaskId
	}

	if errMap, ok := data["error"].(map[string]interface{}); ok {
		if c, ok := errMap["code"].(string); ok {
			res.ErrorCode = c
		}
		if m, ok := errMap["message"].(string); ok {
			res.ErrorMessage = m
		}
	}

	// 解析返回的 usage
	var usage map[string]interface{}
	if u, ok := data["usage"].(map[string]interface{}); ok {
		usage = u
	} else if d, ok := data["data"].(map[string]interface{}); ok {
		// 兼容如果返回被包裹在 data 这个字段内的情况
		if u, ok := d["usage"].(map[string]interface{}); ok {
			usage = u
		}
	}

	if usage != nil {
		res.CompletionTokens = usage["completion_tokens"]
		res.TotalTokens = usage["total_tokens"]
		var model = "text" //文本 或 图文

		if draft_task_id, ok := data["draft_task_id"]; ok && draft_task_id != "" {
			model = "draft_task" //样片信息
		}

		var totalTokens int64
		switch v := res.TotalTokens.(type) {
		case float64:
			totalTokens = int64(v)
		case string:
			totalTokens, _ = strconv.ParseInt(v, 10, 64)
		}

		if modelStr, ok := data["model"].(string); ok {
			res.TotalTokensAmount = ai.Calculate(modelStr, model, totalTokens)

		}
	}

	res.Duration = extractDurationFromTaskData(data)
	res.TotalDuration = res.Duration
	if resRaw, ok := data["resolution"]; ok && resRaw != nil {
		if v, ok := resRaw.(string); ok {
			res.Resolution = v
		} else {
			res.Resolution = fmt.Sprintf("%v", resRaw)
		}
	}

	// 兼容如果它们被包裹在 data 这个字段内的情况
	if data["data"] != nil {
		if d, ok := data["data"].(map[string]interface{}); ok {
			if resRaw, ok := d["resolution"]; ok && resRaw != nil && res.Resolution == "" {
				if v, ok := resRaw.(string); ok {
					res.Resolution = v
				} else {
					res.Resolution = fmt.Sprintf("%v", resRaw)
				}
			}
		}
	}
	if strings.TrimSpace(res.DraftTaskId) != "" {
		if draftDuration, err := queryTaskDuration(res.DraftTaskId, apiKey); err == nil && draftDuration > 0 {
			res.TotalDuration = res.Duration + draftDuration
		}
	}
	res.PanQuAmount = ai.PanQuModelPriceCalculate(data["model"].(string), res.Resolution, res.DraftTaskId, int64(res.Duration))
	res.AmountDiff = res.PanQuAmount - res.TotalTokensAmount
	return res, nil
}
