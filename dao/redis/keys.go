package redis

// 存放redis keys

// 实际业务中利用命名空间来作为所对应的key，方便理解和拆分
const (
	KeyPrefix          = "bluebell:"   // 项目前缀key
	KeyPostTimeZSet    = "post:time"   // 帖子和发表时间-->对应于展示最新帖子的业务
	KeyPostScoreZSet   = "post:score"  //帖子和投票分数 -->对应于展示最热帖子业务
	KeyPostVotedZSetPF = "post:voted:" // 帖子和投票结果 后面需要提供参数    比如每个帖子记录了投票用户id和投票结果
	KeyCommunitySetPF  = "community:"  // 保存帖子的id，每个分区下的帖子的id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
