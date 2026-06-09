package extras

import (
	"cyris/internal/structs"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var gpuDatabase []structs.GPUEntry
var integratedDatabase []structs.GPUEntry

func LoadGPUDatabase() error {
	urls := []string{
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/nvidia/all.json",
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/amd/all.json",
		"https://raw.githubusercontent.com/RightNow-AI/RightNow-GPU-Database/main/data/intel/all.json",
	}
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var entries []structs.GPUEntry
		if err := json.Unmarshal(body, &entries); err != nil {
			return err
		}
		gpuDatabase = append(gpuDatabase, entries...)
	}
	return nil
}

func LoadIntegratedDatabase() error {
	resp, err := http.Get("https://raw.githubusercontent.com/voidful/gpu-info-api/gpu-data/gpu.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var raw map[string]struct {
		Model  string `json:"Model"`
		TFLOPS string `json:"Single-precision TFLOPS"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return err
	}

	for _, entry := range raw {
		tflops, err := strconv.ParseFloat(entry.TFLOPS, 64)
		if err != nil {
			tflops = 0
		}
		integratedDatabase = append(integratedDatabase, structs.GPUEntry{
			Name: entry.Model,
			FP32: tflops,
		})
	}
	return nil
}

func LookupGPU(name string) (structs.GPUEntry, bool) {
	lower := strings.ToLower(strings.TrimSpace(name))
	if lower == "" || lower == "none" {
		return structs.GPUEntry{}, false
	}
	
	// Remove common trademark symbols for better matching
	lower = strings.ReplaceAll(lower, "(r)", "")
	lower = strings.ReplaceAll(lower, "(tm)", "")
	lower = strings.ReplaceAll(lower, "  ", " ") // remove double spaces
	lower = strings.TrimSpace(lower)

	for _, gpu := range gpuDatabase {
		if strings.Contains(strings.ToLower(gpu.Name), lower) ||
			strings.Contains(lower, strings.ToLower(gpu.Name)) {
			return gpu, true
		}
	}

	for _, gpu := range integratedDatabase {
		if strings.Contains(strings.ToLower(gpu.Name), lower) ||
			strings.Contains(lower, strings.ToLower(gpu.Name)) {
			return gpu, true
		}
	}

	// Fallback for generic integrated GPUs
	if strings.Contains(lower, "intel") && strings.Contains(lower, "uhd graphics") {
		return structs.GPUEntry{Name: "Intel UHD Graphics", FP32: 0.4}, true // ~0.4 TFLOPS average
	}
	if strings.Contains(lower, "intel") && strings.Contains(lower, "iris xe") {
		return structs.GPUEntry{Name: "Intel Iris Xe Graphics", FP32: 1.6}, true // ~1.6 TFLOPS average
	}
	if strings.Contains(lower, "intel") && strings.Contains(lower, "hd graphics") {
		return structs.GPUEntry{Name: "Intel HD Graphics", FP32: 0.2}, true // ~0.2 TFLOPS average
	}
	if strings.Contains(lower, "amd") && strings.Contains(lower, "radeon graphics") {
		return structs.GPUEntry{Name: "AMD Radeon Graphics", FP32: 1.5}, true // ~1.5 TFLOPS average for modern APUs
	}
	if strings.Contains(lower, "apple") && strings.Contains(lower, "m1") {
		return structs.GPUEntry{Name: "Apple M1 GPU", FP32: 2.6}, true
	}
	if strings.Contains(lower, "apple") && strings.Contains(lower, "m2") {
		return structs.GPUEntry{Name: "Apple M2 GPU", FP32: 3.6}, true
	}

	return structs.GPUEntry{}, false
}

func ParseRequirements(requirements string) (uint64, uint64, string) {
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

	return ram, disk, gpu
}

func Hystory() {
	// send back the last 3 options searched
}
