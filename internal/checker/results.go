package checker

import (
	"cyris/internal/extras"
	"cyris/internal/structs"
	"strings"
)

func CheckRequirements(systemSpecs structs.SystemSpecs, gameData structs.GameData) structs.RequirementsResult {
	var MeetsMinimumRes bool = true
	var MeetsRecommendedRes bool = true
	var CPUCoresRes bool = true
	var RAMTotalRes bool = true
	var DiskFreeRes bool = true
	var GPURes structs.GPUResult
	var CPURes structs.CPUResult

	minRAM, minDisk, minGPU, minCPU := extras.ParseRequirements(gameData.PCRequirements.Minimum)
	recRAM, recDisk, recGPU, recCPU := extras.ParseRequirements(gameData.PCRequirements.Recommended)

	systemRAMGB := systemSpecs.RAMTotal / 1024

	if systemRAMGB < recRAM {
		MeetsRecommendedRes = false
	}
	if systemRAMGB < minRAM {
		RAMTotalRes = false
		MeetsRecommendedRes = false
	}

	if systemSpecs.DiskFree < recDisk {
		MeetsRecommendedRes = false
	}
	if systemSpecs.DiskFree < minDisk {
		DiskFreeRes = false
		MeetsRecommendedRes = false
	}

	userGPULabel := strings.TrimSpace(systemSpecs.GPUVendor + " " + systemSpecs.GPUName)
	userGPU, foundGPU := extras.LookupGPU(userGPULabel)
	requiredGPU, reqGPUFound := extras.LookupGPU(minGPU)
	requiredRecGPU, recGPUFound := extras.LookupGPU(recGPU)

	GPURes.Found = foundGPU

	isIntegrated := userGPU.Integrated || extras.IsIntegratedGPU(userGPULabel)
	GPURes.IsIntegrated = isIntegrated

	vramThreshold := 1.0

	vramFailMin := isIntegrated && reqGPUFound && requiredGPU.MemoryGB >= vramThreshold
	vramFailRec := isIntegrated && recGPUFound && requiredRecGPU.MemoryGB >= vramThreshold

	if vramFailMin {
		GPURes.Meets = false
		GPURes.VRAMFail = true
		MeetsMinimumRes = false
	} else if foundGPU && reqGPUFound {
		GPURes.Meets = userGPU.FP32 >= requiredGPU.FP32
		if !GPURes.Meets {
			MeetsMinimumRes = false
		}
	} else if foundGPU && !reqGPUFound {
		GPURes.Meets = true
	} else if !foundGPU && reqGPUFound {
		GPURes.Meets = false
		MeetsMinimumRes = false
	} else {
		GPURes.Meets = true
	}

	if vramFailRec {
		GPURes.MeetsRecommended = false
		MeetsRecommendedRes = false
	} else if foundGPU && recGPUFound {
		GPURes.MeetsRecommended = userGPU.FP32 >= requiredRecGPU.FP32
		if !GPURes.MeetsRecommended {
			MeetsRecommendedRes = false
		}
	} else if foundGPU && !recGPUFound {
		GPURes.MeetsRecommended = true
	} else if !foundGPU && recGPUFound {
		GPURes.MeetsRecommended = false
		MeetsRecommendedRes = false
	} else {
		GPURes.MeetsRecommended = true
	}

	userCPU, foundCPU := extras.LookupCPU(systemSpecs.CPUName)
	requiredCPU, reqCPUFound := extras.LookupCPU(minCPU)
	requiredRecCPU, recCPUFound := extras.LookupCPU(recCPU)

	CPURes.Found = foundCPU
	if foundCPU && reqCPUFound {
		CPURes.Meets = userCPU.CPUMark >= requiredCPU.CPUMark
		if !CPURes.Meets {
			MeetsMinimumRes = false
		}
	} else if foundCPU && !reqCPUFound {
		CPURes.Meets = true
	} else if !foundCPU && reqCPUFound {
		CPURes.Meets = false
		MeetsMinimumRes = false
	} else {
		CPURes.Meets = true
	}

	if foundCPU && recCPUFound {
		CPURes.MeetsRecommended = userCPU.CPUMark >= requiredRecCPU.CPUMark
		if !CPURes.MeetsRecommended {
			MeetsRecommendedRes = false
		}
	} else if foundCPU && !recCPUFound {
		CPURes.MeetsRecommended = true
	} else if !foundCPU && recCPUFound {
		CPURes.MeetsRecommended = false
		MeetsRecommendedRes = false
	} else {
		CPURes.MeetsRecommended = true
	}

	if !RAMTotalRes || !CPUCoresRes || !DiskFreeRes {
		MeetsMinimumRes = false
	}

	if !GPURes.Meets {
		GPURes.MeetsRecommended = false
	}
	if !CPURes.Meets {
		CPURes.MeetsRecommended = false
	}

	if !MeetsMinimumRes {
		MeetsRecommendedRes = false
	}

	return structs.RequirementsResult{
		CPUCores:         CPUCoresRes,
		RAMTotal:         RAMTotalRes,
		DiskFree:         DiskFreeRes,
		GPU:              GPURes,
		CPU:              CPURes,
		MeetsMinimum:     MeetsMinimumRes,
		MeetsRecommended: MeetsRecommendedRes,
	}
}
