// Содержит информацию о сервисе

package info

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/pbnjay/memory"
)

var Config map[string]string

var startedAt = time.Now()

var maxMem atomic.Uint64

func init() {
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			v := m.Sys
			if v > maxMem.Load() {
				maxMem.Store(v)
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

type Info struct {
	Service        string
	StartedAt      string
	Uptime         string
	MaxMemoryUsage string
	TotalMemory    string
	FreeMemory     string
	GOOS           string
	GOARCH         string
	Env            map[string]string
}

func GetInfo() (*Info, error) {
	info := &Info{
		Service:        "touchon_server",
		StartedAt:      startedAt.Format("02.01.2006 15:04:05"),
		Uptime:         time.Since(startedAt).Round(time.Second).String(),
		MaxMemoryUsage: fmt.Sprintf("%.1f MiB", float64(maxMem.Load())/1024/1024),
		TotalMemory:    fmt.Sprintf("%.1f MiB", float64(memory.TotalMemory())/1024/1024),
		FreeMemory:     fmt.Sprintf("%.1f MiB", float64(memory.FreeMemory())/1024/1024),
		GOOS:           runtime.GOOS,
		GOARCH:         runtime.GOARCH,
		Env:            Config,
	}

	return info, nil
}
