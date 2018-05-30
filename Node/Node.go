package spider

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type Node struct {
	child []*Node //子目录们

	url        string
	title      string
	statusCode int
	isutf8     bool
	encoding   string //`gbk`
	content    *string
}

var (
	reg = regexp.MustCompile(`<a href="(.+?html)">(.+?)</a>`)
)

func Run(url string) {
	n := &Node{}
	n.setUrl(url).httpGet().getChildsNode()
	GetNodes(n)
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
	n.content = &str
	return n
}

func (n *Node) getChildsNode() *Node {
	if n.statusCode == 200 {
		n.readContent(n.content)
	}
	return n
}

//解析出子url
func (n *Node) readContent(contents *string) {
	matches := reg.FindAllStringSubmatch(*contents, -1)
	childs := make([]*Node, len(matches))

	for i, item := range matches {
		childs[i] = &Node{
			url:   n.url + item[1],
			title: item[2],
		}
	}
	n.child = childs
}
