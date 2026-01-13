package small_program

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"program-collection/models"

	"github.com/samber/lo"
)

// ----------------------------------------------- 处理字符串获取字段 ----------------------------------------------

// 定义操作符的参数规则
type ParamRule struct {
	FieldPositions []int // 哪些位置包含字段（从0开始）
	IgnoreNamed    bool  // 是否忽略命名参数
}

// 移除注释但保留代码
func removeComments(input string) string {
	var result strings.Builder
	inMultiLineComment := false
	inSingleLineComment := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if inMultiLineComment {
			if char == '*' && i+1 < len(input) && input[i+1] == '/' {
				inMultiLineComment = false
				i++ // 跳过 '/'
			}
			continue
		}

		if inSingleLineComment {
			if char == '\n' {
				inSingleLineComment = false
				result.WriteByte(char)
			}
			continue
		}

		// 检查是否开始多行注释
		if char == '/' && i+1 < len(input) && input[i+1] == '*' {
			inMultiLineComment = true
			i++ // 跳过 '*'
			continue
		}

		// 检查是否开始单行注释
		if char == '/' && i+1 < len(input) && input[i+1] == '/' {
			inSingleLineComment = true
			i++ // 跳过第二个 '/'
			continue
		}

		result.WriteByte(char)
	}

	// 清理结果
	cleaned := result.String()
	cleaned = strings.ReplaceAll(cleaned, "\r", " ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	return strings.TrimSpace(cleaned)
}

func extractFields(input string, allOperators []string, functionRules map[string]ParamRule) []string {
	// 首先移除所有注释
	cleanedInput := removeComments(input)

	// 提取所有定义的变量
	definedVars := extractDefinedVariables(cleanedInput, allOperators)

	// 移除赋值语句的左边部分并获取处理后的表达式
	processedInput, rightSides := removeAssignmentLeft(cleanedInput)

	uniqueFields := make(map[string]bool)

	// 递归提取所有字段
	fields := extractFieldsRecursive(processedInput, allOperators, functionRules)
	for _, field := range fields {
		// 过滤掉函数名和定义的变量
		if !isFunctionName(field, allOperators) && !definedVars[field] {
			uniqueFields[field] = true
		}
	}

	// 另外从等号右边部分提取字段（确保不会漏掉）
	for _, expr := range rightSides {
		exprFields := extractFieldsRecursive(expr, allOperators, functionRules)
		for _, field := range exprFields {
			if !isFunctionName(field, allOperators) && !definedVars[field] {
				uniqueFields[field] = true
			}
		}
	}

	// 转换为切片
	result := make([]string, 0, len(uniqueFields))
	for field := range uniqueFields {
		result = append(result, field)
	}

	return result
}

// 提取所有定义的变量（等号左边的标识符）
func extractDefinedVariables(input string, allOperators []string) map[string]bool {
	definedVars := make(map[string]bool)

	// 按分号分割多个表达式
	expressions := strings.Split(input, ";")

	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}

		// 查找第一个不在括号内的等号
		parenDepth := 0
		for i, char := range expr {
			switch char {
			case '(':
				parenDepth++
			case ')':
				parenDepth--
			case '=':
				if parenDepth == 0 {
					// 提取等号左边的变量
					leftSide := strings.TrimSpace(expr[:i])
					if leftSide != "" {
						// 左边可能是一个变量或多个变量，这里简单处理为单个变量
						// 使用正则提取有效的变量名
						varPattern := `[a-zA-Z_][a-zA-Z0-9_]*`
						varRe := regexp.MustCompile(varPattern)
						matches := varRe.FindAllString(leftSide, -1)

						for _, match := range matches {
							if !isFunctionName(match, allOperators) {
								definedVars[match] = true
							}
						}
					}
				}
			}
		}
	}

	return definedVars
}

