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
