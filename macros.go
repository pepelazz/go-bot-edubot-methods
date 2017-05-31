package eduBotMethods

import (
	"regexp"
	"strings"
	"github.com/pepelazz/go-bot-user-session"
	"fmt"
)

func checkIsEduBotMacros(s *userSession.S) (res string, err error) {
	r, _ := regexp.Compile(`{{eduBot-macros\s*\((.*)\)}}`)
	msg := s.AnswerMsg()
	groups := r.FindAllStringSubmatch(msg, 1)
	if len(groups) > 0 && len(groups[0]) > 1 {
		var resStr string
		resStr, err = replaceStr(strings.Split(groups[0][1], ",")[0], groups[0][1])
		if err != nil {
			return
		}
		res = r.ReplaceAllString(msg, resStr)
		return
	}
	return
}

func replaceStr(name string, groups string) (string, error) {
	switch name {
	case "capitalize":
		return macrosCapitalize(groups)
	default:
		return fmt.Sprintf("Function for macros '%s' not found. Add process function.", name), nil
	}
}

func macrosCapitalize(str string) (string, error) {
	groups := strings.Split(str, ",")
	if len(groups) > 1 {
		return strings.Title(groups[1]), nil
	} else {
		return "", nil
	}
}