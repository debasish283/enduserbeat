
package hardware

import (
  "github.com/elastic/beats/libbeat/logp"
  "github.com/elastic/beats/libbeat/common"
  "github.com/elastic/beats/libbeat/beat"
  "os/exec"
  "strings"
  "regexp"
  "time"
  "context"
  "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	// "strings"
)


type InfoStat struct {
    Hostname             string `json:"hostname"`
    Uptime               uint64 `json:"uptime"`
    BootTime             uint64 `json:"bootTime"`
    Procs                uint64 `json:"procs"`           // number of processes
    OS                   string `json:"os"`              // ex: freebsd, linux
    Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
    PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
    PlatformVersion      string `json:"platformVersion"` // version of the complete OS
    KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
    VirtualizationSystem string `json:"virtualizationSystem"`
    VirtualizationRole   string `json:"virtualizationRole"` // guest or host
    HostID               string `json:"hostid"`             // ex: uuid
}


type CpuInfoStat []struct {
	CPU        int      `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int      `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int      `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        int      `json:"mhz"`
	CacheSize  int      `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}


func PostData(client beat.Client) ( error) {
	memory, _ := mem.VirtualMemory()
  out, _ := exec.Command("cmd", "/C", "wmic bios get Manufacturer").Output()
  name, _ := exec.Command("cmd", "/C", "wmic csproduct get name").Output()
  servicepackmajorversion, _ := exec.Command("cmd", "/C", "wmic os get servicepackmajorversion").Output()
  user, _ := exec.Command("cmd", "/C", "wmic netlogin get name").Output()
  lastbootuptime, _ := exec.Command("cmd", "/C", "wmic os get lastbootuptime").Output()
  InstalledOn, _ := exec.Command("cmd", "/C", "wmic qfe get  InstalledOn").Output()
  //////////////UUUUUUUUUUUUUUUU/////////////////////////

  manufacturer := strings.TrimSpace(strings.Trim(string(out), "Manufacturer"))
  modelname := strings.TrimSpace(strings.Trim(string(name), "Name"))
  servicepack := strings.TrimSpace(strings.Trim(string(servicepackmajorversion), "ServicePackMajorVersion"))
  username := strings.TrimSpace(strings.Trim(string(user), "Name"))
  usernamef := strings.TrimSpace(strings.Trim(string(username), "NT AUTHORITY\\SYSTEM"))
  booTime := strings.TrimSpace(strings.Trim(string(lastbootuptime), "LastBootUpTime"))
  re := regexp.MustCompile(".{0,10}\\z")
  match := re.FindStringSubmatch(strings.TrimSpace(string(InstalledOn)))
  ///////////////////////ffffffffffffff/////////////
  InfoStat, _ := host.InfoWithContext(context.Background())
  CpuInfoStat, _ := cpu.InfoWithContext(context.Background())
  UsageStat, _ := disk.Usage("/")

  event := beat.Event{
     Timestamp: time.Now(),
     Fields:common.MapStr{
       "type":       "smartcenterEndUserHardware",
   		 "Hostname":    InfoStat.Hostname,
   		 "Uptime":    InfoStat.Uptime,
   		 "OS":    InfoStat.OS,
   		 "Platform":    InfoStat.Platform,
   		 "KernelVersion":    InfoStat.KernelVersion,
   		 "VirtualizationSystem":    InfoStat.VirtualizationSystem,
   		 "VirtualizationRole":    InfoStat.VirtualizationRole,
   		 "Model": CpuInfoStat[0].Model,
   		 "CoreID": CpuInfoStat[0].CoreID,
   		 "ModelName": CpuInfoStat[0].ModelName,
   		 "Microcode": CpuInfoStat[0].Microcode,
   		 "Path": UsageStat.Path,
   		 "Fstype": UsageStat.Fstype,
   		 "HddTotal": UsageStat.Total,
   		 "HddFree": UsageStat.Free,
   		 "HddUsed": UsageStat.Used,
   		 "HddUsedPercent": UsageStat.UsedPercent,
   		 "Ram": memory.Total,
   		 "FreeRam": memory.Free,
   		 "RamUsedPercent":memory.UsedPercent,
       "manufacturer": manufacturer,
       "modelname": modelname,
       "servicepack": servicepack,
       "username": username,
       "usernamef": usernamef,
       "booTime": booTime,
       "lastPatchUpdateDate":match,
     },
   }

   client.Publish(event)
   logp.Info("sent hardware data.", event)

  return  nil
}
