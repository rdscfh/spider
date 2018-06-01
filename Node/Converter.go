package spider

import (
	"github.com/axgle/mahonia"
)

//gbk to utf8
func ConvertToString(src []byte, srcCode string, tagCode string) (result []byte) {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(string(src))
	tagCoder := mahonia.NewDecoder(tagCode)
	_, result, _ = tagCoder.Translate([]byte(srcResult), true)
	return
}
