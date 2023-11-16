package sanitizer

import (
	"encoding/json"
	"net"
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
