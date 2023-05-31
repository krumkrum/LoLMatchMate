package riotapi

import (
	"LoLMatchMate/api/models"
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

// GetPlayerInfoByPUUID Получить информацию об игроке по его PUUID
func (api *RiotAPI) GetPlayerInfoByPUUID(playerID string) (*models.PlayerInfo, error) {
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-puuid/%s?api_key=%s", api.baseURL, playerID, api.apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playerInfo models.PlayerInfo
	err = json.Unmarshal(body, &playerInfo)
	if err != nil {
		return nil, err
	}

	return &playerInfo, nil
}

func (api *RiotAPI) GetPlayersInfoByPUUID(playerIDs []string) ([]*models.PlayerInfo, error) {

	var playersInfo []*models.PlayerInfo

	for _, playerID := range playerIDs {
		playerInfo, err := api.GetPlayerInfoByPUUID(playerID)
		if err != nil {
			return nil, err
		}
		playersInfo = append(playersInfo, playerInfo)

	}

	return playersInfo, nil
}

// GetMatchHistoryByPUUID Функция GetMatchHistoryByPUUID получает историю матчей игрока по его PUUID
func (api *RiotAPI) GetMatchHistoryByPUUID(puuid string) (MatchHistory, error) {

	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=70", api.baseURL, puuid)
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

func (api *RiotAPI) GetMatchInfo(matchID string) (*models.MatchInfo, error) {
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

	var matchInfo models.MatchInfo

	err = json.Unmarshal(body, &matchInfo)

	if err != nil {
		return nil, err
	}

	return &matchInfo, nil
}

func (api *RiotAPI) GetMatchesInfo(matchIDs []string) ([]*models.MatchInfo, error) {
	var matchesInfo []*models.MatchInfo
	for _, matchID := range matchIDs {
		matchInfo, err := api.GetMatchInfo(matchID)
		if err != nil {
			return nil, err
		}
		matchesInfo = append(matchesInfo, matchInfo)
	}
	return matchesInfo, nil
}

func PrepareMatch(data models.MatchInfo) (*models.MatchTable, error) {

	matchTable := models.MatchTable{MatchPUUID: data.Metadata.MatchID,
		WinningTeam: fmt.Sprintf("%t", data.Info.Teams[0].Win),
		Duration:    int32(data.Info.GameDuration)}

	return &matchTable, nil
}

func PrepareMatches(data []*models.MatchInfo) (models.MatchTables, error) {
	var matchTables models.MatchTables
	for _, matchInfo := range data {
		matchTable, err := PrepareMatch(*matchInfo)
		if err != nil {
			return nil, err
		}
		matchTables = append(matchTables, *matchTable)
	}
	return matchTables, nil
}

// GAME_ID -> 10 PLAYERS ID -> GET_PLAYER_INFO -> PREPARE_PLAYER

func PreparePlayerFromMatch(data models.PlayerInfo, matchPuuid string) (*models.TablePlayerInfo, error) {
	playerInfo := models.TablePlayerInfo{PUUID: data.PUUID,
		Name:          data.Name,
		MatchPUUID:    matchPuuid,
		SummonerLevel: int32(data.SummonerLevel)}

	return &playerInfo, nil
}

func PreparePlayersFromMatch(data []*models.PlayerInfo, matchPuuid string) (models.TablePlayersInfo, error) {
	var tablePlayers models.TablePlayersInfo
	for _, PlayerInfo := range data {

		tablePlayer, err := PreparePlayerFromMatch(*PlayerInfo, matchPuuid)
		if err != nil {
			return nil, err
		}
		tablePlayers = append(tablePlayers, *tablePlayer)
	}
	return tablePlayers, nil
}
