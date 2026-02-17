package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"context"
	"fmt"
	"strings"
	"time"
)

type BrowserAgentDashboardService struct {
	BrowserAgentRepo *repository.BrowserAgentRepo
}

// =========================
// 8. Admin Dashboard APIs
// =========================

func (s *BrowserAgentDashboardService) GetAdminSummary(ctx context.Context) (*response.AdminSummaryResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	todayTasks, _ := s.BrowserAgentRepo.CountAllMessagesByTimeRange(ctx, todayStart, time.Time{})
	yesterdayTasks, _ := s.BrowserAgentRepo.CountAllMessagesByTimeRange(ctx, yesterdayStart, todayStart)

	totalActions, _ := s.BrowserAgentRepo.CountAllActionsByTimeRange(ctx, "", time.Time{}, time.Time{})
	successActions, _ := s.BrowserAgentRepo.CountAllActionsByTimeRange(ctx, "success", time.Time{}, time.Time{})
	successRate := 0
	if totalActions > 0 {
		successRate = int(successActions * 100 / totalActions)
	}

	lastWeekTotal, _ := s.BrowserAgentRepo.CountAllActionsByTimeRange(ctx, "", todayStart.AddDate(0, 0, -7), todayStart)
	lastWeekSuccess, _ := s.BrowserAgentRepo.CountAllActionsByTimeRange(ctx, "success", todayStart.AddDate(0, 0, -7), todayStart)
	lastWeekRate := 0
	if lastWeekTotal > 0 {
		lastWeekRate = int(lastWeekSuccess * 100 / lastWeekTotal)
	}
	successGrowth := calcChange(int64(successRate), int64(lastWeekRate))

	return &response.AdminSummaryResponse{
		TodayTasks:    todayTasks,
		TodayGrowth:   calcChange(todayTasks, yesterdayTasks),
		SuccessRate:   successRate,
		SuccessGrowth: successGrowth,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminWeeklyTaskVolume(ctx context.Context) (*response.VolumeDataResponse, error) {
	now := time.Now()
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = now.AddDate(0, 0, -6+i)
	}

	chartData, _ := s.BrowserAgentRepo.CountAllMessagesByDays(ctx, days)
	thisWeekVolume := int64(0)
	for _, v := range chartData {
		thisWeekVolume += v
	}

	lastWeekDays := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		lastWeekDays[i] = now.AddDate(0, 0, -13+i)
	}
	lastWeekData, _ := s.BrowserAgentRepo.CountAllMessagesByDays(ctx, lastWeekDays)
	lastWeekVolume := int64(0)
	for _, v := range lastWeekData {
		lastWeekVolume += v
	}

	return &response.VolumeDataResponse{
		Volume:    thisWeekVolume,
		Growth:    calcChange(thisWeekVolume, lastWeekVolume),
		ChartData: chartData,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminWeeklyTaskSuccessRate(ctx context.Context) (*response.RateDataResponse, error) {
	return s.getWeeklySuccessRate(ctx, func(ctx context.Context, days []time.Time) ([]dayCountRow, error) {
		rows, err := s.BrowserAgentRepo.CountAllMessageSuccessRateByDays(ctx, days)
		if err != nil {
			return nil, err
		}
		result := make([]dayCountRow, len(rows))
		for i, r := range rows {
			result[i] = dayCountRow{Day: r.Day, Total: r.Total, Success: r.Success}
		}
		return result, nil
	})
}

func (s *BrowserAgentDashboardService) GetAdminTotalTaskVolume(
	ctx context.Context,
) (*response.TotalTaskVolumeResponse, error) {

	stateStats, err := s.BrowserAgentRepo.GetMessageStateStats(ctx)
	if err != nil {
		return nil, err
	}

	// 状态文案映射
	stateNameMap := make(map[string]string)
	stateNameMap[entity.MessageStateRunning] = "进行中"
	stateNameMap[entity.MessageStateFinished] = "已完成"
	stateNameMap[entity.MessageStateError] = "已失败"

	// 哪些状态计入 Total
	countInTotal := map[string]struct{}{
		"running":  {},
		"finished": {},
	}

	distribution := make([]response.DistributionItem, 0, len(stateStats))

	var total, totalAll int64

	for state, count := range stateStats {
		name, ok := stateNameMap[state]
		if !ok {
			name = state
		}

		distribution = append(distribution, response.DistributionItem{
			Name:  name,
			Value: count,
		})

		totalAll += count
		if _, ok := countInTotal[state]; ok {
			total += count
		}
	}

	return &response.TotalTaskVolumeResponse{
		Total:        total,
		TotalAll:     totalAll,
		Distribution: distribution,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminTaskClassification(ctx context.Context) (*response.TaskClassificationResponse, error) {
	tasks, _ := s.BrowserAgentRepo.GetTaskClassification(ctx)

	classificationMap := make(map[string]int64)
	var total int64
	for _, task := range tasks {
		category := classifyTask(task.Content)
		classificationMap[category] += task.Count
		total += task.Count
	}

	distribution := make([]response.ClassificationItem, 0)
	for category, count := range classificationMap {
		percentage := 0
		if total > 0 {
			percentage = int(count * 100 / total)
		}
		distribution = append(distribution, response.ClassificationItem{
			Value: percentage,
			Name:  category,
		})
	}

	return &response.TaskClassificationResponse{
		Distribution: distribution,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminWeeklyOperationVolume(ctx context.Context) (*response.VolumeDataResponse, error) {
	now := time.Now()
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = now.AddDate(0, 0, -6+i)
	}

	chartData, _ := s.BrowserAgentRepo.CountAllActionsByDays(ctx, days)
	thisWeekVolume := int64(0)
	for _, v := range chartData {
		thisWeekVolume += v
	}

	lastWeekDays := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		lastWeekDays[i] = now.AddDate(0, 0, -13+i)
	}
	lastWeekData, _ := s.BrowserAgentRepo.CountAllActionsByDays(ctx, lastWeekDays)
	lastWeekVolume := int64(0)
	for _, v := range lastWeekData {
		lastWeekVolume += v
	}

	return &response.VolumeDataResponse{
		Volume:    thisWeekVolume,
		Growth:    calcChange(thisWeekVolume, lastWeekVolume),
		ChartData: chartData,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminWeeklyOperationSuccessRate(ctx context.Context) (*response.RateDataResponse, error) {
	return s.getWeeklySuccessRate(ctx, func(ctx context.Context, days []time.Time) ([]dayCountRow, error) {
		rows, err := s.BrowserAgentRepo.CountAllByDays(ctx, days)
		if err != nil {
			return nil, err
		}
		result := make([]dayCountRow, len(rows))
		for i, r := range rows {
			result[i] = dayCountRow{Day: r.Day, Total: r.Total, Success: r.Success}
		}
		return result, nil
	})
}

func (s *BrowserAgentDashboardService) GetAdminActiveSessions(
	ctx context.Context,
) (*response.ActiveSessionsResponse, error) {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// ===== 本周 7 天 =====
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = today.AddDate(0, 0, -6+i)
	}

	thisWeekData, err := s.BrowserAgentRepo.CountAllConversationsByDays(ctx, days)
	if err != nil {
		return nil, err
	}

	// 本周总会话数
	var currentCount int64
	for _, c := range thisWeekData {
		currentCount += c
	}

	// ===== 上周 7 天 =====
	lastWeekDays := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		lastWeekDays[i] = today.AddDate(0, 0, -13+i)
	}

	lastWeekData, err := s.BrowserAgentRepo.CountAllConversationsByDays(ctx, lastWeekDays)
	if err != nil {
		return nil, err
	}

	var lastWeekCount int64
	for _, c := range lastWeekData {
		lastWeekCount += c
	}

	return &response.ActiveSessionsResponse{
		Count:     currentCount,                            // 本周总会话数
		Growth:    calcChange(currentCount, lastWeekCount), // 周环比
		ChartData: thisWeekData,                            // 每天会话数
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminAnnualTaskStats(ctx context.Context, year int) (*response.AnnualTaskStatsResponse, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	monthlyData, quarterData, err := s.BrowserAgentRepo.CountAllMessagesByYearWithQuarters(ctx, year)
	if err != nil {
		return nil, err
	}

	return &response.AnnualTaskStatsResponse{
		Year:        year,
		MonthlyData: monthlyData,
		QuarterData: quarterData,
	}, nil
}

func (s *BrowserAgentDashboardService) GetAdminHotTaskList(ctx context.Context, limit int) ([]response.HotTaskItemResponse, error) {
	if limit == 0 || limit > 10 {
		limit = 6
	}

	tasks, _ := s.BrowserAgentRepo.GetHotTasksWithDetails(ctx, limit)

	colors := []string{"primary", "success", "warning", "error", "info", "secondary"}
	result := make([]response.HotTaskItemResponse, len(tasks))
	for i, task := range tasks {
		successRate := 0
		if task.TotalActions > 0 {
			successRate = int(task.SuccessCount * 100 / task.TotalActions)
		}
		result[i] = response.HotTaskItemResponse{
			Name:        truncateString(task.Content, 30),
			Category:    classifyTask(task.Content),
			AvgTime:     int(task.AvgExecTime),
			ExecCount:   task.Count,
			SuccessRate: successRate,
			Color:       colors[i%len(colors)],
		}
	}

	return result, nil
}

// =========================
// 9. User Dashboard APIs
// =========================

func (s *BrowserAgentDashboardService) GetUserSummary(ctx context.Context, userID int64) (*response.UserSummaryResponse, error) {
	thisWeekStart, lastWeekStart, lastWeekEnd := getWeekTimeRanges()

	thisWeekSessions, _ := s.BrowserAgentRepo.CountConversationsByTimeRange(ctx, userID, thisWeekStart, time.Time{})
	lastWeekSessions, _ := s.BrowserAgentRepo.CountConversationsByTimeRange(ctx, userID, lastWeekStart, lastWeekEnd)

	totalActions, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "", time.Time{}, time.Time{})
	successActions, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "success", time.Time{}, time.Time{})
	successRate := 0
	if totalActions > 0 {
		successRate = int(successActions * 100 / totalActions)
	}

	lastWeekActions, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "", lastWeekStart, lastWeekEnd)
	lastWeekSuccess, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "success", lastWeekStart, lastWeekEnd)
	lastWeekRate := 0
	if lastWeekActions > 0 {
		lastWeekRate = int(lastWeekSuccess * 100 / lastWeekActions)
	}

	return &response.UserSummaryResponse{
		SessionCount:  thisWeekSessions,
		SessionGrowth: calcChange(thisWeekSessions, lastWeekSessions),
		SuccessRate:   successRate,
		SuccessGrowth: calcChange(int64(successRate), int64(lastWeekRate)),
	}, nil
}

func (s *BrowserAgentDashboardService) GetUserWeeklyTaskVolume(ctx context.Context, userID int64) (*response.VolumeDataResponse, error) {
	now := time.Now()
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = now.AddDate(0, 0, -6+i)
	}

	chartData, _ := s.BrowserAgentRepo.CountMessagesByDays(ctx, userID, days)
	thisWeekVolume := int64(0)
	for _, v := range chartData {
		thisWeekVolume += v
	}

	lastWeekDays := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		lastWeekDays[i] = now.AddDate(0, 0, -13+i)
	}
	lastWeekData, _ := s.BrowserAgentRepo.CountMessagesByDays(ctx, userID, lastWeekDays)
	lastWeekVolume := int64(0)
	for _, v := range lastWeekData {
		lastWeekVolume += v
	}

	return &response.VolumeDataResponse{
		Volume:    thisWeekVolume,
		Growth:    calcChange(thisWeekVolume, lastWeekVolume),
		ChartData: chartData,
	}, nil
}

