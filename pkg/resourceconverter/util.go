package resourceconverter

import (
	"context"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"reflect"
)

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

func SendAllConverted[T any](ctx context.Context, output chan<- model.Resource, converter ResourceConverter, resources []T) error {
	var converted []model.Resource

	for _, raw := range resources {
		resource, err := converter.ToResource(ctx, raw, nil)
		if err != nil {
			return err
		}

		converted = append(converted, resource)
	}

	return util.SendAllFromSlice(ctx, output, converted)
}

type tagFunc[T any] func(context.Context, T) (model.Tags, error)

func SendAllConvertedTags[T any](ctx context.Context, output chan<- model.Resource, converter ResourceConverter, resources []T, tF tagFunc[T]) error {
	var converted []model.Resource

	for _, raw := range resources {
		tags, err := tF(ctx, raw)
		if err != nil {
			return err
		}
		if tags == nil {
			tags = []model.Tag{}
		}
		resource, err := converter.ToResource(ctx, raw, tags)
		if err != nil {
			return err
		}

		converted = append(converted, resource)
	}

	return util.SendAllFromSlice(ctx, output, converted)
}