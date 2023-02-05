package video

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	"train-tiktok/service/video/videoclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewPublishLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishReq) (resp *types.Resp, err error) {
	var UserId = l.ctx.Value("user_id").(int64)
	var _fileBaseDir = l.svcCtx.VideoTmpPath
	var _fileTypNotSupport = types.Resp{
		Code: 199,
		Msg:  "不支持的文件类型",
	}

	var SystemErrResp = errx.HandleRpcErr(errorx.ErrSystemError)

	// 限制文件大小 150M
	err = l.r.ParseMultipartForm(150 << 20)
	if err != nil {
		logx.Errorf("parse form failed: %v", err)
		return &SystemErrResp, nil
	}

	// 从请求中获取文件句柄
	var file multipart.File
	var header *multipart.FileHeader
	if file, header, err = l.r.FormFile("data"); err != nil {
		logx.Errorf("get form data failed: %v", err)
		return &SystemErrResp, nil
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// 判断 文件目录是否存在
	if _, err := os.Stat(l.svcCtx.VideoTmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(l.svcCtx.VideoTmpPath, 0755); err != nil {
			logx.Errorf("mkdir %s failed: %v", l.svcCtx.VideoTmpPath, err)
			return &SystemErrResp, nil
		}
	}

	// 通过文件 filename 判断是否为视频
	if !tool.IsVideo(header.Filename) {
		logx.Infof("不支持的文件类型: %s", header.Filename)
		return &types.Resp{
			Code: 1,
			Msg:  "不支持的文件类型",
		}, nil
	}

	// 判断文件名是否存在安全风险
	if tool.IsFilenameDangerous(header.Filename) {
		logx.Infof("文件名存在安全风险: %s", header.Filename)
		return &_fileTypNotSupport, nil
	}

	// 生成文件路径
	logx.Info(header.Filename, req.Title)
	_timeMd5 := strconv.Itoa(int(time.Now().UnixNano()))
	_titleMd5 := tool.Md5(req.Title)
	_filenameMd5 := tool.Md5(header.Filename)
	_fileExt := tool.GetFileExt(header.Filename)
	if _fileExt == "" {
		return &_fileTypNotSupport, nil
	}
	_fileTmpPath := fmt.Sprintf("%s/%d_%s_%s_%s.%s", _fileBaseDir, UserId, _filenameMd5, _titleMd5, _timeMd5, _fileExt)

	// 打开临时文件句柄
	var f *os.File
	if f, err = os.OpenFile(_fileTmpPath, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		logx.Errorf("open file failed: %v", err)
		return &SystemErrResp, nil
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// 先获取文件开头字符 判断是否为视频
	var bufHead = make([]byte, 512)
	if n, _ := file.Read(bufHead); n == 0 {
		closeAndRemove(f, _fileTmpPath)
		return &_fileTypNotSupport, nil
	}
	if !tool.IsVideoByHead(bufHead) {
		closeAndRemove(f, _fileTmpPath)
		return &_fileTypNotSupport, nil
	}

	// 1M 分割 写文件
	buf := make([]byte, 1<<20)
	if _, err := f.Write(bufHead); err != nil {
		closeAndRemove(f, _fileTmpPath)
		return &SystemErrResp, nil
	}
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		if _, err := f.Write(buf[:n]); err != nil {
			closeAndRemove(f, _fileTmpPath)
			return &SystemErrResp, nil
		}
	}

	// 请求 video service 存储文件
	request, err := l.svcCtx.VideoRpc.Publish(l.ctx, &videoclient.PublishReq{
		Title:    header.Filename,
		FilePath: _fileTmpPath,
		UserId:   UserId,
	})

	logx.Info(request)
	return &types.Resp{
		Code: 0,
		Msg:  "success",
	}, nil

	// let video service to handle the file

	// delete
	// err = os.Remove("/tmp/" + header.Filename)
	// if err != nil {
	// 	logx.Error(err)
	// 	return err
	// }

}

// closeAndRemove 关闭文件句柄并删除已经创建的临时文件
func closeAndRemove(f *os.File, path string) {
	if err := f.Close(); err != nil {
		logx.Errorf("close file failed: %v", err)
	}
	if err := os.Remove(path); err != nil {
		logx.Errorf("remove tmp file failed: %v", err)
	}
}