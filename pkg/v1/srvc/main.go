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

func (s *Services) Stop() error {
	if runtime.GOOS == "windows" {
		// Stop via Windows Service Control Manager (SCM)
		cmdStr := fmt.Sprintf(`Stop-Service -Name "%s" -Force -ErrorAction SilentlyContinue`, s.Name)
		cmd := exec.Command("powershell", "-Command", cmdStr)
		if err := cmd.Run(); err == nil {
			s.Status.Running = false
			s.Status.Active = false
			return nil
		}

		// Fallback for Docker if service name relates to docker
		if strings.Contains(strings.ToLower(s.Name), "docker") {
			dockerCmd := exec.Command("net", "stop", "com.docker.service")
			if dockerCmd.Run() == nil {
				s.Status.Running = false
				s.Status.Active = false
				return nil
			}
		}

		return fmt.Errorf("failed to stop Windows service: %s", s.Name)
	}

	// Linux / macOS fallback using systemctl
	cmd := exec.Command("systemctl", "stop", s.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop systemctl service %s: %w", s.Name, err)
	}

	s.Status.Running = false
	s.Status.Active = false
	return nil
}

// Restart performs a stop followed by a start operation
func (s *Services) Restart() error {
	if runtime.GOOS == "windows" {
		// Restart via Windows Service Control Manager (SCM)
		cmdStr := fmt.Sprintf(`Restart-Service -Name "%s" -Force -ErrorAction SilentlyContinue`, s.Name)
		cmd := exec.Command("powershell", "-Command", cmdStr)
		if err := cmd.Run(); err == nil {
			s.Status.Running = true
			s.Status.Active = true
			return nil
		}

		// Fallback if Restart-Service fails: explicit Stop then Start
		if err := s.Stop(); err != nil {
			return err
		}
		return s.Start()
	}

	// Linux / macOS fallback using systemctl
	cmd := exec.Command("systemctl", "restart", s.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart systemctl service %s: %w", s.Name, err)
	}

	s.Status.Running = true
	s.Status.Active = true
	return nil
}

// Start initiates the service (included for complete control mapping)
func (s *Services) Start() error {
	if runtime.GOOS == "windows" {
		cmdStr := fmt.Sprintf(`Start-Service -Name "%s" -ErrorAction SilentlyContinue`, s.Name)
		cmd := exec.Command("powershell", "-Command", cmdStr)
		if err := cmd.Run(); err == nil {
			s.Status.Running = true
			s.Status.Active = true
			return nil
		}
		return fmt.Errorf("failed to start Windows service: %s", s.Name)
	}

	cmd := exec.Command("systemctl", "start", s.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start systemctl service %s: %w", s.Name, err)
	}

	s.Status.Running = true
	s.Status.Active = true
	return nil
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
