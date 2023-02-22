package video

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/video/videoclient"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewPublishActionLogic(r *http.Request, ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {

	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq) (resp *types.Resp, err error) {
	var UserId = l.ctx.Value("user_id").(int64)
	var _fileBaseDir = l.svcCtx.PublicPath
	var _videoBaseDir = _fileBaseDir + "/video"
	var _coverBaseDir = _fileBaseDir + "/cover"
	var _fileTypNotSupport = types.Resp{
		Code: 199,
		Msg:  "不支持的文件类型",
	}
	var _fileTitleIllegal = types.Resp{
		Code: 198,
		Msg:  "标题不合法",
	}

	if req.Title == "" {
		return &_fileTitleIllegal, nil
	}

	var SystemErrResp = errx.HandleRpcErr(errorx.ErrSystemError)

	// 从请求中获取文件句柄
	var file multipart.File
	var header *multipart.FileHeader
	logx.WithContext(l.ctx).Debugf("Title: %s", req.Title)

	if file, header, err = l.r.FormFile("data"); err != nil {
		logx.Errorf("get form data failed: %v", err)

		return &SystemErrResp, nil
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	logx.WithContext(l.ctx).Debugf("收到上传请求: %s %s %s", req.Title, header.Filename, UserId)

	// 判断 文件目录是否存在
	if tool.CheckPathOrCreate(_videoBaseDir) != nil {
		logx.Errorf("mkdir %s failed: %v", _videoBaseDir, err)

		return &SystemErrResp, nil
	}
	if tool.CheckPathOrCreate(_coverBaseDir) != nil {
		logx.Errorf("mkdir %s failed: %v", _coverBaseDir, err)

		return &SystemErrResp, nil
	}

	// 通过文件 filename 判断是否为视频
	if !tool.IsVideo(header.Filename) {
		logx.Debugf("不支持的文件类型: %s", header.Filename)

		return &_fileTypNotSupport, nil
	}
	logx.WithContext(l.ctx).Debugf("文件类型: %s", header.Filename)

	// 判断文件名是否存在安全风险
	if tool.IsFilenameDangerous(header.Filename) {
		logx.Debugf("文件名存在安全风险: %s", header.Filename)

		return &_fileTitleIllegal, nil
	}

	// 判断标题合法性
	// TODO

	// 生成文件路径
	logx.WithContext(l.ctx).Debug(header.Filename, req.Title)

	_fileExt := tool.GetFileExt(header.Filename)
	if _fileExt == "" {
		return &_fileTypNotSupport, nil
	}

	_fileFullName := fmt.Sprintf("%d_%s_%s_%s", UserId, tool.Md5(header.Filename)[0:16], tool.Md5(req.Title)[1:16], strconv.Itoa(int(time.Now().UnixNano())))
	_fileTmpPath := fmt.Sprintf("%s/%s.%s", _videoBaseDir, _fileFullName, _fileExt)
	if _fileExt == "" {
		return &_fileTypNotSupport, nil
	}

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
	// 部分视频 无判断 //不知道为什么
	//if !tool.IsVideoByHead(bufHead) {
	//	closeAndRemove(f, _fileTmpPath)
	//
	//	return &_fileTypNotSupport, nil
	//}

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

	// 获取视频的 hash 用于去重
	var videoHash string
	if videoHash, err = tool.GetFileHash(_fileTmpPath); err != nil {
		closeAndRemove(f, _fileTmpPath)
		logx.WithContext(l.ctx).Debugf("获取文件hash失败: %v", err)

		return &SystemErrResp, nil
	}

	// 判断视频是否已经存在
	var hashrpc *videoclient.GetVideoByHashResp
	if hashrpc, err = l.svcCtx.VideoRpc.GetVideoByHash(l.ctx, &videoclient.GetVideoByHashReq{
		Hash: videoHash,
	}); err != nil {
		closeAndRemove(f, _fileTmpPath)
		logx.WithContext(l.ctx).Debugf("获取文件hash失败: %v", err)

		return &SystemErrResp, nil
	}

	// 秒存验证
	if hashrpc.Exists == true {
		closeAndRemove(f, _fileTmpPath)
		if _, err := l.svcCtx.VideoRpc.Publish(l.ctx, &videoclient.PublishReq{
			Title:     req.Title,
			FilePath:  hashrpc.Video.PlayUrl,
			CoverPath: hashrpc.Video.CoverUrl,
			UserId:    UserId,
			Hash:      videoHash,
		}); err != nil {
			return &SystemErrResp, nil
		}

		return &types.Resp{
			Code: 0,
			Msg:  "rapid upload success",
		}, nil
	}

	// 生成封面
	_coverPath := fmt.Sprintf("%s/%s.jpg", _coverBaseDir, _fileFullName)
	if err := tool.GenerateVideoCover(_fileTmpPath, _coverPath); err != nil {
		closeAndRemove(f, _fileTmpPath)

		return &SystemErrResp, nil
	}

	// 是否上传至 cos
	if l.svcCtx.Config.Cos.Enable {
		if remoteVideoUrl, remoteCoverUrl, err := toCos(l, _fileTmpPath, _coverPath); err != nil {
			closeAndRemove(f, _fileTmpPath)
			_ = os.Remove(_coverPath)

			return &SystemErrResp, nil
		} else {
			closeAndRemove(f, _fileTmpPath)
			_ = os.Remove(_coverPath)

			_fileTmpPath = remoteVideoUrl
			_coverPath = remoteCoverUrl
		}
	}

	// 请求 video service 存储文件
	if _, err := l.svcCtx.VideoRpc.Publish(l.ctx, &videoclient.PublishReq{
		Title:     req.Title,
		FilePath:  _fileTmpPath,
		CoverPath: _coverPath,
		UserId:    UserId,
		Hash:      videoHash,
	}); err != nil {
		closeAndRemove(f, _fileTmpPath)
		_ = os.Remove(_coverPath)

		return &SystemErrResp, nil
	}

	return &types.Resp{
		Code: 0,
		Msg:  "upload to server success",
	}, nil
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

func toCos(l *PublishActionLogic, videoPath string, coverPath string) (remoteVideo string, remoteCover string, err error) {
	bucketURL, _ := url.Parse(l.svcCtx.Config.Cos.BucketUrl)
	b := &cos.BaseURL{BucketURL: bucketURL}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  l.svcCtx.Config.Cos.SecretId,
			SecretKey: l.svcCtx.Config.Cos.SecretKey,
		},
	})

	videoKey := fmt.Sprintf("%s/video/%s", l.svcCtx.Config.Cos.Path, filepath.Base(videoPath))
	coverKey := fmt.Sprintf("%s/cover/%s", l.svcCtx.Config.Cos.Path, filepath.Base(coverPath))

	// 上传视频
	if videoResp, err := client.Object.PutFromFile(l.ctx, videoKey, videoPath, nil); err != nil || videoResp.StatusCode != http.StatusOK {
		logx.WithContext(l.ctx).Errorf("上传视频失败: %v %v", err, videoResp)
		return "", "", errorx.ErrSystemError
	}

	if coverResp, err := client.Object.PutFromFile(l.ctx, coverKey, coverPath, nil); err != nil || coverResp.StatusCode != http.StatusOK {
		logx.WithContext(l.ctx).Errorf("上传封面失败: %v %v", err, coverResp)
		_, _ = client.Object.Delete(l.ctx, videoKey)
		return "", "", errorx.ErrSystemError
	}

	remoteVideoFullPath := fmt.Sprintf("%s/%s", l.svcCtx.Config.Cos.BucketUrl, videoKey)
	remoteCoverFullPath := fmt.Sprintf("%s/%s", l.svcCtx.Config.Cos.BucketUrl, coverKey)

	return remoteVideoFullPath, remoteCoverFullPath, nil
}
