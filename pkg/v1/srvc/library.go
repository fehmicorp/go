package srvc

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// VersionCheckResult holds status information about local package installations and version matches
type VersionCheckResult struct {
	Installed       bool
	LocalVersion    string
	RemoteVersion   string
	UpdateAvailable bool
}

// CheckLocalVersion inspects the host system to see if the package is installed,
// retrieves its version, and compares it against the remote target version.
func (s *Services) CheckLocalVersion() (VersionCheckResult, error) {
	if runtime.GOOS == "windows" {
		return s.checkWindowsVersion()
	}
	return s.checkLinuxVersion()
}

// --- Windows Version Checking ---

func (s *Services) checkWindowsVersion() (VersionCheckResult, error) {
	result := VersionCheckResult{
		Installed:     false,
		RemoteVersion: s.Package.Version,
	}

	// Example implementation querying registry or executable paths common for Windows software/Docker
	// For Docker Desktop or general Windows services, check SCM configuration or registry path
	cmdStr := fmt.Sprintf(`Get-Service -Name "%s" -ErrorAction SilentlyContinue`, s.Name)
	cmd := exec.Command("powershell", "-Command", cmdStr)
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return result, nil // Not installed
	}

	result.Installed = true

	// Query product or file version via PowerShell for more precision
	versionCmdStr := fmt.Sprintf(`(Get-WmiObject -Class Win32_Product | Where-Object {$_.Name -like "*%s*"}).Version`, s.Title)
	vCmd := exec.Command("powershell", "-Command", versionCmdStr)
	vOut, vErr := vCmd.Output()

	localVer := strings.TrimSpace(string(vOut))
	if vErr != nil || localVer == "" {
		// Fallback: assume installed with an unknown live version if SCM status succeeded
		localVer = "unknown"
	}

	result.LocalVersion = localVer

	// Version comparison logic (Semantic version check or exact match string comparison)
	if result.LocalVersion != "unknown" && result.RemoteVersion != "latest" {
		if compareVersions(result.LocalVersion, result.RemoteVersion) < 0 {
			result.UpdateAvailable = true
		}
	}

	return result, nil
}

// --- Linux Version Checking ---

func (s *Services) checkLinuxVersion() (VersionCheckResult, error) {
	result := VersionCheckResult{
		Installed:     false,
		RemoteVersion: s.Package.Version,
	}

	// Check if systemd unit or package binary exists / is active
	cmd := exec.Command("systemctl", "is-active", s.Name)
	err := cmd.Run()

	// If systemctl fails or returns inactive, check package manager via dpkg / rpm
	isInstalledViaSystemd := err == nil

	var localVer string
	if strings.Contains(s.Package.Name, "docker") {
		vCmd := exec.Command("docker", "--version")
		vOut, vErr := vCmd.Output()
		if vErr == nil {
			isInstalledViaSystemd = true
			localVer = parseDockerVersion(string(vOut))
		}
	} else {
		// Generic dpkg check for Debian/Ubuntu
		vCmd := exec.Command("dpkg", "-s", s.Package.Name)
		vOut, vErr := vCmd.Output()
		if vErr == nil {
			isInstalledViaSystemd = true
			localVer = parseDpkgVersion(string(vOut))
		}
	}

	if !isInstalledViaSystemd && localVer == "" {
		return result, nil
	}

	result.Installed = true
	if localVer == "" {
		localVer = "unknown"
	}
	result.LocalVersion = localVer

	// Evaluate updates if a specific remote version is set
	if result.LocalVersion != "unknown" && result.RemoteVersion != "latest" {
		if compareVersions(result.LocalVersion, result.RemoteVersion) < 0 {
			result.UpdateAvailable = true
		}
	}

	return result, nil
}

// Helper: Basic version comparison (returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2)
func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		n1, n2 := 0, 0
		if i < len(parts1) {
			_, _ = fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			_, _ = fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}

func parseDockerVersion(output string) string {
	// Example output: "Docker version 24.0.5, build ced0996"
	fields := strings.Fields(output)
	for i, f := range fields {
		if f == "version" && i+1 < len(fields) {
			return strings.TrimRight(fields[i+1], ",")
		}
	}
	return "unknown"
}

func parseDpkgVersion(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Version:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}
