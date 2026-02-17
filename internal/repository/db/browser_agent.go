package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/constant/tablename"
	"Art-Design-Backend/pkg/errors"
	"context"
	"time"

	"gorm.io/gorm"
)

type BrowserAgentDB struct {
	db *gorm.DB
}

func NewBrowserAgentDB(db *gorm.DB) *BrowserAgentDB {
	return &BrowserAgentDB{db: db}
}

// =========================
// Conversation CRUD
// =========================

func (r *BrowserAgentDB) CreateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	if err := DB(ctx, r.db).Create(conv).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetConversationByID(ctx context.Context, id int64) (conv *entity.BrowserAgentConversation, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&conv).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体会话失败")
	}
	return
}

func (r *BrowserAgentDB) ListConversationsByUserID(ctx context.Context, userID int64, queryParam *query.BrowserAgentConversation) (conversations []*entity.BrowserAgentConversation, total int64, err error) {
	db := DB(ctx, r.db).Model(&entity.BrowserAgentConversation{}).Where("created_by = ?", userID)

	if queryParam.Title != "" {
		db = db.Where("title LIKE ?", "%"+queryParam.Title+"%")
	}
	if queryParam.State != "" {
		db = db.Where("state = ?", queryParam.State)
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapDBError(err, "统计浏览器智能体会话数量失败")
	}

	if err = db.Scopes(queryParam.Paginate()).Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, errors.WrapDBError(err, "查询浏览器智能体会话列表失败")
	}

	return
}

func (r *BrowserAgentDB) UpdateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	if err := DB(ctx, r.db).Save(conv).Error; err != nil {
		return errors.WrapDBError(err, "更新浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteConversation(ctx context.Context, id int64) error {
	if err := DB(ctx, r.db).Delete(&entity.BrowserAgentConversation{}, id).Error; err != nil {
		return errors.WrapDBError(err, "删除浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteMessagesByConversationID(ctx context.Context, conversationID int64) error {
	if err := DB(ctx, r.db).Where("conversation_id = ?", conversationID).Delete(&entity.BrowserAgentMessage{}).Error; err != nil {
		return errors.WrapDBError(err, "删除浏览器智能体消息失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteActionsByMessageIDList(ctx context.Context, messageIDList []int64) error {
	if err := DB(ctx, r.db).
		Where("message_id IN (?)", messageIDList).
		Delete(&entity.BrowserAgentAction{}).Error; err != nil {
		return errors.WrapDBError(err, "批量删除浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) UpdateConversationState(ctx context.Context, id int64, state string) error {
	if err := DB(ctx, r.db).Model(&entity.BrowserAgentConversation{}).
		Where("id = ?", id).Update("state", state).Error; err != nil {
		return errors.WrapDBError(err, "更新会话状态失败")
	}
	return nil
}

// =========================
// Message CRUD
// =========================

func (r *BrowserAgentDB) CreateMessage(ctx context.Context, msg *entity.BrowserAgentMessage) error {
	if err := DB(ctx, r.db).Create(msg).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体消息失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetMessageByID(ctx context.Context, id int64) (msg *entity.BrowserAgentMessage, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&msg).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewDBError("消息不存在")
		}
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息失败")
	}
	return
}

func (r *BrowserAgentDB) ListMessagesByConversationID(ctx context.Context, conversationID int64) (messages []*entity.BrowserAgentMessage, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentMessage{}).
		Where("conversation_id = ?", conversationID).
		Find(&messages).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息列表失败")
	}
	return
}

func (r *BrowserAgentDB) ListMessagesIDListByConversationID(ctx context.Context, conversationID int64) (messageIDList []int64, err error) {
	if err = DB(ctx, r.db).
		Model(&entity.BrowserAgentMessage{}).
		Where("conversation_id = ?", conversationID).
		Pluck("id", &messageIDList).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息ID列表失败")
	}
	return
}

func (r *BrowserAgentDB) UpdateMessageState(ctx context.Context, id int64, state string) error {
	if err := DB(ctx, r.db).Model(&entity.BrowserAgentMessage{}).
		Where("id = ?", id).Update("state", state).Error; err != nil {
		return errors.WrapDBError(err, "更新任务状态失败")
	}
	return nil
}

// =========================
// Action CRUD
// =========================

func (r *BrowserAgentDB) CreateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := DB(ctx, r.db).Create(action).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) CreateActions(ctx context.Context, actions []*entity.BrowserAgentAction) error {
	if len(actions) == 0 {
		return nil
	}
	if err := DB(ctx, r.db).Create(&actions).Error; err != nil {
		return errors.WrapDBError(err, "批量创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetActionByID(ctx context.Context, id int64) (action *entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&action).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewDBError("操作不存在")
		}
		return nil, errors.WrapDBError(err, "查询浏览器智能体操作失败")
	}
	return
}

func (r *BrowserAgentDB) ListActionsByMessageID(ctx context.Context, messageID int64) (actions []*entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ?", messageID).
		Find(&actions).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体操作列表失败")
	}
	return
}

