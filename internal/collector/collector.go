package collector

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"mon-pi/internal/model"
)

// Collector gathers system metrics and keeps a small in-memory history.
type Collector struct {
	mu          sync.RWMutex
	latest      *model.Snapshot
	history     []float64
	memHistory  []float64
	prevCPU     map[string]model.CPUSample
	prevNet     map[string][2]uint64
	prevNetTime time.Time
	prevCPUTime time.Time
}

func New() *Collector {
	c := &Collector{
		prevCPU:     make(map[string]model.CPUSample),
		prevNet:     make(map[string][2]uint64),
		prevNetTime: time.Now(),
		prevCPUTime: time.Now(),
	}
	c.readCPUSamples()
	time.Sleep(200 * time.Millisecond)
	return c
}

func (c *Collector) Run() {
	ticker := time.NewTicker(2 * time.Second)
	c.collect()
	for range ticker.C {
		c.collect()
	}
}

func (c *Collector) Latest() *model.Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.latest
}

func (c *Collector) CPUHistory() []float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]float64, len(c.history))
	copy(out, c.history)
	return out
}

func (c *Collector) MemHistory() []float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]float64, len(c.memHistory))
	copy(out, c.memHistory)
	return out
}

func (c *Collector) collect() {
	snap := &model.Snapshot{Timestamp: time.Now().UnixMilli()}
	snap.CPU = c.collectCPU()
	snap.Memory = collectMemory()
	snap.Disks = collectDisks()
	snap.Network = c.collectNetwork()
	snap.Processes = collectProcesses()
	snap.GPU = collectGPU()
	snap.System = collectSystem()

	c.mu.Lock()
	c.latest = snap
	c.history = appendCapped(c.history, snap.CPU.UsagePercent, 60)
	c.memHistory = appendCapped(c.memHistory, snap.Memory.UsedPercent, 60)
	c.mu.Unlock()
}

func appendCapped(s []float64, v float64, max int) []float64 {
	s = append(s, v)
	if len(s) > max {
		s = s[len(s)-max:]
	}
	return s
}

func (c *Collector) readCPUSamples() map[string]model.CPUSample {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil
	}
	defer f.Close()
	samples := make(map[string]model.CPUSample)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}
		key := fields[0]
		var s model.CPUSample
		s.User, _ = strconv.ParseUint(fields[1], 10, 64)
		s.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
		s.System, _ = strconv.ParseUint(fields[3], 10, 64)
		s.Idle, _ = strconv.ParseUint(fields[4], 10, 64)
		s.IOWait, _ = strconv.ParseUint(fields[5], 10, 64)
		s.IRQ, _ = strconv.ParseUint(fields[6], 10, 64)
		s.SoftIRQ, _ = strconv.ParseUint(fields[7], 10, 64)
		samples[key] = s
		c.prevCPU[key] = s
	}
	return samples
}

func cpuPercent(prev, cur model.CPUSample) float64 {
	prevIdle := prev.Idle + prev.IOWait
	curIdle := cur.Idle + cur.IOWait
	prevTotal := prev.User + prev.Nice + prev.System + prev.Idle + prev.IOWait + prev.IRQ + prev.SoftIRQ
	curTotal := cur.User + cur.Nice + cur.System + cur.Idle + cur.IOWait + cur.IRQ + cur.SoftIRQ
	totalDelta := curTotal - prevTotal
	idleDelta := curIdle - prevIdle
	if totalDelta == 0 {
		return 0
	}
	return math.Round((float64(totalDelta-idleDelta)/float64(totalDelta))*1000) / 10
}

