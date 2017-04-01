package eduBotMethods

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pepelazz/BotTarantool_0.1/src/common"
	"strings"
)

func callLua(userId int, args []string) (result interface{}, err error) {

	if len(args) < 1 {
		err = errors.New(fmt.Sprintf("callLua wrong signature: need lua function name, but receive %s", args))
		return
	}

	funcName := args[0]

	_, err = common.Trntl.CallDbFunction(funcName, []interface{}{userId, args[1:]})
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


