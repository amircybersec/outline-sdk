package sanitizer

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

func TestIPAddress_MarshalJSON(t *testing.T) {
	// // Example JSON strings
	// jsonIPv4 := `"192.168.1.1"`
	// jsonIPv6 := `"2001:0db8:85a3:0000:0000:8a2e:0370:7334"`

	// // Unmarshal JSON to IPAddress struct
	// var ipAddr4 IPAddressSanitized
	// var ipAddr6 IPAddressSanitized

	type ConnectionInfo struct {
		Proto     string             `json:"Proto"`
		Transport string             `json:"Transport"`
		IP        IPAddressSanitized `json:"clientIPAddress"`
	}

	var c ConnectionInfo

	jsonStringIPv4 := `{
		"Proto": "udp",
		"Transport": "ShadowsSocks",
		"clientIPAddress": "192.168.1.1"
	}`

	err := json.Unmarshal([]byte(jsonStringIPv4), &c)
	if err != nil {
		t.Errorf("Error unmarshaling IPv4: %v\n", err)
	} else {
		fmt.Printf("Unmarshalled IPv4: %s\n", c.IP.Address)
	}

	// err = json.Unmarshal([]byte(jsonIPv6), &ipAddr6)
	// if err != nil {
	// 	t.Errorf("Error unmarshaling IPv6: %v\n", err)
	// } else {
	// 	fmt.Printf("Unmarshalled IPv6: %s\n", ipAddr6.Address)
	// }
}

func TestReportMarshalJSON(t *testing.T) {
	report := Report{
		Name:       "Example Report",
		IPv4Addr:   net.ParseIP("192.168.1.1"),
		IPv6Addr:   net.ParseIP("fe80::1"),
		OtherField: 12345,
		NestedStruct: Nested{
			NestedIP:  net.ParseIP("10.0.0.1"),
			SomeValue: 42,
		},
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		t.Errorf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
