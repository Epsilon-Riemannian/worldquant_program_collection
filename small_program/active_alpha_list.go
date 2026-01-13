package small_program

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"program-collection/models"

	"gorm.io/gorm"
)

const (
	ModeUpdate = "update" // 更新模式：重新拉取现有数据
	ModeFetch  = "fetch"  // 获取模式：拉取新数据
)

// ActiveAlphaList 活跃Alpha列表
type ActiveAlphaList struct {
	// 主键和核心字段
	ID     string `json:"id" gorm:"column:id;primaryKey;size:50;comment:Alpha ID"`
	Type   string `json:"type" gorm:"column:type;not null;size:20;index:idx_type;comment:Alpha类型"`
	Author string `json:"author" gorm:"column:author;not null;size:50;index:idx_author;comment:作者ID"`

	// Settings字段
	InstrumentType      *string  `json:"instrument_type,omitempty" gorm:"column:instrument_type;size:50;comment:工具类型"`
	Region              *string  `json:"region,omitempty" gorm:"column:region;size:50;index:idx_region;comment:地区"`
	Universe            *string  `json:"universe,omitempty" gorm:"column:universe;size:100;index:idx_universe;comment:股票池"`
	Delay               *int     `json:"delay,omitempty" gorm:"column:delay;comment:延迟参数"`
	Decay               *int     `json:"decay,omitempty" gorm:"column:decay;comment:衰减参数"`
	Neutralization      *string  `json:"neutralization,omitempty" gorm:"column:neutralization;size:50;comment:中性化方式"`
	Truncation          *float64 `json:"truncation,omitempty" gorm:"column:truncation;type:decimal(5,4);comment:截断值"`
	Pasteurization      *string  `json:"pasteurization,omitempty" gorm:"column:pasteurization;size:20;comment:巴氏杀菌设置"`
	UnitHandling        *string  `json:"unit_handling,omitempty" gorm:"column:unit_handling;size:50;comment:单位处理方式"`
	NanHandling         *string  `json:"nan_handling,omitempty" gorm:"column:nan_handling;size:50;comment:NaN处理方式"`
	SelectionHandling   *string  `json:"selection_handling,omitempty" gorm:"column:selection_handling;size:50;comment:选择处理方式"`
	SelectionLimit      *int     `json:"selection_limit,omitempty" gorm:"column:selection_limit;comment:选择限制"`
	MaxTrade            *string  `json:"max_trade,omitempty" gorm:"column:max_trade;size:20;comment:最大交易设置"`
	Language            *string  `json:"language,omitempty" gorm:"column:language;size:50;comment:语言"`
	Visualization       *bool    `json:"visualization,omitempty" gorm:"column:visualization;comment:可视化设置"`
	StartDate           *string  `json:"start_date,omitempty" gorm:"column:start_date;type:date;comment:开始日期"`
	EndDate             *string  `json:"end_date,omitempty" gorm:"column:end_date;type:date;comment:结束日期"`
	ComponentActivation *string  `json:"component_activation,omitempty" gorm:"column:component_activation;size:50;comment:组件激活"`
	TestPeriod          *string  `json:"test_period,omitempty" gorm:"column:test_period;size:20;comment:测试周期"`

	// Alpha代码内容 - SUPER类型
	ComboCode          *string `json:"combo_code,omitempty" gorm:"column:combo_code;type:text;comment:组合代码"`
	ComboDescription   *string `json:"combo_description,omitempty" gorm:"column:combo_description;type:text;comment:组合描述"`
	ComboOperatorCount *int    `json:"combo_operator_count,omitempty" gorm:"column:combo_operator_count;comment:组合运算符数量"`

	SelectionCode          *string `json:"selection_code,omitempty" gorm:"column:selection_code;type:text;comment:选择代码"`
	SelectionDescription   *string `json:"selection_description,omitempty" gorm:"column:selection_description;type:text;comment:选择描述"`
	SelectionOperatorCount *int    `json:"selection_operator_count,omitempty" gorm:"column:selection_operator_count;comment:选择运算符数量"`

	// Alpha代码内容 - REGULAR类型
	RegularCode          *string `json:"regular_code,omitempty" gorm:"column:regular_code;type:text;comment:常规代码"`
	RegularDescription   *string `json:"regular_description,omitempty" gorm:"column:regular_description;type:text;comment:常规描述"`
	RegularOperatorCount *int    `json:"regular_operator_count,omitempty" gorm:"column:regular_operator_count;comment:常规运算符数量"`

	// 基础信息
	DateCreated   *string `json:"date_created,omitempty" gorm:"column:date_created;type:datetime;index:idx_date_created;comment:创建时间"`
	DateSubmitted *string `json:"date_submitted,omitempty" gorm:"column:date_submitted;type:datetime;index:idx_date_submitted;comment:提交时间"`
	DateModified  *string `json:"date_modified,omitempty" gorm:"column:date_modified;type:datetime;comment:修改时间"`
	Name          *string `json:"name,omitempty" gorm:"column:name;size:255;comment:Alpha名称"`
	Favorite      *bool   `json:"favorite,omitempty" gorm:"column:favorite;index:idx_favorite;comment:是否收藏"`
	Hidden        *bool   `json:"hidden,omitempty" gorm:"column:hidden;comment:是否隐藏"`
	Color         *string `json:"color,omitempty" gorm:"column:color;size:50;comment:颜色标签"`
	Category      *string `json:"category,omitempty" gorm:"column:category;size:100;comment:分类"`

	// 标签和分类（JSON存储）
	Tags            *string `json:"tags,omitempty" gorm:"column:tags;type:json;comment:标签数组"`
	Classifications *string `json:"classifications,omitempty" gorm:"column:classifications;type:json;comment:分类信息数组"`

	Grade  *string `json:"grade,omitempty" gorm:"column:grade;size:50;comment:等级"`
	Stage  string  `json:"stage" gorm:"column:stage;not null;size:20;index:idx_stage;comment:阶段"`
	Status string  `json:"status" gorm:"column:status;not null;size:20;index:idx_status;comment:状态"`

	// IS性能指标
	IsPnl             *int     `json:"is_pnl,omitempty" gorm:"column:is_pnl;comment:IS期间PNL"`
	IsBookSize        *int     `json:"is_book_size,omitempty" gorm:"column:is_book_size;comment:IS期间账面大小"`
	IsLongCount       *int     `json:"is_long_count,omitempty" gorm:"column:is_long_count;comment:IS期间多头数量"`
	IsShortCount      *int     `json:"is_short_count,omitempty" gorm:"column:is_short_count;comment:IS期间空头数量"`
	IsTurnover        *float64 `json:"is_turnover,omitempty" gorm:"column:is_turnover;type:decimal(10,4);comment:IS期间换手率"`
	IsReturns         *float64 `json:"is_returns,omitempty" gorm:"column:is_returns;type:decimal(10,4);comment:IS期间收益率"`
	IsDrawdown        *float64 `json:"is_drawdown,omitempty" gorm:"column:is_drawdown;type:decimal(10,4);comment:IS期间回撤"`
	IsMargin          *float64 `json:"is_margin,omitempty" gorm:"column:is_margin;type:decimal(10,6);comment:IS期间保证金"`
	IsSharpe          *float64 `json:"is_sharpe,omitempty" gorm:"column:is_sharpe;type:decimal(10,2);index:idx_sharpe;comment:IS期间夏普比率"`
	IsFitness         *float64 `json:"is_fitness,omitempty" gorm:"column:is_fitness;type:decimal(10,2);index:idx_fitness;comment:IS期间适应度"`
	IsStartDate       *string  `json:"is_start_date,omitempty" gorm:"column:is_start_date;type:date;comment:IS开始日期"`
	IsSelfCorrelation *float64 `json:"is_self_correlation,omitempty" gorm:"column:is_self_correlation;type:decimal(10,4);comment:IS自相关"`
	IsProdCorrelation *float64 `json:"is_prod_correlation,omitempty" gorm:"column:is_prod_correlation;type:decimal(10,4);comment:IS与生产相关"`
	IsChecks          *string  `json:"is_checks,omitempty" gorm:"column:is_checks;type:json;comment:IS检查项数组"`

	// OS信息
	OsStartDate           *string `json:"os_start_date,omitempty" gorm:"column:os_start_date;type:date;comment:OS开始日期"`
	OsIsSharpeRatio       *string `json:"os_is_sharpe_ratio,omitempty" gorm:"column:os_is_sharpe_ratio;type:json;comment:OS IS夏普比率"`
	OsPreCloseSharpeRatio *string `json:"os_pre_close_sharpe_ratio,omitempty" gorm:"column:os_pre_close_sharpe_ratio;type:json;comment:OS前收盘夏普比率"`
	OsChecks              *string `json:"os_checks,omitempty" gorm:"column:os_checks;type:json;comment:OS检查项数组"`

	// Train性能指标（SUPER类型特有）
	TrainPnl        *int     `json:"train_pnl,omitempty" gorm:"column:train_pnl;comment:训练期间PNL"`
	TrainBookSize   *int     `json:"train_book_size,omitempty" gorm:"column:train_book_size;comment:训练期间账面大小"`
	TrainLongCount  *int     `json:"train_long_count,omitempty" gorm:"column:train_long_count;comment:训练期间多头数量"`
	TrainShortCount *int     `json:"train_short_count,omitempty" gorm:"column:train_short_count;comment:训练期间空头数量"`
	TrainTurnover   *float64 `json:"train_turnover,omitempty" gorm:"column:train_turnover;type:decimal(10,4);comment:训练期间换手率"`
	TrainReturns    *float64 `json:"train_returns,omitempty" gorm:"column:train_returns;type:decimal(10,4);comment:训练期间收益率"`
	TrainDrawdown   *float64 `json:"train_drawdown,omitempty" gorm:"column:train_drawdown;type:decimal(10,4);comment:训练期间回撤"`
	TrainMargin     *float64 `json:"train_margin,omitempty" gorm:"column:train_margin;type:decimal(10,6);comment:训练期间保证金"`
	TrainSharpe     *float64 `json:"train_sharpe,omitempty" gorm:"column:train_sharpe;type:decimal(10,2);comment:训练期间夏普比率"`
	TrainFitness    *float64 `json:"train_fitness,omitempty" gorm:"column:train_fitness;type:decimal(10,2);comment:训练期间适应度"`
	TrainStartDate  *string  `json:"train_start_date,omitempty" gorm:"column:train_start_date;type:date;comment:训练开始日期"`

	// Test性能指标（SUPER类型特有）
	TestPnl        *int     `json:"test_pnl,omitempty" gorm:"column:test_pnl;comment:测试期间PNL"`
	TestBookSize   *int     `json:"test_book_size,omitempty" gorm:"column:test_book_size;comment:测试期间账面大小"`
	TestLongCount  *int     `json:"test_long_count,omitempty" gorm:"column:test_long_count;comment:测试期间多头数量"`
	TestShortCount *int     `json:"test_short_count,omitempty" gorm:"column:test_short_count;comment:测试期间空头数量"`
	TestTurnover   *float64 `json:"test_turnover,omitempty" gorm:"column:test_turnover;type:decimal(10,4);comment:测试期间换手率"`
	TestReturns    *float64 `json:"test_returns,omitempty" gorm:"column:test_returns;type:decimal(10,4);comment:测试期间收益率"`
	TestDrawdown   *float64 `json:"test_drawdown,omitempty" gorm:"column:test_drawdown;type:decimal(10,4);comment:测试期间回撤"`
	TestMargin     *float64 `json:"test_margin,omitempty" gorm:"column:test_margin;type:decimal(10,6);comment:测试期间保证金"`
	TestSharpe     *float64 `json:"test_sharpe,omitempty" gorm:"column:test_sharpe;type:decimal(10,2);comment:测试期间夏普比率"`
	TestFitness    *float64 `json:"test_fitness,omitempty" gorm:"column:test_fitness;type:decimal(10,2);comment:测试期间适应度"`
	TestStartDate  *string  `json:"test_start_date,omitempty" gorm:"column:test_start_date;type:date;comment:测试开始日期"`

	// 其他JSON字段
	Prod          *string `json:"prod,omitempty" gorm:"column:prod;type:json;comment:生产数据"`
	Competitions  *string `json:"competitions,omitempty" gorm:"column:competitions;type:json;comment:比赛数据"`
	Themes        *string `json:"themes,omitempty" gorm:"column:themes;type:json;comment:主题数组"`
	Pyramids      *string `json:"pyramids,omitempty" gorm:"column:pyramids;type:json;comment:金字塔数组"`
	PyramidThemes *string `json:"pyramid_themes,omitempty" gorm:"column:pyramid_themes;type:json;comment:金字塔主题"`
	Team          *string `json:"team,omitempty" gorm:"column:team;type:json;comment:团队信息"`
	OsmosisPoints *string `json:"osmosis_points,omitempty" gorm:"column:osmosis_points;type:json;comment:渗透点数"`

	// 三个时间字段
	CreateTime  *time.Time `json:"create_time,omitempty" gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_create_time;comment:创建时间"`
	CreateDate  *string    `json:"create_date,omitempty" gorm:"column:create_date;type:date;default:(CURRENT_DATE);index:idx_create_date;comment:创建日期"`
	CreateMonth *string    `json:"create_month,omitempty" gorm:"column:create_month;size:7;index:idx_create_month;comment:创建月份"`
}

