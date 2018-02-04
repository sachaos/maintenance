Maintenance
===

`maintenance` is little Golang `net/http` middleware to add maintenance feature to your app.

Features
---

* IP address white list.

Usage
---

## Install

```
go get github.com/sachaos/maintenance
```

## Sample

```go
func main() {
	// Create maintenance instance with backend memcached url
	memcachedUrl := os.Getenv("MEMCACHED_SERVER")
	m := maintenance.NewMaintenance(memcachedUrl)

	r := chi.NewRouter()

	r.Use(m.SetMaintenance)            // Set MaintenanceMode in request context
	r.Use(m.AllowByIP)                 // Enable IP white list
	r.Use(m.ResponseIfMaintenanceMode) // If maintenance mode enabled, response specifield message with 503.

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request Succeeded")
	})
    
    http.ListenAndServe(":3000", r)
}
```
