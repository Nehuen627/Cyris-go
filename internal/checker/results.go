package checker

import (
	"cyris/internal/extras"
	"cyris/internal/structs"
)

func CheckRequirements(systemSpecs structs.SystemSpecs, gameData structs.GameData) structs.RequirementsResult {
	var MeetsMinimumRes bool = true
	var MeetsRecommendedRes bool = true
	var CPUCoresRes bool = true
	var RAMTotalRes bool = true
	var DiskTotalRes bool = true
	var GPURes structs.GPUResult

	minRAM, minDisk, minGPU := extras.ParseRequirements(gameData.PCRequirements.Minimum)
	recRAM, recDisk, recGPU := extras.ParseRequirements(gameData.PCRequirements.Recommended)

	GPURes.UserGPU = systemSpecs.GPUName
	GPURes.RequiredGPU = minGPU

	if systemSpecs.RAMTotal < recRAM {
		MeetsRecommendedRes = false
	}
	if systemSpecs.RAMTotal < minRAM {
		RAMTotalRes = false
		MeetsRecommendedRes = false
	}

	if systemSpecs.DiskTotal < recDisk {
		MeetsRecommendedRes = false
	}
	if systemSpecs.DiskTotal < minDisk {
		DiskTotalRes = false
		MeetsRecommendedRes = false
	}

	userGPU, found := extras.LookupGPU(systemSpecs.GPUName)
	requiredGPU, reqFound := extras.LookupGPU(minGPU)
	requiredRecGPU, recGPUFound := extras.LookupGPU(recGPU)

	GPURes.Found = found
	if found && reqFound {
		GPURes.Meets = userGPU.FP32 >= requiredGPU.FP32
		if !GPURes.Meets {
			MeetsMinimumRes = false
		}
	} else if found && !reqFound {
		GPURes.Meets = true // Requirement unknown, assume passes
	} else if !found && reqFound {
		GPURes.Meets = false // User GPU unknown, but requirement known -> Fail to prevent false positives
		MeetsMinimumRes = false
	} else {
		GPURes.Meets = true // Both unknown -> Assume passes
	}

	if found && recGPUFound {
		GPURes.MeetsRecommended = userGPU.FP32 >= requiredRecGPU.FP32
		if !GPURes.MeetsRecommended {
			MeetsRecommendedRes = false
		}
	} else if found && !recGPUFound {
		GPURes.MeetsRecommended = true
	} else if !found && recGPUFound {
		GPURes.MeetsRecommended = false
		MeetsRecommendedRes = false
	} else {
		GPURes.MeetsRecommended = true
	}

	if !RAMTotalRes || !CPUCoresRes || !DiskTotalRes {
		MeetsMinimumRes = false
	}

	return structs.RequirementsResult{
		CPUCores:         CPUCoresRes,
		RAMTotal:         RAMTotalRes,
		DiskTotal:        DiskTotalRes,
		GPU:              GPURes,
		MeetsMinimum:     MeetsMinimumRes,
		MeetsRecommended: MeetsRecommendedRes,
	}
}
