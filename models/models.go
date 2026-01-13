package models

import "time"

// -------------------------------------- 配置文件结构体 -------------------------------------- //
type Config struct {
	Login    Login    `yaml:"login"`
	Third    Third    `yaml:"third"`
	Paths    Paths    `yaml:"path"`
	Database Database `yaml:"database"`
}

type Third struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

type Login struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Paths struct {
	Auth       string `yaml:"auth"`
	Alpha      string `yaml:"alpha"`
	Pnl        string `yaml:"pnl"`
	AlphaList  string `yaml:"alphaList"`
	Operator   string `yaml:"operator"`
	Consultant string `yaml:"consultant"`
}

type Database struct {
	DSN          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
}

// -------------------------------------- 登录返回结构体 -------------------------------------- //
type AuthResponse struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Token struct {
		Expiry float64 `json:"expiry"`
	} `json:"token"`
	Permissions []string `json:"permissions"`
}

// -------------------------- alphaListResponse 获取alpha列表的响应结构 ----------------------- //
type AlphaListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Alpha `json:"results"`
}

// Alpha 表示单个alpha的详细信息
type Alpha struct {
	ID              string           `json:"id"`
	Type            string           `json:"type"`
	Author          string           `json:"author"`
	Settings        Settings         `json:"settings"`
	Combo           *AlphaCode       `json:"combo,omitempty"`
	Selection       *AlphaCode       `json:"selection,omitempty"`
	Regular         *AlphaCode       `json:"regular,omitempty"`
	DateCreated     string           `json:"dateCreated"`
	DateSubmitted   string           `json:"dateSubmitted"`
	DateModified    string           `json:"dateModified"`
	Name            *string          `json:"name"`
	Favorite        bool             `json:"favorite"`
	Hidden          bool             `json:"hidden"`
	Color           string           `json:"color"`
	Category        *string          `json:"category"`
	Tags            []string         `json:"tags"`
	Classifications []Classification `json:"classifications"`
	Grade           *string          `json:"grade"`
	Stage           string           `json:"stage"`
	Status          string           `json:"status"`
	IS              *Performance     `json:"is,omitempty"`
	OS              *OSInfo          `json:"os,omitempty"`
	Train           *Performance     `json:"train,omitempty"`
	Test            *Performance     `json:"test,omitempty"`
	Prod            interface{}      `json:"prod"`         // 根据实际情况可能需要具体类型
	Competitions    []interface{}    `json:"competitions"` // 根据实际情况可能需要具体类型
	Themes          []Theme          `json:"themes"`
	Pyramids        []Pyramid        `json:"pyramids"`
	PyramidThemes   PyramidThemes    `json:"pyramidThemes"`
	Team            interface{}      `json:"team"`          // 根据实际情况可能需要具体类型
	OsmosisPoints   interface{}      `json:"osmosisPoints"` // 根据实际情况可能需要具体类型
}

// 设置结构体
type Settings struct {
	InstrumentType      string  `json:"instrumentType"`
	Region              string  `json:"region"`
	Universe            string  `json:"universe"`
	Delay               int     `json:"delay"`
	Decay               int     `json:"decay"`
	Neutralization      string  `json:"neutralization"`
	Truncation          float64 `json:"truncation"`
	Pasteurization      string  `json:"pasteurization"`
	UnitHandling        string  `json:"unitHandling"`
	NanHandling         string  `json:"nanHandling,omitempty"`
	SelectionHandling   string  `json:"selectionHandling,omitempty"`
	SelectionLimit      int     `json:"selectionLimit,omitempty"`
	MaxTrade            string  `json:"maxTrade"`
	Language            string  `json:"language"`
	Visualization       bool    `json:"visualization"`
	StartDate           string  `json:"startDate"`
	EndDate             string  `json:"endDate"`
	ComponentActivation string  `json:"componentActivation,omitempty"`
	TestPeriod          string  `json:"testPeriod,omitempty"`
}

// Alpha代码结构体
type AlphaCode struct {
	Code          string `json:"code"`
	Description   string `json:"description"`
	OperatorCount *int   `json:"operatorCount"`
}

// 分类结构体
type Classification struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 性能指标结构体
type Performance struct {
	PNL             int     `json:"pnl"`
	BookSize        int     `json:"bookSize"`
	LongCount       int     `json:"longCount"`
	ShortCount      int     `json:"shortCount"`
	Turnover        float64 `json:"turnover"`
	Returns         float64 `json:"returns"`
	Drawdown        float64 `json:"drawdown"`
	Margin          float64 `json:"margin"`
	Sharpe          float64 `json:"sharpe"`
	Fitness         float64 `json:"fitness"`
	StartDate       string  `json:"startDate"`
	SelfCorrelation float64 `json:"selfCorrelation,omitempty"`
	ProdCorrelation float64 `json:"prodCorrelation,omitempty"`
	Checks          []Check `json:"checks,omitempty"`
}

