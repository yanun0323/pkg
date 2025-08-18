# PKG

PKG is an useful develop package collection for golang.

## Requirements

#### _Go 1.21 or higher_

## Feature

- `Config` provides init function to initialize the config yaml.

## Usage

```go
import (
    "github.com/yanun0323/pkg/config"
)

func main() {
    err := config.Init("config", true, "../config", "../../config")
    if err != nil {
        panic("init config")
    }
}
```
