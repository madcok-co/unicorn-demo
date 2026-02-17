# 🚀 Quick Start Guide

## Current Status

⚠️ **This demo is under development**

The Unicorn Demo API is being built to showcase all enterprise features of the Unicorn framework. Some handlers need to be updated to match the latest Unicorn interfaces.

## What Works

✅ Project structure  
✅ Domain models and DTOs  
✅ Database integration (GORM)  
✅ Docker setup  
✅ Documentation  

## What Needs Fixing

🔧 Handler implementations need to be updated for:
- Proper context.Get() usage (returns `(any, bool)`)
- Database interface methods
- Identity/Claims handling

## Simple Working Example

Until the full demo is complete, here's a minimal working example:

```bash
# Create a new project
mkdir my-unicorn-api && cd my-unicorn-api
go mod init github.com/yourusername/my-unicorn-api

# Install Unicorn
go get github.com/madcok-co/unicorn/core@main

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "log"
    httpAdapter "github.com/madcok-co/unicorn/core/pkg/adapters/http"
    "github.com/madcok-co/unicorn/core/pkg/app"
    "github.com/madcok-co/unicorn/core/pkg/context"
)

type HealthResponse struct {
    Status string `json:"status"`
}

func HealthCheck(ctx *context.Context, req struct{}) (*HealthResponse, error) {
    return &HealthResponse{Status: "healthy"}, nil
}

func main() {
    application := app.New(&app.Config{
        Name:       "My API",
        EnableHTTP: true,
        HTTP:       &httpAdapter.Config{Port: 8080},
    })

    application.RegisterHandler(HealthCheck).
        Named("health").
        HTTP("GET", "/health").
        Done()

    log.Fatal(application.Start())
}
EOF

# Run
go run main.go
```

Visit: http://localhost:8080/health

## Contributing

Want to help fix the demo? Check the compilation errors and submit a PR!

## Full Framework Documentation

See the main [Unicorn Framework](https://github.com/madcok-co/unicorn) for:
- Complete examples
- API reference
- Enterprise features documentation
