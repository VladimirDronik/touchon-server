package service

// Регистрируем события
import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	_ "github.com/VladimirDronik/touchon-server/events"
	_ "github.com/VladimirDronik/touchon-server/events/item"
	_ "github.com/VladimirDronik/touchon-server/events/object/port"
	_ "github.com/VladimirDronik/touchon-server/events/object/regulator"
	_ "github.com/VladimirDronik/touchon-server/events/object/sensor"
	_ "github.com/VladimirDronik/touchon-server/events/script"
	"github.com/pbnjay/memory"
)

var Name string

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
	Env            map[string]string
}

func GetInfo() (*Info, error) {
	info := &Info{
		Service:        Name,
		StartedAt:      startedAt.Format("02.01.2006 15:04:05"),
		Uptime:         time.Since(startedAt).Round(time.Second).String(),
		MaxMemoryUsage: fmt.Sprintf("%.1f MiB", float64(maxMem.Load())/1024/1024),
		TotalMemory:    fmt.Sprintf("%.1f MiB", float64(memory.TotalMemory())/1024/1024),
		FreeMemory:     fmt.Sprintf("%.1f MiB", float64(memory.FreeMemory())/1024/1024),
		Env:            Config,
	}

	return info, nil
}
