package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"github.com/axgle/mahonia"
	"sync"
)

var (
	ptnIndexItem    = regexp.MustCompile(`<a target="_blank" href="(.+\.html)" title=".+" >(.+)</a>`)
	ptnContentRough = regexp.MustCompile(`(?s).*<div class="artcontent">(.*)<div id="zhanwei">.*`)
	ptnBrTag        = regexp.MustCompile(`<br>`)
	ptnHTMLTag      = regexp.MustCompile(`(?s)</?.*?>`)
	ptnSpace        = regexp.MustCompile(`(^\s+)|( )`)
	ptnAhref		= regexp.MustCompile(`<a href="(.+?html)">(.+?)</a>`)
	ptnAhrefNext	= regexp.MustCompile(`<a id="book-next" href="(.+?)">(.+?)</a>`)
	ptnRepx			= regexp.MustCompile(`&nbsp;&nbsp;`)
	ptnFooterRegexp = regexp.MustCompile(`<div id="footer">(.+?)</div>`)
	ptnLink 		= regexp.MustCompile(`<div class="link">(.+?)</div>`)

	IndexPage		="http://www.uuxs.la/book/42/42677/"
)

func Get(url string) (content string, statusCode int) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = string(data)
	return
}

type IndexItem struct {
	url   string
	title string
}

//gbk转utf8
func ConvertToString(src string, srcCode string, tagCode string)(result string){
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result = string(cdata)
    return
}

//var chanel=make(chan IndexItem)

func findIndex(content string) (indexs chan IndexItem,length int, err error) {
	content  = ConvertToString(content, "gbk", "utf-8")
	matches:=ptnAhref.FindAllStringSubmatch(content,-1);
	length=len(matches)
	indexs = make(chan IndexItem, length)
	for _, item := range matches {
		indexs <- IndexItem{IndexPage + item[1], item[2]}
	}
	return
}
//读取url中所有标签为
func readContent(url string) (content string) {
	content, statusCode := Get(url)
	if statusCode != 200 {
		fmt.Print("Fail to get the raw data from", url, "\n")
		return
	}
	content  = ConvertToString(content, "gbk", "utf-8")
	dialog := regexp.MustCompile(`<div id="BookText">(.+?)</div>`)
	s:=dialog.FindAllString(content,100)
	content=Join(s)
	content= ptnHTMLTag.ReplaceAllString(content,"\r\n")
	content =ptnRepx.ReplaceAllString(content," ")
	return
}
func Join(s []string)(content string){
	for _,val:=range s{
		content+=val
	}
	return
}

var wg sync.WaitGroup

func main() {
	fmt.Println(`Get index ...`)
	s, statusCode := Get(IndexPage)
	if statusCode != 200 {
		return
	}

	ch,len, _ := findIndex(s)
	wg.Add(len)
	for w := range ch{
		go goContents(w)
	}

	close(ch)
	wg.Wait()
}

func goContents(ch IndexItem)  {
	defer wg.Done();
	fileName := fmt.Sprintf("./m/%s.txt",ch.title)
	content := readContent(ch.url)
	ioutil.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Finish writing to %s.\n", fileName)
}