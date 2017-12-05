# Wego Web FrameWork

Wego is a web framework written in Golang.


##Quick start

```sh
# assume the following codes in example.go file
$ cat example.go
```

```go

package main

import "github.com/dolia/wego"

type HelloController struct {

}

func (this *HelloController)Default(c *wego.Context){

	c.Text(200,"hello")

}

func main() {

	w := wego.New()

	w.GET("/hello", &HelloController,"Default")

	w.Run() // listen and serve on 0.0.0.0:3000
}
```

```
# run example.go and visit 0.0.0.0:3000/hello on browser
$ go run example.go

```

### Dependency Injection


#### How to inject your dependency

```

func main() {

	w := wego.New()


    db, _ := gorm.Open("sqlite3", "test.db") // 1. instance your dependency

    w.Provides(db)                          //  2. inject into wego with 'Provides' method



	w.GET("/hello", &HelloController,"Default")

	w.Run() // listen and serve on 0.0.0.0:3000
}

```

#### How to use you dependency

##### Use it in struct.

1. Use 'service' tag to let  wego know what you want to use.

```
type HelloController struct {

    DB *gorm.DB  `service::`  // 'service' tag tells wego which services you want to use and auto inject  into your struct.

}

```
2. Use DB service

```
func (this *HelloController)Default(){

    fmt.Println(this.DB.find(&Model{}).Value)   // Use DB service


}

```
##### Use it in struct function.

You can also use it like this

```
func (this *HelloController)Default(DB *gorm.DB){

    fmt.Println(this.DB.find(&Model{}).Value)   // Use DB service

}

```
