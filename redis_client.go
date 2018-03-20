package maintenance

import (
	"github.com/go-redis/redis"
)

type RedisClient struct {
	rc *redis.Client
}

func NewRedisClient(opt *redis.Options) Client {
	client := redis.NewClient(opt)
	return &RedisClient{
		rc: client,
	}
}

func (c *RedisClient) getMessage() []byte {
	val, err := c.rc.Get(MaintenanceKey).Result()
	if err != nil {
		return nil
	}
	return []byte(val)
}

func (c *RedisClient) GetMaintenance() *maintenanceMode {
	msg := c.getMessage()
	if msg == nil {
		return &maintenanceMode{
			enabled: false,
		}
	}

	return &maintenanceMode{
		enabled: true,
		message: msg,
	}
}

func (c *RedisClient) GetAllowedIPs() ([]string, error) {
	return c.rc.LRange(AllowedIPsKey, 0, -1).Result()
}

func (c *RedisClient) SetMessage(msg []byte) error {
	return c.rc.Set(MaintenanceKey, msg, 0).Err()
}

func (c *RedisClient) SetAllowedIPs(ips []string) error {
	pipe := c.rc.TxPipeline()

	pipe.Del(AllowedIPsKey)
	for _, ip := range ips {
		pipe.LPush(AllowedIPsKey, ip)
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) Disable() error {
	return c.rc.Del(AllowedIPsKey, MaintenanceKey).Err()
}

func (c *RedisClient) DisableAllowedIPs() error {
	return c.rc.Del(AllowedIPsKey).Err()
}