// 移除赋值语句的左边部分，返回处理后的表达式和等号右边的部分
func removeAssignmentLeft(input string) (string, []string) {
	// 按分号分割多个表达式
	expressions := strings.Split(input, ";")
	var result []string
	var rightSides []string

	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}

		// 查找第一个不在括号内的等号
		parenDepth := 0
		foundEqual := false
		for i, char := range expr {
			switch char {
			case '(':
				parenDepth++
			case ')':
				parenDepth--
			case '=':
				if parenDepth == 0 {
					// 返回等号右边的部分
					rightSide := strings.TrimSpace(expr[i+1:])
					if rightSide != "" {
						result = append(result, rightSide)
						rightSides = append(rightSides, rightSide)
					}
					foundEqual = true
				}
			}
			if foundEqual {
				break
			}
		}

		// 如果没有等号，保留整个表达式
		if !foundEqual {
			result = append(result, expr)
			rightSides = append(rightSides, expr)
		}
	}

	return strings.Join(result, "; "), rightSides
}

// 递归提取字段（带深度限制）
func extractFieldsRecursive(expr string, allOperators []string, functionRules map[string]ParamRule) []string {
	return extractFieldsRecursiveWithDepth(expr, allOperators, functionRules, 0)
}

func extractFieldsRecursiveWithDepth(expr string, allOperators []string, functionRules map[string]ParamRule, depth int) []string {
	// 添加递归深度限制
	if depth > 50 {
		// fmt.Printf("警告: 递归深度超过限制，表达式: %s\n", expr)
		return extractSimpleFields(expr, allOperators)
	}

	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}

	var fields []string

	// 先检查是否包含比较或逻辑运算符
	if hasComparisonOrLogicalOps(expr) {
		// 分割表达式并递归处理每个部分
		parts := splitExpression(expr)
		for _, part := range parts {
			partFields := extractFieldsRecursiveWithDepth(part, allOperators, functionRules, depth+1)
			fields = append(fields, partFields...)
		}
		return fields
	}

	// 修复：使用更强大的函数调用匹配，处理嵌套函数
	funcPattern := `([a-zA-Z_][a-zA-Z0-9_]*)\((.*)\)`
	funcRe := regexp.MustCompile(funcPattern)

	// 尝试匹配最外层的函数调用
	remainingExpr := expr
	for {
		funcMatch := funcRe.FindStringSubmatch(remainingExpr)
		if funcMatch == nil {
			break
		}

		funcName := funcMatch[1]
		paramsStr := funcMatch[2]

		// 检查参数是否平衡（括号匹配）
		if !isBalancedParentheses(paramsStr) {
			// 如果括号不平衡，跳过这个函数匹配
			break
		}

		if rule, exists := functionRules[funcName]; exists {
			// 根据规则处理参数
			paramFields := processParametersByRule(paramsStr, rule, allOperators, functionRules)
			fields = append(fields, paramFields...)
		} else if isFunctionName(funcName, allOperators) {
			// 默认规则：递归处理第一个参数
			params := parseParameters(paramsStr)
			if len(params) > 0 {
				firstParamFields := extractFieldsRecursiveWithDepth(params[0], allOperators, functionRules, depth+1)
				fields = append(fields, firstParamFields...)
			}
		}

		// 移除已处理的函数调用，继续处理剩余部分
		remainingExpr = strings.Replace(remainingExpr, funcMatch[0], "", 1)
	}

	// 如果没有函数调用，提取简单字段和表达式中的字段
	if len(fields) == 0 {
		fields = extractSimpleFields(expr, allOperators)
	}

	return fields
}

// 检查括号是否平衡
func isBalancedParentheses(s string) bool {
	count := 0
	for _, char := range s {
		switch char {
		case '(':
			count++
		case ')':
			count--
			if count < 0 {
				return false
			}
		}
	}
	return count == 0
}

// 检查是否包含比较或逻辑运算符
func hasComparisonOrLogicalOps(expr string) bool {
	operators := []string{">", "<", ">=", "<=", "==", "!=", "&&", "||"}
	for _, op := range operators {
		if strings.Contains(expr, op) {
			return true
		}
	}
	return false
}

