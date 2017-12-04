
package wego

import (
	"net/http"
	"reflect"
)


type Handle interface {}


type InvokeHandler func(w http.ResponseWriter,req *http.Request,handle Handle,ps Params,f string ,tsr bool)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}


type group struct {

	pattern  string
	middleWares []MiddleWare

}

type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {


	trees map[string]*node

	groups     []group

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}


func (r *Router) GROUP(pattern string, fn func(*Router),middleWares ...MiddleWare) {
	r.groups = append(r.groups, group{pattern,middleWares })
	fn(r)
	r.groups = r.groups[:len(r.groups)-1]
}

// GET is a shortcut for router.Handle("GET", path, handle,action)
func (r *Router) GET(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("GET", path, handle,action,middleWares)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle,action)
func (r *Router) HEAD(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("HEAD", path, handle,action,middleWares)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle,action)
func (r *Router) OPTIONS(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("OPTIONS", path, handle,action,middleWares)
}

// POST is a shortcut for router.Handle("POST", path, handle,action)
func (r *Router) POST(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("POST", path, handle,action,middleWares)
}

// PUT is a shortcut for router.Handle("PUT", path, handle,action)
func (r *Router) PUT(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("PUT", path, handle,action,middleWares)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle,action)
func (r *Router) PATCH(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("PATCH", path, handle,action,middleWares)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle,action)
func (r *Router) DELETE(path string, handle Handle, action string,middleWares ...MiddleWare) {
	r.Handle("DELETE", path, handle,action,middleWares)
}


func (r *Router) Handle(method,path string, controller Handle,action string,middleWares []MiddleWare) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}
	//check action exists
	a := reflect.ValueOf(controller).MethodByName(action)
	if !a.IsValid(){
		panic(" controller "+reflect.ValueOf(controller).Type().String()+" has no  action " +action)
	}

	if len(r.groups) > 0 {

		groupPattern := ""

		groupMiddleWares := make([]MiddleWare, 0)

		for _, g := range r.groups {

			groupPattern += g.pattern

			groupMiddleWares = append(groupMiddleWares, g.middleWares...)

		}

		path = groupPattern + path

		middleWares  = append(middleWares,groupMiddleWares...)

	}

	// middleWares dependency...

	for _,m := range middleWares {
		graph.Provide(&Object{Value:m})
	}

	graph.Provide(&Object{Value:controller})

	Log.Info("add route:",path, controller,action,middleWares)

	root.addRoute(path, controller,action,middleWares)
}


// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).


//// Handler is an adapter which allows the usage of an http.Handler as a
//// request handle.
//func (r *Router) Handler(method, path string, handler http.Handler) {
//	r.Handle(method, path,
//		func(w http.ResponseWriter, req *http.Request, _ Params) {
//			handler.ServeHTTP(w, req)
//		},
//	)
//}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
//func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
//	r.Handler(method, path, handler)
//}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
//func (r *Router) ServeFiles(path string, root http.FileSystem) {
//	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
//		panic("path must end with /*filepath in path '" + path + "'")
//	}
//
//	fileServer := http.FileServer(root)
//
//	r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) {
//		req.URL.Path = ps.ByName("filepath")
//		fileServer.ServeHTTP(w, req)
//	})
//}
//
func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (handle Handle, p Params, f string,middleWares []MiddleWare ,tsr bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	middleWares = make([]MiddleWare,0)
	return nil, nil,"",middleWares,false
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _,_, _,_ := r.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request,c *Context) {

	//Log.Info(  req.Method + " " + req.RequestURI +" " + req.Host + " "+req.Header.Get("User-Agent"))

	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	path := req.URL.Path

	if root := r.trees[req.Method]; root != nil {
		if handle, ps,f,middleWares,tsr := root.getValue(path); handle != nil {

			//Log.Info("get middleware :",path,handle,f,middleWares)


			r.handler(c,w,req,handle,ps,middleWares,f,tsr)
			return

		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	if req.Method == "OPTIONS" {
		// Handle OPTIONS requests
		if r.HandleOPTIONS {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				return
			}
		}
	} else {
		// Handle 405
		if r.HandleMethodNotAllowed {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {

		r.NotFound.ServeHTTP(w, req)
	} else {

		Log.Info(  req.Method + " " + req.RequestURI +" " + req.Host + " "+req.Header.Get("User-Agent"))

		http.NotFound(w, req)
	}
}


func (r *Router)handler(c *Context,w http.ResponseWriter,req *http.Request,controller Handle,ps Params,middleWares []MiddleWare,f string ,tsr bool) {
	c.Writer = nil
	c.Writer = w
	c.Req = nil
	c.Req = req
	c.Controller = controller
	c.Action = f
	c.Params = ps
	c.MiddleWares = middleWares
	c.index = 0
	c.handle()
}




