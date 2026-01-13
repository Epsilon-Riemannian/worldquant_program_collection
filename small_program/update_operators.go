package small_program

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"program-collection/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---------------------------------- Operator操作符入库结构体 ------------------------------ //
type Operators struct {
	ID            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `json:"name" gorm:"column:name;type:varchar(255);not null;unique"`
	Category      string `json:"category" gorm:"column:category;type:varchar(100);not null"`
	Scope         string `json:"scope" gorm:"column:scope;type:json"`
	Definition    string `json:"definition" gorm:"column:definition;type:text"`
	EnDescription string `json:"en_description" gorm:"column:en_description;type:text"`
	CnDescription string `json:"cn_description" gorm:"column:cn_description;type:text"`
	Documentation string `json:"documentation" gorm:"column:documentation;type:text"`
	Level         string `json:"level" gorm:"column:level;type:varchar(50)"`
	GeniusLevel   string `json:"genius_level" gorm:"column:genius_level;type:varchar(50)"`
	GeniusQuarter string `json:"genius_quarter" gorm:"column:genius_quarter;type:varchar(50)"`

	// 三个时间字段 - 使用字符串格式
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime;type:timestamp"` // 精确到秒，使用time.Time类型
	CreateDate  string    `json:"create_date" gorm:"column:create_date;type:date"`                     // 精确到日，格式：2025-12-13
	CreateMonth string    `json:"create_month" gorm:"column:create_month;type:varchar(7)"`             // 精确到月，格式：2025-12
}

// TableName 指定表名
func (Operators) TableName() string {
	return "operators"
}

// 三、数据库连接
func ConnectDB(config models.Config) (*gorm.DB, error) {

	// MySQL 连接字符串
	// 格式: "username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	data := config.Database

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别 (静默)
	}

	db, err := gorm.Open(mysql.Open(data.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(data.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(data.MaxOpenConns) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)      // 连接最大生命周期

	// 自动迁移表结构（如果表不存在则创建）
	// err = db.AutoMigrate(&Operators{})
	// if err != nil {
	// 	return nil, fmt.Errorf("自动迁移表结构失败: %v", err)
	// }

	fmt.Println("数据库连接成功!")
	return db, nil
}

// 四、将 Operator 转换为 Operators 并保存到数据库
func SaveOperators(db *gorm.DB, operators []models.Operator, geniusLevel, geniusQuarter string) error {
	var dbOperators []Operators

	for _, op := range operators {
		// 将 Scope 数组转换为 JSON 字符串
		scopeJSON, err := json.Marshal(op.Scope)
		if err != nil {
			return fmt.Errorf("序列化 Scope 失败: %v", err)
		}

		// 处理指针类型的字段
		documentation := ""
		if op.Documentation != nil {
			documentation = *op.Documentation
		}

		level := ""
		if op.Level != nil {
			level = *op.Level
		}

		now := time.Now()

		// 创建数据库模型
		dbOp := Operators{
			Name:          op.Name,                  // 操作符名
			Category:      op.Category,              // 操作符类别
			Scope:         string(scopeJSON),        // 可用类别
			Definition:    op.Definition,            // 操作符定义
			EnDescription: op.Description,           // 英文描述
			CnDescription: "",                       // 中文描述
			Documentation: documentation,            // 文件文档
			Level:         level,                    // 水平
			GeniusLevel:   geniusLevel,              // Genius级别
			GeniusQuarter: geniusQuarter,            // Genius季度
			CreateTime:    now,                      // 精确到秒
			CreateDate:    now.Format("2006-01-02"), // 精确到日
			CreateMonth:   now.Format("2006-01"),    // 精确到月
		}

		dbOperators = append(dbOperators, dbOp)
	}

	// 批量插入数据，使用 CreateInBatches 分批插入，每批100条
	result := db.CreateInBatches(dbOperators, 100)
	if result.Error != nil {
		return fmt.Errorf("插入数据失败: %v", result.Error)
	}

	fmt.Printf("成功插入 %d 条操作符记录\n", result.RowsAffected)
	return nil
}

// 五、检查数据库中是否已有数据
func CheckDataExists(db *gorm.DB) (bool, error) {

	var count int64
	// result := db.Model(&Operators{}).Count(&count)
	// if result.Error != nil {
	// 	return false, fmt.Errorf("查询数据失败: %v", result.Error)
	// }

	return count > 0, nil
}

// 六、清空表数据
func ClearTable(db *gorm.DB) error {
	// 使用 Exec 执行原生 SQL 清空表
	result := db.Exec("DELETE FROM operators")
	if result.Error != nil {
		return fmt.Errorf("清空表数据失败: %v", result.Error)
	}

	fmt.Printf("已清空表，删除 %d 条记录\n", result.RowsAffected)
	return nil
}

// 七、验证和获取Genius等级
func getGeniusLevel() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	validLevels := []string{"Gold", "Expert", "Master", "Grand Master"}

	for {
		// 步骤1: 输入Genius等级
		fmt.Println("\n" + strings.Repeat("-", 50))
		fmt.Println("请输入Genius等级 (必填), 可选值: Gold, Expert, Master, Grand Master")
		fmt.Print("请输入: ")

		if !scanner.Scan() {
			return "", fmt.Errorf("读取输入失败")
		}

		geniusLevel := strings.TrimSpace(scanner.Text())

		// 验证是否为空
		if geniusLevel == "" {
			fmt.Println("❌ 错误: Genius等级不能为空，请重新输入")
			continue
		}

		// 验证是否有效 - 严格大小写匹配
		isValid := false
		for _, validLevel := range validLevels {
			if geniusLevel == validLevel {  // 严格相等，包括大小写
				isValid = true
				break
			}
		}

		if !isValid {
			fmt.Printf("❌ 错误: '%s' 不是有效的Genius等级，请选择: Gold, Expert, Master, Grand Master", geniusLevel)
			fmt.Println("注意: 必须完全匹配大小写")
			continue
		}

		// 步骤2: 确认输入
		fmt.Printf("\n您输入的Genius等级是: %s\n", geniusLevel)
		fmt.Print("确认吗? (y/n): ")

		if !scanner.Scan() {
			return "", fmt.Errorf("读取确认输入失败")
		}

		confirm := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if confirm == "y" || confirm == "yes" || confirm == "是" {
			fmt.Println("✅ Genius等级设置完成")
			return geniusLevel, nil
		} else if confirm == "n" || confirm == "no" || confirm == "否" {
			fmt.Println("重新输入Genius等级...")
			continue
		} else {
			fmt.Println("❌ 无效的确认输入，请输入 y/n 或 是/否")
			continue
		}
	}
}

