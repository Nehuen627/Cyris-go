package extras

import (
	"cyris/internal/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var gpuDatabase []structs.GPUEntry
var integratedDatabase []structs.GPUEntry
var cpuDatabase []structs.CPUEntry

func loadJSON(url string, dst interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body %s: %w", url, err)
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("unmarshal %s: %w", url, err)
	}
	return nil
}

func LoadGPUDatabase() error {
	urls := []string{
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/nvidia/all.json",
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/amd/all.json",
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/intel/all.json",
	}
	for _, url := range urls {
		var entries []structs.GPUEntry
		if err := loadJSON(url, &entries); err != nil {
			return err
		}
		for i := range entries {
			if entries[i].FP32 == 0 && entries[i].PixelRate > 0 {
				entries[i].FP32 = entries[i].PixelRate
			}
		}
		gpuDatabase = append(gpuDatabase, entries...)
	}
	return nil
}

func LoadIntegratedDatabase() error {
	var raw map[string]struct {
		Model       string  `json:"Model"`
		PixelShader float64 `json:"Pixel Shader (MP/s)"`
		TFLOPS      string  `json:"Single-precision TFLOPS"`
		MemSizeMB   string  `json:"Memory Configuration Size (MB)"`
	}
	if err := loadJSON("https://raw.githubusercontent.com/voidful/gpu-info-api/gpu-data/gpu.json", &raw); err != nil {
		return err
	}

	for _, entry := range raw {
		if entry.Model == "" || strings.ToLower(entry.Model) == "nan" {
			continue
		}

		var fp32 float64

		if entry.TFLOPS != "" && entry.TFLOPS != "?" {
			if v, err := strconv.ParseFloat(strings.TrimSpace(entry.TFLOPS), 64); err == nil && v > 0 {
				fp32 = v
			}
		}

		if fp32 == 0 && entry.PixelShader > 0 {
			fp32 = entry.PixelShader / 1000.0
		}

		var memGB float64
		if entry.MemSizeMB != "" {
			for _, tok := range strings.Fields(entry.MemSizeMB) {
				if v, err := strconv.ParseFloat(tok, 64); err == nil && v > 0 {
					memGB = v / 1024.0
					break
				}
			}
		}

		integratedDatabase = append(integratedDatabase, structs.GPUEntry{
			Name:       entry.Model,
			FP32:       fp32,
			MemoryGB:   memGB,
			Integrated: true,
		})
	}
	return nil
}

func LoadCPUDatabase() error {
	var raw struct {
		Data []struct {
			Name    string `json:"name"`
			CPUMark string `json:"cpumark"`
		} `json:"data"`
	}
	if err := loadJSON("https://raw.githubusercontent.com/spapas/cpu-benchmark/master/cpus2.json", &raw); err != nil {
		return err
	}

	for _, entry := range raw.Data {
		cleanMark := strings.ReplaceAll(entry.CPUMark, ",", "")
		mark, err := strconv.ParseFloat(cleanMark, 64)
		if err != nil {
			mark = 0
		}
		cpuDatabase = append(cpuDatabase, structs.CPUEntry{
			Name:    entry.Name,
			CPUMark: mark,
		})
	}
	return nil
}

func normalizeQuery(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	lower = strings.ReplaceAll(lower, "(r)", "")
	lower = strings.ReplaceAll(lower, "(tm)", "")
	lower = strings.ReplaceAll(lower, "  ", " ")
	return strings.TrimSpace(lower)
}

func LookupCPU(name string) (structs.CPUEntry, bool) {
	lower := normalizeQuery(name)
	if lower == "" || lower == "none" {
		return structs.CPUEntry{}, false
	}

	var best structs.CPUEntry
	bestLen := 0
	found := false

	for _, cpu := range cpuDatabase {
		dbName := strings.ToLower(cpu.Name)
		if strings.Contains(lower, dbName) || strings.Contains(dbName, lower) {
			if len(dbName) > bestLen {
				best = cpu
				bestLen = len(dbName)
				found = true
			}
		}
	}

	return best, found
}

func LookupGPU(name string) (structs.GPUEntry, bool) {
	lower := normalizeQuery(name)
	if lower == "" || lower == "none" {
		return structs.GPUEntry{}, false
	}

	if entry, ok := searchAllGPUDbs(lower); ok {
		return entry, true
	}

	if family := extractBracketName(lower); family != "" && family != lower {
		if entry, ok := searchAllGPUDbs(family); ok {
			return entry, true
		}
	}

	return broadFamilyFallback(lower)
}

