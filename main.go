package eduBotMethods

import (
	"regexp"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"github.com/pepelazz/go-bot-user-session"
	"time"
	"github.com/pepelazz/go-bot-telebot"
	"github.com/tarantool/go-tarantool"
)


var (
	methodMap map[string]method
	config *Config
)

func init()  {
	methodMap = map[string]method{}
	Add("flickerPhoto", flickerPhoto)
	Add("getPhoto", getPhoto)
	Add("getSticker", getSticker)
	Add("callLua", callLua)
}

type Config struct {
	TrntlConn *tarantool.Connection
	FlickerKey string
}

func Init(cnf *Config) (err error) {
	config = cnf
	if cnf.TrntlConn == nil {
		err = errors.New("eduBotMethods.Init missed TrntlConn")
	}
	return
}

type method func(int, []string) (interface{}, error)

type PhotoUrlWithCapture struct {
	Url     string
	Capture string
}

type StickerWithText struct {
	FileId string
	Text   string
}

func CheckIsEduBotMethod(s *userSession.S)  {
	res, err := execute(s)
	if err != nil {
		s.SetAnswerMsg(fmt.Sprintf("Ошибка: %s", err))
		return
	}
	if res != nil {
		switch v := res.(type) {
		case PhotoUrlWithCapture:
			s.SetAnswerMsgWithPhoto(v.Capture, "", v.Url)
			break
		case StickerWithText:
			s.SetAnswerWithSticker(v.FileId)
			go sendMsgWithDelay(s, v.Text, 1) // текст сообщения отправляем вслед за стикером
			break
		case string:
			s.SetAnswerMsg(res.(string))
			break
		default:
			s.SetAnswerMsg(fmt.Sprintf("Error: unknown type for interface assertion: %s", res))
		}
	}
}

func execute(s *userSession.S) (res interface{}, err error) {
	r, _ := regexp.Compile(`{{eduBot-call(\s*(\w+)\s*\((.*)\))}}`)
	groups := r.FindAllStringSubmatch(s.AnswerMsg(), 1)

	if len(groups) > 0 && len(groups[0]) > 2 {
		for k, v := range methodMap {
			if strings.ToLower(k) == strings.ToLower(groups[0][2]) {
				var args []string
				for _, argStr := range strings.Split((groups[0][3:])[0], ",") {
					s := strings.TrimSpace(argStr)
					if len(s) > 0 {
						args = append(args, s)
					}
				}
				return v(s.IdInt(), args)
			}
		}
		err = errors.New(fmt.Sprintf("eduBot method not found for function name: %s", groups[0][2]))
	}
	return

}

func Add(funcName string, f method) {
	methodMap[funcName] = f
}


func sendMsgWithDelay(s *userSession.S, text string, delay time.Duration) {
	time.Sleep(delay * time.Second)
	if len(text) > 0 {
		message := telebot.Message{Sender: telebot.User{ID:s.IdInt()}}
		s, err := userSession.New(message)
		if err != nil {
			return
		} else {
			_, err = s.SetAnswerMsg(text).SendMsg()
		}
	}
}