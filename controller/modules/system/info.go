package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	"github.com/reef-pi/reef-pi/controller/utils"
)

const (
	_modelFilePath = "/proc/device-tree/model"
)

//swagger:model systemSummary
type Summary struct {
	Name           string `json:"name"`
	IP             string `json:"ip"`
	CurrentTime    string `json:"current_time"`
	Uptime         string `json:"uptime"`
	CPUTemperature string `json:"cpu_temperature"`
	Version        string `json:"version"`
	Model          string `json:"model"`
}

func (c *Controller) ComputeSummary() Summary {
	ip, err := c.HostIP(c.config.Interface)
	if err != nil {
		log.Println("ERROR: Failed to detect ip for interface '"+c.config.Interface+". Error:", err)
		ip = "unknown"
	}
	temp, err := c.CPUTemperature()
	if err != nil {
		log.Println("ERROR:Failed to get controller temperature. Error:", err)
		temp = "unknown"
	}
	s := Summary{
		Name:           c.config.Name,
		CurrentTime:    time.Now().Format("Mon Jan 2 15:04:05, 2006"),
		IP:             ip,
		Uptime:         c.Uptime(),
		CPUTemperature: string(temp),
		Version:        c.config.Version,
		Model:          c.GetModel(),
	}
	return s
}

func (c *Controller) HostIP(i string) (string, error) {
	iface, err := net.InterfaceByName(i)
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, v := range addrs {
		switch s := v.(type) {
		case *net.IPNet:
			if s.IP.To4() != nil {
				return s.IP.To4().String(), nil
			}
		}
	}
	return "", fmt.Errorf("Cant detect IP of interface:%s", i)
}

// temp=36.9'C
func (c *Controller) CPUTemperature() (string, error) {
	out, err := utils.Command("vcgencmd", "measure_temp").WithDevMode(c.config.DevMode).CombinedOutput()
	if err != nil {
		return "", err
	}
	if c.config.DevMode {
		out = []byte("foo=23.23 ")
	}
	return strings.Split(string(out), "=")[1], nil
}

func (c *Controller) GetModel() string {
	data, err := ioutil.ReadFile("/proc/device-tree/model")
	if err != nil {
		log.Println("Failed to detect Raspberry Pi version. Error:", err)
		return "unknown"
	}
	return strings.TrimSpace(string(data))

}
