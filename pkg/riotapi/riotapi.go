package riotapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RiotAPI struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewRiotAPI Создаем новый RiotAPI клиент
func NewRiotAPI(apiKey, baseURL string) *RiotAPI {
	return &RiotAPI{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type MatchHistory []string

type PlayerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type MatchInfo struct {
	Metadata Metadata `json:"metadata"`
	Info     Info     `json:"info"`
}

type Metadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type Info struct {
	GameCreation       int64  `json:"gameCreation"`
	GameDuration       int    `json:"gameDuration"`
	GameEndTimestamp   int64  `json:"gameEndTimestamp"`
	GameID             int    `json:"gameId"`
	GameMode           string `json:"gameMode"`
	GameName           string `json:"gameName"`
	GameStartTimestamp int64  `json:"gameStartTimestamp"`
	GameType           string `json:"gameType"`
	GameVersion        string `json:"gameVersion"`
	MapID              int    `json:"mapId"`
	PlatformID         string `json:"platformId"`
	QueueID            int    `json:"queueId"`
	Teams              []Team `json:"teams"`
	TournamentCode     string `json:"tournamentCode"`
}

type Team struct {
	Bans       []Ban          `json:"bans"`
	Objectives TeamObjectives `json:"objectives"`
	TeamID     int            `json:"teamId"`
	Win        bool           `json:"win"`
}

type Ban struct {
	// Определение полей для запрещений
}

type TeamObjectives struct {
	Baron      Objective `json:"baron"`
	Champion   Objective `json:"champion"`
	Dragon     Objective `json:"dragon"`
	Inhibitor  Objective `json:"inhibitor"`
	RiftHerald Objective `json:"riftHerald"`
	Tower      Objective `json:"tower"`
}

type Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

// GetPlayerInfoByPUUID Получить информацию об игроке по его PUUID
func (api *RiotAPI) GetPlayerInfoByPUUID(puuid string) (*PlayerInfo, error) {
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-puuid/%s?api_key=%s", api.baseURL, puuid, api.apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playerInfo PlayerInfo
	err = json.Unmarshal(body, &playerInfo)
	if err != nil {
		return nil, err
	}

	return &playerInfo, nil
}

// GetMatchHistoryByPUUID Функция GetMatchHistoryByPUUID получает историю матчей игрока по его PUUID
func (api *RiotAPI) GetMatchHistoryByPUUID(puuid string) (MatchHistory, error) {

	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=100", api.baseURL, puuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var matchHistory MatchHistory

	err = json.Unmarshal(body, &matchHistory)

	if err != nil {
		return nil, err
	}

	return matchHistory, nil
}

func (api *RiotAPI) GetMatchInfo(matchID string) (*MatchInfo, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/%s", api.baseURL, matchID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var matchInfo MatchInfo

	err = json.Unmarshal(body, &matchInfo)

	if err != nil {
		return nil, err
	}

	return &matchInfo, nil
}
