package monitor

import (
	"cmp"
	"context"
	"math"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var ServerController = &serverController{}

type serverController struct {
}

func (c *serverController) GetServerInfo(ctx *gin.Context) {
	server := &Server{}
	server.Sys = InitOS()
	server.Cpu, _ = InitCPU()
	server.Mem, _ = InitRAM()
	server.Disk, _ = InitDisk()
	server.GoInfo = InitGoInfo()
	rsp.OkWithData(server, ctx)
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type Server struct {
	Sys    Sys    `json:"sys"`
	Cpu    Cpu    `json:"cpu"`
	Mem    Mem    `json:"mem"`
	Disk   []Disk `json:"sysFiles"`
	GoInfo GoInfo `json:"goInfo"`
}

type Sys struct {
	OsName       string `json:"osName"`
	OsArch       string `json:"osArch"`
	ComputerName string `json:"computerName"`
	ComputerIp   string `json:"computerIp"`
	UserDir      string `json:"userDir"`
}

type Cpu struct {
	Cores int     `json:"cpuNum"`
	Sys   float64 `json:"sys"`
	Used  float64 `json:"used"`
	Wait  float64 `json:"wait"`
	Free  float64 `json:"free"`
}

type Mem struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Free        float64 `json:"free"`
	UsedPercent int     `json:"usage"`
}

type Disk struct {
	Device      string `json:"sysTypeName"`
	MountPoint  string `json:"dirName"`
	Fstype      string `json:"typeName"`
	Total       string `json:"total"`
	Used        string `json:"used"`
	Free        string `json:"free"`
	UsedPercent int    `json:"usage"`
}

type GoInfo struct {
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
	Mem          int    `json:"mem"`
	NumGC        int    `json:"numGC"`
}

func InitOS() (o Sys) {
	o.OsName = runtime.GOOS
	o.OsArch = runtime.GOARCH
	o.ComputerName, _ = os.Hostname()
	o.ComputerIp = getLocalHost()
	return o
}

func InitCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if stats, err := caclCPU(); err != nil {
		return c, err
	} else {
		c.Used = math.Round(stats[0])
		c.Sys = math.Round(stats[1])
		c.Free = math.Round(stats[2])
	}
	return c, nil
}

func InitRAM() (r Mem, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.Used = math.Round(float64(u.Used)*10/GB) / 10
		r.Total = math.Round(float64(u.Total)*10/GB) / 10
		r.Free = math.Round(float64(u.Available)*10/GB) / 10
		r.UsedPercent = int(u.UsedPercent)
	}
	return r, nil
}

func InitDisk() (d []Disk, err error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return
	}
	for _, p := range parts {
		mp := p.Mountpoint
		if u, err := disk.Usage(mp); err != nil {
			return d, err
		} else {
			d = append(d, Disk{
				Device:      p.Device,
				MountPoint:  mp,
				Fstype:      p.Fstype,
				Total:       util.HumanByteSize(int64(u.Total)),
				Used:        util.HumanByteSize(int64(u.Used)),
				Free:        util.HumanByteSize(int64(u.Free)),
				UsedPercent: int(u.UsedPercent),
			})
		}
	}
	return d, nil
}

func InitGoInfo() GoInfo {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return GoInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		Mem:          int(ms.Sys) / MB,
		NumGC:        int(ms.NumGC),
	}
}

func getLocalHost() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		gb.Logger.Error(err)
		return ""
	}

	pvtIps, pubIps := []string{}, []string{}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			if ipNet.IP.IsPrivate() {
				pvtIps = append(pvtIps, ipNet.IP.String())
			} else {
				pubIps = append(pubIps, ipNet.IP.String())
			}
		}
	}
	pubIps = util.NewList(pubIps).Filter(func(ip string) bool { return ip != "" })
	if len(pubIps) > 0 {
		return pubIps[0]
	}
	pvtIps = util.NewList(pvtIps).Filter(func(ip string) bool { return ip != "" }).
		Order(func(a, b string) int { return cmp.Compare(b, a) })
	if len(pvtIps) > 0 {
		return pvtIps[0]
	}

	return ""
}

func caclCPU() ([]float64, error) {
	// Get CPU usage at the start of the interval.
	cpuTimes1, err := cpu.TimesWithContext(context.Background(), false)
	if err != nil {
		return nil, err
	}
	cpu1 := cpuTimes1[0]

	time.Sleep(time.Second)

	// And at the end of the interval.
	cpuTimes2, err := cpu.TimesWithContext(context.Background(), false)
	if err != nil {
		return nil, err
	}
	cpu2 := cpuTimes2[0]

	result := make([]float64, 0)
	total := cpu2.Total() - cpu1.Total()
	result = append(result, (cpu2.User-cpu1.User)/total*100)
	result = append(result, (cpu2.System-cpu1.System)/total*100)
	result = append(result, (cpu2.Idle-cpu1.Idle)/total*100)
	return result, nil
}