// 八、验证和获取Genius季度
func getGeniusQuarter() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	quarterRegex := regexp.MustCompile(`^\d{4}-Q[1-4]$`)

	// 获取当前季度作为参考
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentQuarter := (currentMonth-1)/3 + 1
	suggestedQuarter := fmt.Sprintf("%d-Q%d", currentYear, currentQuarter)

	for {
		// 步骤1: 输入Genius季度
		fmt.Println("\n" + strings.Repeat("-", 50))
		fmt.Println("请输入Genius季度 (必填):")
		fmt.Printf("格式: YYYY-Q[1-4] (例如: %s)\n", suggestedQuarter)
		fmt.Print("请输入: ")

		if !scanner.Scan() {
			return "", fmt.Errorf("读取输入失败")
		}

		geniusQuarter := strings.TrimSpace(scanner.Text())

		// 验证是否为空
		if geniusQuarter == "" {
			fmt.Println("❌ 错误: Genius季度不能为空，请重新输入")
			continue
		}

		// 验证格式
		if !quarterRegex.MatchString(geniusQuarter) {
			fmt.Println("❌ 错误: 季度格式不正确")
			fmt.Println("正确格式: YYYY-Q[1-4] (例如: 2025-Q1, 2025-Q2, 2025-Q3, 2025-Q4)")
			continue
		}

		// 提取年份和季度
		parts := strings.Split(geniusQuarter, "-Q")
		if len(parts) != 2 {
			fmt.Println("❌ 错误: 季度格式解析失败")
			continue
		}

		yearStr := parts[0]
		quarterStr := parts[1]

		// 验证年份
		if len(yearStr) != 4 {
			fmt.Println("❌ 错误: 年份必须是4位数字")
			continue
		}

		// 检查年份是否合理（2000-2100）
		year := 0
		if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil || year < 2000 || year > 2100 {
			fmt.Println("❌ 错误: 年份必须在 2000 到 2100 之间")
			continue
		}

		// 检查季度是否合理
		quarter := 0
		if _, err := fmt.Sscanf(quarterStr, "%d", &quarter); err != nil || quarter < 1 || quarter > 4 {
			fmt.Println("❌ 错误: 季度必须是 1-4 之间的数字")
			continue
		}

		// 步骤2: 确认输入
		fmt.Printf("\n您输入的Genius季度是: %s\n", geniusQuarter)
		fmt.Print("确认吗? (y/n): ")

		if !scanner.Scan() {
			return "", fmt.Errorf("读取确认输入失败")
		}

		confirm := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if confirm == "y" || confirm == "yes" || confirm == "是" {
			fmt.Println("✅ Genius季度设置完成")
			return geniusQuarter, nil
		} else if confirm == "n" || confirm == "no" || confirm == "否" {
			fmt.Println("重新输入Genius季度...")
			continue
		} else {
			fmt.Println("❌ 无效的确认输入，请输入 y/n 或 是/否")
			continue
		}
	}
}

