package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/layarda-durianpay/go-skeleton/internal/config"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/proto"
)

// ToOtelAttributes converts a label and associated data into OpenTelemetry attributes.
func ToOtelAttributes(label string, data interface{}) []attribute.KeyValue {
	var res []attribute.KeyValue

	if data == nil {
		res = append(res, attribute.String(label, "null"))
		return res
	}

	switch v := data.(type) {
	case proto.Message:
		// If the data is of type proto.Message, manually convert it to a struct type.
		// The proto.Message is expected to have a .String() method to avoid being treated as a Stringer.
		res = handleProtoMessage(label, v)
	case bool:
		res = append(res, attribute.Bool(label, v))
	case int, int8, int16, int32, int64:
		res = append(res, attribute.Int64(label, reflect.ValueOf(v).Int()))
	case float32, float64:
		res = append(res, attribute.Float64(label, reflect.ValueOf(v).Float()))
	case fmt.Stringer:
		res = append(res, handleStringerType(label, v))
	case string:
		res = append(res, attribute.String(label, v))
	default:
		res = handleReflectType(label, v)
	}

	return res
}

func handleProtoMessage(label string, data proto.Message) []attribute.KeyValue {
	var res []attribute.KeyValue

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct &&
		value.Elem().CanInterface() {
		value = value.Elem()
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			val := value.Field(i)
			otelAttr := validateAndGetReflectOTel(
				strings.ToLower(fmt.Sprintf("%s.%s", label, field.Name)),
				val,
			)
			res = append(res, otelAttr...)
		}
	}

	return res
}

// handleReflectType recursively processes data based on its type using reflection.
// It converts pointers, structs, maps, arrays, and slices into OpenTelemetry attributes.
func handleReflectType(label string, data interface{}) []attribute.KeyValue {
	var res []attribute.KeyValue
	var commonConfig = config.ProvideDisbursementConfig()

	allowOtelAttributesTruncate := commonConfig.GetEnableConfigAllowTruncateAttributesOtel()

	switch dataType := reflect.TypeOf(data); dataType.Kind() {
	case reflect.Pointer:
		val := reflect.ValueOf(data).Elem()
		res = append(res, validateAndGetReflectOTel(label, val)...)
	case reflect.Struct:
		structNumField := dataType.NumField()
		for i := 0; i < structNumField; i++ {
			if allowOtelAttributesTruncate && i == 20 {
				res = append(
					res,
					attribute.String(
						label+"."+strconv.Itoa(i),
						" ...struct with total field "+strconv.Itoa(structNumField)+" truncated",
					),
				)

				return res
			}

			field := dataType.Field(i)
			val := reflect.ValueOf(data).Field(i)
			res = append(
				res,
				validateAndGetReflectOTel(
					strings.ToLower(fmt.Sprintf("%s.%s", label, field.Name)),
					val,
				)...)
		}
	case reflect.Map:
		vals := reflect.ValueOf(data)
		for i, key := range vals.MapKeys() {
			if allowOtelAttributesTruncate && i == 20 {
				res = append(
					res,
					attribute.String(
						label+"."+strconv.Itoa(i),
						" ...map truncated",
					),
				)

				return res
			}

			val := vals.MapIndex(key)
			res = append(
				res,
				validateAndGetReflectOTel(
					strings.ToLower(fmt.Sprintf("%s.%s", label, key)),
					val,
				)...)
		}
	case reflect.Array, reflect.Slice:
		if dataType.Elem().Kind() == reflect.Uint8 {
			// It's an array of bytes, so convert it to a string (might be json data)
			val := string(data.([]byte))
			if allowOtelAttributesTruncate && len(val) > 100 {
				val = val[:100] + " ...truncated"
			}

			res = append(res, attribute.String(label, val))

			return res
		}

		vals := reflect.ValueOf(data)

		valsLen := vals.Len()
		for i := 0; i < valsLen; i++ {
			if allowOtelAttributesTruncate && i == 20 {
				res = append(
					res,
					attribute.String(
						label+"."+strconv.Itoa(i),
						"...array with total length "+strconv.Itoa(valsLen)+" truncated",
					),
				)

				return res
			}

			val := vals.Index(i)
			res = append(
				res,
				validateAndGetReflectOTel(
					strings.ToLower(fmt.Sprintf("%s.%d", label, i)),
					val,
				)...)
		}
	default:
		val := fmt.Sprintf("%v", data)
		if allowOtelAttributesTruncate && len(val) > 500 {
			val = val[:500] + " ...truncated"
		}

		res = append(res, attribute.String(label, val))
	}

	return res
}

func handleStringerType(label string, data fmt.Stringer) attribute.KeyValue {
	dataType := reflect.TypeOf(data)
	value := reflect.ValueOf(data)

	if dataType.Kind() == reflect.Ptr && value.IsNil() {
		return attribute.String(label, "null")
	}

	return attribute.Stringer(label, data)
}

func validateAndGetReflectOTel(
	label string,
	val reflect.Value,
) (res []attribute.KeyValue) {
	if !val.IsValid() {
		res = append(res, attribute.String(label, "null"))
		return
	}

	if val.CanInterface() {
		res = ToOtelAttributes(label, val.Interface())
	}

	return
}

// HeadersToAttributes converts HTTP headers into OpenTelemetry attributes.
func HTTPHeadersToAttributes(prefix string, h http.Header) []attribute.KeyValue {
	key := func(k string) attribute.Key {
		k = strings.ToLower(k)
		k = strings.ReplaceAll(k, "-", "_")
		k = fmt.Sprintf("%s.%s", prefix, k)

		return attribute.Key(k)
	}

	attrs := make([]attribute.KeyValue, 0, len(h))
	for k, v := range h {
		attrs = append(attrs, key(k).StringSlice(v))
	}

	return attrs
}
