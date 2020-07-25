package notify

import (
	"bytes"
	"ferry/models/system"
	"ferry/pkg/notify/email"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

/*
  @Author : lanyulei
  @同时发送多种通知方式
*/

type BodyData struct {
	SendTo        interface{} // 接受人
	Subject       string      // 标题
	Classify      []int       // 通知类型
	Id            int         // 工单ID
	Title         string      // 工单标题
	Creator       string      // 工单创建人
	Priority      int         // 工单优先级
	PriorityValue string      // 工单优先级
	CreatedAt     string      // 工单创建时间
	Content       string      // 通知的内容
}

func (b *BodyData) ParsingTemplate() (err error) {
	// 读取模版数据
	var (
		buf bytes.Buffer
	)

	log.Info(os.Getwd())
	tmpl, err := template.ParseFiles("./pkg/notify/template/email.html")
	if err != nil {
		return
	}

	err = tmpl.Execute(&buf, b)
	if err != nil {
		return
	}

	b.Content = buf.String()

	return
}

func (b *BodyData) SendNotify() {
	var (
		emailList []string
		err       error
	)

	switch b.Priority {
	case 1:
		b.PriorityValue = "正常"
	case 2:
		b.PriorityValue = "紧急"
	case 3:
		b.PriorityValue = "非常紧急"
	}

	for _, c := range b.Classify {
		switch c {
		case 1: // 邮件
			for _, user := range b.SendTo.(map[string]interface{})["userList"].([]system.SysUser) {
				emailList = append(emailList, user.Email)
			}
			err = b.ParsingTemplate()
			if err != nil {
				log.Errorf("模版内容解析失败，%v", err.Error())
				return
			}
			go email.SendMail(emailList, b.Subject, b.Content)
		}
	}
}
