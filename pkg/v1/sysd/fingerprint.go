package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func GenerateDeviceFingerprint() (string, error) {
	var builder strings.Builder

	// 1. CPU & OS Architecture Info
	builder.WriteString(fmt.Sprintf("OS:%s;", runtime.GOOS))
	builder.WriteString(fmt.Sprintf("ARCH:%s;", runtime.GOARCH))
	builder.WriteString(fmt.Sprintf("CPUs:%d;", runtime.NumCPU()))

	// 2. Hostname
	hostname, err := os.Hostname()
	if err == nil {
		builder.WriteString(fmt.Sprintf("HOST:%s;", hostname))
	}

	// 3. Network Interfaces & MAC Addresses
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			// Skip loopback and unpopulated MAC addresses
			if iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			mac := iface.HardwareAddr.String()
			if mac != "" {
				builder.WriteString(fmt.Sprintf("MAC:%s;", mac))
				builder.WriteString(fmt.Sprintf("NETNAME:%s;", iface.Name))
			}
		}
	}

	// Additional system-specific metrics can be concatenated here
	// (e.g., disk serial numbers via external system calls or platform-specific libraries).

	// Compute SHA-256 Hash of the combined attributes
	hasher := sha256.New()
	hasher.Write([]byte(builder.String()))
	fingerprintHash := hex.EncodeToString(hasher.Sum(nil))

	return fingerprintHash, nil
}
