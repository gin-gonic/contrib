# pprof

pprof support for gin

## Usage

```go
import "github.com/gin-gonic/gin"
import "github.com/gin-gonic/contrib/pprof"

func main() {
    router := gin.Default()
		pprof.Register(router, nil)
    router.Run(":8080")
}
```

```bash
go tool pprof http://localhost:8080/debug/pprof/profile
```
