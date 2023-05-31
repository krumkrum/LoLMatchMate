package main

import (
	"LoLMatchMate/pkg/database"
	"LoLMatchMate/pkg/riotapi"
	"fmt"
	"log"
)

func main() {
	const matchBaseURL = "https://europe.api.riotgames.com"
	const playerBaseURL = "https://ru.api.riotgames.com"

	var riotAPIKey = "RGAPI-bfe05208-f21c-4702-b45b-602763d1ced5"
	var puuidUser = "gfaji4Vf6jeWX37y02PUpWaTaKeRyGm-j0oQLoxw6oBPIzGkGLZy-vjoVbFNokbDnaRPUQwz4qvmTA"

	api := riotapi.NewRiotAPI(riotAPIKey, matchBaseURL)
	playerApi := riotapi.NewRiotAPI(riotAPIKey, playerBaseURL)

	db, err := database.NewDatabase("tcp://31.207.44.222:9000?username=develop&password=popok@31&database=LoLMatchMate")
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	matchIDs, err := api.GetMatchHistoryByPUUID(puuidUser)
	if err != nil {
		log.Fatalf("Failed to get match history: %v", err)
	}

	matchesInfo, err := api.GetMatchesInfo(matchIDs)
	if err != nil {
		log.Fatalf("Failed to get matches info: %v", err)
	}

	matchTables, err := riotapi.PrepareMatches(matchesInfo)
	if err != nil {
		log.Fatalf("Failed to prepare match tables: %v", err)
	}

	err = db.SaveMatches(matchTables)
	if err != nil {
		return
	}

	for _, match := range matchesInfo {
		exists, err := db.IsMatchPUUIDInPlayerMatches(match.Metadata.MatchID)
		if err != nil {
			log.Fatalf("Failed to check if match exists: %v", err)
		}
		if exists {
			fmt.Printf("MatchPUUID: %s already exists. Skipping...\n", match.Metadata.MatchID)
			continue
		}

		playersInfo, err := playerApi.GetPlayersInfoByPUUID(match.Metadata.Participants)
		if err != nil {
			log.Fatalf("Failed to get players info: %v", err)
		}

		playersReadyInfo, err := riotapi.PreparePlayersFromMatch(playersInfo, match.Metadata.MatchID)
		if err != nil {
			return
		}

		err = db.SavePlayerInfos(playersReadyInfo)
		if err != nil {
			return
		}
	}
	err = db.PrintPlayerStats()
	if err != nil {
		return
	}
}