func (c *Collector) collectCPU() model.CPUInfo {
	cur := c.readCPUSamples()
	info := model.CPUInfo{NumCores: runtime.NumCPU()}

	if prev, ok := c.prevCPU["cpu"]; ok {
		if cu, ok2 := cur["cpu"]; ok2 {
			info.UsagePercent = cpuPercent(prev, cu)
		}
	}

	for i := 0; i < info.NumCores; i++ {
		key := fmt.Sprintf("cpu%d", i)
		if prev, ok := c.prevCPU[key]; ok {
			if cu, ok2 := cur[key]; ok2 {
				info.CoreUsages = append(info.CoreUsages, cpuPercent(prev, cu))
			}
		}
	}
	for k, v := range cur {
		c.prevCPU[k] = v
	}

	info.TempCelsius = readCPUTemp()
	info.Frequency = readCPUFreq()

	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			info.LoadAvg1, _ = strconv.ParseFloat(fields[0], 64)
			info.LoadAvg5, _ = strconv.ParseFloat(fields[1], 64)
			info.LoadAvg15, _ = strconv.ParseFloat(fields[2], 64)
		}
	}

	info.ModelName = readCPUModel()
	info.Throttled = checkThrottled()
	return info
}

func readCPUTemp() float64 {
	paths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/devices/virtual/thermal/thermal_zone0/temp",
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			if err == nil {
				return v / 1000.0
			}
		}
	}
	return 0
}

func readCPUFreq() float64 {
	if out, err := exec.Command("vcgencmd", "measure_clock", "arm").Output(); err == nil {
		parts := strings.Split(strings.TrimSpace(string(out)), "=")
		if len(parts) == 2 {
			v, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				return v / 1e6
			}
		}
	}
	if data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq"); err == nil {
		v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err == nil {
			return v / 1000.0
		}
	}
	return 0
}

func readCPUModel() string {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return runtime.GOARCH
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "Model name") || strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "Hardware") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return runtime.GOARCH
}

func checkThrottled() bool {
	out, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return false
	}
	parts := strings.Split(strings.TrimSpace(string(out)), "=")
	if len(parts) == 2 {
		return parts[1] != "0x0"
	}
	return false
}

func collectMemory() model.MemInfo {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return model.MemInfo{}
	}
	defer f.Close()

	vals := make(map[string]uint64)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			key := strings.TrimSuffix(parts[0], ":")
			v, _ := strconv.ParseUint(parts[1], 10, 64)
			vals[key] = v
		}
	}

	total := float64(vals["MemTotal"]) / 1024
	free := float64(vals["MemFree"]) / 1024
	buffers := float64(vals["Buffers"]) / 1024
	cached := float64(vals["Cached"]+vals["SReclaimable"]) / 1024
	available := float64(vals["MemAvailable"]) / 1024
	used := total - available

	var usedPct float64
	if total > 0 {
		usedPct = math.Round(used/total*1000) / 10
	}

	swapTotal := float64(vals["SwapTotal"]) / 1024
	swapFree := float64(vals["SwapFree"]) / 1024
	swapUsed := swapTotal - swapFree
	var swapPct float64
	if swapTotal > 0 {
		swapPct = math.Round(swapUsed/swapTotal*1000) / 10
	}

	return model.MemInfo{
		TotalMB:     math.Round(total*10) / 10,
		UsedMB:      math.Round(used*10) / 10,
		FreeMB:      math.Round(free*10) / 10,
		CachedMB:    math.Round(cached*10) / 10,
		BuffersMB:   math.Round(buffers*10) / 10,
		UsedPercent: usedPct,
		SwapTotalMB: math.Round(swapTotal*10) / 10,
		SwapUsedMB:  math.Round(swapUsed*10) / 10,
		SwapPercent: swapPct,
	}
}

