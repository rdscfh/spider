package spider

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=68957423")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Node{})
}

/*
type Nodes struct {
	gorm.Model
	sync.Mutex
	Ch       []*Nodes
	Url      string `json:"url"`
	Tittle   string `json:"Tittle"`
	Contents string `json:"Contents"`
}

func NewNodes(url string, title string, contents string) *Nodes {
	return &Nodes{
		Url:      url,
		Tittle:   title,
		Contents: contents,
	}
}
func (n *Nodes) Push(c *Nodes) {
	n.Lock()
	if len(n.Ch) < 100 {
		n.Ch = append(n.Ch, c)
	} else {

	}
	n.Unlock()
}

func (n *Nodes) saveBatch() {
	//for
	//db.Create()
}

func (n *Node) GetNodes() {

	lens := len(n.child)
	var wg sync.WaitGroup
	wg.Add(lens)

	for _, item := range n.child {
		go goGetContent(item.url, item.title, &wg)
	}
	wg.Wait()
	db.Close()
}

func goGetContent(url string, t string, wg *sync.WaitGroup) {
	defer wg.Done()
	n := &Node{}
	n.setUrl(url).httpGet()
	if n.statusCode == 200 {
		content := readContent2(n.content)
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
*/