// ------------------------------------------------ 更新或加载新赛季操作符 -----------------------------------------------

func UpdateOperators(config models.Config, token string) {
	fmt.Println("\n====================== 执行更新操作符 ======================")

	// 1. 连接数据库
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 2. 检查是否已有数据
	hasData, err := CheckDataExists(db)
	if err != nil {
		log.Fatal("检查数据失败:", err)
	}

	if hasData {
		fmt.Print("数据库已有数据，是否清空重新导入？(y/n): ")
		var answer string
		fmt.Scanln(&answer)

		if strings.ToLower(answer) == "y" {
			err = ClearTable(db)
			if err != nil {
				log.Fatal("清空表失败:", err)
			}
		} else {
			fmt.Println("已取消操作")
			return
		}
	}

	// 3. 获取操作符列表
	// fmt.Println("正在获取操作符列表...")
	allOperators, err := FetchOperators(config, token)
	if err != nil {
		log.Fatal("获取操作符失败:", err)
	}

	fmt.Printf("成功获取 %d 个操作符\n", len(allOperators))

	// 4. 打印前3个操作符信息
	// fmt.Println("\n前3个操作符信息:")
	// for i, op := range allOperators {
	// 	if i >= 3 {
	// 		break
	// 	}
	// 	fmt.Printf("%d. %s (%s)\n", i+1, op.Name, op.Category)
	// }

	// 5. 获取Genius等级和Genius季度（必填，带验证和确认）
    geniusLevel, err := getGeniusLevel()
    if err != nil {
        log.Fatalf("获取Genius等级失败: %v", err)
    }
    geniusQuarter, err := getGeniusQuarter()
    if err != nil {
        log.Fatalf("获取Genius季度失败: %v", err)
    }
    
    // 6. 显示最终配置确认
    fmt.Println("\n" + strings.Repeat("=", 60))
    fmt.Println("                  最终配置确认")
    fmt.Println(strings.Repeat("=", 60))
    fmt.Printf("Genius等级: %s\n", geniusLevel)
    fmt.Printf("Genius季度: %s\n", geniusQuarter)
    fmt.Println(strings.Repeat("-", 60))
    
    // 7. 最终确认
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("\n确认使用以上配置保存到数据库吗? (y/n): ")
        
        if !scanner.Scan() {
            log.Fatal("读取最终确认失败")
        }
        
        finalConfirm := strings.TrimSpace(strings.ToLower(scanner.Text()))
        
        if finalConfirm == "y" || finalConfirm == "yes" || finalConfirm == "是" {
            break
        } else if finalConfirm == "n" || finalConfirm == "no" || finalConfirm == "否" {
            fmt.Println("❌ 操作已取消")
            return
        } else {
            fmt.Println("❌ 无效输入，请输入 y/n 或 是/否")
        }
    }

	// 8. 保存到数据库
	fmt.Println("\n正在保存到数据库...")
	err = SaveOperators(db, allOperators, geniusLevel, geniusQuarter)
	if err != nil {
		log.Fatal("保存到数据库失败:", err)
	}

	fmt.Println("操作完成!")

	fmt.Println("更新操作符功能正在执行...")
	// TODO: 实现实际的更新操作符逻辑
	fmt.Println("更新操作符完成！")
}
