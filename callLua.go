package eduBotMethods

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"github.com/tidwall/gjson"
)

func callLua(userId int, args []string) (result interface{}, err error) {

	if len(args) < 1 {
		err = errors.New(fmt.Sprintf("callLua wrong signature: need lua function name, but receive %s", args))
		return
	}

	funcName := args[0]

	_, err = callDbFunction(funcName, []interface{}{userId, args[1:]})
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), fmt.Sprintf("Procedure '%s' is not defined", funcName)) {
			err = fmt.Errorf("Функция '%s' не найдена в lua.", funcName)
		}
		return
	}

	// записываем строчку "nil" чтобы сообщение не отправлялось. Ответ будет отправлен через jsonRpc
	result = "nil"

	return
}

func callDbFunction(functionName string, args []interface{}) (result []byte, err error) {

	resp, err := config.TrntlConn.Call17(functionName, args)
	if err != nil {
		err = fmt.Errorf("call '%s': %s", functionName, err)
	}

	if len(resp.Data) == 0 {
		result = []byte{}
		return
	}

	queryRes := resp.Data[0].(string)

	//fmt.Printf("rsp %s\n", resp)
	code := gjson.Get(queryRes, "ok")
	if !code.Bool() {
		msg := gjson.Get(queryRes, "message").Str
		err = fmt.Errorf(msg)
	}
	result = []byte(gjson.Get(queryRes, "result").Raw)

	return
}