func searchAllGPUDbs(query string) (structs.GPUEntry, bool) {
	var best structs.GPUEntry
	bestLen := 0
	found := false

	for _, db := range [][]structs.GPUEntry{gpuDatabase, integratedDatabase} {
		for _, gpu := range db {
			dbName := strings.ToLower(gpu.Name)
			if strings.Contains(query, dbName) || strings.Contains(dbName, query) {
				if len(dbName) > bestLen {
					best = gpu
					bestLen = len(dbName)
					found = true
				}
			}
		}
	}
	return best, found
}

func extractBracketName(s string) string {
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start != -1 && end > start+1 {
		return strings.TrimSpace(s[start+1 : end])
	}
	return ""
}

func IsIntegratedGPU(name string) bool {
	lower := normalizeQuery(name)
	if lower == "" || lower == "none" {
		return false
	}

	for _, gpu := range integratedDatabase {
		dbName := strings.ToLower(gpu.Name)
		if strings.Contains(lower, dbName) || strings.Contains(dbName, lower) {
			return true
		}
	}

	if family := extractBracketName(lower); family != "" && family != lower {
		for _, gpu := range integratedDatabase {
			dbName := strings.ToLower(gpu.Name)
			if strings.Contains(family, dbName) || strings.Contains(dbName, family) {
				return true
			}
		}
	}

	integratedPatterns := []string{
		"iris xe", "uhd graphics", "hd graphics", "radeon graphics",
		"iris plus", "iris pro", "radeon vega",
	}
	for _, pattern := range integratedPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

func broadFamilyFallback(lower string) (structs.GPUEntry, bool) {
	switch {
	case strings.Contains(lower, "iris xe"):
		return structs.GPUEntry{Name: "Intel Iris Xe Graphics", FP32: 1.6, Integrated: true}, true
	case strings.Contains(lower, "uhd graphics"):
		return structs.GPUEntry{Name: "Intel UHD Graphics", FP32: 0.4, Integrated: true}, true
	case strings.Contains(lower, "hd graphics"):
		return structs.GPUEntry{Name: "Intel HD Graphics", FP32: 0.2, Integrated: true}, true
	case strings.Contains(lower, "radeon graphics"):
		return structs.GPUEntry{Name: "AMD Radeon Graphics", FP32: 1.5, Integrated: true}, true
	}
	return structs.GPUEntry{}, false
}

func ParseRequirements(requirements string) (uint64, uint64, string, string) {
	lower := strings.ToLower(requirements)

	var ram uint64
	ramIdx := strings.Index(lower, "memory:</strong>")
	if ramIdx != -1 {
		after := strings.TrimSpace(lower[ramIdx+len("memory:</strong>"):])
		parts := strings.Fields(after)
		if len(parts) > 0 {
			num, err := strconv.ParseUint(parts[0], 10, 64)
			if err == nil {
				ram = num
				if len(parts) > 1 && strings.ToLower(parts[1]) == "mb" {
					ram = ram / 1024
				}
			}
		}
	}

	var disk uint64
	diskIdx := strings.Index(lower, "storage:</strong>")
	if diskIdx != -1 {
		after := strings.TrimSpace(lower[diskIdx+len("storage:</strong>"):])
		parts := strings.Fields(after)
		if len(parts) > 0 {
			num, err := strconv.ParseUint(parts[0], 10, 64)
			if err == nil {
				disk = num
				if len(parts) > 1 && strings.ToLower(parts[1]) == "mb" {
					disk = disk / 1024
				}
			}
		}
	}

	var gpu string
	gpuIdx := strings.Index(lower, "graphics:</strong>")
	if gpuIdx != -1 {
		after := strings.TrimSpace(requirements[gpuIdx+len("graphics:</strong>"):])
		endIdx := strings.Index(after, "<")
		if endIdx != -1 {
			gpu = strings.TrimSpace(after[:endIdx])
		} else {
			gpu = strings.TrimSpace(after)
		}
	}

	var cpu string
	cpuIdx := strings.Index(lower, "processor:</strong>")
	if cpuIdx != -1 {
		after := strings.TrimSpace(requirements[cpuIdx+len("processor:</strong>"):])
		endIdx := strings.Index(after, "<")
		if endIdx != -1 {
			cpu = strings.TrimSpace(after[:endIdx])
		} else {
			cpu = strings.TrimSpace(after)
		}
	}

	return ram, disk, gpu, cpu
}
