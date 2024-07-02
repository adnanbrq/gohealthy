# GoHealthy

GoHealthy is a handy way to run Health Checks (like Ping) against any target and check the system's health.\
Tests run in a sequence with a fail-fast mechanism. If the Health Check says the process is unhealthy, it stops right there and returns the result.\
It is designed not to handle HTTP calls or anything else. It's mainly used to create a response that describes the system's health. You can decide how to communicate this. 👍

## Contents

- [Installation](#installation)
- [Usage](#usage)
- [Dependencies](#dependencies)

## Installation

```sh
$ go get -u github.com/adnanbrq/gohealthy
```

## Usage

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/adnanbrq/gohealthy/gohealthy"
)

func main() {
	gh := gohealthy.New([]gohealthy.HealthCheck{
        // Simulated fail
		gohealthy.NewTimeoutHealthCheck("timeout", time.Microsecond*-1),
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if health := gh.GetHealth(r.Context()); !health.IsHealthy {
			http.Error(w, health.UnhealthyReason, http.StatusServiceUnavailable)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

Visit http://localhost:8081 with your Browser or use the following command to see the health status.

```console
$ curl -v localhost:8081
```

## Dependencies

- No dependencies
