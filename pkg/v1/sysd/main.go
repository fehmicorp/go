package main

import "fmt"

func main() {
	fingerprint, err := GenerateDeviceFingerprint()
	if err != nil {
		fmt.Printf("Error generating fingerprint: %v\n", err)
		return
	}

	fmt.Printf("Unique Device Fingerprint: %s\n", fingerprint)
}
