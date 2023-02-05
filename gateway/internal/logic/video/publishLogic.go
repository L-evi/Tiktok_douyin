package video

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

	// check if _fileTmpPath exists
	if _, err := os.Stat(l.svcCtx.VideoTmpPath); os.IsNotExist(err) {
		if err := os.Mkdir(l.svcCtx.VideoTmpPath, 0755); err != nil {
			logx.Errorf("mkdir %s failed", l.svcCtx.VideoTmpPath)
			return &SystemErrResp, nil
		}
	}

	// 生成文件路径
	logx.Info(header.Filename)
	_tmp := l.svcCtx.VideoTmpPath
	_FilenameMd5 := tool.Md5(header.Filename)
	_fileTmpPath := _tmp + "/" + _FilenameMd5

	// save file to /tmp
	var f *os.File
	if f, err = os.OpenFile(_fileTmpPath, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		logx.Errorf("open file failed: %v", err)
		return &SystemErrResp, nil
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
			return &SystemErrResp, nil
		}
		log.Printf("read %d bytes, percent %d%%\n", n, n*100/1024)
	}

	request, err := l.svcCtx.VideoRpc.Publish(l.ctx, &videoclient.PublishReq{
		Title:    header.Filename,
		FilePath: _fileTmpPath,
		UserId:   l.ctx.Value("userId").(int64),
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
