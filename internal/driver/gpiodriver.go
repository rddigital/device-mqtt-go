package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type GPIODev struct {
	lc   logger.LoggingClient
	gpio int
}

func NewGPIODev(lc logger.LoggingClient) *GPIODev {
	return &GPIODev{lc: lc, gpio: -1}
}

func (dev *GPIODev) ExportGPIO(gpio int) error {
	err := exportgpio(gpio)
	if err == nil {
		dev.gpio = gpio
	}
	return err
}

func (dev *GPIODev) UnexportGPIO(gpio int) error {
	err := unexportgpio(gpio)
	if err == nil {
		dev.gpio = -1
	}
	return err
}

func (dev *GPIODev) SetDirection(direction string) error {
	if dev.gpio == -1 {
		return errors.New("Please export gpio first")
	}
	return setgpiodirection(dev.gpio, direction)
}

func (dev *GPIODev) GetDirection() (string, error) {
	if dev.gpio == -1 {
		return "", errors.New("Please export gpio first")
	}
	direction, err := getgpiodirection(dev.gpio)
	if err != nil {
		return "", err
	} else {
		res, _ := json.Marshal(map[string]interface{}{"gpio": dev.gpio, "direction": direction})
		return string(res), err
	}
}

func (dev *GPIODev) SetGPIO(gpio int, value bool) error {
	if gpio < 0 || gpio > 255 {
		return errors.New("gpio out of range")
	}

	exportgpio(gpio)
	err := setgpiodirection(gpio, "out")
	if err == nil {
		return setgpiovalue(gpio, value)
	}

	return err
}

func (dev *GPIODev) GetGPIO(gpio int) (string, error) {

	if gpio < 0 || gpio > 255 {
		return "", errors.New("gpio out of range")
	}

	exportgpio(gpio)
	err := setgpiodirection(gpio, "in")
	if err == nil {
		gpiovalue, err := getgpiovalue(gpio)
		if err != nil {
			return "", err
		} else {
			if gpiovalue == 0 {
				return "false", nil
			} else {
				return "true", nil
			}
		}
	}
	return "", err
}

func exportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	} else {
		return ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	}
}

func unexportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ioutil.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	} else {
		return nil
	}
}

func setgpiodirection(gpioNum int, direction string) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var way string
		if direction == "in" {
			way = "in"
		} else {
			way = "out"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum), []byte(way), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiodirection(gpioNum int) (string, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		direction, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum))
		if err != nil {
			return "", err
		} else {
			return strings.Replace(string(direction), "\n", "", -1), err
		}
	} else {
		return "", errors.New("Please export gpio first")
	}
}

func setgpiovalue(gpioNum int, value bool) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value == false {
			tmp = "0"
		} else {
			tmp = "1"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum), []byte(tmp), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiovalue(gpioNum int) (int, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum))
		if err != nil {
			return 0, err
		} else {
			value, _ := strconv.Atoi(strings.Replace(string(ret), "\n", "", -1))
			return value, err
		}
	} else {
		return -1, errors.New("Please export gpio first")
	}
}
