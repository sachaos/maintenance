package maintenance

type MaintenanceMode interface {
	IsEnabled() bool
	Disable()
}

type maintenanceMode struct {
	enabled bool
	message []byte
}

func (mm *maintenanceMode) IsEnabled() bool {
	return mm.enabled
}

func (mm *maintenanceMode) Disable() {
	mm.enabled = false
}
