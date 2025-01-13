package wxwork

import (
	"fmt"
	"testing"
)

func TestPushBotMarkdown(t *testing.T) {
	fmt.Println(PushToBot(
		"API_KEY",
		map[string]any{
			"msgtype": "markdown",
			"markdown": map[string]any{
				"content": `实时新增用户反馈<font color="warning">132例</font>，请相关同事注意。
> 类型:<font color="comment">用户反馈</font>
> 普通用户反馈:<font color="comment">117例</font>
> VIP用户反馈:<font color="comment">15例</font>`,
			},
		},
		"",
	))
}
