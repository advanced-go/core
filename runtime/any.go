package runtime

import "reflect"

type Nillable *struct{}

// IsNil - determine if the interface{} is nil, or if it holds a nil pointer
func IsNil(a any) bool {
	if a == nil {
		return true
	}
	if reflect.TypeOf(a).Kind() != reflect.Pointer {
		return false
	}
	return reflect.ValueOf(a).IsNil()
}

func TypeName(a any) string {
	if a == nil {
		return "<nil>"
	}
	// TO DO: determine underlying type name of a pointer
	if IsPointer(a) {
		k := reflect.TypeOf(a).Kind()
		return k.String()
	}
	return reflect.TypeOf(a).Name()
}

func IsPointer(a any) bool {
	if a == nil {
		return false
	}
	if reflect.TypeOf(a).Kind() != reflect.Pointer {
		return false
	}
	return true
}

/*
func IsNillable(a any) bool {
	return IsPointer(a) || IsPointerType(a)
}



func IsPointerType(a any) bool {
	if a == nil {
		return false
	}
	switch reflect.ValueOf(a).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map:
		return true
	}
	return false
}

*/
