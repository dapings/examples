package collect

// RuleTree 采集规则树
type RuleTree struct {
	Root  func() []*Request // 根节点，执行入口
	Trunk map[string]*Rule  // 规则哈希表
}

// Rule 采集规则节点
type Rule struct {
	ParseFunc func(*Context) ParseResult // 内容解析函数
}
