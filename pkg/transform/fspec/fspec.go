package fspec

import (
	"reflect"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

// ComputeAvailableFields returns the names of non-nil pointer fields on the Record.
func ComputeAvailableFields(msg *asterix.AsterixMessage) []string {
	if msg.Record == nil {
		return nil
	}

	v := reflect.ValueOf(msg.Record)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	var fields []string
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr && !f.IsNil() {
			tag := t.Field(i).Tag.Get("json")
			if tag != "" {
				// strip ",omitempty"
				for j := 0; j < len(tag); j++ {
					if tag[j] == ',' {
						tag = tag[:j]
						break
					}
				}
				fields = append(fields, tag)
			}
		}
	}
	return fields
}