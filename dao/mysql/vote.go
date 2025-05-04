package mysql

import (
	"database/sql"
)

// SavePostVote 保存投票记录
func SavePostVote(postID, userID int64, voteType int8) error {
	// 使用REPLACE INTO可以实现存在则更新，不存在则插入
	sqlStr := `REPLACE INTO post_vote (post_id, user_id, vote_type) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStr, postID, userID, voteType)
	if err != nil {
		return err
	}

	// 更新投票统计
	return updatePostVoteCount(postID)
}

// updatePostVoteCount 更新帖子投票统计
func updatePostVoteCount(postID int64) error {
	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 统计赞成票和反对票数量
	sqlStr := `SELECT 
		COUNT(CASE WHEN vote_type = 1 THEN 1 END) as up_count,
		COUNT(CASE WHEN vote_type = -1 THEN 1 END) as down_count
	FROM post_vote 
	WHERE post_id = ?`

	var upCount, downCount int
	err = tx.QueryRow(sqlStr, postID).Scan(&upCount, &downCount)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 更新统计表
	sqlStr = `REPLACE INTO post_vote_count (post_id, up_vote_count, down_vote_count) VALUES (?, ?, ?)`
	_, err = tx.Exec(sqlStr, postID, upCount, downCount)
	return err
}

// GetPostVoteCount 获取帖子投票统计
func GetPostVoteCount(postID int64) (upCount, downCount int, err error) {
	sqlStr := `SELECT up_vote_count, down_vote_count FROM post_vote_count WHERE post_id = ?`
	err = db.QueryRow(sqlStr, postID).Scan(&upCount, &downCount)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	return
}

// GetUserPostVote 获取用户对特定帖子的投票记录
func GetUserPostVote(postID, userID int64) (voteType int8, err error) {
	sqlStr := `SELECT vote_type FROM post_vote WHERE post_id = ? AND user_id = ?`
	err = db.QueryRow(sqlStr, postID, userID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}
