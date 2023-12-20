package doubanbook

import (
	"regexp"
	"strconv"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"go.uber.org/zap"
)

var (
	cookie  = "bid=-UXUw--yL5g; push_doumail_num=0; __utmv=30149280.21428; __utmc=30149280; __gads=ID=c6eaa3cb04d5733a-2259490c18d700e1:T=1666111347:RT=1666111347:S=ALNI_MaonVB4VhlZG_Jt25QAgq-17DGDfw; frodotk_db=\\\"17dfad2f83084953479f078e8918dbf9\\\"; gr_user_id=cecf9a7f-2a69-4dfd-8514-343ca5c61fb7; __utmc=81379588; _vwo_uuid_v2=D55C74107BD58A95BEAED8D4E5B300035|b51e2076f12dc7b2c24da50b77ab3ffe; __yadk_uid=BKBuETKRjc2fmw3QZuSw4rigUGsRR4wV; ct=y; ll=\\\"108288\\\"; viewed=\\\"36104107\\\"; ap_v=0,6.0; __gpi=UID=000008887412003e:T=1666111347:RT=1668851750:S=ALNI_MZmNsuRnBrad4_ynFUhTl0Hi0l5oA; __utma=30149280.2072705865.1665849857.1668851747.1668854335.25; __utmz=30149280.1668854335.25.4.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/misc/sorry; __utma=81379588.990530987.1667661846.1668852024.1668854335.8; __utmz=81379588.1668854335.8.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/misc/sorry; _pk_ref.100001.3ac3=[\\\"\\\",\\\"\\\",1668854335,\\\"https://www.douban.com/misc/sorry?original-url=https%3A%2F%2Fbook.douban.com%2Ftag%2F%25E5%25B0%258F%25E8%25AF%25B4\\\"]; _pk_ses.100001.3ac3=*; gr_cs1_5f43ac5c-3e30-4ffd-af0e-7cd5aadeb3d1=user_id:0; __utmt=1; dbcl2=\\\"214281202:GLkwnNqtJa8\\\"; ck=dBZD; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03=ca04de17-2cbf-4e45-914a-428d3c26cfe3; gr_cs1_ca04de17-2cbf-4e45-914a-428d3c26cfe3=user_id:1; __utmt_douban=1; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03_ca04de17-2cbf-4e45-914a-428d3c26cfe3=true; __utmb=30149280.10.10.1668854335; __utmb=81379588.9.10.1668854335; _pk_id.100001.3ac3=02339dd9cc7d293a.1667661846.8.1668855011.1668852362.; push_noty_num=0"
	bookURL = "https://book.douban.com"
)

const (
	BookListTaskName = "douban_book_list"
	bookListReg      = `<a.*?href="([^"]+)" title="([^"]+)"`
	tagReg           = `<a href="([^"]+)" class="tag">([^<]+)</a>`
)

var (
	autoReg   = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
	publicReg = regexp.MustCompile(`<span class="pl">出版社:</span>[\d\D]*?<a.*?>([^<]+)</a>`)
	pageReg   = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
	priceReg  = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
	scoreReg  = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
	intoReg   = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)
)

var DoubanBookTask = &spider.Task{
	Options: spider.Options{Name: BookListTaskName},
	Rule: spider.RuleTree{
		Root: func() ([]*spider.Request, error) {
			roots := []*spider.Request{
				{
					Priority: 1,
					URL:      bookURL,
					Method:   "GET",
					RuleName: "数据tag",
				},
			}

			return roots, nil
		},
		Trunk: map[string]*spider.Rule{
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

func ParseTag(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(tagReg)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := spider.ParseResult{}

	for _, m := range matches {
		result.Requests = append(
			result.Requests, &spider.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				URL:      bookURL + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	zap.S().Debugln("parse book tag, request count:", len(result.Requests))
	// result.Requests = result.Requests[:1]
	return result, nil
}

func ParseBookList(ctx *spider.Context) (spider.ParseResult, error) {
	re := regexp.MustCompile(bookListReg)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := spider.ParseResult{}

	for _, m := range matches {
		req := &spider.Request{
			Priority: 100,
			Method:   "GET",
			Task:     ctx.Req.Task,
			URL:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &spider.Temp{}

		err := req.TmpData.Set("book_name", string(m[2]))
		if err != nil {
			zap.L().Error("set tmp data failed", zap.Error(err))

			return spider.ParseResult{}, err
		}

		result.Requests = append(result.Requests, req)
	}
	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	zap.S().Debugln("parse book list, request count:", len(result.Requests))
	// result.Requests = result.Requests[:3]
	return result, nil
}

func ParseBookDetail(ctx *spider.Context) (spider.ParseResult, error) {
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

	zap.S().Debugln("parse book detail", book)

	return spider.ParseResult{Items: []any{ctx.Output(book)}}, nil
}

func ExtraString(contents []byte, reg *regexp.Regexp) string {
	match := reg.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}

	return ""
}
