package tool

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"path"
	"sort"
	"strings"
)

var LegalVideoExt = []string{
	"mp4",
	"avi",
	"rmvb",
	"rm",
	"flv",
	"mkv",
	"wmv",
	"mpg",
	"mpeg",
	"3gp",
	"mov",
	"dat",
	"asf",
	"vob",
	"m4v",
}

var IllegalFileNames = []string{
	"..",
	"/",
	"\\",
	":",
	"*",
	"?",
	"\"",
	"<",
	">",
	"|",
}

// IsVideo 通过文件名判断是否是视频
func IsVideo(filename string) bool {
	ext := GetFileExt(filename)
	if inArray(ext, LegalVideoExt) {
		return true
	} else {
		return false
	}
}

// IsVideoByHead 通过文件头判断是否是视频
func IsVideoByHead(buf []byte) bool {
	filetype := http.DetectContentType(buf)
	if strings.Contains(filetype, "video") {
		logx.Infof("支持的文件类型: %s", filetype)
		return true
	} else {
		logx.Infof("不支持的文件类型: %s", filetype)
		return false
	}
}

func IsFilenameDangerous(filename string) bool {
	for _, illegal := range IllegalFileNames {
		if strings.Contains(filename, illegal) {
			return true
		}
	}
	return false
}

// GetFileExt 获取文件扩展名
func GetFileExt(filename string) string {
	// 防止没有后缀的文件
	if !strings.Contains(filename, ".") {
		return ""
	}
	return strings.ToLower(path.Ext(filename)[1:])
}

// inArray 判断字符串是否在字符串数组中
func inArray(str string, strArray []string) bool {
	sort.Strings(strArray)
	i := sort.SearchStrings(strArray, str)
	if i < len(strArray) && strArray[i] == str {
		return true
	} else {
		return false
	}
}
