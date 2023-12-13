package doubangroup

import (
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
)

const (
	JSFindSumRoomTaskName = "js_find_douban_sun_room"
)

var DoubangroupJSTask = &collect.TaskModel{
	Property: collect.Property{
		Name:     JSFindSumRoomTaskName,
		Cookie:   cookie,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
	},
	Root: `
				var arr = new Array();
		        for (var i = 25; i <= 25; i+=25) {
					var obj = {
					   Url: "https://www.douban.com/group/szsh/discussion?start=" + i,
					   Priority: 1,
					   RuleName: "解析网站URL",
					   Method: "GET",
				   };
					arr.push(obj);
				};
				console.log(arr[0].Url);
				AddJsReq(arr);
	`,
	Rules: []collect.RuleModel{
		{
			Name: "解析网站URL",
			ParseFunc: `
			ctx.ParseJSReg("解析阳台房","(https://www.douban.com/group/topic/[0-9a-z]+/)\"[^>]*>([^<]+)</a>");
			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
			//console.log("parse output");
			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
			`,
		},
	},
}
