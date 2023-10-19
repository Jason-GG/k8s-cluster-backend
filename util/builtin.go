package util

import (
	"net"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func UintPtr(s uint) *uint                   { return &s }
func Uint8Ptr(s uint8) *uint8                { return &s }
func Uint16Ptr(s uint16) *uint16             { return &s }
func Uint32Ptr(s uint32) *uint32             { return &s }
func Uint64Ptr(s uint64) *uint64             { return &s }
func Float32Ptr(s float32) *float32          { return &s }
func Float64Ptr(s float64) *float64          { return &s }
func Complex64Ptr(s complex64) *complex64    { return &s }
func Complex128Ptr(s complex128) *complex128 { return &s }
func BoolPtr(s bool) *bool                   { return &s }
func StringPtr(s string) *string             { return &s }
func IntPtr(i int) *int                      { return &i }
func Int16Ptr(i int16) *int16                { return &i }
func Int32Ptr(i int32) *int32                { return &i }
func Int64Ptr(i int64) *int64                { return &i }

func NilString(str *string) bool {
	p := (*reflect.StringHeader)(unsafe.Pointer(str))
	return unsafe.Pointer(p.Data) == nil
}

func ByteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ParseBool(str string) bool {
	switch str {
	case "TRUE", "true", "True", "Yes", "yes", "YES", "1", "t", "T":
		return true
	case "FALSE", "false", "False", "No", "no", "NO", "0", "f", "F":
		return false
	default:
	}
	return false
}

func Connectable(host string, port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
