# PKG
PKG is an useful develop package collection for golang.

## Requirements
Go 1.21 or higher

## Feature
- Config provide config function to initializes the config yaml.
- Logger provide logger to print/output the log. 

## Usage
```go
import (
    "github.com/yanun0323/pkg/logs"
    "github.com/yanun0323/pkg/config"
)

func main() {
    l := logs.New("example_service", 0)

    err := config.Init("config", true, "../config", "../../config")
    if err != nil {
        l.WithError(err).Fatal("init config")
    }

    l.Info("init config succeed")
}
```
