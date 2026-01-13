package small_program

import (
	"fmt"
	"log"
	"time"

	"program-collection/models"

	"gorm.io/gorm/clause"
)

type QueryResult struct {
	TotalCount int `gorm:"column:total_count"`
	IsToday    int `gorm:"column:is_today"`
}

// Consultant 研究顾问的 Weight | Value_factor信息
type WeightValueFactor struct {
	ID       int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID   string    `gorm:"column:user_id;size:50;not null;index:idx_user_id" json:"userId"`
	StatDate time.Time `gorm:"column:stat_date;not null;index:idx_stat_date" json:"statDate"`

	// 因子相关（当天数据）
	WeightFactor float64 `gorm:"column:weight_factor;type:decimal(5,2)" json:"weightFactor"`
	ValueFactor  float64 `gorm:"column:value_factor;type:decimal(5,2)" json:"valueFactor"`

	// 变化量
	WeightFactorChange float64 `gorm:"column:weight_factor_change;type:decimal(6,3)" json:"weightFactorChange"`
	ValueFactorChange  float64 `gorm:"column:value_factor_change;type:decimal(6,3)" json:"valueFactorChange"`

	// 变化率
	WeightFactorChangeRate float64 `gorm:"column:weight_factor_change_rate;type:decimal(8,4)" json:"weightFactorChangeRate"`
	ValueFactorChangeRate  float64 `gorm:"column:value_factor_change_rate;type:decimal(8,4)" json:"valueFactorChangeRate"`

	// 其他数据
	DataFieldsUsed                int        `gorm:"column:data_fields_used;default:0" json:"dataFieldsUsed"`
	SubmissionsCount              int        `gorm:"column:submissions_count;default:0" json:"submissionsCount"`
	SuperAlphaSubmissionsCount    int        `gorm:"column:super_alpha_submissions_count;default:0" json:"superAlphaSubmissionsCount"`
	MeanProdCorrelation           float64    `gorm:"column:mean_prod_correlation;type:decimal(5,4)" json:"meanProdCorrelation"`
	MeanSelfCorrelation           float64    `gorm:"column:mean_self_correlation;type:decimal(5,4)" json:"meanSelfCorrelation"`
	SuperAlphaMeanProdCorrelation float64    `gorm:"column:super_alpha_mean_prod_correlation;type:decimal(5,4)" json:"superAlphaMeanProdCorrelation"`
	SuperAlphaMeanSelfCorrelation float64    `gorm:"column:super_alpha_mean_self_correlation;type:decimal(5,4)" json:"superAlphaMeanSelfCorrelation"`
	University                    *string    `gorm:"column:university;size:255" json:"university"`
	Country                       string     `gorm:"column:country;size:100" json:"country"`
	DateStarted                   *time.Time `gorm:"column:date_started" json:"dateStarted"`

	// 时间字段
	CreateTime *time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// TableName 指定表名
func (WeightValueFactor) TableName() string {
	return "weight_value_factor"
}

// --------------------------------------- 保存研究顾问 wf 和 vf 变化 -----------------------------------------
func SaveWeightValueFactor(config models.Config, token string) error {
	fmt.Println("\n====================== 执行Weight和Value_factor更新 ======================")

	// 获取研究顾问Consultant的wf和vf数据
	resp, err := FetchConsultant(config, token)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("consultant response is nil")
	}

	// 连接数据库
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	var result QueryResult
	err = db.Raw(`
        SELECT 
            (SELECT COUNT(*) FROM worldquant.active_alpha_list) as total_count,
            MAX(date_submitted) as max_date,
            IF(DATE(MAX(date_submitted)) = CURDATE(), 1, 0) as is_today
        FROM worldquant.active_alpha_list
        WHERE type = 'SUPER'
    `).Scan(&result).Error

	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	// 获取当前日期
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 解析开始日期
	var dateStarted *time.Time
	if resp.DateStarted != "" {
		parsedDate, err := time.Parse("2006-01-02", resp.DateStarted)
		if err == nil {
			dateStarted = &parsedDate
		}
	}

	// 1. 首先检查数据库中是否有该用户的任何数据
	var totalCount int64
	err = db.Model(&WeightValueFactor{}).Where("user_id = ?", resp.Leaderboard.User).Count(&totalCount).Error
	if err != nil {
		return fmt.Errorf("查询用户数据总数失败: %v", err)
	}

	// 如果是首次加入数据
	if totalCount == 0 {
		fmt.Printf("老辈子些，这是你第一次存Weight_factor和Value_factor数据哦。")

		// 构建每日统计数据（首次添加，没有变化量）
		dailyStat := WeightValueFactor{
			UserID:   resp.Leaderboard.User,
			StatDate: currentDate,

			WeightFactor: resp.Leaderboard.WeightFactor,
			ValueFactor:  resp.Leaderboard.ValueFactor,

			WeightFactorChange:     0,
			ValueFactorChange:      0,
			WeightFactorChangeRate: 0,
			ValueFactorChangeRate:  0,

			DataFieldsUsed:                resp.Leaderboard.DataFieldsUsed,
			SubmissionsCount:              result.TotalCount,
			SuperAlphaSubmissionsCount:    result.IsToday,
			MeanProdCorrelation:           resp.Leaderboard.MeanProdCorrelation,
			MeanSelfCorrelation:           resp.Leaderboard.MeanSelfCorrelation,
			SuperAlphaMeanProdCorrelation: resp.Leaderboard.SuperAlphaMeanProdCorrelation,
			SuperAlphaMeanSelfCorrelation: resp.Leaderboard.SuperAlphaMeanSelfCorrelation,
			University:                    resp.Leaderboard.University,
			Country:                       resp.Leaderboard.Country,
			DateStarted:                   dateStarted,
		}

		// 保存数据
		return db.Create(&dailyStat).Error
	}

	// 2. 获取数据库中该用户的最新数据（按CreateTime排序）
	var latestStat WeightValueFactor
	err = db.Where("user_id = ?", resp.Leaderboard.User).Order("create_time DESC").First(&latestStat).Error
	if err != nil {
		return fmt.Errorf("获取最新数据失败: %v", err)
	}

	latestCreateTime := latestStat.CreateTime

	// 获取最新数据的日期（按StatDate，用于显示）
	latestDate := latestStat.StatDate
	currentDateOnly := currentDate

	// 5. 检查今天是否已经记录过数据
	// 如果最新数据的日期是今天，说明今天已经记录过了
	var hasPreviousData bool
	var previousWeightFactor, previousValueFactor float64

	// 判断今天是否已经记录过
	if latestDate.Equal(currentDateOnly) {
		fmt.Printf("老辈子些，今天已经记录过数据了哦(记录时间%s)\n", latestCreateTime.Format("15:04:05"))
		previousWeightFactor = latestStat.WeightFactor
		previousValueFactor = latestStat.ValueFactor
		hasPreviousData = true
	} else {
		// 使用最新的历史数据作为基准
		previousWeightFactor = latestStat.WeightFactor
		previousValueFactor = latestStat.ValueFactor
		hasPreviousData = true
	}

	// 6. 计算变化量和变化率
	weightFactorChange := 0.0
	valueFactorChange := 0.0
	weightFactorChangeRate := 0.0
	valueFactorChangeRate := 0.0

	if hasPreviousData {
		// 计算变化量
		weightFactorChange = resp.Leaderboard.WeightFactor - previousWeightFactor
		valueFactorChange = resp.Leaderboard.ValueFactor - previousValueFactor

		// 计算变化率（避免除零）
		if previousWeightFactor != 0 {
			weightFactorChangeRate = weightFactorChange / previousWeightFactor
		}
		if previousValueFactor != 0 {
			valueFactorChangeRate = valueFactorChange / previousValueFactor
		}

		// 输出变化信息
		fmt.Printf("与最近记录(日期: %s, 记录时间: %s)比较:\n",
			latestDate.Format("2006-01-02"),
			latestCreateTime.Format("15:04:05"))
		fmt.Printf("权重因子: %.4f -> %.4f (变化: %+.4f, 变化率: %+.2f%%)\n",
			previousWeightFactor, resp.Leaderboard.WeightFactor,
			weightFactorChange, weightFactorChangeRate*100)
		fmt.Printf("价值因子: %.4f -> %.4f (变化: %+.4f, 变化率: %+.2f%%)\n",
			previousValueFactor, resp.Leaderboard.ValueFactor,
			valueFactorChange, valueFactorChangeRate*100)

	}

	// 7. 构建每日统计数据
	dailyStat := WeightValueFactor{
		UserID:   resp.Leaderboard.User,
		StatDate: currentDate,

		WeightFactor: resp.Leaderboard.WeightFactor,
		ValueFactor:  resp.Leaderboard.ValueFactor,

		WeightFactorChange:     weightFactorChange,
		ValueFactorChange:      valueFactorChange,
		WeightFactorChangeRate: weightFactorChangeRate,
		ValueFactorChangeRate:  valueFactorChangeRate,

		DataFieldsUsed:                resp.Leaderboard.DataFieldsUsed,
		SubmissionsCount:              resp.Leaderboard.SubmissionsCount,
		SuperAlphaSubmissionsCount:    resp.Leaderboard.SuperAlphaSubmissionsCount,
		MeanProdCorrelation:           resp.Leaderboard.MeanProdCorrelation,
		MeanSelfCorrelation:           resp.Leaderboard.MeanSelfCorrelation,
		SuperAlphaMeanProdCorrelation: resp.Leaderboard.SuperAlphaMeanProdCorrelation,
		SuperAlphaMeanSelfCorrelation: resp.Leaderboard.SuperAlphaMeanSelfCorrelation,
		University:                    resp.Leaderboard.University,
		Country:                       resp.Leaderboard.Country,
		DateStarted:                   dateStarted,
	}

	// 8. 使用 Upsert（如果当天已存在数据则更新）
	upsertErr := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "stat_date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"weight_factor", "value_factor", "weight_factor_change", "value_factor_change",
			"weight_factor_change_rate", "value_factor_change_rate", "data_fields_used",
			"submissions_count", "super_alpha_submissions_count", "mean_prod_correlation",
			"mean_self_correlation", "super_alpha_mean_prod_correlation",
			"super_alpha_mean_self_correlation", "university", "country", "update_time",
		}),
	}).Create(&dailyStat).Error

	if upsertErr == nil {
		fmt.Println("数据保存成功!")
	}

	return upsertErr
}
