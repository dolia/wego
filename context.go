package wego

import (
	"net/http"
	"reflect"
)

type MiddleWare interface {

	Handle(Context *Context)
	
}

type Context struct {

	*Render `service:""`

	Req *http.Request

	Writer http.ResponseWriter

	MiddleWares []MiddleWare

	Controller   Handle
	Action 		string
	Params Params

	Keys map[string]interface{}

	index    int
}


func (c *Context) middleware() MiddleWare {
	if c.index < len(c.MiddleWares) {
		return c.MiddleWares[c.index]
	}
	panic("invalid index for Context middleware")
}

func (c *Context) Next() {
	c.index += 1
	c.handle()
}


func (c *Context) handle() {

	for c.index <= len(c.MiddleWares) {

		if c.index == len(c.MiddleWares) {
			c.callController(c.Controller,c.Action)
		}else{
			c.callMiddleware(c.middleware())
		}
		c.index += 1

	}
}


func (c *Context) callController(controller Handle,action string)  {

	cv := reflect.ValueOf(controller)

	fv := cv.MethodByName(action)

	t := fv.Type()

	var in = make([]reflect.Value, t.NumIn()) //Panic if t is not kind of Func
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i)

		val := c.findValueFromGraph(argType)

		if !val.IsValid() {

			//panic("Value not found for type " + argType)
		}

		in[i] = val
	}

	fv.Call(in)
}

func (c *Context) callMiddleware(middleware MiddleWare)  {

	middleware.Handle(c)

}

func (c *Context) findValueFromGraph(p reflect.Type) reflect.Value  {

	if (p ==reflect.TypeOf(&Context{})){
		return reflect.ValueOf(c)
	}

	val := graph.GetService(p)

	return val

}