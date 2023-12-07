# Crawler

## 4种网页文本处理手段

1. 正则表达式

    ```text
    // regexp.MustCompile函数会在编译时，提前解析好正则表达式内容，在一定程度上加速程序的运行。
    // [\s\S]*?，[\s\S] 任意字符串，*将前面任意字符匹配0次或无数次，?非贪婪匹配，找到第一次出现的地方，就认定匹配成功。
    // 由于回溯的原因，复杂的正则表达式，可能比较消耗CPU资源。
    var headerReg = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
    
    // FindAllSubmatch返回一个三维字节数组，第三层是字符实际对应的字节数组
    matches := headerReg.FindAllSubmatch(body, -1)
    ```

2. XPath(XML Path Language)

    定义了一种遍历XML文档中节点层次结构，并返回匹配元素的灵活方法。
    第三方库`github.com/antchfx/htmlquery`提供了在HTML中通过XPath匹配XML节点的引擎。
    ```text
    // XPath语法
    var xpathReg = `//div[@class="news_li"]/h2/a[@target="_blank"]`
    // 解析HTML文本
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	
	// 通过XPath语法查找符合条件的节点
	nodes := htmlquery.Find(doc, xpathReg)
    ```

3. CSS选择器
4. 标准库：strings,bytes,text/encoding,html/charset


- https://github.com/dreamerjackson/crawler “聚沙万塔-Go语言构建高性能、分布式爬虫项目”