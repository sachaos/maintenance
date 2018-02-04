package maintenance

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tomasen/realip"
)

const MaintenanceKey = "maintenance"
const AllowedIPsKey = "maintenance_allowed_ips"

type Maintenance interface {
	SetMaintenance(next http.Handler) http.Handler
	AllowByIP(next http.Handler) http.Handler
	ResponseIfMaintenanceMode(next http.Handler) http.Handler
}

type maintenance struct {
	URL string
	mc  *memcache.Client
}

func NewMaintenance(url string) Maintenance {
	return &maintenance{
		URL: url,
		mc:  memcache.New(url),
	}
}

func (ms *maintenance) SetMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), MaintenanceKey, ms.getMaintenance())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ms *maintenance) AllowByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode := r.Context().Value(MaintenanceKey).(MaintenanceMode)
		if !mode.IsEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		ips, err := ms.getAllowedIPs()
		if err != nil {
			mode.Disable()
			next.ServeHTTP(w, r)
			return
		}

		ip := realip.FromRequest(r)
		for _, allowedIP := range ips {
			if ip == allowedIP {
				mode.Disable()
				next.ServeHTTP(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
		return
	})
}

func (ms *maintenance) ResponseIfMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode := r.Context().Value(MaintenanceKey).(*maintenanceMode)
		if mode.IsEnabled() {
			http.Error(w, string(mode.message), 503)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ms *maintenance) getMessage() []byte {
	item, err := ms.mc.Get(MaintenanceKey)
	if err != nil {
		return nil
	}
	return item.Value
}

func (ms *maintenance) getMaintenance() *maintenanceMode {
	msg := ms.getMessage()
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

func (ms *maintenance) getAllowedIPs() ([]string, error) {
	var ips []string
	m, err := ms.mc.Get(AllowedIPsKey)
	if err != nil {
		return ips, nil
	}

	if err := json.Unmarshal(m.Value, &ips); err != nil {
		return ips, err
	}
	return ips, nil
}