// TableName 指定表名
func (ActiveAlphaList) TableName() string {
	return "active_alpha_list"
}

// 获取用户输入
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// 显示 ActiveAlpha 管理子菜单
func showActiveAlphaMenu() {
	fmt.Println("\n--------------------------------------------")
	fmt.Println("         ActiveAlpha 管理")
	fmt.Println("--------------------------------------------")
	fmt.Println("1. 获取新的 Alpha")
	fmt.Println("2. 更新现有 Alpha")
	fmt.Println("3. 返回主菜单")
	fmt.Println("--------------------------------------------")
	fmt.Print("请选择操作 (1-3): ")
}

// 10. 运行 ActiveAlpha 管理
func RunActiveAlphaManagement(config models.Config, token string) {

	// 1. 连接数据库
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	for {
		showActiveAlphaMenu()
		choice := getUserInput()

		switch choice {
		case "1":
			err := FetchNewAlphas(config, token, db)
			if err != nil {
				log.Printf("获取新的 Alpha 失败: %v", err)
			} else {
				fmt.Println("获取新的 Alpha 成功！")
			}
		case "2":
			err := UpdateExistingAlphas(config, token, db)
			if err != nil {
				log.Printf("更新现有 Alpha 失败: %v", err)
			} else {
				fmt.Println("更新现有 Alpha 成功！")
			}

		case "3":
			return
		default:
			fmt.Println("无效的选择，请输入 1-3 之间的数字！")
		}
	}
}

