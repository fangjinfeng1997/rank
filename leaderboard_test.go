package rank

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLeaderboardImpl(t *testing.T) {
	lbImpl := NewLeaderboardImpl()

	now := time.Now()
	for i := 0; i < 50; i++ {
		plaerId := fmt.Sprintf("player%03d", i)
		score := rand.Int63n(200)
		timestamp := now.Add(time.Duration(rand.Int63n(120)-60) * time.Second)
		lbImpl.UpdateScore(plaerId, score, timestamp)
	}

	fmt.Println(lbImpl.sl.String())
	fmt.Println(lbImpl.GetTopN(10))
	fmt.Println(lbImpl.GetPlayerRankRange("player000", 10))
	fmt.Println(lbImpl.GetPlayerRank("player000"))
}
