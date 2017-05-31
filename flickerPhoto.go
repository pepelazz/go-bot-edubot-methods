package eduBotMethods

import (
	"strings"
	"fmt"
	"time"
	"github.com/azer/go-flickr"
	"encoding/json"
	"github.com/pkg/errors"
	"math/rand"
)

func flickerPhoto(userId int, args []string) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New(fmt.Sprintf("flickerPhoto wrong signature: need (string, string), but receive %s", args))
	}

	url, err := flickrGetPhoto(args[0])
	if err != nil {
		return "", err
	}
	var captureMsg, riveVarName, riveVarValue string
	if len(args) > 1 {
		captureMsg = args[1]
	}
	if len(args)>3 {
		riveVarName = args[2]
		riveVarValue = args[3]
	}
	return PhotoUrlWithCapture{url, captureMsg, riveVarName, riveVarValue}, nil
}

func flickrGetPhoto(tag string) (url string, err error) {
	if len(config.FlickerKey) == 0 {
		err = errors.New("Missed FlickerKey in config. Write FlickerKey snd restart function.")
		return
	}
	client := &flickr.Client{
		Key: config.FlickerKey,
		//Token: "token", // optional
		//Sig: "sig", // optional
	}

	res := struct {
		Stat    string
		Message string
		Photos  struct {
				Photo []struct {
					Id       string `json:"id"`
					Title    string `json:"title"`
					Secret   string `json:"secret"`
					Server   string `json:"server"`
					Farm     int `json:"farm"`
					Ispublic int `json:"ispublic"`
				}
			} `json:"photos"`
	}{}

	response, err := client.Request("photos.search", flickr.Params{"tags": strings.TrimSpace(tag) })
	if err != nil {
		return
	}

	//https://farm{farm-id}.staticflickr.com/{server-id}/{id}_{secret}.jpg
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	if res.Stat == "fail" {
		err = errors.New(fmt.Sprintf("<b>flickr error</b>: %s", res.Message))
		return
	}

	if len(res.Photos.Photo) > 1 {
		rand.Seed(time.Now().Unix())
		i := rand.Intn(len(res.Photos.Photo) - 1)
		v := res.Photos.Photo[i]
		url = fmt.Sprintf("https://farm%v.staticflickr.com/%s/%s_%s.jpg\n", v.Farm, v.Server, v.Id, v.Secret)
	} else {
		err = errors.New("ничего похожего не найдено")
	}

	//
	//for i, v := range res.Photos.Photo {
	//	fmt.Println(i, ") ", v)
	//	fmt.Printf("%v) https://farm%v.staticflickr.com/%s/%s_%s.jpg\n", i, v.Farm, v.Server, v.Id, v.Secret)
	//}
	//fmt.Printf("flicker url %s\n", url)

	return
}