func (s *BrowserAgentDashboardService) GetUserWeeklyTaskSuccessRate(ctx context.Context, userID int64) (*response.RateDataResponse, error) {
	now := time.Now()
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = now.AddDate(0, 0, -6+i)
	}

	chartData, _ := s.BrowserAgentRepo.CountActionsSuccessRateByDays(ctx, userID, days)
	currentRate := int64(0)
	if len(chartData) > 0 {
		currentRate = chartData[len(chartData)-1]
	}

	lastWeekDays := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		lastWeekDays[i] = now.AddDate(0, 0, -13+i)
	}
	lastWeekData, _ := s.BrowserAgentRepo.CountActionsSuccessRateByDays(ctx, userID, lastWeekDays)
	lastWeekRate := int64(0)
	if len(lastWeekData) > 0 {
		lastWeekRate = lastWeekData[len(lastWeekData)-1]
	}

	return &response.RateDataResponse{
		Rate:      int(currentRate),
		Growth:    calcChange(currentRate, lastWeekRate),
		ChartData: chartData,
	}, nil
}

func (s *BrowserAgentDashboardService) GetUserTaskOverview(ctx context.Context, userID int64) (*response.UserTaskOverviewResponse, error) {
	sessionCount, _ := s.BrowserAgentRepo.CountConversationsByTimeRange(ctx, userID, time.Time{}, time.Time{})
	taskCount, _ := s.BrowserAgentRepo.CountMessagesByTimeRange(ctx, userID, time.Time{}, time.Time{})

	totalActions, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "", time.Time{}, time.Time{})
	successActions, _ := s.BrowserAgentRepo.CountActionsByTimeRange(ctx, userID, "success", time.Time{}, time.Time{})
	successRate := 0
	if totalActions > 0 {
		successRate = int(successActions * 100 / totalActions)
	}

	_, lastWeekStart, lastWeekEnd := getWeekTimeRanges()
	lastWeekTasks, _ := s.BrowserAgentRepo.CountMessagesByTimeRange(ctx, userID, lastWeekStart, lastWeekEnd)

	now := time.Now()
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = now.AddDate(0, 0, -6+i)
	}
	chartData, _ := s.BrowserAgentRepo.CountMessagesByDays(ctx, userID, days)

	return &response.UserTaskOverviewResponse{
		SessionCount: sessionCount,
		TaskCount:    taskCount,
		SuccessRate:  successRate,
		WeekGrowth:   calcChange(taskCount, lastWeekTasks),
		ChartData: response.UserTaskChart{
			XAxis: []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
			Data:  chartData,
		},
	}, nil
}

