package sdm630

import (
	"reflect"
	"regexp"
)

type pathSlice []string

// StructureAccessor manages access to arbitrary nested structures
// by access pattern
type StructureAccessor struct {
	pathes map[string]pathSlice
	regex  *regexp.Regexp
}

func NewStructureAccessor(pattern string) *StructureAccessor {
	return &StructureAccessor{
		pathes: make(map[string]pathSlice),
		regex:  regexp.MustCompile(pattern),
	}
}

func (sa *StructureAccessor) getPath(key string) (pathSlice, bool) {
	path, ok := sa.pathes[key]
	return path, ok
}

func (sa *StructureAccessor) setPath(key string, path pathSlice) {
	sa.pathes[key] = path
}

func (sa *StructureAccessor) value(obj interface{}, key string) (reflect.Value, bool) {
	// first, get the cached access path for the key
	path, ok := sa.getPath(key)
	// log.Printf("sa.value %v %b", path, ok)

	// path not yet determined
	if !ok {
		// log.Printf("sa.value did not get path from cache - determining")

		path = sa.regex.FindStringSubmatch(key)
		// log.Printf("parts %s %v\n", key, path)

		if len(path) <= 0 {
			sa.setPath(key, nil) // remeber this struct doesn't exist
			return reflect.Value{}, false
		}

		path = path[1:]       // remove entire match
		sa.setPath(key, path) // remember path
	}

	// key found buth path nil - struct access path does not exist
	if path == nil {
		// log.Printf("sa.value got nil path from cache")
		return reflect.Value{}, false
	}

	// start inspecting object under reflection
	v := reflect.ValueOf(obj)

	// descent access path
	for _, field := range path {
		// log.Printf("Kind %s", v.Kind())
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			// log.Printf("Kind %s", v.Kind())
		}

		// log.Printf("FieldByName %s %v", field, v)
		v = v.FieldByName(field)

		if v.Kind() == reflect.Invalid {
			// could not find structure
			// log.Printf("sa.value struct does not exist")
			sa.setPath(key, nil) // remeber this struct doesn't exist
			return reflect.Value{}, false
		}
	}

	return v, true
}

func (sa *StructureAccessor) SetFloat(obj interface{}, key string, val float64) bool {
	if v, ok := sa.value(obj, key); ok {
		// log.Printf("Kind %s", v.Kind())
		if v.Kind() == reflect.Ptr {
			// create pointer
			v.Set(reflect.New(v.Type().Elem()))
			v = v.Elem()
		}

		v.SetFloat(val)
		return true
	}

	return false
}
