package onlinejudge_test

import (
	"testing"

	"github.com/woshilaixuex/third/onlinejudge"
)

func TestGetExamRank(t *testing.T) {
	tools := onlinejudge.NewOjTools()
	params := onlinejudge.GetRankParams{
		Offset:    "0",
		Limit:     "10",
		ContestId: "10",
	}
	rank, err := tools.GetExamRank(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rank)
}
func TestPushUser(t *testing.T) {
	tools := onlinejudge.NewOjTools(
		onlinejudge.WithCsrfToken("your_csrf_token_here"),
		onlinejudge.WithSessionId("your_session_id_here"),
	)
	accountData := []onlinejudge.User{
		{
			Account:  "csd00052",
			Password: "xHu6fL32",
			Email:    "csd00051@exam.com",
			Name:     "csd00051",
		},
		{
			Account:  "csd00053",
			Password: "yJt7gK91",
			Email:    "csd00052@exam.com",
			Name:     "csd00052",
		},
	}
	err := tools.PushAccount(accountData)
	if err != nil {
		t.Error(err)
		return
	}
}