func (r *BrowserAgentDB) UpdateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := DB(ctx, r.db).Save(action).Error; err != nil {
		return errors.WrapDBError(err, "更新浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) UpdateActionStatus(ctx context.Context, id int64, status string, errMsg *string, execTime *int) error {
	updates := map[string]interface{}{"status": status}
	if errMsg != nil {
		updates["error_message"] = *errMsg
	}
	if execTime != nil {
		updates["execution_time"] = *execTime
	}

	if err := DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.WrapDBError(err, "更新操作状态失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetPendingActionsByMessageID(ctx context.Context, messageID int64) (actions []*entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ? AND status = ?", messageID, entity.ActionStatusPending).
		Order("sequence ASC").
		Find(&actions).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询待执行操作失败")
	}
	return
}

// =========================
// Dashboard - 用户维度统计
// =========================

func (r *BrowserAgentDB) CountConversationsByTimeRange(ctx context.Context, userID int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).
		Table(tablename.BrowserAgentConversationTableName).
		Where("created_by = ?", userID)
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

func (r *BrowserAgentDB) CountMessagesByTimeRange(ctx context.Context, userID int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table("browser_agent_message m").
		Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
		Where("c.created_by = ?", userID)
	if !startTime.IsZero() {
		queryCond = queryCond.Where("m.created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("m.created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

func (r *BrowserAgentDB) CountActionsByTimeRange(ctx context.Context, userID int64, status string, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table("browser_agent_action a").
		Joins("JOIN browser_agent_message m ON a.message_id = m.id").
		Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
		Where("c.created_by = ?", userID)
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
	var count int64
	queryCond := DB(ctx, r.db).Table(tablename.BrowserAgentConversationTableName)
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

func (r *BrowserAgentDB) CountAllMessagesByTimeRange(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table(tablename.BrowserAgentMessageTableName)
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
}

func (r *BrowserAgentDB) CountAllActionsByTimeRange(ctx context.Context, status string, startTime, endTime time.Time) (int64, error) {
	var count int64
	queryCond := DB(ctx, r.db).Table(tablename.BrowserAgentActionTableName)
	if status != "" {
		queryCond = queryCond.Where("status = ?", status)
	}
	if !startTime.IsZero() {
		queryCond = queryCond.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		queryCond = queryCond.Where("created_at < ?", endTime)
	}
	return count, queryCond.Count(&count).Error
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

	// 计算整体时间范围
	start := time.Date(days[0].Year(), days[0].Month(), days[0].Day(), 0, 0, 0, 0, days[0].Location())
	end := time.Date(days[len(days)-1].Year(), days[len(days)-1].Month(), days[len(days)-1].Day(), 0, 0, 0, 0, days[len(days)-1].Location()).
		Add(24 * time.Hour)

	// 查询结果结构
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

	// 构建 map 方便按天填充
	countMap := make(map[string]int64, len(rows))
	for _, row := range rows {
		key := row.Day.Format("2006-01-02")
		countMap[key] = row.Count
	}

	// 按传入 days 顺序填充结果（保证顺序稳定）
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

// GetHotTasksWithDetails 查询热门任务及执行详情
//
// 返回结果包含每条任务内容、执行次数、平均执行时间、成功次数以及总动作数
// limit: 返回条数上限
func (r *BrowserAgentDB) GetHotTasksWithDetails(
	ctx context.Context,
	limit int,
) ([]HotTaskDetailRow, error) {

	var results []HotTaskDetailRow

	// 构建 SQL 查询
	// 说明：
	//  - m.content: 任务内容
	//  - COUNT(DISTINCT m.id) as count: 该任务出现次数（去重 message id）
	//  - COALESCE(AVG(a.execution_time), 0) as avg_exec_time: 平均执行时间，如果没有动作则为 0
	//  - SUM(CASE WHEN a.status = 'success' THEN 1 ELSE 0 END) as success_count: 成功动作数
	//  - COUNT(a.id) as total_actions: 总动作数（包括失败动作）
	err := DB(ctx, r.db).
		Table("browser_agent_message m").
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

func (r *BrowserAgentDB) CountActionsSuccessRateByDays(ctx context.Context, userID int64, days []time.Time) ([]int64, error) {
	result := make([]int64, len(days))
	for i, day := range days {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		var total, success int64
		if err := DB(ctx, r.db).Table("browser_agent_action a").
			Joins("JOIN browser_agent_message m ON a.message_id = m.id").
			Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
			Where("c.created_by = ? AND a.created_at >= ? AND a.created_at < ?", userID, dayStart, dayEnd).
			Count(&total).Error; err != nil {
			return nil, err
		}
		if err := DB(ctx, r.db).Table("browser_agent_action a").
			Joins("JOIN browser_agent_message m ON a.message_id = m.id").
			Joins("JOIN browser_agent_conversation c ON m.conversation_id = c.id").
			Where("c.created_by = ? AND a.created_at >= ? AND a.created_at < ? AND a.status = ?", userID, dayStart, dayEnd, "success").
			Count(&success).Error; err != nil {
			return nil, err
		}
		if total > 0 {
			result[i] = success * 100 / total
		}
	}
	return result, nil
}
