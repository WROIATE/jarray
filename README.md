# Self Use Array Utils

## install

```shell
go get github.com/WROIATE/jarray
```

## Code Example
### list
```go
package main

import "github.com/WROIATE/jarray"

func main() {
	l := jarray.NewSimpleList[int]()
	l.Add(1)
	l.Append([]int{2, 3, 4})
	l.Contains(2, func(a, b int) bool {
		return a == b
	})
	l.Sort(func(i, j int, sortSlice []int) bool {
		return sortSlice[i] < sortSlice[j]
	})
}
```
### stream
```go
package main

import (
	"fmt"
	"github.com/WROIATE/jarray"
)

func main() {
	jarray.NewMapStream[int, string]([]int{1, 2, 3, 4, 1}).
		Map(func(i int) string {
			return fmt.Sprintf("%d", i)
		}).
		Filter(func(a int) bool {
			return a == 1 || a == 2 || a == 3
		}).
		Distinct(func(i int) int64 {
			return int64(i)
		}).
		Range(func(k int, v int) {
			fmt.Println(v)
		}).
		Limit(2)
}

```