func (s *BrowserAgentDashboardService) GetUserTaskTrend(ctx context.Context, userID int64, year int) (*response.UserTaskTrendResponse, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	monthlyData, err := s.BrowserAgentRepo.CountMessagesByMonths(ctx, userID, year)
	if err != nil {
		return nil, err
	}

	lastYearData, _ := s.BrowserAgentRepo.CountMessagesByMonths(ctx, userID, year-1)

	var thisYearTotal, lastYearTotal int64
	for _, v := range monthlyData {
		thisYearTotal += v
	}
	for _, v := range lastYearData {
		lastYearTotal += v
	}

	return &response.UserTaskTrendResponse{
		Year:        year,
		Growth:      calcChange(thisYearTotal, lastYearTotal),
		MonthlyData: monthlyData,
	}, nil
}

// =========================
// 辅助函数
// =========================

type dayCountRow struct {
	Day     time.Time
	Total   int64
	Success int64
}

type dayCountFetcher func(ctx context.Context, days []time.Time) ([]dayCountRow, error)

func makeWeekDays(today time.Time, startOffset int) []time.Time {
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = today.AddDate(0, 0, startOffset+i)
	}
	return days
}

func calcWeekRate(rows []dayCountRow, days []time.Time) (rate int64, chart []int64) {
	rowMap := make(map[string]dayCountRow)
	for _, r := range rows {
		rowMap[r.Day.Format("2006-01-02")] = r
	}

	chart = make([]int64, 7)
	var total, success int64

	for i, d := range days {
		if r, ok := rowMap[d.Format("2006-01-02")]; ok {
			total += r.Total
			success += r.Success
			if r.Total > 0 {
				chart[i] = r.Success * 100 / r.Total
			}
		}
	}

	if total > 0 {
		rate = success * 100 / total
	}
	return
}

