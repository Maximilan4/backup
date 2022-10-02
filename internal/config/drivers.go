package config

import (
	"backup/internal/drivers"
	"fmt"
	"reflect"
	"strings"
)

type (
	Drivers struct {
		S3 map[string]*drivers.S3DriverConfig `mapstructure:"s3"`
		Fs map[string]*drivers.FsDriverConfig `mapstructure:"fs"`
	}
)

func (d Drivers) Get(info drivers.DriverInfo) (any, error) {
	rVal := reflect.ValueOf(d)
	fieldValue := rVal.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(s, info.Type())
	})

	if !fieldValue.IsValid() || fieldValue.IsNil() || fieldValue.IsZero() {
		return nil, fmt.Errorf("unable to find %s driver named configuration: %s", info.Type(), info.Name())
	}

	value := fieldValue.MapIndex(reflect.ValueOf(info.Name()))
	if value.IsNil() {
		return nil, fmt.Errorf("unable to find %s driver named configuration: %s", info.Type(), info.Name())
	}

	return value.Interface(), nil
}
