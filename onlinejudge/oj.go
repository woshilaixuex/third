package onlinejudge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OjOptions 配置项结构体
type OjOptions struct {
	CsrfToken string
	Origin    string
	SessionId string
}

// OjTools 主结构体
type OjTools struct {
	option OjOptions
	client *http.Client
}

// Option 定义函数选项类型
type Option func(*OjOptions)

// DefaultOjOptions 提供默认配置选项
func DefaultOjOptions() OjOptions {
	return OjOptions{
		CsrfToken: "u4RJHKE1Qc6AuO6t6vonnLR0VQbujeen5vuz9zNeVlMJzGIYIvgKqkBTeHMwzgLg",
		Origin:    "oj.dsstudio.tech",
		SessionId: "t2b0p2840gb5mbgpkt0focoirpri4lti",
	}
}

// NewOjTools 构造函数，接受可变的 Option 参数,默认参数会过期，到时候可以这样写
// tools := onlinejudge.NewOjTools(
//
//	onlinejudge.WithCsrfToken("your_csrf_token_here"),
//	onlinejudge.WithSessionId("your_session_id_here"),
//
// )
func NewOjTools(opts ...Option) *OjTools {
	options := DefaultOjOptions()

	for _, opt := range opts {
		opt(&options)
	}

	return &OjTools{
		option: options,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// WithCsrfToken 设置 CsrfToken
func WithCsrfToken(csrfToken string) Option {
	return func(o *OjOptions) {
		o.CsrfToken = csrfToken
	}
}

// WithOrigin 设置 Origin
func WithOrigin(origin string) Option {
	return func(o *OjOptions) {
		o.Origin = origin
	}
}

// WithSessionId 设置 SessionId
func WithSessionId(sessionId string) Option {
	return func(o *OjOptions) {
		o.SessionId = sessionId
	}
}

// 添加账号给oj平台
func (tools OjTools) PushAccount(usersData []User) error {
	client := NewPushAccountClient(tools.option)
	accountData := AddDictWithUser(usersData)
	data := ojData{
		dataType: BODY_JSON,
		data:     accountData,
	}
	request, err := client.GetRequest(data)
	if err != nil {
		return err
	}
	bodyBytes, err := tools.doRequset(request)
	if err != nil {
		return err
	}
	return returnErrResponse(bodyBytes)
}

// 获取对应考试排名信息
func (tools OjTools) GetExamRank(parmas GetRankParams) (rankData RankData, err error) {
	client := NewGetRankClient(tools.option)
	parmasMap := make(map[string]interface{})
	parmasMap["offset"] = parmas.Offset
	parmasMap["limit"] = parmas.Limit
	parmasMap["contest_id"] = parmas.ContestId
	data := ojData{
		dataType: PARAMS,
		data:     parmasMap,
	}
	request, err := client.GetRequest(data)
	if err != nil {
		return
	}
	bodyBytes, err := tools.doRequset(request)
	if err != nil {
		return
	}
	var apiResponse struct {
		Error *string  `json:"error"`
		Data  RankData `json:"data"` // 直接使用RankData解析data字段
	}
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return RankData{}, returnErrResponse(bodyBytes)
	}
	if apiResponse.Error != nil && *apiResponse.Error != "" {
		return RankData{}, fmt.Errorf("api returned error: %s", *apiResponse.Error)
	}

	return apiResponse.Data, nil
}
func (tools OjTools) doRequset(request *http.Request) ([]byte, error) {
	bodyBytes := make([]byte, 0)
	response, err := tools.client.Do(request)
	if err != nil {
		return bodyBytes, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close() // 确保响应体被关闭

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return bodyBytes, fmt.Errorf("request failed with status code: %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(response.Body)
	if err != nil {
		return bodyBytes, fmt.Errorf("failed to read response body: %w", err)
	}
	return bodyBytes, nil
}
func returnErrResponse(bodyBytes []byte) error {
	errResponse := ErrResponse{}
	err := json.Unmarshal(bodyBytes, &errResponse)
	if err != nil {
		return err
	}
	if errResponse.Error == nil {
		return nil
	}
	return fmt.Errorf("err cause: %s", string(bodyBytes))
}
