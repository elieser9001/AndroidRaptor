package android

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func ExecuteCommand(cmd string) (string, error) {
	result, err := exec.Command("sh", "-c", cmd).Output()

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func ExecuteCommandToFile(cmd string, fPath string) error {
	cmd += " > " + fPath
	_, err := exec.Command("sh", "-c", cmd).Output()

	if err != nil {
		return err
	}

	return nil
}

func GetNotifications(fPath string) error {
	err := ExecuteCommandToFile("dumpsys notification --noredact | grep tickerText | grep -v null | cut -d= -f2", fPath)
	return err
}

func StartCamera() error {
	_, err := ExecuteCommand("am start -a android.media.action.IMAGE_CAPTURE")
	return err
}

func GetUserAppsList(fPath string) error {
	err := ExecuteCommandToFile("pm list packages -3", fPath)
	return err
}

func GetScreenshot(capturePath string) error {
	_, err := ExecuteCommand("screencap -p > " + capturePath)

	return err
}

func OpenURL(url string) error {
	_, err := ExecuteCommand("am start -a android.intent.action.VIEW -d " + url)

	return err
}

func TurnScreenOFF() error {
	_, err := ExecuteCommand("input keyevent KEYCODE_POWER")

	return err
}

func GetGoogleContacts(fPath string) error {
	err := ExecuteCommandToFile("content query --uri content://contacts/phones/  --projection display_name:number:notes", fPath)
	return err
}

func GetWifiIP() (string, error) {
	result, err := ExecuteCommand("ip addr show wlan0 | grep 'inet ' | cut -d ' ' -f 6 | cut -d / -f 1 | tr -d '\n'")

	if err != nil {
		return "", err
	}

	return result, nil
}

func GetWifiSSID() (string, error) {
	result, err := ExecuteCommand("dumpsys wifi | grep -o 'SSID=\".*\"' | cut -d '\"' -f 2 | tr -d '\n'")

	if err != nil {
		return "", err
	}

	return result, nil
}

func GetWifiMAC() (string, error) {
	result, err := ExecuteCommand("dumpsys wifi | grep -o 'BSSID=.*duration' | cut -d '=' -f 2 | cut -d ',' -f 1 | tr -d '\n'")

	if err != nil {
		return "", err
	}

	return result, nil
}

func GetCurrentActivity() (string, error) {
	result, err := ExecuteCommand("dumpsys activity a . | grep -E 'ResumedActivity' | cut -d ' ' -f 8")

	if err != nil {
		return "", err
	}

	return result, nil
}

func GetCurrentPacketName() (string, error) {
	c, err := GetCurrentActivity()

	if err != nil {
		return "", err
	}

	s := strings.Split(c, "/")

	if len(s) <= 0 {
		return "", fmt.Errorf("no current activity detected")
	}

	return s[0], nil
}

func CloseApp(packageName string) error {
	_, err := ExecuteCommand("am force-stop " + packageName)
	return err
}

func CloseCurrentApp() error {
	ca, err := GetCurrentPacketName()

	if err != nil {
		return err
	}

	if len(ca) <= 0 {
		return fmt.Errorf("no current app detected")
	}

	if err := CloseApp(ca); err != nil {
		return err
	}

	return nil
}

func SendClick(x, y string) error {
	_, err := ExecuteCommand("input tap " + x + " " + y)
	return err
}

func GetDir(path string, fPath string) error {
	err := ExecuteCommandToFile("ls -la "+path, fPath)
	return err
}

func IsScreenOn() (bool, error) {
	result, err := ExecuteCommand("dumpsys input_method | grep -c 'mInteractive=true' | tr -d '\n'")

	if err != nil {
		return false, err
	}

	if result == "1" {
		return true, nil
	} else {
		return false, nil
	}
}

func TurnScreenOn() error {
	screenAlreadyOn, err := IsScreenOn()

	if err != nil {
		return err
	}

	if !screenAlreadyOn {
		_, err := ExecuteCommand("input keyevent KEYCODE_WAKEUP")

		if err != nil {
			return err
		}
	}

	return nil
}

func SwipeUp() error {
	_, err := ExecuteCommand("input touchscreen swipe 530 1420 530 1120")

	return err
}

func GetBatteryLevel() (int, error) {
	result, err := ExecuteCommand("dumpsys battery | sed -n '/level: /p' | tr -d '[a-z] | : | \n'")

	if err != nil {
		return 0, err
	}

	r, err := strconv.Atoi(result)

	if err != nil {
		return 0, err
	}

	return r, nil
}

func Shutdown() error {
	_, err := ExecuteCommand("reboot -p")

	return err
}

func Reboot() error {
	_, err := ExecuteCommand("reboot")

	return err
}