// 分割表达式（处理比较和逻辑运算符）
func splitExpression(expr string) []string {
	var parts []string
	var current strings.Builder
	parenDepth := 0

	for i := 0; i < len(expr); i++ {
		char := expr[i]

		switch char {
		case '(':
			parenDepth++
			current.WriteByte(char)
		case ')':
			parenDepth--
			current.WriteByte(char)
		case '>', '<', '!', '&', '|':
			if parenDepth == 0 {
				// 检查是否是复合运算符
				if i+1 < len(expr) {
					nextChar := expr[i+1]
					compoundOps := []string{">=", "<=", "==", "!=", "&&", "||"}
					compound := string(char) + string(nextChar)

					// 检查是否是逻辑运算符 and, or 的一部分
					if i+2 < len(expr) {
						thirdChar := expr[i+2]
						if (char == 'a' && nextChar == 'n' && thirdChar == 'd') ||
							(char == 'o' && nextChar == 'r') {
							current.WriteByte(char)
							continue
						}
					}

					foundCompound := false
					for _, op := range compoundOps {
						if op == compound {
							if current.Len() > 0 {
								parts = append(parts, strings.TrimSpace(current.String()))
								current.Reset()
							}
							i++ // 跳过下一个字符
							foundCompound = true
							break
						}
					}
					if foundCompound {
						continue
					}
				}

				// 单个运算符
				if current.Len() > 0 {
					parts = append(parts, strings.TrimSpace(current.String()))
					current.Reset()
				}
			} else {
				current.WriteByte(char)
			}
		default:
			current.WriteByte(char)
		}
	}

	// 添加最后一部分
	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// 根据规则处理参数
func processParametersByRule(paramsStr string, rule ParamRule, allOperators []string, functionRules map[string]ParamRule) []string {
	params := parseParameters(paramsStr)
	var allFields []string

	for _, pos := range rule.FieldPositions {
		if pos < len(params) {
			param := strings.TrimSpace(params[pos])
			// 只有 IgnoreNamed 为 true 时才跳过命名参数
			if rule.IgnoreNamed && strings.Contains(param, "=") {
				continue
			}
			// 递归处理参数（可能包含嵌套函数或简单字段）
			paramFields := extractFieldsRecursive(param, allOperators, functionRules)
			allFields = append(allFields, paramFields...)
		}
	}

	return allFields
}

// 提取简单表达式中的字段（过滤掉命名参数）
func extractSimpleFields(expr string, allOperators []string) []string {
	var fields []string

	// 使用正则匹配所有标识符
	fieldPattern := `[a-zA-Z_][a-zA-Z0-9_]*`
	fieldRe := regexp.MustCompile(fieldPattern)
	matches := fieldRe.FindAllString(expr, -1)

	for _, match := range matches {
		// 过滤掉函数名和命名参数
		if !isFunctionName(match, allOperators) && !isNamedParameter(expr, match) {
			fields = append(fields, match)
		}
	}

	return fields
}

// 检查是否是命名参数（如 std=4.0 中的 std）
func isNamedParameter(expr, identifier string) bool {
	// 查找标识符后面是否有等号
	pattern := regexp.QuoteMeta(identifier) + `\s*=`
	re := regexp.MustCompile(pattern)
	return re.MatchString(expr)
}

// 解析参数字符串，处理嵌套函数
func parseParameters(paramsStr string) []string {
	var params []string
	var currentParam strings.Builder
	parenDepth := 0
	bracketDepth := 0
	braceDepth := 0

	for i, char := range paramsStr {
		switch char {
		case '(':
			parenDepth++
			currentParam.WriteRune(char)
		case ')':
			parenDepth--
			currentParam.WriteRune(char)
		case '[':
			bracketDepth++
			currentParam.WriteRune(char)
		case ']':
			bracketDepth--
			currentParam.WriteRune(char)
		case '{':
			braceDepth++
			currentParam.WriteRune(char)
		case '}':
			braceDepth--
			currentParam.WriteRune(char)
		case ',':
			if parenDepth == 0 && bracketDepth == 0 && braceDepth == 0 {
				param := strings.TrimSpace(currentParam.String())
				if param != "" {
					params = append(params, param)
				}
				currentParam.Reset()
			} else {
				currentParam.WriteRune(char)
			}
		default:
			currentParam.WriteRune(char)
		}

		// 如果是最后一个字符，添加当前参数
		if i == len(paramsStr)-1 {
			param := strings.TrimSpace(currentParam.String())
			if param != "" {
				params = append(params, param)
			}
		}
	}

	return params
}

// 判断是否为函数名
func isFunctionName(s string, operatorNmae []string) bool {
	return lo.Contains(operatorNmae, s)
}

// ------------------------------------------------------------------- 获取操作符参数位置 -------------------------------------------------------------------

// GenerateFunctionRules 从操作符列表生成函数规则映射
func GenerateFunctionRules(allOperators []models.Operator, specialOperatorNames []string) map[string]ParamRule {
	functionRules := make(map[string]ParamRule, len(allOperators))

	for _, operator := range allOperators {
		// 解析definition获取参数位置
		fieldPositions := parseDefinition(operator.Definition)

		// 如果解析失败或者需要特殊处理的操作符，跳过此操作符
		if len(fieldPositions) == 0 && lo.Contains(specialOperatorNames, operator.Name) {
			continue
		}

		functionRules[operator.Name] = ParamRule{
			FieldPositions: fieldPositions,
			IgnoreNamed:    true, // 默认忽略命名参数
		}
	}

	return functionRules
}

// 解析definition字符串，提取x,y,z等单个字母参数、input参数和alpha参数的位置
func parseDefinition(definition string) []int {
	// 查找第一个括号内的内容
	start := strings.Index(definition, "(")
	end := strings.Index(definition, ")")

	if start == -1 || end == -1 || start >= end {
		return []int{}
	}

	// 提取括号内的参数列表
	paramsStr := definition[start+1 : end]

	// 分割参数
	params := strings.Split(paramsStr, ",")

	var positions []int
	position := 0

	for _, param := range params {
		param = strings.TrimSpace(param)

		// 跳过空参数和省略号
		if param == "" || param == "..." {
			continue
		}

		// 移除默认值部分
		if equalsIndex := strings.Index(param, "="); equalsIndex != -1 {
			param = strings.TrimSpace(param[:equalsIndex])
		}

		// 识别单个字母参数：x, y, z 等、input 参数和 alpha 参数
		if isFieldParam(param) {
			positions = append(positions, position)
		}
		position++
	}

	return positions
}

// 判断是否是字段参数（只有x,y,z单个字母、input、alpha）
func isFieldParam(param string) bool {
	// 单个字母参数：只有 x, y, z
	if len(param) == 1 && (param == "x" || param == "y" || param == "z") {
		return true
	}

	// input 参数（input、input 2、input1等）
	if strings.HasPrefix(param, "input") {
		// 如果是单纯的 "input"
		if param == "input" {
			return true
		}

		// 处理 "input 2" 这种带空格的情况
		if len(param) > 5 && param[5] == ' ' {
			// 检查空格后面的内容是否是数字
			_, err := strconv.Atoi(strings.TrimSpace(param[5:]))
			return err == nil
		}

		// 处理 "input2" 这种不带空格的情况
		if len(param) > 5 {
			_, err := strconv.Atoi(param[5:])
			return err == nil
		}
	}

	// alpha 参数
	if param == "alpha" {
		return true
	}

	return false
}

// ------------------------------------------------- 字段使用情况检测 ----------------------------------------------

func GetFieldData(config models.Config, token string) ([]string, map[string]ParamRule, map[string][]string) {

	// 获取操作符列表
	allOperators, err := FetchOperators(config, token)
	if err != nil {
		log.Fatal("获取操作符失败:", err)
	}

	var allOperatorName []string
	for _, operator := range allOperators {
		allOperatorName = append(allOperatorName, operator.Name)
	}

	// 手动添加特殊操作符的规则
	specialRules := map[string]ParamRule{
		"greater_equal": {FieldPositions: []int{0, 1}, IgnoreNamed: true},
		"multiply":      {FieldPositions: []int{0, 1, 2, 3, 4, 5, 6}, IgnoreNamed: true},
		"max":           {FieldPositions: []int{0, 1, 2, 3, 4, 5, 6}, IgnoreNamed: true},
		"min":           {FieldPositions: []int{0, 1, 2, 3, 4, 5, 6}, IgnoreNamed: true},
	}

	// 获取特殊操作符的key放入切片
	specialOperatorNames := make([]string, 0, len(specialRules))
	for key := range specialRules {
		specialOperatorNames = append(specialOperatorNames, key)
	}

	// 生成函数规则（自动跳过特殊操作符）
	functionRules := GenerateFunctionRules(allOperators, specialOperatorNames)

	// 将特殊操作符规则合并到functionRules中
	for name, rule := range specialRules {
		functionRules[name] = rule
	}

	// fmt.Printf("\n总共生成 %d 个函数规则\n\n", len(functionRules))

	beginDate, _ := ConvertToUTCPlus5("2025-09-01")
	endDate, _ := ConvertToUTCPlus5("2025-10-01")

	// 获取 alpha 列表信息
	alphaLists, _ := GetAllAlphas(config, token, models.GetAlphasRequest{
		Limit:    50,
		Offset:   0,
		DateFrom: beginDate,
		DateTo:   endDate,
		Order:    "-dateSubmitted",
		Type:     "REGULAR",
	})

	// var alphaFields []string
	// var i int
	alphaIDFieldsMap := make(map[string][]string, len(alphaLists))

	for _, alpha := range alphaLists {
		fields := extractFields(alpha.Regular.Code, allOperatorName, functionRules)
		alphaIDFieldsMap[alpha.ID] = fields

		// fmt.Printf("表达式: %s\n", alpha.Regular.Code)
		// fmt.Printf("字段: %v, 数量: %d\n", fields, len(fields))
		// fmt.Printf("如有疑问请访问: %s/alphas/%s\n", config.Paths.Auth, alpha.ID)
		// fmt.Printf("或访问: %s/alpha/%s\n\n", config.Third.Addr, alpha.ID)
		// i += 1
		// fmt.Printf("%d\n", i)
	}

	return allOperatorName, functionRules, alphaIDFieldsMap
}

// ExtractContent 从字符串中提取内容
// 如果字符串包含 https://platform.worldquantbrain.com，则提取最后一个斜杠后的内容
// 否则，返回整个字符串
func ExtractContent(config models.Config, input string) (string, bool) {

	prefix := config.Third.Addr
	// 检查是否包含指定前缀
	if !strings.Contains(input, prefix) {
		return input, false
	}

	// 找到前缀的位置
	prefixIndex := strings.Index(input, prefix)
	if prefixIndex == -1 {
		return input, true
	}

	// 获取前缀之后的部分
	afterPrefix := input[prefixIndex+len(prefix):]

	// 处理可能的查询参数
	// 先分割掉查询参数（如果有的话）
	if questionMarkIndex := strings.Index(afterPrefix, "?"); questionMarkIndex != -1 {
		afterPrefix = afterPrefix[:questionMarkIndex]
	}

	// 分割路径部分
	parts := strings.Split(afterPrefix, "/")

	// 获取最后一个非空的部分
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i], true
		}
	}

	// 如果没有任何内容，返回空字符串或原始输入
	return input, true
}

