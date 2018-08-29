
package software

import (
  "time"
  "fmt"
  "github.com/elastic/beats/libbeat/common"
  "github.com/elastic/beats/libbeat/beat"
  "github.com/shirou/gopsutil/host"
  "github.com/elastic/beats/libbeat/logp"
  "context"
  wapi "github.com/Microland/go-win64api"
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

func PostData(client beat.Client) ( error) {
    InfoStat, _ := host.InfoWithContext(context.Background())
		 SoftwareInfo, err := wapi.InstalledSoftwareList()
			 if err != nil {
					 fmt.Printf("%s\r\n", err.Error())
			 }

			 for _, s := range SoftwareInfo {


         event := beat.Event{

           Timestamp: time.Now(),
           Fields:common.MapStr{
             "Hostname":    InfoStat.Hostname,
             "type":       "smartcenterEndUserSoftware",
             "name": s.Name(),
             "Architecture": s.Architecture(),
             "Version": s.Version(),
             "Publisher": s.Publisher(),
             "InstallDate": s.InstallDate(),
             "EstimatedSize": s.EstimatedSize(),
             "Contact": s.Contact(),
             "HelpLink": s.HelpLink(),
             "InstallSource": s.InstallSource(),
             "VersionMajor": s.VersionMajor(),
             "VersionMinor": s.VersionMinor(),
           },
         }

         client.Publish(event)
         logp.Info("sent software data.", event)

			 }

  return  nil
}
