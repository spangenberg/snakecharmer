package snakecharmer

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const tagName = "flag"

var (
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
)

type Flag struct {
	Name  string
	Value string
	Usage string
}

func GenerateFlags(cmd *cobra.Command, cfg Config) *cobra.Command {
	parseFlags(cmd.Flags(), reflect.TypeOf(cfg).Elem(), "")
	return cmd
}

func parseFlags(flags *pflag.FlagSet, r reflect.Type, path string) {
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		tag := f.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		s := strings.Split(tag, ";")
		var name, shorthand, value, usage string
		switch len(s) {
		case 4:
			name, shorthand, value, usage = s[0], s[1], s[2], s[3]
		case 3:
			name, shorthand, value = s[0], s[1], s[2]
		case 2:
			name, shorthand = s[0], s[1]
		case 1:
			name = s[0]
		}
		kind := f.Type.Kind()
		if kind == reflect.Struct {
			parseFlags(flags, f.Type, buildKey(path, name))
		} else {
			key := buildKey(path, name)
			name = strings.Replace(key, ".", "-", -1)
			name = kebabCase(name)
			addFlag(flags, kind, value, name, shorthand, usage)
			if err := viper.BindPFlag(key, flags.Lookup(name)); err != nil {
				panic(err)
			}
		}
	}
}

func addFlag(flags *pflag.FlagSet, kind reflect.Kind, value, name, shorthand, usage string) {
	switch kind {
	case reflect.Bool:
		var boolValue bool
		if value != "" {
			var err error
			if boolValue, err = strconv.ParseBool(value); err != nil {
				panic(err)
			}
		}
		flags.BoolP(name, shorthand, boolValue, usage)
	case reflect.Int:
		var intValue int64
		if value != "" {
			var err error
			intValue, err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.IntP(name, shorthand, int(intValue), usage)
	case reflect.Int8:
		var intValue int64
		if value != "" {
			var err error
			intValue, err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Int8P(name, shorthand, int8(intValue), usage)
	case reflect.Int16:
		var intValue int64
		if value != "" {
			var err error
			intValue, err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Int16P(name, shorthand, int16(intValue), usage)
	case reflect.Int32:
		var intValue int64
		if value != "" {
			var err error
			intValue, err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Int32P(name, shorthand, int32(intValue), usage)
	case reflect.Int64:
		var intValue int64
		if value != "" {
			var err error
			intValue, err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Int64P(name, shorthand, intValue, usage)
	case reflect.Uint:
		var uintValue uint64
		if value != "" {
			var err error
			uintValue, err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.UintP(name, shorthand, uint(uintValue), usage)
	case reflect.Uint8:
		var uintValue uint64
		if value != "" {
			var err error
			uintValue, err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Uint8P(name, shorthand, uint8(uintValue), usage)
	case reflect.Uint16:
		var uintValue uint64
		if value != "" {
			var err error
			uintValue, err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Uint16P(name, shorthand, uint16(uintValue), usage)
	case reflect.Uint32:
		var uintValue uint64
		if value != "" {
			var err error
			uintValue, err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Uint32P(name, shorthand, uint32(uintValue), usage)
	case reflect.Uint64:
		var uintValue uint64
		if value != "" {
			var err error
			uintValue, err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Uint64P(name, shorthand, uintValue, usage)
	case reflect.Float32:
		var floatValue float64
		if value != "" {
			var err error
			floatValue, err = strconv.ParseFloat(value, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Float32P(name, shorthand, float32(floatValue), usage)
	case reflect.Float64:
		var floatValue float64
		if value != "" {
			var err error
			floatValue, err = strconv.ParseFloat(value, 0)
			if err != nil {
				panic(err)
			}
		}
		flags.Float64P(name, shorthand, floatValue, usage)
	// case reflect.Array:
	// flags.StringArray()
	// NOTE: Not just strings possible
	case reflect.Slice:
		sliceValue := strings.Split(value, ",")
		flags.StringSliceP(name, shorthand, sliceValue, usage)
	case reflect.String:
		flags.StringP(name, shorthand, value, usage)
	// time.Duration
	default:
		panic(fmt.Sprintf("flag type %v not supported", kind))
	}
}

func buildKey(a, b string) string {
	if a == "" {
		return b
	}
	return fmt.Sprintf("%s.%s", a, b)
}

func kebabCase(str string) string {
	str = matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	str = matchAllCap.ReplaceAllString(str, "${1}-${2}")
	return strings.ToLower(str)
}

func tagInformation(r reflect.Type, path []string, previous string) string {
	current := path[0]
	structField, ok := r.FieldByName(current)
	if !ok {
		panic(fmt.Sprintf("field %s not found", structField.Name))
	}
	tag := structField.Tag.Get(tagName)
	s := strings.Split(tag, ";")
	usage := kebabCase(s[0])
	if previous != "" {
		usage = previous+"-"+usage
	}
	if len(path) > 1 {
		return tagInformation(structField.Type, path[1:], usage)
	}
	return usage
}
