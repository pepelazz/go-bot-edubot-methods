package eduBotMethods

import (
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
	"github.com/pepelazz/utils"
	"os"
	"html/template"
	"github.com/pepelazz/go-bot-utils"
)

func getLearningMap(userId int, args []string) (res interface{}, err error) {
	fmt.Printf("step1 args %s\n", args)
	if len(args) < 1 {
		err = errors.New(fmt.Sprintf("getLearningMap wrong signature: need (string, string), but receive %s", args))
		return
	}

	result := struct {
		Tasks struct {
			     Title string `json:"title"`
			     Tasks []struct {
				     Score      float64 `json:"score"`
				     Topic      string `json:"topic"`
				     Rivescript string `json:"rivescript"`
				     Progress   float64 `json:"progress"`
				     State      string `json:"state"`
			     } `json:"tasks"`
		     }`json:"tasks"`
		Rating []struct{
			Num int64 `json:"num"`
			User string `json:"user"`
			Score float64 `json:"score"`

		} `json:"rating"`
	}{}

	mapTitle := args[0]

	jsonStr, _ := json.Marshal(map[string]interface{}{"user_id": userId, "map_title": mapTitle})

	err = utils.CallPgFunc(Pg, "report_learning_map_for_user", jsonStr, &result, nil)
	if err != nil {
		return
	}

	t, err := template.New("learningMap").Parse(learningMapTmpl) // Create a template.
	if err != nil {
		return
	}

	path := fmt.Sprintf("tmp/learningMap_%s_%v.html", mapTitle, userId)
	err = goBotUtils.CreateFile(path)
	if err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return
	}
	defer f.Close()
	err = t.Execute(f, result)
	if err != nil {
		fmt.Println(err)
	}

	res = LearningMapFile{path}

	return
}

var learningMapTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <title>Карта обучения {{.Tasks.Title}}</title>
</head>
<body>
<div class="container-fluid">
    <a href="http://nl-a.ru/"><img src="http://nl-a.ru/wp-content/uploads/2016/08/NLA-title-2.png" alt="" style="margin-top: 20px"></a>
    <h3>Список задач: {{.Tasks.Title}}</h3>


    <div class="row">
	<div class="col-md-12">
	    <div class="panel panel-info" style="margin-top: 20px">
		<div class="panel-heading">
		    <h3 class="panel-title">Пройденные задания</h3>
		</div>
		<div class="panel-body">
		{{range .Tasks.Tasks}}
			{{if eq .State "finished"}}
				<pre><b>{{.Rivescript}}</b> {{if .Score}}<br>баллов:{{.Score}}{{end}} {{if .Progress}}<br>пройдено:{{.Progress}}%{{end}}
				</pre>
			{{end}}
		{{end}}
		</div>
	    </div>
	</div>
    </div>
    <div class="row">
	<div class="col-md-12">
	    <div class="panel panel-info" style="margin-top: 20px">
		<div class="panel-heading">
		    <h3 class="panel-title">Открытые задания</h3>
		</div>
		<div class="panel-body">
		{{range .Tasks.Tasks}}
			{{if eq .State "inProcess"}}
				<pre><b>{{.Rivescript}}</b> {{if .Score}}<br>баллов:{{.Score}}{{end}} {{if .Progress}}<br>пройдено:{{.Progress}}%{{end}}
				</pre>
			{{end}}
		{{end}}
		</div>
	    </div>
	</div>
    </div>

    <div class="row">
	<div class="col-md-12">
	    <div class="panel panel-info" style="margin-top: 20px">
		<div class="panel-heading">
		    <h3 class="panel-title">Рейтинг</h3>
		</div>
		<div class="panel-body">
			<ul class="list-group">
			{{range .Rating}}
				<li class="list-group-item">
				<span class="badge">{{.Score}}</span>
				{{.Num}}) {{.User}}
				</li>
			{{end}}
			</ul>
		</div>
	    </div>
	</div>
    </div>

</div>
</body>
</html>
`
