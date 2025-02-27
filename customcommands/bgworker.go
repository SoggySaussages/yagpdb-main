package customcommands

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/botlabs-gg/sgpdb/v2/common/backgroundworkers"
)

var _ backgroundworkers.BackgroundWorkerPlugin = (*Plugin)(nil)

// RunBackgroundWorker implements backgroundworkers.BackgroundWorkerPlugin
func (p *Plugin) RunBackgroundWorker() {
	t := time.NewTicker(time.Second * 30)
	for range t.C {
		go p.pullDirectories()
	}
}

// StopBackgroundWorker implements backgroundworkers.BackgroundWorkerPlugin
func (p *Plugin) StopBackgroundWorker(wg *sync.WaitGroup) {
	wg.Done()
}

func (p *Plugin) pullDirectories() {
	dir, err := os.Open("cc-github")
	if err != nil {
		return
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return
	}

	for _, f := range entries {
		cmd := exec.Command("git", "pull")
		cmd.Dir = "cc-github/" + f.Name()
		go runCmdLogErr(cmd)
	}
	logger.Info("DONE pulling any recent CC GitHub changes")
}
