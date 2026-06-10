package main

import (
	"context"
	"cyris/internal/checker"
	"cyris/internal/extras"
	"cyris/internal/steamspy"
	"cyris/internal/structs"
	"cyris/internal/system"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := extras.LoadGPUDatabase(); err != nil {
		fmt.Println("Warning: failed to load GPU database:", err)
	}
	if err := extras.LoadIntegratedDatabase(); err != nil {
		fmt.Println("Warning: failed to load integrated GPU database:", err)
	}
	if err := extras.LoadCPUDatabase(); err != nil {
		fmt.Println("Warning: failed to load CPU database:", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetHardwareInfo() structs.SystemSpecs {
	return system.GetHardwareInfo()
}

func (a *App) SearchGame(gameTerm string) []structs.SteamApp {
	return steamspy.SearchGame(gameTerm)
}

func (a *App) GetGameRequirements(gameId int) structs.GameData {
	return steamspy.GetGameID(gameId)
}

func (a *App) CheckRequirements(systemSpecs structs.SystemSpecs, gameData structs.GameData) structs.RequirementsResult {
	return checker.CheckRequirements(systemSpecs, gameData)
}
