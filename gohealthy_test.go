package gohealthy_test

import (
	"context"
	"testing"
	"time"

	"github.com/adnanbrq/gohealthy"
)

func TestNewGoHealthy(t *testing.T) {
	g1 := gohealthy.New([]gohealthy.HealthCheck{
		gohealthy.NewTimeoutHealthCheck("timeout", time.Second),
		{
			Name:  "yo",
			Check: func(ctx context.Context) (bool, string) { return true, "Yo" },
		},
		{
			Name:  "whops",
			Check: func(ctx context.Context) (bool, string) { return false, "whops" },
		},
	})

	g2 := gohealthy.New([]gohealthy.HealthCheck{
		// Negative duration will always return a Unhealthy response
		gohealthy.NewTimeoutHealthCheck("timeout", time.Microsecond*-1),
	})

	g3 := gohealthy.New([]gohealthy.HealthCheck{})

	health := g1.GetHealth(context.Background())

	if health.IsHealthy {
		t.Fatal("IsHealthy should be false")
	}

	if health.Origin != "whops" {
		t.Fatal("GetHealth Origin should have been \"whops\"")
	}

	if health.UnhealthyReason != "whops" {
		t.Fatal("UnhealthyReason should equal \"whops\"")
	}

	health = g2.GetHealth(context.Background())
	if health.IsHealthy {
		t.Fatal("IsHealthy should be false")
	}

	if health.Origin != "timeout" {
		t.Fatal("GetHealth Origin should have been \"timeout\"")
	}

	if health.UnhealthyReason != "We exceeded the maximum time to run" {
		t.Fatal("UnhealthyReason should equal \"We exceeded the maximum time to run\"")
	}

	health = g3.GetHealth(context.Background())
	if !health.IsHealthy {
		t.Fatal("IsHealthy should be true")
	}
}
