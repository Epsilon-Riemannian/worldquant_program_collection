package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"program-collection/models"
	sp "program-collection/small_program"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// 1. 加载配置文件
func loadConfig() models.Config {
	file, err := os.Open("configs/config.yaml")
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	var config models.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}
	return config
}

// 2. 登录并获取 token
func globalSignIn(config models.Config) string {
	// 用户登录
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", config.Third.Addr, config.Paths.Auth), nil)
	req.SetBasicAuth(config.Login.Username, config.Login.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 201 {
		log.Fatalf("登录失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Fatalf("解析响应失败: %v", err)
	}

	// 打印结果
	fmt.Println("Login to BRAIN successfully.")
	data, err := json.Marshal(authResp)
	if err != nil {
		log.Fatalf("序列化 JSON 失败: %v", err)
	}
	fmt.Println(string(data))

	// 从 Set-Cookie 中查找名为 "t" 的 token
	var token string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "t" {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		log.Fatal("未从 Set-Cookie 中获取到名为 t 的 token")
	}

	return strings.TrimPrefix(token, "Bearer ")
}

// 3. 显示菜单
func showMenu() {
	fmt.Println("\n============================================")
	fmt.Println("         程序集合控制中心")
	fmt.Println("============================================")
	fmt.Println("1. 字段使用情况检查 (FieldCheck)")
	fmt.Println("2. 相似度检测 (ProdCorrCheck)")
	fmt.Println("3. 更新操作符 (UpdateOperators)")
	fmt.Println("4. 阿尔法管理 (RunActiveAlphaManagement)")
	fmt.Println("5. 权重|因子价值差分 (SaveWeightValueFactor)")
	fmt.Println("6. 优先推金字塔 (PyramidAlphaInfo)")
	fmt.Println("7. 运行所有程序")
	fmt.Println("8. 自定义选择多个程序")
	fmt.Println("0. 退出")
	fmt.Println("============================================")
	fmt.Print("请选择要执行的操作 (0-8): ")
}

// 4. 获取用户输入
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// 5. 运行所有程序
func runAllPrograms(config models.Config, token string) {
	fmt.Println("\n>>>>>>>>>>>>>>>> 开始执行所有程序 <<<<<<<<<<<<<<<<")

	sp.FieldCheck(config, token)
	fmt.Println()

	sp.ProdCorrCheck(config, token)
	fmt.Println()

	sp.UpdateOperators(config, token)
	fmt.Println()

	sp.RunActiveAlphaManagement(config, token)
	fmt.Println()

	sp.SaveWeightValueFactor(config, token)
	fmt.Println()

	sp.PyramidAlphaInfo(config, token)

	fmt.Println("\n>>>>>>>>>>>>>>>> 所有程序执行完毕 <<<<<<<<<<<<<<<<")
}

// 6. 运行选择的程序
func runSelectedPrograms(config models.Config, token string, selections []int) {
	for i, selection := range selections {
		switch selection {
		case 1:
			sp.FieldCheck(config, token)
		case 2:
			sp.ProdCorrCheck(config, token)
		case 3:
			sp.UpdateOperators(config, token)
		case 4:
			sp.RunActiveAlphaManagement(config, token)
		case 5:
			sp.SaveWeightValueFactor(config, token)
		case 6:
			sp.PyramidAlphaInfo(config, token)
		}
		if i != len(selections)-1 {
			fmt.Println() // 在程序之间添加空行
		}
	}
}

// 7. 获取多个选择
func getMultipleSelections() []int {
	fmt.Println("\n请选择要运行的程序（输入数字，用空格分隔）:")
	fmt.Println("示例: 1 2 3 4 或 1  3")
	fmt.Print("你的选择: ")

	input := getUserInput()
	if input == "" {
		return []int{}
	}

	parts := strings.Fields(input)
	var selections []int

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 1 || num > 7 {
			fmt.Printf("无效的选择: %s，已跳过\n", part)
			continue
		}

		// 去重
		found := false
		for _, s := range selections {
			if s == num {
				found = true
				break
			}
		}

		if !found {
			selections = append(selections, num)
		}
	}

	return selections
}

// 8. 确认运行
func confirmRun(programName string) bool {
	fmt.Printf("确定要运行 %s 吗？(y/n): ", programName)
	input := strings.ToLower(getUserInput())
	return input == "y" || input == "yes" || input == "是" || input == "1"
}

// 9. 初始化数据库连接
func initDatabase(config models.Config) *gorm.DB {

	// db, err := sp.ConnectDB(config)
	// if err != nil {
	// 	log.Fatal("数据库连接失败:", err)
	// }
	// defer func() {
	// 	sqlDB, _ := db.DB()
	// 	sqlDB.Close()
	// }()
	return nil
}

// ---------------------------------------------- main-主程序 --------------------------------------------

func main() {

	// 加载配置
	// fmt.Println("正在加载配置文件...")
	config := loadConfig()
	// fmt.Println("配置文件加载成功！")

	// 初始化数据库
	// db := initDatabase(config)

	// 登录获取token
	// fmt.Println("\n正在登录获取token...")
	token := globalSignIn(config)
	// fmt.Printf("Token 获取成功！\n")

	// 主循环
	for {
		showMenu()
		choice := getUserInput()

		switch choice {
		case "0":
			fmt.Println("感谢使用，再见！")
			return

		case "1":
			if confirmRun("字段使用情况检查 (FieldCheck)") {
				sp.FieldCheck(config, token)
			}

		case "2":
			if confirmRun("相似度检测 (ProdCorrCheck)") {
				sp.ProdCorrCheck(config, token)
			}

		case "3":
			if confirmRun("更新操作符 (UpdateOperators)") {
				sp.UpdateOperators(config, token)
			}

		case "4":
			if confirmRun("阿尔法管理 (RunActiveAlphaManagement)") {
				sp.RunActiveAlphaManagement(config, token)
			}

		case "5":
			if confirmRun("权重|因子价值差分 (SaveWeightValueFactor)") {
				sp.SaveWeightValueFactor(config, token)
			}

		case "6":
			if confirmRun("优先推金字塔 (PyramidAlphaInfo)") {
				sp.PyramidAlphaInfo(config, token)
			}

		case "7":
			if confirmRun("所有程序") {
				runAllPrograms(config, token)
			}

		case "8":
			selections := getMultipleSelections()
			if len(selections) == 0 {
				fmt.Println("未选择任何程序，返回菜单。")
				continue
			}

			fmt.Println("\n你选择了以下程序:")
			for _, s := range selections {
				switch s {
				case 1:
					fmt.Println("  - 字段使用情况检查 (FieldCheck)")
				case 2:
					fmt.Println("  - 相似度检测 (ProdCorrCheck)")
				case 3:
					fmt.Println("  - 更新操作符 (UpdateOperators)")
				case 4:
					fmt.Println("  - 阿尔法管理 (RunActiveAlphaManagement)")
				case 5:
					fmt.Println("  - 权重|因子价值差分 (SaveWeightValueFactor)")
				}
			}

			if confirmRun("以上程序") {
				runSelectedPrograms(config, token, selections)
			}

		default:
			fmt.Println("无效的选择，请输入 0-7 之间的数字！")
		}

		// 询问是否继续
		fmt.Print("\n是否继续运行其他程序？(y/n): ")
		cont := strings.ToLower(getUserInput())
		if cont != "y" && cont != "yes" && cont != "是" && cont != "1" {
			fmt.Println("感谢使用，再见！")
			break
		}
	}
}
