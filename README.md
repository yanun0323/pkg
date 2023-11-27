# PKG
PKG is an useful develop package collection for golang.

## Requirements
#### _Go 1.21 or higher_

## Feature
- `Config` provides init function to initialize the config yaml.
- `Logger` provides logger to print/output the log. 

## Usage
```go
import (
    "github.com/yanun0323/pkg/logs"
    "github.com/yanun0323/pkg/config"
)

func main() {
    l := logs.New(logs.LevelInfo, logs.OutputStd())

    err := config.Init("config", true, "../config", "../../config")
    if err != nil {
        l.WithError(err).Fatal("init config")
    }

    l.Info("init config succeed")
}
```
