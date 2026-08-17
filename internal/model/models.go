package model

// CPUSample stores the raw per-CPU counters from /proc/stat.
type CPUSample struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
}

// CPUInfo contains the current CPU stats and system-level observations.
type CPUInfo struct {
	UsagePercent float64   `json:"usage_percent"`
	CoreUsages   []float64 `json:"core_usages"`
	TempCelsius  float64   `json:"temp_celsius"`
	Frequency    float64   `json:"freq_mhz"`
	LoadAvg1     float64   `json:"load_avg_1"`
	LoadAvg5     float64   `json:"load_avg_5"`
	LoadAvg15    float64   `json:"load_avg_15"`
	NumCores     int       `json:"num_cores"`
	ModelName    string    `json:"model_name"`
	Throttled    bool      `json:"throttled"`
}

// MemInfo describes memory and swap statistics.
type MemInfo struct {
	TotalMB     float64 `json:"total_mb"`
	UsedMB      float64 `json:"used_mb"`
	FreeMB      float64 `json:"free_mb"`
	CachedMB    float64 `json:"cached_mb"`
	BuffersMB   float64 `json:"buffers_mb"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotalMB float64 `json:"swap_total_mb"`
	SwapUsedMB  float64 `json:"swap_used_mb"`
	SwapPercent float64 `json:"swap_percent"`
}

// DiskInfo contains filesystem usage metrics.
type DiskInfo struct {
	Path        string  `json:"path"`
	Device      string  `json:"device"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	FreeGB      float64 `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
	FSType      string  `json:"fs_type"`
}

// NetIface stores network traffic and rates for a single interface.
type NetIface struct {
	Name      string  `json:"name"`
	RxBytes   uint64  `json:"rx_bytes"`
	TxBytes   uint64  `json:"tx_bytes"`
	RxRate    float64 `json:"rx_rate_kbps"`
	TxRate    float64 `json:"tx_rate_kbps"`
	RxPackets uint64  `json:"rx_packets"`
	TxPackets uint64  `json:"tx_packets"`
}

// ProcessInfo contains process resource usage information.
type ProcessInfo struct {
	PID   int     `json:"pid"`
	Name  string  `json:"name"`
	CPU   float64 `json:"cpu_percent"`
	MemMB float64 `json:"mem_mb"`
	State string  `json:"state"`
	User  string  `json:"user"`
}

// GPUInfo contains GPU temperature and memory split information.
type GPUInfo struct {
	TempCelsius float64 `json:"temp_celsius"`
	MemSplitMB  int     `json:"mem_split_mb"`
	Available   bool    `json:"available"`
}

// SystemInfo describes host-level metadata.
type SystemInfo struct {
	Hostname   string `json:"hostname"`
	Uptime     string `json:"uptime"`
	UptimeSecs uint64 `json:"uptime_secs"`
	OS         string `json:"os"`
	Kernel     string `json:"kernel"`
	PiModel    string `json:"pi_model"`
	GoVersion  string `json:"go_version"`
	ServerTime string `json:"server_time"`
}

// Snapshot is the complete latest system data set.
type Snapshot struct {
	Timestamp int64         `json:"ts"`
	CPU       CPUInfo       `json:"cpu"`
	Memory    MemInfo       `json:"memory"`
	Disks     []DiskInfo    `json:"disks"`
	Network   []NetIface    `json:"network"`
	Processes []ProcessInfo `json:"processes"`
	GPU       GPUInfo       `json:"gpu"`
	System    SystemInfo    `json:"system"`
}