// FindKeysForSliceElements 找出切片元素在map的值中出现的所有keys
func FindKeysForSliceElements(m map[string][]string, slice []string) []string {
	result := make([]string, 0)

	// 用于去重，避免重复添加相同的key
	keySet := make(map[string]bool)

	for _, element := range slice {
		for key, values := range m {
			if contains(values, element) && !keySet[key] {
				keySet[key] = true
				result = append(result, key)
			}
		}
	}

	return result
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// CheckAndPrintResults 统一处理检查和打印结果
func CheckAndPrintResults(config models.Config, token string, checkFields []string, alphaIDFieldsMap map[string][]string, alphaID string) {

	// 查找包含这些字段的alphaIDs
	matchingAlphaIDs := FindKeysForSliceElements(alphaIDFieldsMap, checkFields)

	if len(matchingAlphaIDs) > 0 {
		fmt.Println("\n🔍 找到相关Alpha:")
		for _, id := range matchingAlphaIDs {
			fmt.Printf("   - Alpha ID: %s\n", id)
			fmt.Printf("     如需查看详情: %s/alphas/%s\n", config.Paths.Auth, id)
			fmt.Printf("     或访问: %s/alpha/%s\n", config.Third.Addr, id)
		}
	} else if alphaID != "" {
		fmt.Printf("\n⚠️  未找到包含这些字段的其他Alpha。\n")
		fmt.Printf("   当前Alpha: %s/alphas/%s\n", config.Paths.Auth, alphaID)
	} else {
		fmt.Println("\n❌ 未找到包含这些字段的Alpha。")
	}
}

// GetUserInput 获取用户输入
func GetUserInput() string {
	fmt.Print("\n请输入（输入 'quit' 退出）: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时出错:", err)
		return ""
	}

	// 去除首尾空白字符
	return strings.TrimSpace(input)
}

