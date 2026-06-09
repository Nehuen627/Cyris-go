package steamspy

import (
	"cyris/internal/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func SearchGame(gameTerm string) []structs.SteamApp {
	searchUrl := "https://steamcommunity.com/actions/SearchApps/" + url.QueryEscape(gameTerm)
	resp, err := http.Get(searchUrl)
	if err != nil {
		return []structs.SteamApp{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []structs.SteamApp{}
	}
	var result []structs.SteamApp
	json.Unmarshal(body, &result)
	return result
}

func GetGameID(gameId int) structs.GameData {
	gameUrl := "https://store.steampowered.com/api/appdetails?appids=" + fmt.Sprintf("%d", gameId)
	resp, err := http.Get(gameUrl)
	if err != nil {
		return structs.GameData{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return structs.GameData{}
	}
	type steamResponse struct {
		Success bool             `json:"success"`
		Data    structs.GameData `json:"data"`
	}

	var result map[string]steamResponse
	json.Unmarshal(body, &result)

	for _, v := range result {
		return v.Data
	}

	return structs.GameData{}
}
