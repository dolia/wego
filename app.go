package wego

import (
	"net/http"
	"github.com/unrolled/render"
)

func Provide(values ...interface{}) error {

	for _, v := range values {
		if err := graph.Provide(&Object{Value: v}); err != nil {
			Log.Info("App Provide err",err)
			return err
		}
	}

	return graph.Populate()
}

func New() *App  {

	return &App{}
}

var graph Graph

type Service interface {}

type App struct {

	Config

	Router

	context *Context

	routerPrefix string

}

func (d *App)bootstrap()  {

	c := &Context{}

	render := render.New()

	d.Provides(d,render,&Render{},c)

	if err := graph.Populate(); err != nil {

		Log.Info("Service Populate err",err)
	}

	Log.Info("Service Loaded in container :",graph.Objects())

	d.context = c

}

func (d *App) Middleware(middleware MiddleWare) {
	d.Provides(middleware)
	//d.middleWares = append(d.middleWares, middleware)
}

func (d *App) MiddleWares(middleWares ...MiddleWare) {

	for _, v := range middleWares {
		d.Middleware(v)
	}
}

func (d *App)Provides(values ...interface{})  {
	for _, v := range values {
		if err := graph.Provide(&Object{Value: v}); err != nil {
			Log.Info("App Provide err",err)
		}
	}
}


func (d *App)Run()  {

	d.bootstrap()

	host := d.GetString("HTTP_SERVER_HOST","127.0.0.1")
	port := d.GetString("HTTP_SERVER_PORT","3000")
	addr := host+":"+port
	Log.Info("Http server is listening on  "+addr)
	http.ListenAndServe(addr,d)
}



func (d *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.Router.ServeHTTP(w,req,d.context)
}


