package doubangroup

import (
	"fmt"
	"regexp"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
)

var (
	cookie        = "bid=-UXUw--yL5g; dbcl2=\\\"214281202:q0BBm9YC2Yg\\\"; __yadk_uid=jigAbrEOKiwgbAaLUt0G3yPsvehXcvrs; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1665849857.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.21428; ck=SAvm; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1665925405%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=*; __utma=30149280.2072705865.1665849857.1665849857.1665925407.2; __utmc=30149280; __utmt=1; __utmb=30149280.23.5.1665925419338; _pk_id.100001.8cb4=fc1581490bf2b70c.1665849856.2.1665925421.1665849856."
	discussionURL = "https://www.douban.com/group/szsh/discussion?start=%d"
)

const (
	urlListRe           = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
	ContentRe           = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`
	FindSumRoomTaskName = "find_douban_sun_room"
)

var DoubangroupTask = &spider.Task{
	Property: spider.Property{
		Name:     FindSumRoomTaskName,
		Cookie:   cookie,
		WaitTime: 2,
		MaxDepth: 5,
	},
	Rule: spider.RuleTree{
		Root: func() ([]*spider.Request, error) {
			var roots []*spider.Request
			for i := 0; i < 25; i += 25 {
				roots = append(roots, &spider.Request{
					URL:      fmt.Sprintf(discussionURL, i),
					Method:   "GET",
					Priority: 1,
					RuleName: "解析网站URL",
				})
			}

			return roots, nil
		},
		Trunk: map[string]*spider.Rule{
			"解析网站URL": {ParseFunc: ParseURL},
			"解析阳台房":   {ParseFunc: GetSunRoom},
		},
	},
}

func GetSunRoom(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(ContentRe)

	if ok := re.Match(ctx.Body); !ok {
		return spider.ParseResult{Items: make([]any, 0)}, nil
	}

	return spider.ParseResult{Items: []any{ctx.Req.URL}}, nil
}

func ParseURL(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(urlListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := spider.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &spider.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			URL:      u,
			Depth:    ctx.Req.Depth + 1,
			RuleName: "解析阳台房",
		})
	}

	return result, nil
}
