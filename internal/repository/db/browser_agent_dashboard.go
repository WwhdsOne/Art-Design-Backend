package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/tablename"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"time"
)

// countConversations 统计会话数量，userID 为 nil 时统计平台维度
func (r *BrowserAgentDB) countConversations(ctx context.Context, userID *int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table(tablename.BrowserAgentConversationTableName)
	if userID != nil {
		queryCond = queryCond.Where("created_by = ?", *userID)
	}
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

// countMessages 统计消息数量，userID 为 nil 时统计平台维度
func (r *BrowserAgentDB) countMessages(ctx context.Context, userID *int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table("browser_agent_message m")
	if userID != nil {
		queryCond = queryCond.
			Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
			Where("c.created_by = ?", *userID)
	}
	if !startTime.IsZero() {
		queryCond = queryCond.Where("m.created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("m.created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

// countActions 统计动作数量，userID 为 nil 时统计平台维度
func (r *BrowserAgentDB) countActions(ctx context.Context, userID *int64, status string, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table("browser_agent_action a")
	if userID != nil {
		queryCond = queryCond.
			Joins("JOIN browser_agent_message m ON a.message_id = m.id").
			Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
			Where("c.created_by = ?", *userID)
	}
	if status != "" {
		queryCond = queryCond.Where("a.status = ?", status)
	}
	if !startTime.IsZero() {
		queryCond = queryCond.Where("a.created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("a.created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

// =========================
// Dashboard - 用户维度统计
// =========================

func (r *BrowserAgentDB) CountConversationsByTimeRange(ctx context.Context, userID int64, startTime, endTime time.Time) (int64, error) {
	return r.countConversations(ctx, &userID, startTime, endTime)
}

func (r *BrowserAgentDB) CountMessagesByTimeRange(ctx context.Context, userID int64, startTime, endTime time.Time) (int64, error) {
	return r.countMessages(ctx, &userID, startTime, endTime)
}

func (r *BrowserAgentDB) CountActionsByTimeRange(ctx context.Context, userID int64, status string, startTime, endTime time.Time) (int64, error) {
	return r.countActions(ctx, &userID, status, startTime, endTime)
}

func (r *BrowserAgentDB) CountMessagesByDays(ctx context.Context, userID int64, days []time.Time) ([]int64, error) {
	result := make([]int64, len(days))
	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		count, err := r.CountMessagesByTimeRange(ctx, userID, dayStart, dayEnd)
		if err != nil {
			return nil, err
		}
		result[i] = count
	}
	return result, nil
}

func (r *BrowserAgentDB) CountMessagesByMonths(ctx context.Context, userID int64, year int) ([]int64, error) {
	result := make([]int64, 12)
	for month := 1; month <= 12; month++ {
		monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0)
		count, err := r.CountMessagesByTimeRange(ctx, userID, monthStart, monthEnd)
		if err != nil {
			return nil, err
		}
		result[month-1] = count
	}
	return result, nil
}

func (r *BrowserAgentDB) GetRecentConversations(ctx context.Context, userID int64, limit int) ([]*entity.BrowserAgentConversation, error) {
	var conversations []*entity.BrowserAgentConversation
	err := DB(ctx, r.db).Model(&entity.BrowserAgentConversation{}).
		Where("created_by = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&conversations).Error
	return conversations, errors.WrapDBError(err, "查询最近会话失败")
}

func (r *BrowserAgentDB) CountConversationsByDays(ctx context.Context, userID int64, days []time.Time) ([]int64, error) {
	result := make([]int64, len(days))
	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		var count int64
		if err := DB(ctx, r.db).Table(tablename.BrowserAgentConversationTableName).
			Where("created_by = ? AND created_at >= ? AND created_at < ?", userID, dayStart, dayEnd).
			Count(&count).Error; err != nil {
			return nil, err
		}
		result[i] = count
	}
	return result, nil
}

func (r *BrowserAgentDB) CountActionsSuccessRateByDays(
	ctx context.Context, userID int64, days []time.Time,
) (result []int64, totalSum int64, successSum int64, err error) {
	result = make([]int64, len(days))

	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		var total, success int64
		db := DB(ctx, r.db).
			Table(fmt.Sprintf("%s a", tablename.BrowserAgentActionTableName)).
			Joins("JOIN browser_agent_message m ON a.message_id = m.id").
			Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
			Where("c.created_by = ? AND a.created_at >= ? AND a.created_at < ?", userID, dayStart, dayEnd)

		if err := db.Count(&total).Error; err != nil {
			return nil, 0, 0, err
		}
		if err := db.Where("a.status = ?", "success").Count(&success).Error; err != nil {
			return nil, 0, 0, err
		}

		if total > 0 {
			result[i] = success * 100 / total
		}
		totalSum += total
		successSum += success
	}

	return result, totalSum, successSum, nil
}

// =========================
// Dashboard - 平台维度统计
// =========================

func (r *BrowserAgentDB) CountUsersByTimeRange(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table("user")
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

func (r *BrowserAgentDB) CountAllConversationsByTimeRange(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	return r.countConversations(ctx, nil, startTime, endTime)
}

func (r *BrowserAgentDB) CountAllMessagesByTimeRange(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	return r.countMessages(ctx, nil, startTime, endTime)
}

func (r *BrowserAgentDB) CountAllActionsByTimeRange(ctx context.Context, status string, startTime, endTime time.Time) (int64, error) {
	return r.countActions(ctx, nil, status, startTime, endTime)
}

func (r *BrowserAgentDB) CountTotalConversations(ctx context.Context) (int64, error) {
	var count int64
	err := DB(ctx, r.db).Table(tablename.BrowserAgentConversationTableName).Count(&count).Error
	return count, err
}

func (r *BrowserAgentDB) GetMessageStateStats(ctx context.Context) (map[string]int64, error) {
	var items []struct {
		State string
		Count int64
	}
	err := DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
		Select("state, COUNT(*) as count").
		Group("state").
		Scan(&items).Error
	if err != nil {
		return nil, errors.WrapDBError(err, "查询消息状态统计失败")
	}
	result := make(map[string]int64)
	for _, item := range items {
		result[item.State] = item.Count
	}
	return result, nil
}

func (r *BrowserAgentDB) CountAllMessagesByMonths(ctx context.Context, months []time.Time) (thisYear []int64, lastYear []int64, err error) {
	thisYear = make([]int64, len(months))
	lastYear = make([]int64, len(months))
	for i, month := range months {
		monthStart := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0)

		var count int64
		if err = DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&count).Error; err != nil {
			return
		}
		thisYear[i] = count

		lastYearStart := monthStart.AddDate(-1, 0, 0)
		lastYearEnd := monthEnd.AddDate(-1, 0, 0)
		if err = DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
			Where("created_at >= ? AND created_at < ?", lastYearStart, lastYearEnd).
			Count(&count).Error; err != nil {
			return nil, nil, err
		}
		lastYear[i] = count
	}
	return
}

func (r *BrowserAgentDB) CountAllMessagesByDays(ctx context.Context, days []time.Time) ([]int64, error) {
	result := make([]int64, len(days))
	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		var count int64
		if err := DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&count).Error; err != nil {
			return nil, err
		}
		result[i] = count
	}
	return result, nil
}

type UserRankingRow struct {
	UserID    int64
	Username  string
	TaskCount int64
}

func (r *BrowserAgentDB) GetUserTaskRanking(ctx context.Context, limit int) ([]UserRankingRow, error) {
	var results []UserRankingRow
	err := DB(ctx, r.db).
		Table(`"user" AS u`).
		Select(`u.id AS user_id, u.username, COUNT(m.id) AS task_count`).
		Joins(`JOIN browser_agent_conversation c ON c.created_by = u.id`).
		Joins(`JOIN browser_agent_message m ON m.conversation_id = c.id`).
		Group(`u.id, u.username`).
		Order(`task_count DESC`).
		Limit(limit).
		Scan(&results).Error
	return results, err
}

type HotTaskRow struct {
	Content string
	Count   int64
}

func (r *BrowserAgentDB) GetHotTaskContents(ctx context.Context, limit int) ([]HotTaskRow, error) {
	var results []HotTaskRow
	err := DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
		Select("content, COUNT(*) as count").
		Group("content").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

func (r *BrowserAgentDB) CountAllActionsByDays(ctx context.Context, days []time.Time) ([]int64, error) {
	result := make([]int64, len(days))
	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		var count int64
		if err := DB(ctx, r.db).Table(tablename.BrowserAgentActionTableName).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&count).Error; err != nil {
			return nil, err
		}
		result[i] = count
	}
	return result, nil
}

