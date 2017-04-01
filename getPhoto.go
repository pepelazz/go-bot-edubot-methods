package eduBotMethods

import (
	"fmt"
	"github.com/pkg/errors"
)

func getPhoto(userId int, args []string) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New(fmt.Sprintf("getPhoto wrong signature: need (string, string), but receive %s", args))
	}

	var captureMsg string
	if len(args)>1 {
		captureMsg = args[1]
	}
	return PhotoUrlWithCapture{args[0], captureMsg}, nil
}


func getSticker(userId int, args []string) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New(fmt.Sprintf("getSticker wrong signature: need (string, string), but receive %s", args))
	}

	var captureMsg string
	if len(args)>1 {
		captureMsg = args[1]
	}
	return StickerWithText{args[0], captureMsg}, nil
}


