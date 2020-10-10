package ndriver

import (
	"fmt"
	"reflect"
)

func Call(Iface interface{}, MethodName string, params ...interface{}) error {

	ValueIface := reflect.ValueOf(Iface)

	//
	// Check if the passed interface is a pointer
	//
	if ValueIface.Type().Kind() != reflect.Ptr {
		//
		// Create a new type of Iface, so we have a pointer to work with
		//
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	//
	// Get the method by name
	//
	Method1 := ValueIface.MethodByName(MethodName)

	if !Method1.IsValid() {
		return fmt.Errorf("Method not found `%s` `%s`", MethodName, ValueIface.Type())
	}

	in := make([]reflect.Value, len(params))

	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	res := Method1.Call(in)
	err, _ := res[0].Interface().(error)

	return err
}
