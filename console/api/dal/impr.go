package dal

import (
	"context"
	"errors"
	"sort"

	"eigenflux_server/pkg/impr"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var redisClient *redis.Client

type AgentImprRecord struct {
	ItemIDs  []int64
	GroupIDs []int64
	URLs     []string
	Items    []ItemWithProcessed
}

func InitRedis(addr, password string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	ctx := context.Background()
	return redisClient.Ping(ctx).Err()
}

func GetAgentImprRecord(ctx context.Context, db *gorm.DB, agentID int64) (*AgentImprRecord, error) {
	if redisClient == nil {
		return nil, errors.New("redis client not initialized")
	}

	seen, err := impr.GetSeenItems(ctx, redisClient, agentID)
	if err != nil {
		return nil, err
	}

	sort.Slice(seen.ItemIDs, func(i, j int) bool { return seen.ItemIDs[i] > seen.ItemIDs[j] })
	sort.Slice(seen.GroupIDs, func(i, j int) bool { return seen.GroupIDs[i] < seen.GroupIDs[j] })
	sort.Strings(seen.URLs)

	items, err := ListItemsByIDs(db, seen.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &AgentImprRecord{
		ItemIDs:  seen.ItemIDs,
		GroupIDs: seen.GroupIDs,
		URLs:     seen.URLs,
		Items:    items,
	}, nil
}
