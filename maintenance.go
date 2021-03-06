package maintenance

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
)

type Maintenance interface {
	SetMaintenance(next http.Handler) http.Handler
	AllowByIP(next http.Handler) http.Handler
	ResponseIfMaintenanceMode(next http.Handler) http.Handler
}

type maintenance struct {
	client Client
}

func NewMaintenanceRedisBackend(opt *redis.Options) Maintenance {
	client := NewRedisClient(opt)
	return &maintenance{
		client: client,
	}
}
func NewMaintenanceMemcachedBackend(url string) Maintenance {
	client := NewMemcachedClient(url)
	return &maintenance{
		client: client,
	}
}

func (ms *maintenance) SetMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), MaintenanceKey, ms.client.GetMaintenance())
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

		ips, err := ms.client.GetAllowedIPs()
		if err != nil {
			mode.Disable()
			next.ServeHTTP(w, r)
			return
		}

		ip := strings.Split(r.RemoteAddr, ":")[0]

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
