package spider

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type Node struct {
	child      []*Node //子目录们
	url        string
	title      string
	statusCode int
	isutf8     bool
	encoding   string //`gbk`
	content    string
}

var (
	reg        = regexp.MustCompile(`<a href="(.+?html)">(.+?)</a>`)
	ptnRepx    = regexp.MustCompile(`&nbsp;&nbsp;`)
	ptnHTMLTag = regexp.MustCompile(`(?s)</?.*?>`)
)

func Run(url string) {
	n := &Node{}
	n.setUrl(url).httpGet().getChildsNode().getChildsContent()
}

func (n *Node) setUrl(url string) *Node {
	n.url = url
	return n
}

func (n *Node) httpGet() *Node {
	resp, err := http.Get(n.url)
	if err != nil {
		n.statusCode = -100
		return n
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		n.statusCode = -200
		return n
	}
	n.statusCode = resp.StatusCode
	rs := ConvertToString(data, "gbk", "utf-8")
	str := string(rs)
	n.content = str
	return n
}

func (n *Node) getChildsNode() *Node {
	if n.statusCode == 200 {
		n.readContent()
	}
	return n
}

//解析出子url
func (n *Node) readContent() {
	matches := reg.FindAllStringSubmatch(n.content, -1)
	childs := make([]*Node, len(matches))

	for i, item := range matches {
		childs[i] = &Node{
			url:   n.url + item[1],
			title: item[2],
		}
	}
	n.child = childs
}

func (n *Node) getChildsContent() {
	lens := len(n.child)
	var wg sync.WaitGroup
	wg.Add(lens)
	for _, item := range n.child {
		go item.goGetContent(&wg)
	}
	wg.Wait()
	db.Close()
}

func (n *Node) goGetContent(wg *sync.WaitGroup) {
	defer wg.Done()
	n.httpGet()
	if n.statusCode == 200 {
		n.readContent2()
		n.savetoPG()
	}
}

func (n *Node) readContent2() {
	dialog := regexp.MustCompile(`<div id="BookText">(.+?)</div>`)
	s := dialog.FindAllString(n.content, 100)
	contents := join(s)
	contents = ptnHTMLTag.ReplaceAllString(contents, "\r\n")
	contents = ptnRepx.ReplaceAllString(contents, " ")
	n.content = contents
}

func join(s []string) (content string) {
	for _, val := range s {
		content += val
	}
	return
}

func (n *Node) savetoPG() {
	//db.Create(n)
	ioutil.WriteFile(n.title+".txt", []byte(n.content), 0666) //写入文件(字节数组)
}
