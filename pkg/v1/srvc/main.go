package srvc

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

// Install registers the service with the OS Service Manager (SCM on Windows, Systemd on Linux)
func (s *Services) Install(executablePath string) error {
	if runtime.GOOS == "windows" {
		return s.installWindows(executablePath)
	}
	return s.installLinux(executablePath)
}

// Uninstall removes the service registration from the OS
func (s *Services) Uninstall() error {
	if runtime.GOOS == "windows" {
		return s.uninstallWindows()
	}
	return s.uninstallLinux()
}

// Start launches the service
func (s *Services) Start() error {
	if runtime.GOOS == "windows" {
		return s.controlWindows("start")
	}
	return exec.Command("systemctl", "start", s.Name).Run()
}

// Stop halts the service
func (s *Services) Stop() error {
	if runtime.GOOS == "windows" {
		return s.controlWindows("stop")
	}
	return exec.Command("systemctl", "stop", s.Name).Run()
}

// Restart reloads the service
func (s *Services) Restart() error {
	if runtime.GOOS == "windows" {
		if err := s.Stop(); err != nil {
			// Ignore stop error if already stopped, then try starting
		}
		return s.Start()
	}
	return exec.Command("systemctl", "restart", s.Name).Run()
}

// --- Windows Internal Implementations ---

func (s *Services) installWindows(execPath string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	cfg := mgr.Config{
		DisplayName: s.Title,
		Description: s.Desc,
		StartType:   mgr.StartAutomatic,
	}

	ins, err := m.CreateService(s.Name, execPath, cfg)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil
}

func (s *Services) uninstallWindows() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(s.Name)
	if err != nil {
		return err
	}
	defer srv.Close()

	return srv.Delete()
}

func (s *Services) controlWindows(action string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(s.Name)
	if err != nil {
		return err
	}
	defer srv.Close()

	status, err := srv.Query()
	if err != nil {
		return err
	}

	switch action {
	case "start":
		if status.State == 4 { // Running
			return nil
		}
		_, err = srv.Control(15) // Stop code or start logic via custom command / pass arguments
		// Standard windows start control:
		return srv.Start()
	case "stop":
		_, err = srv.Control(1) // 1 = svc.ControlRequest(Stop)
		return err
	}
	return fmt.Errorf("unknown control action: %s", action)
}

func (s *Services) refreshStatusWindows() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(s.Name)
	if err != nil {
		s.Status = Status{Active: false, Running: false, Startup: false}
		return nil // Service might not be installed yet
	}
	defer srv.Close()

	status, err := srv.Query()
	if err != nil {
		return err
	}

	cfg, err := srv.Config()
	if err != nil {
		return err
	}

	s.Status = Status{
		Active:  status.State != 1, // 1 = Stopped
		Running: status.State == 4, // 4 = Running
		Startup: cfg.StartType == mgr.StartAutomatic,
	}
	return nil
}

// --- Linux Internal Implementations (Systemd) ---

func (s *Services) installLinux(execPath string) error {
	// Generates standard systemd unit file dynamically
	serviceContent := fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
ExecStart=%s
WorkingDirectory=%s
Restart=always
User=root

[Install]
WantedBy=multi-user.target
`, s.Desc, execPath, s.Package.WorkDir)

	unitPath := fmt.Sprintf("/etc/systemd/system/%s.service", s.Name)
	if err := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' > %s", serviceContent, unitPath)).Run(); err != nil {
		return err
	}

	return exec.Command("systemctl", "enable", s.Name).Run()
}

func (s *Services) uninstallLinux() error {
	_ = exec.Command("systemctl", "stop", s.Name).Run()
	_ = exec.Command("systemctl", "disable", s.Name).Run()

	unitPath := fmt.Sprintf("/etc/systemd/system/%s.service", s.Name)
	if err := exec.Command("rm", "-f", unitPath).Run(); err != nil {
		return err
	}

	return exec.Command("systemctl", "daemon-reload").Run()
}

func (s *Services) refreshStatusLinux() error {
	// Check if running
	runningCmd := exec.Command("systemctl", "is-active", s.Name)
	runningOut, _ := runningCmd.Output()
	isRunning := strings.TrimSpace(string(runningOut)) == "active"

	// Check if enabled (startup)
	enabledCmd := exec.Command("systemctl", "is-enabled", s.Name)
	enabledOut, _ := enabledCmd.Output()
	isStartup := strings.TrimSpace(string(enabledOut)) == "enabled"

	s.Status = Status{
		Active:  isRunning,
		Running: isRunning,
		Startup: isStartup,
	}
	return nil
}
