package dal

import (
	"context"
	"testing"
	"time"

	"eigenflux_server/pkg/milestone"
	milestonedal "eigenflux_server/pkg/milestone/dal"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMilestoneRuleDALDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&milestonedal.MilestoneRule{}))
	return db
}

func setupMilestoneRuleDALRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisClient = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = redisClient.Close()
		redisClient = nil
		mr.Close()
	})

	return mr
}

func TestCreateMilestoneRulePublishesInvalidation(t *testing.T) {
	db := setupMilestoneRuleDALDB(t)
	mr := setupMilestoneRuleDALRedis(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	received := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		errCh <- milestone.SubscribeRuleInvalidation(ctx, redisClient, func(metricKey string) {
			received <- metricKey
			cancel()
		})
	}()

	require.Eventually(t, func() bool {
		return len(mr.PubSubChannels("")) == 1
	}, time.Second, 10*time.Millisecond)

	rule, err := CreateMilestoneRule(context.Background(), db, CreateMilestoneRuleParams{
		MetricKey:       milestone.MetricConsumed,
		Threshold:       2,
		RuleEnabled:     true,
		ContentTemplate: `Your Content "{{.ItemSummary}}" reached {{.CounterValue}} consumptions. Item Id {{.ItemID}}`,
	})
	require.NoError(t, err)
	require.NotNil(t, rule)

	select {
	case metricKey := <-received:
		assert.Equal(t, milestone.MetricConsumed, metricKey)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for milestone invalidation publish")
	}

	require.NoError(t, <-errCh)
}

func TestReplaceMilestoneRulePublishesOldAndNewMetricInvalidation(t *testing.T) {
	db := setupMilestoneRuleDALDB(t)
	mr := setupMilestoneRuleDALRedis(t)

	oldRule := milestonedal.MilestoneRule{
		RuleID:          1,
		MetricKey:       milestone.MetricConsumed,
		Threshold:       2,
		RuleEnabled:     true,
		ContentTemplate: `old`,
		CreatedAt:       1000,
		UpdatedAt:       1000,
	}
	require.NoError(t, db.Create(&oldRule).Error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	received := make(chan string, 2)
	errCh := make(chan error, 1)
	go func() {
		errCh <- milestone.SubscribeRuleInvalidation(ctx, redisClient, func(metricKey string) {
			received <- metricKey
			if len(received) == 2 {
				cancel()
			}
		})
	}()

	require.Eventually(t, func() bool {
		return len(mr.PubSubChannels("")) == 1
	}, time.Second, 10*time.Millisecond)

	oldRuleResp, newRuleResp, err := ReplaceMilestoneRule(context.Background(), db, 1, ReplaceMilestoneRuleParams{
		MetricKey:       milestone.MetricScore2,
		Threshold:       3,
		RuleEnabled:     true,
		ContentTemplate: `new`,
	})
	require.NoError(t, err)
	require.NotNil(t, oldRuleResp)
	require.NotNil(t, newRuleResp)

	var metrics []string
	require.Eventually(t, func() bool {
		for len(received) > 0 {
			metrics = append(metrics, <-received)
		}
		return len(metrics) == 2
	}, time.Second, 10*time.Millisecond)

	assert.ElementsMatch(t, []string{milestone.MetricConsumed, milestone.MetricScore2}, metrics)
	require.NoError(t, <-errCh)
}
