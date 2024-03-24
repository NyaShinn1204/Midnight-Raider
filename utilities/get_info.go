package get_info

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func GetHwid() string {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		cmd := "cat /etc/machine-id"
		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			fmt.Println("getHWID error:", err)
			return "unknown"
		}
		return strings.TrimSpace(string(out))
	} else if runtime.GOOS == "windows" {
		cmd := "powershell -Command (Get-WmiObject -Class Win32_ComputerSystemProduct).UUID"
		out, err := exec.Command("cmd", "/C", cmd).Output()
		if err != nil {
			fmt.Println("getHWID error:", err)
			return "unknown"
		}
		return strings.TrimSpace(string(out))
	} else {
		fmt.Println("Unsupported OS")
		return "unknown"
	}
}
