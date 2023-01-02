package redis

import (
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 432
	PostPerAge               = 20
)

/*
投票算法：http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html
*/

/* PostVote 为帖子投票
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况
	1.之前没投过票，现在要投赞成票
	2.之前投过反对票，现在要改为赞成票
v=0时，有两种情况
	1.之前投过赞成票，现在要取消
	2.之前投过反对票，现在要取消
v=-1时，有两种情况
	1.之前没投过票，现在要投反对票
	2.之前投过赞成票，现在要改为反对票
*/

func PostVote(postID, userID string, v float64) (err error) {
	// 1. 取帖子发布时间
	postTime := rdb.ZScore(KeyPostTimeZSet, postID).Val()
	if float64(time.Now().Unix())-postTime > float64(OneWeekInSeconds) {
		// 不允许投票了
		return ErrorVoteTimeExpire
	}
	// 判断是否已经投过票
	key := KeyPostVotedZSetPrefix + postID
	ov := rdb.ZScore(key, userID).Val() // 获取当前分数

	diffAbs := math.Abs(ov - v)
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(key, redis.Z{ // 记录已投票
		Score:  v,
		Member: userID,
	})
	pipeline.ZIncrBy(KeyPostScoreZSet, VoteScore*diffAbs*v, postID) // 更新分数

	switch math.Abs(ov) - math.Abs(v) {
	case 1:
		// 取消投票 ov=1/-1 v=0
		// 投票数-1
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -1)
	case 0:
		// 反转投票 ov=-1/1 v=1/-1
		// 投票数不用更新
	case -1:
		// 新增投票 ov=0 v=1/-1
		// 投票数+1
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 1)
	default:
		// 已经投过票了
		return ErrorVoted
	}
	_, err = pipeline.Exec()
	return
}

// CreatePost 使用hash存储帖子信息
func CreatePost(postID, userID, title, summary, communityName string) (err error) {
	// 时间戳
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + postID
	communityKey := KeyCommunityPostSetPrefix + communityName
	postInfo := map[string]interface{}{
		"title":    title,
		"summary":  summary,
		"post:id":  postID,
		"user:id":  userID,
		"time":     now,
		"votes":    1,
		"comments": 0,
	}

	// 事务操作
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(votedKey, redis.Z{ // 作者默认投赞成票
		Score:  1,
		Member: userID,
	})
	pipeline.Expire(votedKey, time.Second*OneWeekInSeconds) // 一周时间

	pipeline.HMSet(KeyPostInfoHashPrefix+postID, postInfo)
	pipeline.ZAdd(KeyPostScoreZSet, redis.Z{ // 添加到分数的ZSet
		Score:  now + VoteScore,
		Member: postID,
	})
	pipeline.ZAdd(KeyPostTimeZSet, redis.Z{ // 添加到时间的ZSet
		Score:  now,
		Member: postID,
	})
	pipeline.SAdd(communityKey, postID) // 添加到对应版块
	_, err = pipeline.Exec()
	return
}
