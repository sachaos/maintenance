package maintenance

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedClient struct {
	URL string
	mc  *memcache.Client
}

func NewMemcachedClient(url string) Client {
	return &MemcachedClient{
		URL: url,
		mc:  memcache.New(url),
	}
}

func (c *MemcachedClient) getMessage() []byte {
	item, err := c.mc.Get(MaintenanceKey)
	if err != nil {
		return nil
	}
	return item.Value
}

func (c *MemcachedClient) GetMaintenance() *maintenanceMode {
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

func (c *MemcachedClient) GetAllowedIPs() ([]string, error) {
	var ips []string
	m, err := c.mc.Get(AllowedIPsKey)
	if err != nil {
		return ips, nil
	}

	if err := json.Unmarshal(m.Value, &ips); err != nil {
		return ips, err
	}
	return ips, nil
}

func (c *MemcachedClient) SetMessage(msg []byte) error {
	return c.mc.Set(&memcache.Item{Key: MaintenanceKey, Value: msg})
}

func (c *MemcachedClient) SetAllowedIPs(ips []string) error {
	allowedIPs, err := json.Marshal(ips)
	if err != nil {
		return err
	}
	return c.mc.Set(&memcache.Item{Key: AllowedIPsKey, Value: allowedIPs})
}

func (c *MemcachedClient) Disable() error {
	return c.mc.DeleteAll()
}

func (c *MemcachedClient) DisableAllowedIPs() error {
	return c.mc.Delete(AllowedIPsKey)
}
