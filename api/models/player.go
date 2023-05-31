package models

// PlayerInfo содержит информацию об игроке
type PlayerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type PlayerStat struct {
	Name             string
	MatchCount       int
	MaxSummonerLevel int
}

type TablePlayerInfo struct {
	PUUID         string
	Name          string
	MatchPUUID    string
	SummonerLevel int32
}

type TablePlayersInfo []TablePlayerInfo