func (s *BrowserAgentDashboardService) getWeeklySuccessRate(
	ctx context.Context,
	fetcher dayCountFetcher,
) (*response.RateDataResponse, error) {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())

	thisWeekDays := makeWeekDays(today, -6)
	lastWeekDays := makeWeekDays(today, -13)

	thisWeekRows, err := fetcher(ctx, thisWeekDays)
	if err != nil {
		return nil, err
	}

	lastWeekRows, err := fetcher(ctx, lastWeekDays)
	if err != nil {
		return nil, err
	}

	thisWeekRate, thisWeekChart := calcWeekRate(thisWeekRows, thisWeekDays)
	lastWeekRate, _ := calcWeekRate(lastWeekRows, lastWeekDays)

	return &response.RateDataResponse{
		Rate:      int(thisWeekRate),
		Growth:    calcChange(thisWeekRate, lastWeekRate),
		ChartData: thisWeekChart,
	}, nil
}

func getWeekTimeRanges() (thisWeekStart, lastWeekStart, lastWeekEnd time.Time) {
	now := time.Now()
	thisWeekStart = now.AddDate(0, 0, -int(now.Weekday())-7)
	lastWeekStart = thisWeekStart.AddDate(0, 0, -7)
	lastWeekEnd = thisWeekStart
	return
}

// calcChange 计算本周相比上周的百分比变化
// 返回格式示例：+20%、-15%、+100%
//
// 计算公式：
//
//	(thisWeek - lastWeek) / lastWeek * 100
//
// 特殊情况处理：
//   - 如果 lastWeek == 0 且 thisWeek > 0，认为增长 100%
//   - 如果 lastWeek == 0 且 thisWeek == 0，返回 +0%
func calcChange(thisWeek, lastWeek int64) string {
	// 如果上周为 0，需要特殊处理（避免除 0）
	if lastWeek == 0 {
		// 上周为 0，本周有数据，视为 100% 增长
		if thisWeek > 0 {
			return "+100%"
		}
		// 两周都为 0
		return "+0%"
	}

	// 计算变化百分比
	diff := float64(thisWeek-lastWeek) / float64(lastWeek) * 100

	// 正数加 "+" 号
	if diff >= 0 {
		return fmt.Sprintf("+%.0f%%", diff)
	}

	// 负数自动带 "-"
	return fmt.Sprintf("%.0f%%", diff)
}

func classifyTask(content string) string {
	// todo 添加更多分类
	keywords := map[string][]string{
		"电商平台":  {"淘宝", "京东", "拼多多", "电商"},
		"办公自动化": {"表单", "填写", "提交"},
		"信息检索":  {"搜索", "查询", "检索"},
	}
	for category, words := range keywords {
		for _, word := range words {
			if strings.Contains(content, word) {
				return category
			}
		}
	}
	return "其他任务"
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