// 1. 更新模式：重新拉取数据库中已有数据
func UpdateExistingAlphas(config models.Config, token string, db *gorm.DB) error {
	log.Println("=== 开始更新模式：重新拉取数据库中已有数据 ===")

	// 获取数据库中所有Alpha的ID
	var alphaIDs []string
	result := db.Model(&ActiveAlphaList{}).Pluck("id", &alphaIDs)
	if result.Error != nil {
		return fmt.Errorf("failed to get alpha IDs: %v", result.Error)
	}

	log.Printf("数据库中共有 %d 条Alpha数据需要更新", len(alphaIDs))

	// 分批处理，避免一次性请求过多
	batchSize := 20
	totalUpdated := 0

	for i := 0; i < len(alphaIDs); i += batchSize {
		end := i + batchSize
		if end > len(alphaIDs) {
			end = len(alphaIDs)
		}

		batchIDs := alphaIDs[i:end]
		log.Printf("处理批次 %d-%d", i+1, end)

		// 1.1 为这批ID获取最新数据
		updatedCount, err := updateBatchAlphas(config, token, db, batchIDs)
		if err != nil {
			log.Printf("批次 %d-%d 更新失败: %v", i+1, end, err)
			continue
		}

		totalUpdated += updatedCount
		log.Printf("批次 %d-%d 更新完成，更新了 %d 条", i+1, end, updatedCount)

		// 避免请求过于频繁
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("=== 更新模式完成，总共更新了 %d 条数据 ===", totalUpdated)
	return nil
}

// 1.1 更新一批Alpha数据
func updateBatchAlphas(config models.Config, token string, db *gorm.DB, alphaIDs []string) (int, error) {
	updatedCount := 0

	// 逐个更新
	for _, alphaID := range alphaIDs {
		alpha, err := GetAlphaByID(config, token, alphaID)
		if err != nil {
			log.Printf("获取Alpha %s 失败: %v", alphaID, err)
			continue
		}

		// 转换为数据库结构
		dbAlpha := convertAlphaToDB(alpha)

		// 更新数据库（只更新，不创建）
		result := db.Model(&ActiveAlphaList{}).Where("id = ?", alphaID).Updates(dbAlpha)
		if result.Error != nil {
			log.Printf("更新Alpha %s 到数据库失败: %v", alphaID, result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			updatedCount++
			log.Printf("已更新Alpha: %s", alphaID)
		}
	}

	return updatedCount, nil
}

// 2. 获取模式：从数据库最大日期拉到今天当前，获取新数据
func FetchNewAlphas(config models.Config, token string, db *gorm.DB) error {
	log.Println("=== 开始获取模式：拉取新数据 ===")

	// 获取数据库中最大的日期
	// 使用指针来处理 NULL 值
	var maxDate *string
	result := db.Model(&ActiveAlphaList{}).Select("MAX(date_submitted) as max_date").Scan(&maxDate)
	if result.Error != nil {
		return fmt.Errorf("failed to get max date: %v", result.Error)
	}

	var startDateStr string
	// 检查指针是否为 nil（表示数据库返回 NULL）
	if maxDate == nil {
		// 数据库为空，设置一个较远的开始日期，比如5年前
		fiveYearsAgo := time.Now().AddDate(-5, 0, 0)
		startDateStr = fiveYearsAgo.Format("2006-01-02 15:04:05")
		log.Printf("数据库为空，从 %s 开始获取", startDateStr)
	} else {
		startDateStr = *maxDate
		log.Printf("数据库中最新日期: %s", startDateStr)
	}

	// 转换为time.Time
	dateFrom, err := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
	if err != nil {
		dateFrom, err = time.Parse("2006-01-02T15:04:05-07:00", startDateStr)
		if err != nil {
			return fmt.Errorf("failed to parse max date: %v", err)
		}
	}

	// 今天的现在日期
	now := time.Now()

	// 重要：确保 dateFrom 在 now 之前
	if dateFrom.After(now) || dateFrom.Equal(now) {
		log.Println("数据已是最新，无需获取")
		return nil
	}

	// 分页获取数据
	limit := 50
	offset := 0
	totalFetched := 0

	// 设置结束日期为当前时间
	endDate := now

	// 只尝试一次，使用正确的格式
	for {
		log.Printf("获取数据: %s 到 %s, 偏移量: %d", dateFrom.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"), offset)

		beginISO, _ := ConvertToUTCPlus5(dateFrom.Format("2006-01-02 15:04:05"))
		endISO, _ := ConvertToUTCPlus5(endDate.Format("2006-01-02 15:04:05"))

		// 调用API获取数据
		alphaLists, _ := GetAllAlphas(config, token, models.GetAlphasRequest{
			Limit:    limit,
			Offset:   offset,
			DateFrom: beginISO,
			DateTo:   endISO,
			Order:    "dateSubmitted", // 按提交日期升序，确保获取完整
		})

		if len(alphaLists) == 0 {
			log.Println("没有更多数据")
			break
		}

		// 转换并保存数据
		var dbAlphas []ActiveAlphaList
		for _, alpha := range alphaLists {
			dbAlpha := convertAlphaToDB(alpha)
			dbAlphas = append(dbAlphas, dbAlpha)
		}

		// 批量插入（使用FirstOrCreate避免重复）
		insertedCount, err := batchInsertOrIgnore(db, dbAlphas)
		if err != nil {
			log.Printf("批量插入失败: %v，尝试逐个插入", err)
			// 逐个插入
			insertedCount = 0
			for _, dbAlpha := range dbAlphas {
				result := db.Where("id = ?", dbAlpha.ID).FirstOrCreate(&dbAlpha)
				if result.Error == nil && result.RowsAffected > 0 {
					insertedCount++
				}
			}
		}

		totalFetched += insertedCount
		log.Printf("批次获取 %d 条，插入 %d 条，累计 %d 条", len(alphaLists), insertedCount, totalFetched)

		offset += limit

		// 如果返回的数据少于请求的数量，说明已经到达末尾
		if len(alphaLists) < limit {
			break
		}

		// 避免请求过于频繁
		time.Sleep(300 * time.Millisecond)
	}

	log.Printf("=== 获取模式完成，总共获取了 %d 条新数据 ===", totalFetched)
	return nil
}

// 辅助函数：创建指针
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

// 批量插入，忽略重复（使用INSERT IGNORE）
func batchInsertOrIgnore(db *gorm.DB, alphas []ActiveAlphaList) (int, error) {
	if len(alphas) == 0 {
		return 0, nil
	}

	// 使用原生SQL进行批量INSERT IGNORE
	// 注意：这里假设你使用MySQL
	sql := `INSERT IGNORE INTO active_alpha_list 
		(id, type, author, instrument_type, region, universe, delay, decay, 
		neutralization, truncation, pasteurization, unit_handling, max_trade, 
		language, start_date, end_date, visualization, date_created, date_submitted, 
		date_modified, name, favorite, hidden, color, category, grade, stage, status,
		tags, classifications, is_pnl, is_book_size, is_long_count, is_short_count,
		is_turnover, is_returns, is_drawdown, is_margin, is_sharpe, is_fitness,
		is_start_date, is_self_correlation, is_prod_correlation, is_checks,
		os_start_date, os_is_sharpe_ratio, os_pre_close_sharpe_ratio, os_checks,
		train_pnl, train_book_size, train_long_count, train_short_count, train_turnover,
		train_returns, train_drawdown, train_margin, train_sharpe, train_fitness,
		train_start_date, test_pnl, test_book_size, test_long_count, test_short_count,
		test_turnover, test_returns, test_drawdown, test_margin, test_sharpe,
		test_fitness, test_start_date, prod, competitions, themes, pyramids,
		pyramid_themes, team, osmosis_points, create_time, create_date, create_month)
		VALUES `

	var valueStrings []string
	var valueArgs []interface{}

	for _, alpha := range alphas {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs,
			alpha.ID, alpha.Type, alpha.Author,
			alpha.InstrumentType, alpha.Region, alpha.Universe,
			alpha.Delay, alpha.Decay, alpha.Neutralization,
			alpha.Truncation, alpha.Pasteurization, alpha.UnitHandling,
			alpha.MaxTrade, alpha.Language, alpha.StartDate, alpha.EndDate,
			alpha.Visualization, alpha.DateCreated, alpha.DateSubmitted,
			alpha.DateModified, alpha.Name, alpha.Favorite, alpha.Hidden,
			alpha.Color, alpha.Category, alpha.Grade, alpha.Stage, alpha.Status,
			alpha.Tags, alpha.Classifications, alpha.IsPnl, alpha.IsBookSize,
			alpha.IsLongCount, alpha.IsShortCount, alpha.IsTurnover,
			alpha.IsReturns, alpha.IsDrawdown, alpha.IsMargin, alpha.IsSharpe,
			alpha.IsFitness, alpha.IsStartDate, alpha.IsSelfCorrelation,
			alpha.IsProdCorrelation, alpha.IsChecks, alpha.OsStartDate,
			alpha.OsIsSharpeRatio, alpha.OsPreCloseSharpeRatio, alpha.OsChecks,
			alpha.TrainPnl, alpha.TrainBookSize, alpha.TrainLongCount,
			alpha.TrainShortCount, alpha.TrainTurnover, alpha.TrainReturns,
			alpha.TrainDrawdown, alpha.TrainMargin, alpha.TrainSharpe,
			alpha.TrainFitness, alpha.TrainStartDate, alpha.TestPnl,
			alpha.TestBookSize, alpha.TestLongCount, alpha.TestShortCount,
			alpha.TestTurnover, alpha.TestReturns, alpha.TestDrawdown,
			alpha.TestMargin, alpha.TestSharpe, alpha.TestFitness,
			alpha.TestStartDate, alpha.Prod, alpha.Competitions, alpha.Themes,
			alpha.Pyramids, alpha.PyramidThemes, alpha.Team, alpha.OsmosisPoints,
			alpha.CreateTime, alpha.CreateDate, alpha.CreateMonth,
		)
	}

	sql += join(valueStrings, ",")

	result := db.Exec(sql, valueArgs...)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// 辅助函数：拼接字符串
func join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// 转换函数
func convertAlphaToDB(alpha models.Alpha) ActiveAlphaList {

	// 如果是更新模式，不修改create_time
	// 如果是获取模式，设置新的create_time

	dbAlpha := ActiveAlphaList{
		// 主键和核心字段
		ID:     alpha.ID,
		Type:   alpha.Type,
		Author: alpha.Author,

		// Settings字段转换
		InstrumentType: stringPtr(alpha.Settings.InstrumentType),
		Region:         stringPtr(alpha.Settings.Region),
		Universe:       stringPtr(alpha.Settings.Universe),
		Delay:          intPtr(alpha.Settings.Delay),
		Decay:          intPtr(alpha.Settings.Decay),
		Neutralization: stringPtr(alpha.Settings.Neutralization),
		Truncation:     float64Ptr(alpha.Settings.Truncation),
		Pasteurization: stringPtr(alpha.Settings.Pasteurization),
		UnitHandling:   stringPtr(alpha.Settings.UnitHandling),
		MaxTrade:       stringPtr(alpha.Settings.MaxTrade),
		Language:       stringPtr(alpha.Settings.Language),
		StartDate:      stringPtr(alpha.Settings.StartDate),
		EndDate:        stringPtr(alpha.Settings.EndDate),
	}

	// 可选字段处理
	if alpha.Settings.NanHandling != "" {
		dbAlpha.NanHandling = stringPtr(alpha.Settings.NanHandling)
	}
	if alpha.Settings.SelectionHandling != "" {
		dbAlpha.SelectionHandling = stringPtr(alpha.Settings.SelectionHandling)
	}
	if alpha.Settings.SelectionLimit > 0 {
		dbAlpha.SelectionLimit = intPtr(alpha.Settings.SelectionLimit)
	}
	if alpha.Settings.ComponentActivation != "" {
		dbAlpha.ComponentActivation = stringPtr(alpha.Settings.ComponentActivation)
	}
	if alpha.Settings.TestPeriod != "" {
		dbAlpha.TestPeriod = stringPtr(alpha.Settings.TestPeriod)
	}

	dbAlpha.Visualization = boolPtr(alpha.Settings.Visualization)

	// Alpha代码内容 - 根据类型处理
	switch alpha.Type {
	case "SUPER":
		if alpha.Combo != nil {
			dbAlpha.ComboCode = stringPtr(alpha.Combo.Code)
			dbAlpha.ComboDescription = stringPtr(alpha.Combo.Description)
			dbAlpha.ComboOperatorCount = alpha.Combo.OperatorCount
		}
		if alpha.Selection != nil {
			dbAlpha.SelectionCode = stringPtr(alpha.Selection.Code)
			dbAlpha.SelectionDescription = stringPtr(alpha.Selection.Description)
			dbAlpha.SelectionOperatorCount = alpha.Selection.OperatorCount
		}
	case "REGULAR":
		if alpha.Regular != nil {
			dbAlpha.RegularCode = stringPtr(alpha.Regular.Code)
			dbAlpha.RegularDescription = stringPtr(alpha.Regular.Description)
			dbAlpha.RegularOperatorCount = alpha.Regular.OperatorCount
		}
	}

	// 基础信息
	dbAlpha.DateCreated = stringPtr(alpha.DateCreated)
	dbAlpha.DateSubmitted = stringPtr(alpha.DateSubmitted)
	dbAlpha.DateModified = stringPtr(alpha.DateModified)
	dbAlpha.Name = alpha.Name
	dbAlpha.Favorite = boolPtr(alpha.Favorite)
	dbAlpha.Hidden = boolPtr(alpha.Hidden)
	dbAlpha.Color = stringPtr(alpha.Color)
	dbAlpha.Category = alpha.Category
	dbAlpha.Grade = alpha.Grade
	dbAlpha.Stage = alpha.Stage
	dbAlpha.Status = alpha.Status

	// 标签和分类（转换为JSON字符串）
	if len(alpha.Tags) > 0 {
		tagsJSON, _ := json.Marshal(alpha.Tags)
		dbAlpha.Tags = stringPtr(string(tagsJSON))
	}
	if len(alpha.Classifications) > 0 {
		classificationsJSON, _ := json.Marshal(alpha.Classifications)
		dbAlpha.Classifications = stringPtr(string(classificationsJSON))
	}

	// IS性能指标
	if alpha.IS != nil {
		dbAlpha.IsPnl = intPtr(alpha.IS.PNL)
		dbAlpha.IsBookSize = intPtr(alpha.IS.BookSize)
		dbAlpha.IsLongCount = intPtr(alpha.IS.LongCount)
		dbAlpha.IsShortCount = intPtr(alpha.IS.ShortCount)
		dbAlpha.IsTurnover = float64Ptr(alpha.IS.Turnover)
		dbAlpha.IsReturns = float64Ptr(alpha.IS.Returns)
		dbAlpha.IsDrawdown = float64Ptr(alpha.IS.Drawdown)
		dbAlpha.IsMargin = float64Ptr(alpha.IS.Margin)
		dbAlpha.IsSharpe = float64Ptr(alpha.IS.Sharpe)
		dbAlpha.IsFitness = float64Ptr(alpha.IS.Fitness)
		dbAlpha.IsStartDate = stringPtr(alpha.IS.StartDate)
		dbAlpha.IsSelfCorrelation = float64Ptr(alpha.IS.SelfCorrelation)
		dbAlpha.IsProdCorrelation = float64Ptr(alpha.IS.ProdCorrelation)

		if len(alpha.IS.Checks) > 0 {
			checksJSON, _ := json.Marshal(alpha.IS.Checks)
			dbAlpha.IsChecks = stringPtr(string(checksJSON))
		}
	}

	// OS信息
	if alpha.OS != nil {
		dbAlpha.OsStartDate = stringPtr(alpha.OS.StartDate)

		if alpha.OS.OsISSharpeRatio != nil {
			osSharpeJSON, _ := json.Marshal(alpha.OS.OsISSharpeRatio)
			dbAlpha.OsIsSharpeRatio = stringPtr(string(osSharpeJSON))
		}

		if alpha.OS.PreCloseSharpeRatio != nil {
			preSharpeJSON, _ := json.Marshal(alpha.OS.PreCloseSharpeRatio)
			dbAlpha.OsPreCloseSharpeRatio = stringPtr(string(preSharpeJSON))
		}

		if len(alpha.OS.Checks) > 0 {
			osChecksJSON, _ := json.Marshal(alpha.OS.Checks)
			dbAlpha.OsChecks = stringPtr(string(osChecksJSON))
		}
	}

	// Train性能指标
	if alpha.Train != nil {
		dbAlpha.TrainPnl = intPtr(alpha.Train.PNL)
		dbAlpha.TrainBookSize = intPtr(alpha.Train.BookSize)
		dbAlpha.TrainLongCount = intPtr(alpha.Train.LongCount)
		dbAlpha.TrainShortCount = intPtr(alpha.Train.ShortCount)
		dbAlpha.TrainTurnover = float64Ptr(alpha.Train.Turnover)
		dbAlpha.TrainReturns = float64Ptr(alpha.Train.Returns)
		dbAlpha.TrainDrawdown = float64Ptr(alpha.Train.Drawdown)
		dbAlpha.TrainMargin = float64Ptr(alpha.Train.Margin)
		dbAlpha.TrainSharpe = float64Ptr(alpha.Train.Sharpe)
		dbAlpha.TrainFitness = float64Ptr(alpha.Train.Fitness)
		dbAlpha.TrainStartDate = stringPtr(alpha.Train.StartDate)
	}

	// Test性能指标
	if alpha.Test != nil {
		dbAlpha.TestPnl = intPtr(alpha.Test.PNL)
		dbAlpha.TestBookSize = intPtr(alpha.Test.BookSize)
		dbAlpha.TestLongCount = intPtr(alpha.Test.LongCount)
		dbAlpha.TestShortCount = intPtr(alpha.Test.ShortCount)
		dbAlpha.TestTurnover = float64Ptr(alpha.Test.Turnover)
		dbAlpha.TestReturns = float64Ptr(alpha.Test.Returns)
		dbAlpha.TestDrawdown = float64Ptr(alpha.Test.Drawdown)
		dbAlpha.TestMargin = float64Ptr(alpha.Test.Margin)
		dbAlpha.TestSharpe = float64Ptr(alpha.Test.Sharpe)
		dbAlpha.TestFitness = float64Ptr(alpha.Test.Fitness)
		dbAlpha.TestStartDate = stringPtr(alpha.Test.StartDate)
	}

	// 其他JSON字段
	if alpha.Prod != nil {
		prodJSON, _ := json.Marshal(alpha.Prod)
		dbAlpha.Prod = stringPtr(string(prodJSON))
	}

	if len(alpha.Competitions) > 0 {
		competitionsJSON, _ := json.Marshal(alpha.Competitions)
		dbAlpha.Competitions = stringPtr(string(competitionsJSON))
	}

	if len(alpha.Themes) > 0 {
		themesJSON, _ := json.Marshal(alpha.Themes)
		dbAlpha.Themes = stringPtr(string(themesJSON))
	}

	if len(alpha.Pyramids) > 0 {
		pyramidsJSON, _ := json.Marshal(alpha.Pyramids)
		dbAlpha.Pyramids = stringPtr(string(pyramidsJSON))
	}

	pyramidThemesJSON, _ := json.Marshal(alpha.PyramidThemes)
	dbAlpha.PyramidThemes = stringPtr(string(pyramidThemesJSON))

	if alpha.Team != nil {
		teamJSON, _ := json.Marshal(alpha.Team)
		dbAlpha.Team = stringPtr(string(teamJSON))
	}

	if alpha.OsmosisPoints != nil {
		pointsJSON, _ := json.Marshal(alpha.OsmosisPoints)
		dbAlpha.OsmosisPoints = stringPtr(string(pointsJSON))
	}

	// 设置创建时间为当前时间（只对新数据）
	now := time.Now()
	dbAlpha.CreateTime = &now
	dateStr := now.Format("2006-01-02")
	dbAlpha.CreateDate = &dateStr
	monthStr := now.Format("2006-01")
	dbAlpha.CreateMonth = &monthStr

	return dbAlpha
}
