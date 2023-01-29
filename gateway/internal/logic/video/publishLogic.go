package video

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
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

func (l *PublishLogic) Publish(req *types.PublishReq) (resp *videoclient.PublishResp, err error) {
	// 限制文件大小 150M
	err = l.r.ParseMultipartForm(150 << 20)
	if err != nil {
		logx.Error(err)
		return &videoclient.PublishResp{}, errorx.ErrSystemError
	}

	// 从请求中获取文件句柄
	var file multipart.File
	var header *multipart.FileHeader
	if file, header, err = l.r.FormFile("data"); err != nil {
		logx.Error(err)
		return &videoclient.PublishResp{}, errorx.ErrSystemError
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// 生成文件路径
	logx.Info(header.Filename)
	_tmp := l.svcCtx.VideoTmpPath
	_FilenameMd5 := tool.Md5(header.Filename)
	_fileTmpPath := _tmp + "/" + _FilenameMd5

	// save file to /tmp
	var f *os.File
	if f, err = os.OpenFile(_fileTmpPath, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		logx.Error(err)
		return &videoclient.PublishResp{}, errorx.ErrSystemError
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// 1M 分割 读取
	buf := make([]byte, 1<<20)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		if _, err := f.Write(buf[:n]); err != nil {
			return &videoclient.PublishResp{}, errorx.ErrSystemError
		}
		log.Printf("read %d bytes, percent %d%%\n", n, n*100/1024)
	}

	//request, err := l.svcCtx.VideoRpc.Publish(l.ctx, &videoclient.PublishReq{
	//	Title: header.Filename,
	//	Path:  _fileTmpPath,
	//})
	return &videoclient.PublishResp{Response: &videoclient.Resp{
		StatusCode: 0,
		StatusMsg:  "success",
	}}, nil

	// let video service to handle the file

	// delete
	// err = os.Remove("/tmp/" + header.Filename)
	// if err != nil {
	// 	logx.Error(err)
	// 	return err
	// }

}
