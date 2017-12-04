#Wego Web FrameWork

Wego is a web framework written in Golang.


```sh
# assume the following codes in example.go file
$ cat example.go
```


```go
package main

import "github.com/dolia/wego"

type HelloController struct {

	Db *gorm.DB  `service:""`

}

func (this *Category)Default(c *wego.Context){

	c.Text(200,"hello")

}



func main() {

	w := wego.New()

	w.GET("/hello", &HelloController,"Default")

	r.Run() // listen and serve on 0.0.0.0:3000
}
```

```
# run example.go and visit 0.0.0.0:8080/hello on browser
$ go run example.go

```