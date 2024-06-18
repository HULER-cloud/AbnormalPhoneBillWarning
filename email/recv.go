package email

import (
	"AbnormalPhoneBillWarning/config"
	"encoding/base64"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"net/mail"
	"time"
)

func Recv(cfg config.Config) {

	// 连接imap邮箱服务器
	imap_addr := fmt.Sprintf("%s:%d", cfg.Email.RecvHost, cfg.Email.RecvPort)
	c, err := client.DialTLS(imap_addr, nil)
	if err != nil {
		log.Fatal("imap服务器连接失败")
	}
	log.Println("imap服务器已连接")

	// 登录
	if err := c.Login(cfg.Email.User, cfg.Email.Password); err != nil {
		log.Fatal(err)
	}
	log.Println("登陆成功")
	defer c.Logout()

	// 邮箱文件夹列表
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	// async函数可能需要这种方式来提升效率？
	go func() {
		done <- c.List("", "*", mailboxes)
	}()
	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("邮箱文件夹:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("收件箱中共有%d封邮件 ", mbox.Messages)

	// 到这里为止，上面的只需要做一次
	// 下面的应该会塞到定时器里面循环

	for {
		// 获取到未读邮件
		unseen_emails := GetUnseenEmails(c)

		for _, email := range unseen_emails {
			// 这里放一些处理的逻辑
			fmt.Println(email)
			// 另外需要注意的是，发送完之后，应当回复一封邮件来确认
			// 但是好像由于团队沟通不畅，这部分内容已经没有用了？？？
			Process(email)
		}

		time.Sleep(60 * time.Second)
	}

}

func GetUnseenEmails(c *client.Client) []UnseenEmail {

	var unseen_emails []UnseenEmail

	// 获取未读邮件列表
	// 新建规则
	criteria := imap.NewSearchCriteria()
	// 未读标记
	criteria.WithoutFlags = []string{imap.SeenFlag}
	// 检索并获取到符合“未读”规则的邮件的uid
	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("未读邮件数：", len(uids))
	fmt.Println("未读邮件列表:")

	// 打印未读邮件列表
	for _, uid := range uids {
		// 获取到当前处理的这封未读邮件的序列号
		seqset := new(imap.SeqSet)
		seqset.AddNum(uid)

		// 获取邮件内容

		// 创建用于读取的channel
		messages := make(chan *imap.Message, 1)

		section := &imap.BodySectionName{}
		messages = make(chan *imap.Message, 1)
		done := make(chan error, 1)
		// 同理async获取
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
		}()
		if err := <-done; err != nil {
			log.Fatal(err)
		}
		// 将某封邮件的信息获取过来
		msg := <-messages

		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("服务器没有返回消息体")
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}

		// 获取到相应的属性
		header := m.Header
		subject, _ := decodeMIMEWord(header.Get("Subject"))
		from, _ := decodeMIMEWord(header.Get("From"))
		fmt.Println("Subject:", subject)
		fmt.Println("From:", from)

		// 读取邮件主体
		body, err := decodeBody(m.Header.Get("Content-Type"), m.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 读取到邮件正文的base64编码数据
		main_text_base64 := readMainText(body)
		// 解码为字节流
		decoded, err := base64.StdEncoding.DecodeString(main_text_base64)
		if err != nil {
			fmt.Println("解码出错:", err)
		}

		// 将解码后的字节转换为字符串
		decodedStr := string(decoded)
		fmt.Println("Body:", decodedStr)

		// 将邮件标记为已读
		flags := []interface{}{imap.SeenFlag}
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Fatal(err)
		}
		// 追加到邮件列表并返回
		unseen_emails = append(unseen_emails, UnseenEmail{
			Subject: subject,
			From:    from,
			Body:    body,
		})
	}
	fmt.Println("操作完成，未读邮件已标记为已读")
	return unseen_emails
}
