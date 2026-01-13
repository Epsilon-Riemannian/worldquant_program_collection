package small_program

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"program-collection/models"

	"gorm.io/gorm"
)

type PyramidAlphas struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        string    `gorm:"column:user_id;size:100;not null;index:idx_user_id" json:"userId"`
	CategoryID    string    `gorm:"column:category_id;size:50;index:idx_category_id" json:"categoryId"`
	CategoryName  string    `gorm:"column:category_name;size:255" json:"categoryName"`
	Region        string    `gorm:"column:region;size:100;index:idx_region" json:"region"`
	Delay         int       `gorm:"column:delay;default:0" json:"delay"`
	AlphaCount    int       `gorm:"column:alpha_count;default:0" json:"alphaCount"`
	QuarterTag    string    `gorm:"column:quarter_tag;size:10;index:idx_quarter_tag" json:"quarterTag"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime;index:idx_create_time" json:"createTime"`
	CreateDate    time.Time `gorm:"column:create_date;autoCreateTime;index:idx_create_date" json:"createDate"`
	CreateMonth   string    `gorm:"-" json:"createMonth"` // 添加 gorm:"-" 忽略此字段
	StatStartDate time.Time `gorm:"column:stat_start_date;index:idx_stat_date_range" json:"statStartDate"`
	StatEndDate   time.Time `gorm:"column:stat_end_date;index:idx_stat_date_range" json:"statEndDate"`
}

func (PyramidAlphas) TableName() string {
	return "pyramid_alphas"
}

func PyramidAlphaInfo(config models.Config, token string) error {
	// 创建控制台读取器
	reader := bufio.NewReader(os.Stdin)

	// 1. 交互式获取季度信息
	quarter, err := getQuarterInput(reader)
	if err != nil {
		return err
	}

	// 2. 根据季度计算日期范围
	startDate, endDate, err := calculateQuarterDates(quarter)
	if err != nil {
		return fmt.Errorf("季度格式错误: %v", err)
	}

	// 3. 交互式获取用户ID
	userID, err := getUserIDInput(reader)
	if err != nil {
		return err
	}

	// 4. 显示确认信息
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("请确认以下信息:")
	fmt.Printf("季度: %s\n", quarter)
	fmt.Printf("日期范围: %s 至 %s\n", startDate, endDate)
	fmt.Printf("用户ID: %s\n", userID)
	fmt.Println(strings.Repeat("=", 50))

	// 5. 确认是否继续
	if !confirmAction(reader, "是否继续执行数据获取和保存？") {
		fmt.Println("操作已取消")
		return nil
	}

	fmt.Printf("\n正在获取 %s 的数据...\n", quarter)

	// 6. 调用API获取数据
	pyramids, err := PyramidInfo(config, token, startDate, endDate)
	if err != nil {
		return fmt.Errorf("获取金字塔数据失败: %v", err)
	}

	// 7. 连接数据库
	db, err := ConnectDB(config)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 8. 检查是否已存在该季度数据
	if shouldCheckExistingData(reader, quarter, userID) {
		exists, err := checkExistingQuarterData(db, quarter, userID)
		if err != nil {
			return fmt.Errorf("检查已有数据失败: %v", err)
		}

		if exists {
			if !confirmAction(reader, "该季度数据已存在，是否更新？") {
				fmt.Println("操作已取消")
				return nil
			}

			// 删除现有数据
			if err := deleteExistingQuarterData(db, quarter, userID); err != nil {
				return fmt.Errorf("删除现有数据失败: %v", err)
			}
			fmt.Println("已删除现有数据，准备重新插入...")
		}
	}

	// 9. 解析日期范围
	statStart, _ := time.Parse("2006-01-02", startDate)
	statEnd, _ := time.Parse("2006-01-02", endDate)

	// 10. 遍历数据并转换为数据库模型
	var records []PyramidAlphas
	for _, p := range pyramids {
		record := PyramidAlphas{
			UserID:        userID,
			CategoryID:    p.Category.ID,
			CategoryName:  p.Category.Name,
			Region:        p.Region,
			Delay:         p.Delay,
			AlphaCount:    p.AlphaCount,
			StatStartDate: statStart,
			StatEndDate:   statEnd,
			QuarterTag:    quarter,
		}
		records = append(records, record)
	}

	// 11. 批量插入数据库
	tx := db.Begin()
	if err := tx.CreateInBatches(&records, 100).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("批量插入数据失败: %v", err)
	}
	tx.Commit()

	fmt.Printf("\n✅ 成功保存 %d 条金字塔Alpha记录\n", len(records))
	fmt.Printf("   季度: %s\n", quarter)
	fmt.Printf("   用户: %s\n", userID)
	fmt.Printf("   时间: %s 至 %s\n", startDate, endDate)

	return nil
}