func collectDisks() []model.DiskInfo {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil
	}
	defer f.Close()

	seen := make(map[string]bool)
	var disks []model.DiskInfo
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		device, mount, fstype := fields[0], fields[1], fields[2]
		skip := []string{"tmpfs", "devtmpfs", "sysfs", "proc", "devpts", "cgroup",
			"pstore", "debugfs", "securityfs", "fusectl", "hugetlbfs",
			"mqueue", "tracefs", "configfs", "bpf", "autofs", "ramfs"}
		skip2 := false
		for _, s := range skip {
			if fstype == s {
				skip2 = true
				break
			}
		}
		if skip2 || seen[mount] {
			continue
		}
		if !strings.HasPrefix(device, "/dev/") {
			continue
		}
		seen[mount] = true

		di := diskUsage(mount)
		di.Device = filepath.Base(device)
		di.FSType = fstype
		if di.TotalGB*1024 < 100 {
			continue
		}
		disks = append(disks, di)
	}
	sort.Slice(disks, func(i, j int) bool { return disks[i].Path < disks[j].Path })
	return disks
}

func diskUsage(path string) model.DiskInfo {
	out, err := exec.Command("df", "-B1", path).Output()
	if err != nil {
		return model.DiskInfo{Path: path}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return model.DiskInfo{Path: path}
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return model.DiskInfo{Path: path}
	}
	total, _ := strconv.ParseFloat(fields[1], 64)
	used, _ := strconv.ParseFloat(fields[2], 64)
	free, _ := strconv.ParseFloat(fields[3], 64)
	var pct float64
	if total > 0 {
		pct = math.Round(used/total*1000) / 10
	}
	gb := 1024.0 * 1024 * 1024
	return model.DiskInfo{
		Path:        path,
		TotalGB:     math.Round(total/gb*100) / 100,
		UsedGB:      math.Round(used/gb*100) / 100,
		FreeGB:      math.Round(free/gb*100) / 100,
		UsedPercent: pct,
	}
}

func (c *Collector) collectNetwork() []model.NetIface {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil
	}
	defer f.Close()

	now := time.Now()
	elapsed := now.Sub(c.prevNetTime).Seconds()
	if elapsed < 0.1 {
		elapsed = 2
	}
	c.prevNetTime = now

	var ifaces []model.NetIface
	sc := bufio.NewScanner(f)
	sc.Scan()
	sc.Scan()
	for sc.Scan() {
		line := sc.Text()
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		name := strings.TrimSpace(line[:colonIdx])
		if name == "lo" {
			continue
		}
		fields := strings.Fields(line[colonIdx+1:])
		if len(fields) < 9 {
			continue
		}
		rxBytes, _ := strconv.ParseUint(fields[0], 10, 64)
		rxPkts, _ := strconv.ParseUint(fields[1], 10, 64)
		txBytes, _ := strconv.ParseUint(fields[8], 10, 64)
		txPkts, _ := strconv.ParseUint(fields[9], 10, 64)

		iface := model.NetIface{
			Name:      name,
			RxBytes:   rxBytes,
			TxBytes:   txBytes,
			RxPackets: rxPkts,
			TxPackets: txPkts,
		}

		if prev, ok := c.prevNet[name]; ok {
			rxDelta := float64(rxBytes - prev[0])
			txDelta := float64(txBytes - prev[1])
			if rxDelta < 0 {
				rxDelta = 0
			}
			if txDelta < 0 {
				txDelta = 0
			}
			iface.RxRate = math.Round(rxDelta/elapsed/1024*10) / 10
			iface.TxRate = math.Round(txDelta/elapsed/1024*10) / 10
		}
		c.prevNet[name] = [2]uint64{rxBytes, txBytes}
		ifaces = append(ifaces, iface)
	}
	return ifaces
}

func collectProcesses() []model.ProcessInfo {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	var procs []model.ProcessInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid <= 0 {
			continue
		}
		p := readProc(pid)
		if p != nil {
			procs = append(procs, *p)
		}
	}

	sort.Slice(procs, func(i, j int) bool {
		if procs[i].CPU != procs[j].CPU {
			return procs[i].CPU > procs[j].CPU
		}
		return procs[i].MemMB > procs[j].MemMB
	})

	if len(procs) > 15 {
		procs = procs[:15]
	}
	return procs
}

