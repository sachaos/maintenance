package maintenance

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"fmt"

	"github.com/go-chi/chi"
)

func TestMaintenance(t *testing.T) {
	m := NewMaintenanceMemcachedBackend(os.Getenv("MEMCACHED_SERVER"))
	mc := NewMemcachedClient(os.Getenv("MEMCACHED_SERVER"))

	middlewares := chi.Chain(m.SetMaintenance, m.ResponseIfMaintenanceMode)

	ts := httptest.NewTLSServer(middlewares.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request Succeeded")
	})))
	defer ts.Close()

	client := ts.Client()

	t.Run("When maintenance mode was enabled", func(t *testing.T) {
		if err := mc.SetMessage([]byte("Service Unavailable")); err != nil {
			t.Error(err.Error())
		}
		defer mc.Disable()

		res, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("Expected status code was %v, but acctually: %v", http.StatusServiceUnavailable, res.StatusCode)
		}
	})

	t.Run("When maintenance mode was disabled", func(t *testing.T) {
		res, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code was %v, but acctually: %v", http.StatusOK, res.StatusCode)
		}
	})
}

func TestAllowByIPFeature(t *testing.T) {
	m := NewMaintenanceMemcachedBackend(os.Getenv("MEMCACHED_SERVER"))

	mc := NewMemcachedClient(os.Getenv("MEMCACHED_SERVER"))
	if err := mc.SetMessage([]byte("Service Unavailable")); err != nil {
		t.Error(err.Error())
	}
	defer mc.Disable()

	middlewares := chi.Chain(m.SetMaintenance, m.AllowByIP, m.ResponseIfMaintenanceMode)

	ts := httptest.NewTLSServer(middlewares.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request Succeeded")
	})))
	defer ts.Close()

	client := ts.Client()

	t.Run("When the white list was not exist", func(t *testing.T) {
		res, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("Expected status code was %v, but acctually: %v", http.StatusServiceUnavailable, res.StatusCode)
		}
	})

	t.Run("When requested by IP specified in white list", func(t *testing.T) {
		allowedIPs := []string{
			"127.0.0.1",
		}
		if err := mc.SetAllowedIPs(allowedIPs); err != nil {
			t.Error(err.Error())
		}
		defer mc.DisableAllowedIPs()

		res, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code was %v, but acctually: %v", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("When requested by IP not specified in white list", func(t *testing.T) {
		allowedIPs := []string{
			"127.0.0.2",
		}
		if err := mc.SetAllowedIPs(allowedIPs); err != nil {
			t.Error(err.Error())
		}
		defer mc.DisableAllowedIPs()

		res, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("Expected status code was %v, but acctually: %v", http.StatusServiceUnavailable, res.StatusCode)
		}
	})
}
