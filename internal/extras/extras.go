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

// loadJSON fetches a URL and unmarshals the JSON body into dst.
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

// LoadGPUDatabase loads discrete GPU data from the RightNow-AI GPU database.
// The RightNow JSON entries have no "fp32" field; we use "pixelRate" (Gpixels/s)
// as a proportional performance proxy and store it in FP32 so the comparisons work.
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
			// If fp32 is not set (RightNow DB doesn't have it), fall back to pixelRate.
			if entries[i].FP32 == 0 && entries[i].PixelRate > 0 {
				entries[i].FP32 = entries[i].PixelRate
			}
		}
		gpuDatabase = append(gpuDatabase, entries...)
	}
	return nil
}

// LoadIntegratedDatabase loads integrated GPU data from voidful/gpu-info-api.
// That DB stores performance as "Pixel Shader (MP/s)" — we convert to a rough
// TFLOPS equivalent by dividing by 1000 (MP/s → GFLOPS ≈ TFLOPS at low scale).
// Entries with Model == "nan" (Python NaN) are skipped.
func LoadIntegratedDatabase() error {
	// The voidful DB is a map of key → entry object.
	var raw map[string]struct {
		Model       string  `json:"Model"`
		PixelShader float64 `json:"Pixel Shader (MP/s)"`
		// Some modern entries do have this field directly.
		TFLOPS string `json:"Single-precision TFLOPS"`
	}
	if err := loadJSON("https://raw.githubusercontent.com/voidful/gpu-info-api/gpu-data/gpu.json", &raw); err != nil {
		return err
	}

	for _, entry := range raw {
		// Skip bogus entries where the model was serialized as Python NaN.
		if entry.Model == "" || strings.ToLower(entry.Model) == "nan" {
			continue
		}

		var fp32 float64

		// Prefer explicit TFLOPS field if present and parseable.
		if entry.TFLOPS != "" && entry.TFLOPS != "?" {
			if v, err := strconv.ParseFloat(strings.TrimSpace(entry.TFLOPS), 64); err == nil && v > 0 {
				fp32 = v
			}
		}

		// Fall back to Pixel Shader rate converted to approximate TFLOPS.
		// 1 MP/s pixel shading ≈ several GFLOPS; dividing by 1000 gives a rough TFLOPS proxy.
		if fp32 == 0 && entry.PixelShader > 0 {
			fp32 = entry.PixelShader / 1000.0
		}

		integratedDatabase = append(integratedDatabase, structs.GPUEntry{
			Name: entry.Model,
			FP32: fp32,
		})
	}
	return nil
}

// LoadCPUDatabase loads CPU benchmark data from spapas/cpu-benchmark.
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

// normalizeQuery cleans up common trademark noise from a GPU/CPU name string.
func normalizeQuery(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	lower = strings.ReplaceAll(lower, "(r)", "")
	lower = strings.ReplaceAll(lower, "(tm)", "")
	lower = strings.ReplaceAll(lower, "  ", " ")
	return strings.TrimSpace(lower)
}

// LookupCPU finds a CPU entry by name using fuzzy substring matching.
// Prefers the longest DB name that matches to avoid false positives on short names.
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

// LookupGPU finds a GPU entry by name. It searches discrete GPUs first, then
// integrated, then tries to extract a GPU family name from bracket notation
// (e.g. "Raptor Lake-P [UHD Graphics]" → search for "UHD Graphics"), and only
// falls back to broad family estimates as a last resort.
// Matching always prefers the longest (most specific) DB name.
func LookupGPU(name string) (structs.GPUEntry, bool) {
	lower := normalizeQuery(name)
	if lower == "" || lower == "none" {
		return structs.GPUEntry{}, false
	}

	if entry, ok := searchAllGPUDbs(lower); ok {
		return entry, true
	}

	// ghw often returns codename-based names like "Intel Raptor Lake-P [UHD Graphics]"
	// where the real GPU family is inside the brackets. Extract it and retry.
	if family := extractBracketName(lower); family != "" && family != lower {
		if entry, ok := searchAllGPUDbs(family); ok {
			return entry, true
		}
	}

	// Last resort: broad-family estimates. These are intentionally generic —
	// just enough to avoid treating an unknown integrated GPU as having 0 FP32.
	return broadFamilyFallback(lower)
}

// searchAllGPUDbs searches both the discrete and integrated GPU databases for
// the given (already-normalized) query string. Returns the longest match found.
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

// extractBracketName returns the text inside the first pair of square brackets,
// e.g. "raptor lake-p [uhd graphics]" → "uhd graphics". Returns "" if none.
func extractBracketName(s string) string {
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start != -1 && end > start+1 {
		return strings.TrimSpace(s[start+1 : end])
	}
	return ""
}

// broadFamilyFallback is the genuine last resort — only reached when neither the
// databases nor bracket extraction produced a result. Values are intentionally
// rough averages per GPU family, not per model.
func broadFamilyFallback(lower string) (structs.GPUEntry, bool) {
	switch {
	case strings.Contains(lower, "iris xe"):
		return structs.GPUEntry{Name: "Intel Iris Xe Graphics", FP32: 1.6}, true
	case strings.Contains(lower, "uhd graphics"):
		return structs.GPUEntry{Name: "Intel UHD Graphics", FP32: 0.4}, true
	case strings.Contains(lower, "hd graphics"):
		return structs.GPUEntry{Name: "Intel HD Graphics", FP32: 0.2}, true
	case strings.Contains(lower, "radeon graphics"):
		return structs.GPUEntry{Name: "AMD Radeon Graphics", FP32: 1.5}, true
	}
	return structs.GPUEntry{}, false
}

// ParseRequirements parses a Steam PC requirements HTML blob and extracts
// RAM (in GB), Disk (in GB), GPU name, and CPU name.
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
				// Convert MB to GB if the next token says "mb"
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
				// Convert MB to GB if the next token says "mb"
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

func Hystory() {
	// send back the last 3 options searched
}
