package main

import (
	"errors"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var errUsageOfCPU error = errors.New("Unable to measure CPU usage.")

func main() {

}

func cpuUsage() ([]float64, error) {
	percent, err := cpu.Percent(time.Duration(1)*time.Second, false)
	if err != nil {
		return nil, errUsageOfCPU
	}
	return percent, nil
}