// 检查项结构体
type Check struct {
	Name       string      `json:"name"`
	Result     string      `json:"result"`
	Limit      interface{} `json:"limit,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Date       string      `json:"date,omitempty"`
	Year       int         `json:"year,omitempty"`
	StartDate  string      `json:"startDate,omitempty"`
	EndDate    string      `json:"endDate,omitempty"`
	Effective  int         `json:"effective,omitempty"`
	Multiplier float64     `json:"multiplier,omitempty"`
	Pyramids   []Pyramid   `json:"pyramids,omitempty"`
	Themes     []Theme     `json:"themes,omitempty"`
}

// OS信息结构体
type OSInfo struct {
	StartDate           string      `json:"startDate"`
	OsISSharpeRatio     interface{} `json:"osISSharpeRatio"`
	PreCloseSharpeRatio interface{} `json:"preCloseSharpeRatio"`
	Checks              []Check     `json:"checks"`
}

// 主题结构体
type Theme struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Multiplier float64 `json:"multiplier"`
}

// 金字塔结构体
type Pyramid struct {
	Name       string  `json:"name"`
	Multiplier float64 `json:"multiplier"`
}

// 金字塔主题结构体
type PyramidThemes struct {
	Effective int       `json:"effective"`
	Pyramids  []Pyramid `json:"pyramids"`
}

// GetAlphasRequest 获取alpha列表的请求参数
type GetAlphasRequest struct {
	Limit        int       // 最大 50
	Offset       int       //
	StatusFilter string    // 例如: "UNSUBMITTED,IS-FAIL"
	DateFrom     time.Time // 开始日期
	DateTo       time.Time // 结束日期
	Order        string    // 排序字段，如: "-dateSubmitted"
	Type         string    // alpha类型
	Hidden       *bool
}

// PnLResponse 完整的响应结构
type PnLResponse struct {
	Schema  Schema      `json:"schema"`
	Records [][]interface{} `json:"records"`
}

// Schema 定义数据结构
type Schema struct {
	Name       string     `json:"name"`
	Title      string     `json:"title"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// PnLRecord 单条记录
type PnLRecord struct {
	Date    time.Time         `json:"date"`
	Values  map[string]float64 `json:"values"`
	Schema  Schema            `json:"schema,omitempty"`
}

// --------------------------------------- Operator操作符函数结构体 -------------------------------------- //
type Operator struct {
	Name          string   `json:"name"`
	Category      string   `json:"category"`
	Scope         []string `json:"scope"`
	Definition    string   `json:"definition"`
	Description   string   `json:"description"`
	Documentation *string  `json:"documentation,omitempty"`
	Level         *string  `json:"level,omitempty"`
}

// --------------------------------- consultant | weight/value-factor 结构体 -------------------------- //
type Leaderboard struct {
	User                          string  `json:"user"`
	WeightFactor                  float64 `json:"weightFactor"`
	ValueFactor                   float64 `json:"valueFactor"`
	DataFieldsUsed                int     `json:"dataFieldsUsed"`
	SubmissionsCount              int     `json:"submissionsCount"`
	MeanProdCorrelation           float64 `json:"meanProdCorrelation"`
	MeanSelfCorrelation           float64 `json:"meanSelfCorrelation"`
	SuperAlphaSubmissionsCount    int     `json:"superAlphaSubmissionsCount"`
	SuperAlphaMeanProdCorrelation float64 `json:"superAlphaMeanProdCorrelation"`
	SuperAlphaMeanSelfCorrelation float64 `json:"superAlphaMeanSelfCorrelation"`
	University                    *string `json:"university"`
	Country                       string  `json:"country"`
}

type ConsultantResponse struct {
	DateStarted string      `json:"dateStarted"`
	Leaderboard Leaderboard `json:"leaderboard"`
}

// ---------------------------------------------- 金字塔优先推塔情况 -----------------------------------------------
type PyramidsResponse struct {
	Pyramids []Pyramids `json:"pyramids"`
}

type Pyramids struct {
	Category   Categorys `json:"category"`
	Region     string    `json:"region"`
	Delay      int       `json:"delay"`
	AlphaCount int       `json:"alphaCount"`
}

type Categorys struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
