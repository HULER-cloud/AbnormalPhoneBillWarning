package email

import (
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"mime"
	"strings"
)

type UnseenEmail struct {
	Subject string
	From    string
	Body    string
}

//var GlobalUnseenEmails []UnseenEmail

// decodeMIMEWord 解码 MIME encoded-word 字符串
func decodeMIMEWord(s string) (string, error) {
	decoder := new(mime.WordDecoder)
	decoded, err := decoder.DecodeHeader(s)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

// decodeBody 解码邮件体，处理字符集
func decodeBody(contentType string, body io.Reader) (string, error) {
	// 获取字符集
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	chars := params["charset"]
	if chars == "" {
		chars = "utf-8"
	}
	// 创建新读取器
	reader, err := charset.NewReaderLabel(chars, body)
	if err != nil {
		return "", err
	}
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// readMainText 读取邮件正文的base64编码串
func readMainText(body string) string {
	// 按行拆分文本
	lines := strings.Split(body, "\n")
	// 选择第8行（索引为7，因为索引从0开始）
	lineIndex := 7
	res := ""
	for i := lineIndex; i < len(lines); i++ {

		if lines[i][0] != 13 {
			res += lines[i]
		} else {
			break
		}
	}
	return res
}
