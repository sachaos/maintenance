package maintenance

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
)

const MaintenanceKey = "maintenance"
const AllowedIPsKey = "maintenance_allowed_ips"

type Client struct {
	URL string
	mc  *memcache.Client
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
		mc:  memcache.New(url),
	}
}

func (c *Client) getMessage() []byte {
	item, err := c.mc.Get(MaintenanceKey)
	if err != nil {
		return nil
	}
	return item.Value
}

func (c *Client) getMaintenance() *maintenanceMode {
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

func (c *Client) getAllowedIPs() ([]string, error) {
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

func (c *Client) SetMessage(msg []byte) error {
	return c.mc.Set(&memcache.Item{Key: MaintenanceKey, Value: msg})
}

func (c *Client) SetAllowedIPs(ips []string) error {
	allowedIPs, err := json.Marshal(ips)
	if err != nil {
		return err
	}
	return c.mc.Set(&memcache.Item{Key: AllowedIPsKey, Value: allowedIPs})
}

func (c *Client) DeleteAll() error {
	return c.mc.DeleteAll()
}

func (c *Client) Delete(key string) error {
	return c.mc.Delete(key)
}
