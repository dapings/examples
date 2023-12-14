package doubanbook

import (
	"regexp"
	"strconv"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
)

var (
	cookie  = "bid=-UXUw--yL5g; dbcl2=\\\"214281202:q0BBm9YC2Yg\\\"; __yadk_uid=jigAbrEOKiwgbAaLUt0G3yPsvehXcvrs; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1665849857.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.21428; ck=SAvm; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1665925405%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=*; __utma=30149280.2072705865.1665849857.1665849857.1665925407.2; __utmc=30149280; __utmt=1; __utmb=30149280.23.5.1665925419338; _pk_id.100001.8cb4=fc1581490bf2b70c.1665849856.2.1665925421.1665849856."
	bookUrl = "https://book.douban.com"
)

const (
	BookListTaskName = "douban_book_list"
	bookListReg      = `<a.*?href="([^"]+)" title="([^"]+)"`
	tagReg           = `<a href="([^"]+)" class="tag">([^<]+)</a>`
)

var (
	autoReg   = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
	publicReg = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
	pageReg   = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
	priceReg  = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
	scoreReg  = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
	intoReg   = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)
)

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     BookListTaskName,
		Cookie:   cookie,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				{
					Url:      bookUrl,
					Method:   "GET",
					Priority: 1,
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": {ParseFunc: ParseTag},
			"书籍列表":  {ParseFunc: ParseBookList},
			"书籍简介": {
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(tagReg)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		result.Requests = append(
			result.Requests, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      bookUrl + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:1]
	return result, nil
}

func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(bookListReg)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &collect.Temp{}
		err := req.TmpData.Set("book_name", string(m[2]))
		if err != nil {
			return collect.ParseResult{}, err
		}
		result.Requests = append(result.Requests, req)
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:1]
	return result, nil
}

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageReg))
	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoReg),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, publicReg),
		"得分":  ExtraString(ctx.Body, scoreReg),
		"价格":  ExtraString(ctx.Body, priceReg),
		"简介":  ExtraString(ctx.Body, intoReg),
	}
	return collect.ParseResult{Items: []any{ctx.Output(book)}}, nil
}

func ExtraString(contents []byte, reg *regexp.Regexp) string {
	match := reg.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}
