package rediskeyutil

import "fmt"

type Keyer interface {
	GetVideoFavoriteKey(videoId int64) string
	GetUserKey(userId int64) string
	GetPublisherFavoriteKey(videoPublisherId int64) string
	GetCommentCount(videoId int64) string
}

type Keys struct {
	RedisPrefix string
}

func NewKeys(rPrefix string) *Keys {
	return &Keys{
		RedisPrefix: rPrefix,
	}
}

// GetVideoFavoriteKey 记录视频点赞数
// 数据类型: SET
// Incr / Decr
func (k *Keys) GetVideoFavoriteKey(videoId int64) string {
	return fmt.Sprintf("%s:favorite_count:%d", k.RedisPrefix, videoId)
}

// GetUserKey 记录用户是否点赞该视频
// 数据类型: ZSET
// ZScore/ ZIncr / ZDecr / ZCard
func (k *Keys) GetUserKey(userId int64) string {
	return fmt.Sprintf("%s:favorite_user:%d", k.RedisPrefix, userId)
}

// GetPublisherFavoriteKey 记录视频发布用户的 获赞数
// 数据类型: SET
// Incr / Decr
func (k *Keys) GetPublisherFavoriteKey(videoPublisherId int64) string {
	return fmt.Sprintf("%s:user_favorited:%d", k.RedisPrefix, videoPublisherId)
}

// GetCommentCount 视频评论量 key
// 数据类型: SET
// Incr / Decr
func (k *Keys) GetCommentCount(videoId int64) string {
	return fmt.Sprintf("%s:comment_count:%d", k.RedisPrefix, videoId)
}