// 获取季度输入的辅助函数
func getQuarterInput(reader *bufio.Reader) (string, error) {
	for {
		fmt.Println("\n请输入季度（格式: 年份-QN，例如: 2025-Q3）:")
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("读取输入失败: %v", err)
		}

		// 清理输入
		input = strings.TrimSpace(input)

		// 验证季度格式
		if isValidQuarterFormat(input) {
			return input, nil
		}

		fmt.Printf("❌ 格式错误: %s\n", input)
		fmt.Println("正确格式示例: 2025-Q3, 2024-Q1, 2023-Q4")
		fmt.Println("Q1: 1-3月, Q2: 4-6月, Q3: 7-9月, Q4: 10-12月")
	}
}

// 验证季度格式
func isValidQuarterFormat(input string) bool {
	if len(input) != 7 && len(input) != 6 {
		return false
	}

	// 检查格式：YYYY-QN
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return false
	}

	year := parts[0]
	quarter := parts[1]

	// 检查年份
	if len(year) != 4 {
		return false
	}

	// 检查季度
	if len(quarter) < 2 || quarter[0] != 'Q' {
		return false
	}

	// 检查季度数字 (1-4)
	if len(quarter) == 2 {
		qNum := quarter[1]
		if qNum < '1' || qNum > '4' {
			return false
		}
	}

	return true
}

// 根据季度计算日期范围
func calculateQuarterDates(quarter string) (startDate, endDate string, err error) {
	parts := strings.Split(quarter, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效的季度格式")
	}

	year := parts[0]
	quarterNum := parts[1][1] - '0' // 提取Q后面的数字

	var startMonth, endMonth int

	switch quarterNum {
	case 1: // Q1: 1月-3月
		startMonth, endMonth = 1, 4
	case 2: // Q2: 4月-6月
		startMonth, endMonth = 4, 7
	case 3: // Q3: 7月-9月
		startMonth, endMonth = 7, 10
	case 4: // Q4: 10月-12月
		startMonth, endMonth = 10, 1
	default:
		return "", "", fmt.Errorf("季度数字必须在1-4之间")
	}

	// 计算开始日期
	startDate = fmt.Sprintf("%s-%02d-01", year, startMonth)

	// 计算结束日期（如果是Q4，年份要加1）
	endYear := year
	if quarterNum == 4 {
		endYearInt := 0
		fmt.Sscanf(year, "%d", &endYearInt)
		endYear = fmt.Sprintf("%d", endYearInt+1)
	}

	endDate = fmt.Sprintf("%s-%02d-01", endYear, endMonth)

	return startDate, endDate, nil
}

// 获取用户ID输入
func getUserIDInput(reader *bufio.Reader) (string, error) {
	for {
		fmt.Println("\n请输入用户ID:")
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("读取输入失败: %v", err)
		}

		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("❌ 用户ID不能为空")
			continue
		}

		// 确认用户ID
		fmt.Printf("您输入的用户ID是: %s\n", input)
		if confirmAction(reader, "是否确认？") {
			return input, nil
		}
	}
}

// 确认操作的辅助函数
func confirmAction(reader *bufio.Reader, prompt string) bool {
	for {
		fmt.Printf("%s (y/n): ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes", "是", "确认":
			return true
		case "n", "no", "否", "取消":
			return false
		default:
			fmt.Println("请输入 y/是 或 n/否")
		}
	}
}

// 检查是否应该检查已有数据
func shouldCheckExistingData(reader *bufio.Reader, quarter, userID string) bool {
	fmt.Printf("\n是否检查 %s 用户 %s 季度的现有数据？\n", userID, quarter)
	return confirmAction(reader, "检查并提示是否更新？")
}

// 检查数据库中是否已存在该季度的数据
func checkExistingQuarterData(db *gorm.DB, quarter, userID string) (bool, error) {
	var count int64
	err := db.Model(&PyramidAlphas{}).
		Where("user_id = ? AND quarter_tag = ?", userID, quarter).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	if count > 0 {
		fmt.Printf("⚠️  数据库中年发现 %d 条已存在的 %s 季度数据\n", count, quarter)
		return true, nil
	}

	fmt.Printf("✅ 数据库中未找到 %s 季度的现有数据\n", quarter)
	return false, nil
}

// 删除现有季度数据
func deleteExistingQuarterData(db *gorm.DB, quarter, userID string) error {
	result := db.Where("user_id = ? AND quarter_tag = ?", userID, quarter).
		Delete(&PyramidAlphas{})

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("已删除 %d 条现有记录\n", result.RowsAffected)
	return nil
}
