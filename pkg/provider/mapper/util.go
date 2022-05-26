package mapper

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
)

func findImplMethod(v reflect.Value, impl string) *reflect.Value {
	method := v.MethodByName(impl)
	if reflect.ValueOf(method).IsZero() {
		panic(fmt.Errorf("could not find a method called '%v' on '%T'", impl, v.Interface()))
	}

	t := method.Type()

	if isFetchMethodAsync(t) {
		return &method
	}

	panic(fmt.Errorf("method %v has invalid signature; expecting func(context.Context, chan<- T) error", impl))
}

func findTagMethod(v reflect.Value, impl string) *reflect.Value {
	method := v.MethodByName(impl)
	if reflect.ValueOf(method).IsZero() {
		panic(fmt.Errorf("could not find a method called '%v' on '%T'", impl, v.Interface()))
	}

	t := method.Type()
	if isFetchTagSync(t) {
		return &method
	}

	panic(fmt.Errorf("method %v has invalid signature; expecting func(context.Context, T) (model.Tags, error)", impl))
}

func isFetchMethodAsync(t reflect.Type) bool {
	if t.Kind() != reflect.Func {
		return false
	}

	if t.NumIn() != 2 {
		return false
	}

	var ctx context.Context
	ctxType := reflect.TypeOf(&ctx).Elem()
	if !ctxType.AssignableTo(t.In(0)) {
		return false
	}

	in1 := t.In(1)
	if in1.Kind() != reflect.Chan {
		return false
	}

	if in1.ChanDir() != reflect.SendDir {
		return false
	}

	if t.NumOut() != 1 {
		return false
	}
	// TODO: Make sure this is the builtin error type
	if t.Out(0).Name() != "error" {
		return false
	}

	return true
}

func isFetchTagSync(t reflect.Type) bool {
	if t.Kind() != reflect.Func {
		return false
	}

	if t.NumIn() != 2 {
		return false
	}

	var ctx context.Context
	ctxType := reflect.TypeOf(&ctx).Elem()
	if !ctxType.AssignableTo(t.In(0)) {
		return false
	}

	if t.NumOut() != 2 {
		return false
	}

	if t.Out(0) != reflect.TypeOf(make(model.Tags, 0)) {
		return false
	}

	// TODO: Make sure this is the builtin error type
	if t.Out(1).Name() != "error" {
		return false
	}

	return true
}

func getProperties(name string, v reflect.Value, ignoredFields []string, maxRecursion int) []model.Property {
	if util.Contains(ignoredFields, name) || maxRecursion <= 0 {
		//ignore this field
		return nil
	}

	emptyProp := []model.Property{{Name: name, Value: ""}}

	switch v.Kind() {
	case reflect.Invalid:
		//ignore this field
		return nil
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			//empty pointer
			return emptyProp
		}
		//display pointer value
		return getProperties(name, v.Elem(), ignoredFields, maxRecursion)
	case reflect.Slice:
		if v.IsZero() {
			//empty slice
			return emptyProp
		}
		//return a distinct property for each slice element
		//ex: Subnets=[a,b] -> Subnets[0]=a Subnets[1]=b
		var props []model.Property
		for i := 0; i < v.Len(); i++ {
			props = append(props,
				getProperties(
					fmt.Sprintf("%v[%d]", name, i), v.Index(i), ignoredFields, maxRecursion-1)...)
		}
		return props
	case reflect.Struct:
		defaultFormatVal := fmt.Sprintf("%v", v)
		if !strings.Contains(defaultFormatVal, "{") {
			// this looks like a custom format, use it
			// ex: a datetime might have a nice formatting already implemented
			return []model.Property{{Name: name, Value: defaultFormatVal}}
		}
		//return a distinct property for each struct element
		//ex: IamInstanceProfile{Arn:, Id:} -> IamInstanceProfile["Arn"]=... IamInstanceProfile["Id"]=...
		t := v.Type()
		var props []model.Property
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			props = append(props,
				getProperties(
					fmt.Sprintf("%v[%v]", name, field.Name), v.Field(i), ignoredFields, maxRecursion-1)...)
		}
		return props
	default:
		return []model.Property{{Name: name, Value: fmt.Sprintf("%v", v)}}
	}
}

func getTags(v reflect.Value, tagField TagField) []model.Tag {
	switch v.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			//empty pointer
			return nil
		}
		//display pointer value
		return getTags(v.Elem(), tagField)
	case reflect.Slice:
		if v.IsZero() {
			//empty slice
			return nil
		}
		//return a distinct Tag for each slice element
		//ex: Tags=[a,b] -> Tag=a Tag=b
		var tags []model.Tag
		for i := 0; i < v.Len(); i++ {
			tags = append(tags,
				getTags(v.Index(i), tagField)...)
		}
		return tags
	case reflect.Struct:
		//expects a Key and Value
		key := getPtrVal(v.FieldByName(tagField.Key))
		value := getPtrVal(v.FieldByName(tagField.Value))
		keyStr := fmt.Sprintf("%v", key)
		valStr := fmt.Sprintf("%v", value)
		// we have a tag
		return []model.Tag{{Key: keyStr, Value: valStr}}
	default:
		return nil
	}

}

func getPtrVal(v reflect.Value) reflect.Value {
	if v.IsValid() && !v.IsZero() && v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}
