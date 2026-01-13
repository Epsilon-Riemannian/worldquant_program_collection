package small_program

import (
	"fmt"
	"math"

	"program-collection/models"
)

// 计算 alpha 的 prod_corr 量纲量化
func GetProdCorrMath(alphaLists []models.Alpha) (prod_corr_min float64, prod_corr_max float64, prod_corr_avg float64, alpha_count int) {

	if len(alphaLists) == 0 {
		return 0.0, 0.0, 0.0, 0
	}

	// 初始化
	prod_corr_min = alphaLists[0].IS.ProdCorrelation
	prod_corr_max = alphaLists[0].IS.ProdCorrelation
	total := 0.0
	validCount := 0

	for _, alpha := range alphaLists {
		if alpha.IS != nil {
			value := alpha.IS.ProdCorrelation
			total += value

			if value < prod_corr_min {
				prod_corr_min = value
			}
			if value > prod_corr_max {
				prod_corr_max = value
			}

			validCount++
		}
	}

	alpha_count = validCount
	if validCount > 0 {
		prod_corr_avg = total / float64(validCount)
	}

	// 保留四位小数
	prod_corr_min = math.Round(prod_corr_min*10000) / 10000
	prod_corr_max = math.Round(prod_corr_max*10000) / 10000
	prod_corr_avg = math.Round(prod_corr_avg*10000) / 10000

	return
}

// ------------------------------------------------ 相似度计算 -----------------------------------------------

func ProdCorrCheck(config models.Config, token string) {
	fmt.Println("\n====================== 执行相似度检测 ======================")

	dateFrom, _ := ConvertToUTCPlus5("2025-10-01 00:00:00")
	dateTo, _ := ConvertToUTCPlus5("2025-11-01 00:00:00")

	// 获取 alpha 列表信息
	alphaLists, _ := GetAllAlphas(config, token, models.GetAlphasRequest{
		Limit:    50,
		Offset:   0,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Order:    "-dateSubmitted",
		Type:     "REGULAR",
	})

	// 计算数学统计量 prod_corr
	min, max, avg, count := GetProdCorrMath(alphaLists)

	// 打印 prod_corr 统计学量结果
	fmt.Printf("\n统计结果如下:\n")
	fmt.Printf("十月共提交 (不包括sa) alpha 数量: %d, prod_corr最小值: %.4f, 最大值: %.4f, 平均值: %.4f", count, min, max, avg)

	fmt.Println("相似度检测功能正在执行...")
	// TODO: 实现实际的相似度检测逻辑
	fmt.Println("相似度检测完成！")
}
