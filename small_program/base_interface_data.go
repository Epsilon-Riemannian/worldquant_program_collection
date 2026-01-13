package small_program

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"program-collection/models"
)

// 1.1 获取 alpha 列表数据
func GetAlphas(config models.Config, token string, req models.GetAlphasRequest) (*models.AlphaListResponse, error) {

	// 构建URL
	urL := fmt.Sprintf("%s%s", config.Third.Addr, config.Paths.AlphaList)

	// 构建查询参数
	params := url.Values{}
	params.Add("limit", strconv.Itoa(req.Limit))
	params.Add("offset", strconv.Itoa(req.Offset))

	if req.StatusFilter != "" {
		params.Add("status!=", req.StatusFilter)
	}

	if !req.DateFrom.IsZero() {
		params.Add("dateSubmitted>", req.DateFrom.UTC().Format("2006-01-02T15:04:05.000Z"))
	}

	if !req.DateTo.IsZero() {
		params.Add("dateSubmitted<", req.DateTo.UTC().Format("2006-01-02T15:04:05.000Z"))
	}

	if req.Order != "" {
		params.Add("order", req.Order)
	}

	if req.Type != "" {
		params.Add("type", req.Type)
	}

	// 修改这里：只有当 Hidden 参数明确设置时才添加
	// 注意：req.Hidden 应该是 *bool 类型才能区分"未设置"和"false"
	// 如果 req.Hidden 是 bool 类型，我们无法区分"未设置"和"false"
	if req.Hidden != nil {
		params.Add("hidden", strconv.FormatBool(*req.Hidden))
	}

	if len(params) > 0 {
		urL += "?" + params.Encode()
	}

	// 创建请求
	httpReq, err := http.NewRequest("GET", urL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置认证头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	// 解析JSON
	var alphaResponse models.AlphaListResponse
	if err := json.Unmarshal(body, &alphaResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &alphaResponse, nil
}

// 1.2 GetAllAlphas 分页获取 alpha
func GetAllAlphas(config models.Config, token string, req models.GetAlphasRequest) ([]models.Alpha, error) {

	var allAlphas []models.Alpha
	offset := req.Offset

	for {
		req.Offset = offset
		response, err := GetAlphas(config, token, req)
		if err != nil {
			return nil, err
		}

		allAlphas = append(allAlphas, response.Results...)

		// 检查是否还有下一页
		if response.Next == nil || *response.Next == "" {
			break
		}

		offset += req.Limit

		// 防止无限循环
		if offset >= response.Count {
			break
		}
	}

	return allAlphas, nil
}

// 1.3 按照 alpha_id 获取 alpha信息
func GetAlphaByID(config models.Config, token, alphaID string) (alpha models.Alpha, err error) {

	url := fmt.Sprintf("%s%s/%s", config.Third.Addr, config.Paths.Alpha, alphaID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Alpha{}, fmt.Errorf("create request failed: %v", err)
	}

	// 设置认证头（使用Cookie方式）
	req.Header.Set("Cookie", fmt.Sprintf("t=%s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Alpha{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Alpha{}, fmt.Errorf("read response body failed: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return models.Alpha{}, fmt.Errorf("fetch alpha failed: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var alphaInfo models.Alpha
	if err := json.Unmarshal(bodyBytes, &alphaInfo); err != nil {
		return models.Alpha{}, fmt.Errorf("decode failed: %v, response: %s", err, string(bodyBytes))
	}

	return alphaInfo, nil
}

// 2.1 获取操作符列表
func FetchOperators(config models.Config, token string) ([]models.Operator, error) {

	url := fmt.Sprintf("%s%s", config.Third.Addr, config.Paths.Operator)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// 设置认证头（使用Cookie方式）
	req.Header.Set("Cookie", fmt.Sprintf("t=%s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch operators failed: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var operators []models.Operator
	if err := json.Unmarshal(bodyBytes, &operators); err != nil {
		return nil, fmt.Errorf("decode failed: %v, response: %s", err, string(bodyBytes))
	}

	return operators, nil
}

// 3.1 获取研究顾问的 Weight | Value_factor信息
func FetchConsultant(config models.Config, token string) (*models.ConsultantResponse, error) {

	url := fmt.Sprintf("%s%s", config.Third.Addr, config.Paths.Consultant)

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// 设置认证头（使用Cookie方式）
	req.Header.Set("Cookie", fmt.Sprintf("t=%s", token))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch consultant failed: status %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// 解析JSON响应
	var consultantResp models.ConsultantResponse
	if err := json.Unmarshal(bodyBytes, &consultantResp); err != nil {
		return nil, fmt.Errorf("decode failed: %v, response: %s", err, string(bodyBytes))
	}

	return &consultantResp, nil
}

// 4.1 PyramidInfo点金字塔内容数据
func PyramidInfo(config models.Config, token string, beginDate, endDate string) ([]models.Pyramids, error) {
	url := fmt.Sprintf("%s/users/self/activities/pyramid-alphas?startDate=%s&endDate=%s", config.Third.Addr, beginDate, endDate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch pyramid failed: %s\n%s", resp.Status, string(bodyBytes))
	}

	var alphaInfo models.PyramidsResponse
	if err := json.NewDecoder(resp.Body).Decode(&alphaInfo); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}

	return alphaInfo.Pyramids, nil
}

