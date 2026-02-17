package middleware

import (
	"net/http"
	"strings"

	"github.com/madcok-co/unicorn/contrib/multitenancy"
)

// TenantResolver middleware resolves tenant from request
func TenantResolver(mt *multitenancy.Driver) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tenant, err := mt.GetTenantFromRequest(r.Context(), r)
			if err != nil {
				// Try to get from header as fallback
				tenantID := r.Header.Get("X-Tenant-ID")
				if tenantID == "" {
					// Extract from host
					host := r.Host
					if idx := strings.Index(host, "."); idx != -1 {
						tenantID = host[:idx]
					} else {
						tenantID = "default"
					}
				}

				// Store in request context (will be accessible in handler via ctx.Get("tenant_id"))
				// For now, we'll pass it through header
				r.Header.Set("X-Tenant-ID", tenantID)
			} else {
				r.Header.Set("X-Tenant-ID", tenant.ID)
			}

			next(w, r)
		}
	}
}
