package gohealthy

import (
	"context"
	"time"
)

type contextKey struct{}

type healthyContext struct {
	start time.Time
}

// HealthCheck describes a test
type HealthCheck struct {
	// Name of the health check
	Name string
	// Check is a function that is performed when requesting the health status.
	Check func(ctx context.Context) (bool, string)
}

type health struct {
	// IsHealthy describes the state of health as a bool.
	IsHealthy bool
	// Origin is the name of the health check and is only set/not empty if IsHealth is false.
	Origin string
	// UnhealthyReason contains the reason returned by the original health check.
	UnhealthyReason string
}

type goHealthy struct {
	healthChecks []HealthCheck
}

// NewTimeoutHealthCheck creates a HealthCheck that ensures that the previous tests did not last longer than the specified time.
func NewTimeoutHealthCheck(name string, duration time.Duration) HealthCheck {
	return HealthCheck{
		Name: name,
		Check: func(ctx context.Context) (bool, string) {
			if healthyContext, ok := ctx.Value(contextKey{}).(healthyContext); ok {
				if time.Since(healthyContext.start) >= duration {
					return false, "We exceeded the maximum time to run"
				}
			}
			return true, ""
		},
	}
}

// GetHealth performs each given health check and returns a health response.
func (gh goHealthy) GetHealth(ctx context.Context) health {
	healthyContext := context.WithValue(ctx, contextKey{}, healthyContext{time.Now()})

	for _, healthCheck := range gh.healthChecks {
		if healthy, reason := healthCheck.Check(healthyContext); !healthy {
			return health{IsHealthy: false, UnhealthyReason: reason, Origin: healthCheck.Name}
		}
	}

	return health{IsHealthy: true}
}

// New creates a new goHealthy instance.
func New(checks []HealthCheck) goHealthy {
	return goHealthy{checks}
}
