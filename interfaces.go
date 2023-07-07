package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"unsafe"
)

type HSLInterface struct {
	Name string `json:"name"`
}

type HSLAddress struct {
	Address string `json:"address"`
}

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

//export interfaces
func interfaces(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	interfaces := []HSLInterface{}

	ifs, err := net.Interfaces()
	if err != nil {
		SetException(hhc, err.Error())
		return
	}
	for _, i := range ifs {
		interfaces = append(interfaces, HSLInterface{Name: i.Name})
	}

	SetReturnValueToAny(ret, interfaces)
}

//export local_ips
func local_ips(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	x, err := GetArgumentAsString(args, 0, false)
	if err != nil {
		SetException(hhc, err.Error())
		return
	}

	ips := []HSLAddress{}

	ifs, err := net.Interfaces()
	if err != nil {
		SetException(hhc, err.Error())
		return
	}
	for _, i := range ifs {
		if x != "" && i.Name != x {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			SetException(hhc, err.Error())
			return
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.IsGlobalUnicast() {
					ips = append(ips, HSLAddress{Address: v.IP.String()})
				}
			case *net.IPAddr:
				if v.IP.IsGlobalUnicast() {
					ips = append(ips, HSLAddress{Address: v.IP.String()})
				}
			}
		}
	}
	SetReturnValueToAny(ret, ips)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	C.HalonMTA_hsl_module_register_function(hhrc, C.CString("interfaces"), nil)
	C.HalonMTA_hsl_module_register_function(hhrc, C.CString("local_ips"), nil)
	return true
}

func SetReturnValueToAny(ret *C.HalonHSLValue, val interface{}) error {
	x, err := json.Marshal(val)
	if err != nil {
		return err
	}
	y := C.CString(string(x))
	defer C.free(unsafe.Pointer(y))
	var z *C.char
	if !(C.HalonMTA_hsl_value_from_json(ret, y, &z, nil)) {
		if z != nil {
			err = errors.New(C.GoString(z))
			C.free(unsafe.Pointer(z))
		} else {
			err = errors.New("failed to parse return value")
		}
		return err
	}
	return nil
}

func GetArgumentAsString(args *C.HalonHSLArguments, pos uint64, required bool) (string, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return "", fmt.Errorf("missing argument at position %d", pos)
		} else {
			return "", nil
		}
	}
	var y *C.char
	if C.HalonMTA_hsl_value_get(x, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&y), nil) {
		return C.GoString(y), nil
	} else {
		return "", fmt.Errorf("invalid argument at position %d", pos)
	}
}

func SetException(hhc *C.HalonHSLContext, msg string) {
	x := C.CString(msg)
	y := unsafe.Pointer(x)
	defer C.free(y)
	exception := C.HalonMTA_hsl_throw(hhc)
	C.HalonMTA_hsl_value_set(exception, C.HALONMTA_HSL_TYPE_EXCEPTION, y, 0)
}

func main() {}
