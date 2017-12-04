package wego

import (
	"github.com/unrolled/render"
)


type Render struct {
	Context *Context `service:""`
	Render *render.Render `service:""`

}

func (r *Render) Text(status int,v string) error {
	return r.Render.Text(r.Context.Writer,status,v)
}

func (r *Render) XML(status int,v string) error {
	return r.Render.XML(r.Context.Writer,status,v)
}

func (r *Render) JSON(status int,v interface{}) error{

	return r.Render.JSON(r.Context.Writer,status,v)
}

func (r *Render) JSONP(status int,callback string,v interface{}) error{

	return r.Render.JSONP(r.Context.Writer,status,callback,v)
}

func (r *Render) Data(status int,v []byte) error{

	return r.Render.Data(r.Context.Writer,status,v)
}

func (r *Render) HTML(status int, name string, binding interface{}, layoutOption ...string) error{

	if len(layoutOption)>0 {

		layout := layoutOption[0]

		return r.Render.HTML(r.Context.Writer,status,name,binding,render.HTMLOptions{Layout:layout})

	}else{
		return r.Render.HTML(r.Context.Writer,status,name,binding)

	}




}




