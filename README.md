# GenerateGoStruct

```go
import (
	"github.com/Seann-Moser/GenerateGoStruct"
	"log"
)

func main() {
  structs, err := GenerateGoStruct.ConvertJson(``)
  if err != nil{
	  log.Fatal(err)
  }
  for _,s := range{
    fmt.Println(s)	  
  }
}
```