func readProc(pid int) *model.ProcessInfo {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(statPath)
	if err != nil {
		return nil
	}
	fields := strings.Fields(string(data))
	if len(fields) < 24 {
		return nil
	}

	raw := string(data)
	start := strings.Index(raw, "(")
	end := strings.LastIndex(raw, ")")
	name := ""
	if start >= 0 && end > start {
		name = raw[start+1 : end]
	}

	state := fields[2]
	utime, _ := strconv.ParseFloat(fields[13], 64)
	stime, _ := strconv.ParseFloat(fields[14], 64)
	clkTck := 100.0
	uptimeData, _ := os.ReadFile("/proc/uptime")
	var uptime float64
	if len(uptimeData) > 0 {
		fmt.Sscanf(string(uptimeData), "%f", &uptime)
	}
	starttime, _ := strconv.ParseFloat(fields[21], 64)
	elapsed := uptime - starttime/clkTck
	var cpuPct float64
	if elapsed > 0 {
		cpuPct = math.Round((utime+stime)/clkTck/elapsed*1000) / 10
	}

	var memMB float64
	statusData, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err == nil {
		sc := bufio.NewScanner(strings.NewReader(string(statusData)))
		for sc.Scan() {
			line := sc.Text()
			if strings.HasPrefix(line, "VmRSS:") {
				f := strings.Fields(line)
				if len(f) >= 2 {
					v, _ := strconv.ParseFloat(f[1], 64)
					memMB = math.Round(v/1024*10) / 10
				}
				break
			}
		}
	}

	user := readProcUser(pid)

	return &model.ProcessInfo{
		PID:   pid,
		Name:  name,
		CPU:   cpuPct,
		MemMB: memMB,
		State: state,
		User:  user,
	}
}

func readProcUser(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/loginuid", pid))
	if err != nil {
		return ""
	}
	uid := strings.TrimSpace(string(data))
	if uid == "4294967295" {
		return "root"
	}
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return uid
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), ":")
		if len(parts) >= 3 && parts[2] == uid {
			return parts[0]
		}
	}
	return uid
}

func collectGPU() model.GPUInfo {
	info := model.GPUInfo{}
	out, err := exec.Command("vcgencmd", "measure_temp").Output()
	if err != nil {
		return info
	}
	info.Available = true
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "temp=")
	s = strings.TrimSuffix(s, "'C")
	s = strings.TrimSuffix(s, "°C")
	info.TempCelsius, _ = strconv.ParseFloat(s, 64)

	out2, err2 := exec.Command("vcgencmd", "get_mem", "gpu").Output()
	if err2 == nil {
		parts := strings.Split(strings.TrimSpace(string(out2)), "=")
		if len(parts) == 2 {
			v := strings.TrimSuffix(parts[1], "M")
			info.MemSplitMB, _ = strconv.Atoi(v)
		}
	}
	return info
}

func collectSystem() model.SystemInfo {
	info := model.SystemInfo{
		OS:         runtime.GOOS,
		GoVersion:  runtime.Version(),
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	info.Hostname, _ = os.Hostname()

	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		var secs float64
		fmt.Sscanf(string(data), "%f", &secs)
		info.UptimeSecs = uint64(secs)
		info.Uptime = formatUptime(uint64(secs))
	}

	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(out))
	}

	if data, err := os.ReadFile("/proc/device-tree/model"); err == nil {
		info.PiModel = strings.TrimRight(string(data), "\x00\n")
	} else if data, err := os.ReadFile("/sys/firmware/devicetree/base/model"); err == nil {
		info.PiModel = strings.TrimRight(string(data), "\x00\n")
	}
	if info.PiModel == "" {
		info.PiModel = "Generic Linux (" + runtime.GOARCH + ")"
	}

	return info
}

func formatUptime(secs uint64) string {
	days := secs / 86400
	hours := (secs % 86400) / 3600
	mins := (secs % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm %ds", mins, secs%60)
}

func Log(msg string) {
	log.Println(msg)
}