// 主处理函数
func FieldCheck(config models.Config, token string) {
	fmt.Println("\n====================== 执行字段检查 ======================")
	fmt.Println("🚀 字段检查功能正在执行...")
	fmt.Println("📝 支持的输入格式:")
	fmt.Println("   1. 完整URL: https://platform.worldquantbrain.com/alpha/1Y5Nj28K")
	fmt.Println("   2. Alpha ID: 1Y5Nj28K")
	fmt.Println("   3. Alpha表达式: (rank(correlation(close, volume, 10)))")
	fmt.Println("   ------------------------------------------------------")

	for {
		input := GetUserInput()

		// 检查是否退出
		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("👋 再见！")
			break
		}

		if input == "" {
			fmt.Println("⚠️  输入不能为空，请重新输入。")
			continue
		}

		// 获取全部操作符以及分解的字段和alphaID数据
		allOperatorName, functionRules, alphaIDFieldsMap := GetFieldData(config, token)

		// 处理输入
		alphaInfo, isAlphaID := ExtractContent(config, input)

		if isAlphaID {
			// 输入是URL或Alpha ID
			fmt.Printf("🔍 检测到Alpha ID: %s\n", alphaInfo)

			// 尝试获取Alpha详情
			alpha, err := GetAlphaByID(config, token, alphaInfo)
			if err != nil {
				fmt.Printf("❌ 无法获取Alpha '%s' 的详情: %v\n", alphaInfo, err)
				fmt.Println("📝 尝试将其作为Alpha表达式处理...")

				// 作为表达式处理
				checkFields := extractFields(input, allOperatorName, functionRules)
				CheckAndPrintResults(config, token, checkFields, alphaIDFieldsMap, "")
			} else {
				// 成功获取Alpha，提取字段并检查
				checkFields := extractFields(alpha.Regular.Code, allOperatorName, functionRules)
				fmt.Printf("📊 从Alpha代码中提取到 %d 个字段\n", len(checkFields))
				CheckAndPrintResults(config, token, checkFields, alphaIDFieldsMap, alphaInfo)
			}
		} else {
			// 输入是Alpha表达式
			fmt.Println("📝 检测到Alpha表达式")
			checkFields := extractFields(input, allOperatorName, functionRules)
			fmt.Printf("📊 从表达式中提取到 %d 个字段\n", len(checkFields))
			CheckAndPrintResults(config, token, checkFields, alphaIDFieldsMap, "")
		}

		fmt.Println("\n" + strings.Repeat("-", 50))
	}

	fmt.Println("✅ 字段检查完成！")
}
