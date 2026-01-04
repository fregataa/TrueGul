package mq

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	client     *redis.Client
	streamName string
}

func NewRedisPublisher(redisURL, streamName string) (*RedisPublisher, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisPublisher{
		client:     client,
		streamName: streamName,
	}, nil
}

func (p *RedisPublisher) Publish(ctx context.Context, task AnalysisTask) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamName,
		Values: map[string]interface{}{
			"task": string(taskJSON),
		},
	}).Err()
}

func (p *RedisPublisher) Close() error {
	return p.client.Close()
}
