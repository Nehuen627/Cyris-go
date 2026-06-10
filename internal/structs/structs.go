package structs

type Requirements struct {
	Minimum     string `json:"minimum"`
	Recommended string `json:"recommended"`
}

type GameData struct {
	PCRequirements    Requirements `json:"pc_requirements"`
	MacRequirements   Requirements `json:"mac_requirements"`
	LinuxRequirements Requirements `json:"linux_requirements"`
}

type SteamApp struct {
	AppID string `json:"appid"`
	Name  string `json:"name"`
}

type SystemSpecs struct {
	CPUName   string `json:"cpu_name"`
	CPUCores  int    `json:"cpu_cores"`
	RAMTotal  uint64 `json:"ram_total_mb"`
	DiskFree  uint64 `json:"disk_free_gb"`
	GPUName   string `json:"gpu_name"`
	GPUVendor string `json:"gpu_vendor"`
	OS        string `json:"os"`
}

type GPUResult struct {
	Meets            bool
	MeetsRecommended bool
	Found            bool
	IsIntegrated     bool
	VRAMFail         bool
	UserGPU          string
	RequiredGPU      string
}

type RequirementsResult struct {
	CPUCores bool
	RAMTotal bool
	DiskFree bool
	GPU      GPUResult
	CPU      CPUResult

	MeetsMinimum     bool
	MeetsRecommended bool
}

type CPUResult struct {
	Meets            bool
	MeetsRecommended bool
	Found            bool
	UserCPU          string
	RequiredCPU      string
}

type CPUEntry struct {
	Name    string  `json:"name"`
	CPUMark float64 `json:"cpumark,string"`
}

type GPUEntry struct {
	Name       string  `json:"name"`
	FP32       float64 `json:"fp32"`
	PixelRate  float64 `json:"pixelRate"`
	MemoryGB   float64 `json:"memorySize"`
	Integrated bool    `json:"-"`
}
