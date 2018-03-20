package maintenance

type Client interface {
	GetMaintenance() *maintenanceMode
	GetAllowedIPs() ([]string, error)
	SetMessage([]byte) error
	SetAllowedIPs([]string) error
	Disable() error
	DisableAllowedIPs() error
}
