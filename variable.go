package wego

type Variable struct {
	Keys     map[string]interface{}
}

// Set is used to store a new key/value pair exclusivelly for this context.
// It also lazy initializes  v.Keys if it was not used previously.
func (v *Variable) Set(key string, value interface{}) {
	if v.Keys == nil {
		v.Keys = make(map[string]interface{})
	}
	v.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (v *Variable) Get(key string) (value interface{}, exists bool) {
	if v.Keys != nil {
		value, exists = v.Keys[key]
	}
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (v *Variable) MustGet(key string) interface{} {
	if value, exists := v.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}


// GetString retrieves the string-typed configuration value corresponding to the specified key.
// Please refer to Get for the detailed usage explanation.
func (v *Variable) GetString(key string, defaultValue ...string) string {
	var d string
	if len(defaultValue) > 0 {
		d = defaultValue[0]
	}
	if value,exists:= v.Get(key);exists==true{
		return value.(string)
	}else{
		return d
	}
}

// GetString retrieves the int-typed configuration value corresponding to the specified key.
// Please refer to Get for the detailed usage explanation.
func (v *Variable) GetInt(key string, defaultValue ...int) int {
	var d int
	if len(defaultValue) > 0 {
		d = defaultValue[0]
	}
	if value,exists:= v.Get(key);exists==true{
		return value.(int)
	}else{
		return d
	}

}

// GetString retrieves the int64-typed configuration value corresponding to the specified key.
// Please refer to Get for the detailed usage explanation.
func (v *Variable) GetInt64(key string, defaultValue ...int64) int64 {
	var d int64
	if len(defaultValue) > 0 {
		d = defaultValue[0]
	}
	if value,exists:= v.Get(key);exists==true{
		return value.(int64)
	}else{
		return d
	}
}

// GetString retrieves the float64-typed configuration value corresponding to the specified key.
// Please refer to Get for the detailed usage explanation.
func (v *Variable) GetFloat(key string, defaultValue ...float64) float64 {
	var d float64
	if len(defaultValue) > 0 {
		d = defaultValue[0]
	}
	if value,exists:= v.Get(key);exists==true{
		return value.(float64)
	}else{
		return d
	}
}

// GetString retrieves the bool-typed configuration value corresponding to the specified key.
// Please refer to Get for the detailed usage explanation.
func (v *Variable) GetBool(key string, defaultValue ...bool) bool {
	d := false
	if len(defaultValue) > 0 {
		d = defaultValue[0]
	}
	if value,exists:= v.Get(key);exists==true{
		return value.(bool)
	}else{
		return d
	}
}