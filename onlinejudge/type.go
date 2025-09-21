package onlinejudge

// ErrResponse
type ErrResponse struct {
	Error *string `json:"error"` // 使用指针可以区分null和空字符串
	Data  *string `json:"data"`
}
type AccountData struct {
	Users [][]string `json:"users"`
}

func AddDictWithUser(usersData []User) *AccountData {
	users := make([][]string, len(usersData))
	for i, user := range usersData {
		users[i] = []string{user.Account, user.Password, user.Email, user.Name}
	}
	return &AccountData{
		Users: users,
	}
}

type User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Email    string `json:"email"` // 通过邮箱标识考生信息
	Name     string `json:"name"`
}

//	GetRankParams struct {
//		Offset       string `json:"offset"`        // 游标
//		Limit        string `json:"limit"`         // 返回数据量
//		ContestId    string `json:"contest_id"`    // 比赛id
//		ForceRefresh string `json:"force_refresh"` // 是否是最新值 （默认0就行）
//	}
type GetRankParams struct {
	Offset       string `json:"offset"`        // 游标
	Limit        string `json:"limit"`         // 返回数据量
	ContestId    string `json:"contest_id"`    // 比赛id
	ForceRefresh string `json:"force_refresh"` // 是否是最新值 （默认0就行）
}

// RankData 包含排名数据和总数
type RankData struct {
	Results []RankRecord `json:"results"`
	Total   int          `json:"total"`
}

// RankRecord 单个用户的排名记录
type RankRecord struct {
	ID               int            `json:"id"`
	User             UserInfo       `json:"user"`
	SubmissionNumber int            `json:"submission_number"`
	TotalScore       int            `json:"total_score"`
	SubmissionInfo   map[string]int `json:"submission_info"` // 键：题目ID，值：该题得分
	Contest          int            `json:"contest"`
}

// UserInfo 考试用户基本信息
type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
}

const (
	BODY_JSON = "jsonbody"
	PARAMS    = "params"
)
