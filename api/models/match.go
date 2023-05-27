package models

// MatchInfo содержит информацию о матче
type MatchTable struct {
	MatchPUUID  string
	WinningTeam string
	Duration    int32
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
	GameCreation       int64    `json:"gameCreation"`
	GameDuration       int      `json:"gameDuration"`
	GameEndTimestamp   int64    `json:"gameEndTimestamp"`
	GameID             int      `json:"gameId"`
	GameMode           string   `json:"gameMode"`
	GameName           string   `json:"gameName"`
	GameStartTimestamp int64    `json:"gameStartTimestamp"`
	GameType           string   `json:"gameType"`
	GameVersion        string   `json:"gameVersion"`
	MapID              int      `json:"mapId"`
	Participants       []string `json:"participants"`
	PlatformID         string   `json:"platformId"`
	QueueID            int      `json:"queueId"`
	Teams              []Team   `json:"teams"`
	TournamentCode     string   `json:"tournamentCode"`
}

type Team struct {
	Bans       []interface{} `json:"bans"` // Данные о банах отсутствуют, поэтому используется пустой интерфейс
	Objectives Objectives    `json:"objectives"`
	TeamID     int           `json:"teamId"`
	Win        bool          `json:"win"`
}

type Objectives struct {
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
