package spider

import (
	"log"
	"regexp"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ptnRepx    = regexp.MustCompile(`&nbsp;&nbsp;`)
	ptnHTMLTag = regexp.MustCompile(`(?s)</?.*?>`)
)

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=68957423")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Nodes{})
}

type Nodes struct {
	gorm.Model
	Url      string `json:"url"`
	Tittle   string `json:"Tittle"`
	Contents string `json:"Contents"`
}

func GetNodes(P *Node) {

	lens := len(P.child)
	var wg sync.WaitGroup
	wg.Add(lens)

	for _, item := range P.child {
		go goGetContent(item.url, item.title, &wg)
	}
	wg.Wait()
}

func goGetContent(url string, t string, wg *sync.WaitGroup) {
	defer wg.Done()
	n := &Node{}
	n.setUrl(url).httpGet()
	if n.statusCode == 200 {
		content := readContent2(n.content)
		//ioutil.WriteFile(t, []byte(content), 0644)
		db.Create(&Nodes{Url: url, Tittle: t, Contents: content})
	}
}

func readContent2(content *string) (contents string) {
	dialog := regexp.MustCompile(`<div id="BookText">(.+?)</div>`)
	s := dialog.FindAllString(*content, 100)
	contents = join(s)
	contents = ptnHTMLTag.ReplaceAllString(contents, "\r\n")
	contents = ptnRepx.ReplaceAllString(contents, " ")
	return
}

func join(s []string) (content string) {
	for _, val := range s {
		content += val
	}
	return
}
