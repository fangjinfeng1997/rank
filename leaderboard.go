package rank

import (
	"sync"
	"time"
)

const (
	less    int = -1
	equal   int = 0
	greater int = 1
)

type KeyInfo struct {
	PlayerId  string
	Score     int64
	TimeStamp time.Time
}

type RankInfo struct {
	Rank int
	KeyInfo
}

type LeaderboardService interface {
	// 更新玩家分数
	UpdateScore(playerId string, score int64, timestamp time.Time)

	// 获取玩家当前排名
	GetPlayerRank(playerId string) RankInfo

	// 获取排名前N名
	GetTopN(n int) []RankInfo

	// 获取玩家周边排名
	GetPlayerRankRange(playerId string, rangeNum int) []RankInfo
}

func Compare(a, b interface{}) int {
	lv, rv := a.(KeyInfo), b.(KeyInfo)
	if lv.Score > rv.Score {
		return less
	}
	if lv.Score < rv.Score {
		return greater
	}
	if lv.TimeStamp.Unix() < rv.TimeStamp.Unix() {
		return less
	}
	if lv.TimeStamp.Unix() > rv.TimeStamp.Unix() {
		return greater
	}
	return equal
}

type LeaderboardImpl struct {
	sl             *SkipList
	player2KeyInfo map[string]KeyInfo
	mu             sync.Mutex
}

func NewLeaderboardImpl() *LeaderboardImpl {
	return &LeaderboardImpl{
		sl:             NewSkipList(Compare),
		player2KeyInfo: make(map[string]KeyInfo),
	}
}

// 更新玩家分数
func (l *LeaderboardImpl) UpdateScore(playerId string, score int64, timestamp time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ki := KeyInfo{
		PlayerId:  playerId,
		Score:     score,
		TimeStamp: timestamp,
	}
	l.player2KeyInfo[playerId] = ki
	l.sl.Insert(ki)
}

// 获取玩家当前排名
func (l *LeaderboardImpl) GetPlayerRank(playerId string) RankInfo {
	l.mu.Lock()
	defer l.mu.Unlock()

	ki := l.player2KeyInfo[playerId]
	rank, _ := l.sl.Search(ki)
	return RankInfo{
		rank, ki,
	}
}

// 获取排名前N名
func (l *LeaderboardImpl) GetTopN(n int) []RankInfo {
	l.mu.Lock()
	defer l.mu.Unlock()

	var result []RankInfo
	l.sl.Range(1, n, func(rank int, value interface{}) bool {
		v := value.(KeyInfo)
		result = append(result, RankInfo{rank, v})
		return true
	})
	return result
}

// 获取玩家周边排名
func (l *LeaderboardImpl) GetPlayerRankRange(playerId string, rangeNum int) []RankInfo {
	if rangeNum < 0 {
		return nil
	}
	ri := l.GetPlayerRank(playerId)

	l.mu.Lock()
	defer l.mu.Unlock()

	halfCnt := (rangeNum - 1) / 2
	start := ri.Rank - halfCnt
	if start <= 0 {
		start = 1
	}
	var result []RankInfo
	l.sl.Range(start, rangeNum, func(rank int, value interface{}) bool {
		v := value.(KeyInfo)
		result = append(result, RankInfo{rank, v})
		return true
	})
	return result
}
