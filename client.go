package maintenance

const MaintenanceKey = "maintenance"
const AllowedIPsKey = "maintenance_allowed_ips"

type Client interface {
	GetMaintenance() *maintenanceMode
	GetAllowedIPs() ([]string, error)
	SetMessage([]byte) error
	SetAllowedIPs([]string) error
	Disable() error
	DisableAllowedIPs() error
}
