package beater

import (
	"fmt"
	"time"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/Manjukb/enduserbeat/config"
	"github.com/Manjukb/enduserbeat/modules/hardware"
  "github.com/Manjukb/enduserbeat/modules/software"
)

type Enduserbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}


// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Enduserbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Enduserbeat) Run(b *beat.Beat) error {
	logp.Info("enduserbeat is running! Hit CTRL-C to stop it.")
	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		  hardware.PostData(bt.client)
		  software.PostData(bt.client)
		counter++
	}
}

func (bt *Enduserbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
