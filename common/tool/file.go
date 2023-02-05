package tool

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"os/exec"
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

// GenerateVideoCover 生成视频封面 / 使用ffmpeg取视频第一帧
func GenerateVideoCover(videoPath string, coverPath string) error {
	// ffmpeg -i 1.mp4 -y -f mjpeg -ss 1 -t 0.001 // -s 320x240 1.jpg
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-y", "-f", "mjpeg", "-ss", "1", "-t", "0.001", coverPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// IsFilenameDangerous 判断文件名是否有危险字符
func IsFilenameDangerous(filename string) bool {
	for _, illegal := range IllegalFileNames {
		if strings.Contains(filename, illegal) {
			return true
		}
	}
	return false
}

// CheckPathOrCreate 检查路径是否存在，不存在则创建
func CheckPathOrCreate(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
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