type DayCount struct {
	Day     time.Time
	Total   int64
	Success int64
}

func (r *BrowserAgentDB) CountAllByDays(ctx context.Context, days []time.Time) ([]DayCount, error) {
	start := days[0]
	end := days[len(days)-1].Add(24 * time.Hour)

	var rows []DayCount
	err := DB(ctx, r.db).
		Table(tablename.BrowserAgentActionTableName).
		Select(`
			DATE(created_at) as day,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success
		`).
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Scan(&rows).Error

	return rows, err
}

func (r *BrowserAgentDB) CountAllMessageSuccessRateByDays(ctx context.Context, days []time.Time) ([]DayCount, error) {
	start := days[0]
	end := days[len(days)-1].Add(24 * time.Hour)

	var rows []DayCount
	err := DB(ctx, r.db).
		Table(tablename.BrowserAgentMessageTableName).
		Select(`
			DATE(created_at) as day,
			COUNT(*) as total,
			SUM(CASE WHEN state = 'finished' THEN 1 ELSE 0 END) as success
		`).
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Scan(&rows).Error

	return rows, err
}

func (r *BrowserAgentDB) CountAllConversationsByDays(
	ctx context.Context,
	days []time.Time,
) ([]int64, error) {

	if len(days) == 0 {
		return nil, nil
	}

	start := time.Date(days[0].Year(), days[0].Month(), days[0].Day(), 0, 0, 0, 0, days[0].Location())
	end := time.Date(days[len(days)-1].Year(), days[len(days)-1].Month(), days[len(days)-1].Day(), 0, 0, 0, 0, days[len(days)-1].Location()).
		Add(24 * time.Hour)

	type dayCount struct {
		Day   time.Time
		Count int64
	}

	var rows []dayCount

	err := DB(ctx, r.db).
		Table(tablename.BrowserAgentConversationTableName).
		Select("DATE(created_at) as day, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int64, len(rows))
	for _, row := range rows {
		key := row.Day.Format("2006-01-02")
		countMap[key] = row.Count
	}

	result := make([]int64, len(days))
	for i, day := range days {
		key := day.Format("2006-01-02")
		result[i] = countMap[key]
	}

	return result, nil
}

func (r *BrowserAgentDB) CountAllMessagesByYearWithQuarters(ctx context.Context, year int) (monthly []int64, quarterly []int64, err error) {
	monthly = make([]int64, 12)
	quarterly = make([]int64, 4)

	for month := 1; month <= 12; month++ {
		monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0)
		var count int64
		if err = DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&count).Error; err != nil {
			return
		}
		monthly[month-1] = count
		quarterly[(month-1)/3] += count
	}
	return
}

type TaskClassificationRow struct {
	Content string
	Count   int64
}

func (r *BrowserAgentDB) GetTaskClassification(ctx context.Context) ([]TaskClassificationRow, error) {
	var results []TaskClassificationRow
	err := DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName).
		Select("content, COUNT(*) as count").
		Group("content").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error
	return results, err
}

type HotTaskDetailRow struct {
	Content      string
	Count        int64
	AvgExecTime  float64
	SuccessCount int64
	TotalActions int64
}

func (r *BrowserAgentDB) GetHotTasksWithDetails(
	ctx context.Context,
	limit int,
) ([]HotTaskDetailRow, error) {

	var results []HotTaskDetailRow

	err := DB(ctx, r.db).
		Table(tablename.BrowserAgentMessageTableName + " AS m").
		Select(`
			m.content,
			COUNT(DISTINCT m.id) AS count,
			COALESCE(AVG(a.execution_time), 0) AS avg_exec_time,
			SUM(CASE WHEN a.status = 'success' THEN 1 ELSE 0 END) AS success_count,
			COUNT(a.id) AS total_actions
		`).
		Joins("LEFT JOIN browser_agent_action a ON a.message_id = m.id").
		Group("m.content").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
