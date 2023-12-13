package collect

type (
	TaskModel struct {
		Property
		Root  string      `json:"root_script"`
		Rules []RuleModel `json:"rule"`
	}

	RuleModel struct {
		Name      string `json:"name"`
		ParseFunc string `json:"parse_script"`
	}
)

// 用于动态规则添加请求。

func generateReq(jsReq map[string]any) *Request {
	req := &Request{}
	u, ok := jsReq["Url"].(string)
	if !ok {
		return nil
	}
	req.Url = u
	req.RuleName, _ = jsReq["RuleName"].(string)
	req.Method, _ = jsReq["Method"].(string)
	req.Priority, _ = jsReq["Priority"].(int64)
	return req
}

func AddJsReq(jsReq map[string]any) []*Request {
	reqs := make([]*Request, 0)
	reqs = append(reqs, generateReq(jsReq))
	return reqs
}

func AddJsReqs(jsReqs []map[string]any) []*Request {
	reqs := make([]*Request, 0)
	for _, jsReq := range jsReqs {
		reqs = append(reqs, generateReq(jsReq))
	}
	return reqs
}
