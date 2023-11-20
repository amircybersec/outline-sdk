package sanitizer

import (
	"encoding/json"
	"net"
	"reflect"
)

// IPAddress is a custom type for IP addresses
type IPAddressSanitized struct {
	Address net.IP
	// Private bool: Flag to control Marshaller behavior?
}

// MarshalJSON customizes the JSON representation of IPAddress
func (ip IPAddressSanitized) MarshalJSON() ([]byte, error) {
	// Check if the IP address is an IPv4 or IPv6 address
	if ip.Address.To4() != nil {
		// It's an IPv4 address, marshal as "0.0.0.0"
		return json.Marshal("0.0.0.0")
	} else {
		// It's an IPv6 address, marshal as "::"
		return json.Marshal("::")
	}
}

// UnmarshalJSON customizes the JSON unmarshaling of IPAddress
func (ip *IPAddressSanitized) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data to a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Check if the unmarshaled string is an IPv4 or IPv6 address
	parsedIP := net.ParseIP(str)
	if parsedIP.To4() != nil {
		// It's an IPv4 address, set to "0.0.0.0"
		ip.Address = net.ParseIP("0.0.0.0")
	} else {
		// It's an IPv6 address, set to "::"
		ip.Address = net.ParseIP("::")
	}

	return nil
}

// second approach is to implement MarshalJSON on the Report struct
// instead of the IPAddress struct (see below):

type Report struct {
	Name         string
	IPv4Addr     net.IP
	IPv6Addr     net.IP
	OtherField   int
	NestedStruct Nested
}

type Nested struct {
	NestedIP  net.IP
	SomeValue int
}

func (r Report) MarshalJSON() ([]byte, error) {
	return marshalStruct(reflect.ValueOf(r))
}

func marshalStruct(val reflect.Value) ([]byte, error) {
	typ := val.Type()
	sanitizedData := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check if the field is of type net.IP
		if field.Type() == reflect.TypeOf(net.IP{}) {
			ip := field.Interface().(net.IP)
			sanitizedData[fieldType.Name] = sanitizeIPAddress(ip)
		} else if field.Kind() == reflect.Struct {
			// Handle nested structs
			nestedJSON, err := marshalStruct(field)
			if err != nil {
				return nil, err
			}
			sanitizedData[fieldType.Name] = json.RawMessage(nestedJSON)
		} else {
			sanitizedData[fieldType.Name] = field.Interface()
		}
	}

	return json.Marshal(sanitizedData)
}

func sanitizeIPAddress(ip net.IP) string {
	if ip.To4() != nil {
		return "0.0.0.0"
	}
	return "::"
}
