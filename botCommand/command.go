package botcommand

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jeanhua/PinBot/config"
	messageChain "github.com/jeanhua/PinBot/messagechain"
	"github.com/jeanhua/PinBot/model"
	"github.com/jeanhua/PinBot/utils"
)

var (
	EnableAIAudio = false
	CommandMu     sync.RWMutex
)

func DealGroupCommand(com string, msg *model.GroupMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		config.ConfigInstance_mu.RLock()
		chain.Text(config.ConfigInstance.HelpWords.Group)
		config.ConfigInstance_mu.RUnlock()
		messageChain.SendMessage(chain)
		return true
	case "/enable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		CommandMu.Lock()
		EnableAIAudio = true
		CommandMu.Unlock()
		chain.Text(" 已开启AI语音")
		messageChain.SendMessage(chain)
		return true
	case "/disable AI语音":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		CommandMu.Lock()
		EnableAIAudio = false
		CommandMu.Unlock()
		chain.Text(" 已关闭AI语音")
		messageChain.SendMessage(chain)
		return true
	case "/zanao post":
		zanao := &utils.Zanao{}
		resp, err := zanao.GetNewest()
		if err != nil {
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text("我遇到了一点错误，请稍后再试")
			messageChain.SendMessage(chain)
			return true
		}
		groupForward := messageChain.GroupForward(msg.GroupId, "集市最新帖子")
		for _, v := range resp.Data.List {
			groupForward.Text(fmt.Sprintf("%s\n%s", v.Title, v.Content), msg.SelfId, "江颦")
		}
		groupForward.Send()
		return true
	case "/zanao hot":
		chain := messageChain.Group(msg.GroupId)
		chain.Reply(msg.MessageId)
		chain.Mention(msg.UserId)
		zanao := &utils.Zanao{}
		resp, err := zanao.GetHot()
		if err != nil {
			utils.LogErr(err.Error())
			chain.Text("我遇到了一点错误，请稍后再试")
			messageChain.SendMessage(chain)
			return true
		}
		text := "实时热帖：\n"
		for i, v := range resp.Data.List {
			text += fmt.Sprintf("[%d]%s\n", i+1, v.Title)
		}
		text = strings.TrimSpace(text)
		chain.Text(text)
		messageChain.SendMessage(chain)
		return true
	}

	// 查课
	if strings.HasPrefix(com, "/查课") {
		param := strings.Split(com, " ")
		if len(param) < 2 {
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text(" 请输入参数哦! 例如 /查课 高等数学 page1")
			messageChain.SendMessage(chain)
			return true
		}
		page := 1
		if len(param) == 3 {
			tp, err := strconv.Atoi(strings.TrimLeft(param[2], "page"))
			if err != nil || !strings.HasPrefix(param[2], "page") {
				chain := messageChain.Group(msg.GroupId)
				chain.Reply(msg.MessageId)
				chain.Mention(msg.UserId)
				chain.Text(" 请正确输入页码哦! 例如 /查课 高等数学 page1")
				messageChain.SendMessage(chain)
				return true
			}
			page = tp
		}

		course := utils.Course{}
		result, err := course.Search(param[1], page)
		if err != nil {
			utils.LogErr(err.Error())
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text(" 查询失败，请稍后再试!")
			messageChain.SendMessage(chain)
			return true
		}
		var responseStr strings.Builder
		responseStr.WriteString(" 查询到以下内容：\n")
		for _, v := range result.Data {
			responseStr.WriteString(fmt.Sprintf("[%d]%s %d-%s %s %s\n", v.Kid, v.CourseName, v.Kch, v.Kxh, v.ExamTypeName, v.TeachersName))
		}
		responseStr.WriteString("\n@我发送 /课程详情 课程id 可查看课程详情，比如 /课程详情 12154")
		groupForward := messageChain.GroupForward(msg.GroupId, "查询结果")
		groupForward.Text(responseStr.String(), msg.SelfId, "江颦")
		groupForward.Send()
		return true
	}

	// 课程详情
	if strings.HasPrefix(com, "/课程详情") {
		param := strings.Split(com, " ")
		if len(param) != 2 {
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text(" 请正确输入参数哦! 例如 /课程详情 12154")
			messageChain.SendMessage(chain)
			return true
		}
		kid, err := strconv.Atoi(param[1])
		if err != nil {
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text(" 请正确输入课程id哦! 例如 /课程详情 12154")
			messageChain.SendMessage(chain)
			return true
		}
		course := utils.Course{}
		result, err := course.GetDetail(kid)
		if err != nil {
			utils.LogErr(err.Error())
			chain := messageChain.Group(msg.GroupId)
			chain.Reply(msg.MessageId)
			chain.Mention(msg.UserId)
			chain.Text(" 查询失败，请稍后再试!")
			messageChain.SendMessage(chain)
			return true
		}
		groupForward := messageChain.GroupForward(msg.GroupId, "课程详情")
		groupForward.Text(fmt.Sprintf("课程名:%s\n任课教师:%s\n课程号:%d\n课序号:%s\n课程学分:%.1f\n考察类型:%s", result.Data.CourseName, result.Data.TeachersName, result.Data.Kch, result.Data.Kxh, result.Data.Credit, result.Data.ExamTypeName), msg.SelfId, "江颦")
		groupForward.Text(fmt.Sprintf("统计总表:\n统计人数:%d\n平均分:%.2f\n最高分:%.2f\n最低分:%.2f\n90~100分:%d人\n80~89分:%d人\n70~79分:%d人\n60~69分:%d人\n0~59分:%d人", result.Data.Count, result.Data.Average, result.Data.Max, result.Data.Min, result.Data.A_levelCount, result.Data.B_levelCount, result.Data.C_levelCount, result.Data.D_levelCount, result.Data.E_levelCount), msg.SelfId, "江颦")
		for _, v := range result.Data.History {
			groupForward.Text(fmt.Sprintf("考试时间:%d:\n统计人数:%d\n平均分:%.2f\n最高分:%.2f\n最低分:%.2f\n90~100分:%d人\n80~89分:%d人\n70~79分:%d人\n60~69分:%d人\n0~59分:%d人", v.ExamTime, v.Count, v.Average, v.Max, v.Min, v.A_levelCount, v.B_levelCount, v.C_levelCount, v.D_levelCount, v.E_levelCount), msg.SelfId, "江颦")
		}
		groupForward.Send()
		return true
	}

	return false
}

func DealFriendCommand(com string, msg *model.FriendMessage) bool {
	com = "/" + strings.TrimLeft(com, "/")
	switch com {
	case "/help", "/帮助":
		chain := messageChain.Friend(msg.UserId)
		config.ConfigInstance_mu.RLock()
		chain.Text(config.ConfigInstance.HelpWords.Friend)
		config.ConfigInstance_mu.RUnlock()
		messageChain.SendMessage(chain)
		return true
	default:
		return false
	}
}
