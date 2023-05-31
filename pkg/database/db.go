package database

import (
	"LoLMatchMate/api/models"
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(dsn string) (*Database, error) {
	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db := &Database{
		conn: conn,
	}
	return db, nil
}

func (db *Database) SaveMatches(matches models.MatchTables) error {
	// Проверка соединения с базой данных
	if err := db.conn.Ping(); err != nil {
		return fmt.Errorf("connection to the database failed: %w", err)
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO Matches (MatchPUUID, WinningTeam, Duration) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	for _, match := range matches {
		// Проверка существования MatchPUUID в базе данных
		var existingMatchID string
		err = db.conn.QueryRow(`SELECT MatchPUUID FROM Matches WHERE MatchPUUID = ?`, match.MatchPUUID).Scan(&existingMatchID)
		if err == nil {
			// MatchPUUID уже существует, переходим к следующему
			fmt.Printf("MatchPUUID: %s already exists. Skipping...\n", match.MatchPUUID)
			continue
		} else if err != sql.ErrNoRows {
			// Если ошибка не связана с отсутствием результата, возвращаем ошибку
			return err
		}

		fmt.Printf("Inserting matchID: %s\n", match.MatchPUUID)
		if _, err := stmt.Exec(match.MatchPUUID, match.WinningTeam, match.Duration); err != nil {
			fmt.Printf("Failed to insert matchID: %s with error: %v\n", match.MatchPUUID, err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) SavePlayerInfos(playerInfos models.TablePlayersInfo) error {
	// Проверка соединения с базой данных
	if err := db.conn.Ping(); err != nil {
		return fmt.Errorf("connection to the database failed: %w", err)
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO PlayerMatches (PUUID, Name, MatchPUUID, SummonerLevel) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	for _, playerInfo := range playerInfos {
		if playerInfo.Name == "" {
			fmt.Println("Player name is empty. ", playerInfo.PUUID, playerInfo.MatchPUUID)
			continue
		}
		// Проверка существования PUUID в базе данных
		var existingPUUID string
		err = db.conn.QueryRow(`SELECT PUUID FROM PlayerMatches WHERE PUUID = ? AND MatchPUUID = ?`, playerInfo.PUUID, playerInfo.MatchPUUID).Scan(&existingPUUID)
		if err == nil {
			// PUUID уже существует для данного MatchPUUID, переходим к следующему
			fmt.Printf("PUUID: %s already exists for MatchPUUID: %s. Skipping...\n", playerInfo.PUUID, playerInfo.MatchPUUID)
			continue
		} else if err != sql.ErrNoRows {
			// Если ошибка не связана с отсутствием результата, возвращаем ошибку
			return err
		}

		fmt.Printf("Inserting Player: %s\n", playerInfo.Name)
		if _, err := stmt.Exec(playerInfo.PUUID, playerInfo.Name, playerInfo.MatchPUUID, playerInfo.SummonerLevel); err != nil {
			fmt.Printf("Failed to insert PUUID: %s with error: %v\n", playerInfo.PUUID, err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) IsMatchPUUIDInPlayerMatches(matchPUUID string) (bool, error) {
	// Проверка соединения с базой данных
	if err := db.conn.Ping(); err != nil {
		return false, fmt.Errorf("connection to the database failed: %w", err)
	}

	var existingMatchPUUID string
	err := db.conn.QueryRow(`SELECT MatchPUUID FROM PlayerMatches WHERE MatchPUUID = ?`, matchPUUID).Scan(&existingMatchPUUID)

	if err == sql.ErrNoRows {
		// Нет результатов, значит, MatchPUUID не существует в таблице PlayerMatches
		return false, nil
	} else if err != nil {
		// Если ошибка не связана с отсутствием результата, возвращаем ошибку
		return false, err
	}

	// MatchPUUID существует в таблице PlayerMatches
	return true, nil
}

func (db *Database) ListPlayerStats() ([]models.PlayerStat, error) {
	rows, err := db.conn.Query(`
		SELECT Name, COUNT(*), MAX(SummonerLevel) 
		FROM PlayerMatches 
		GROUP BY Name
		ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.PlayerStat
	for rows.Next() {
		var stat models.PlayerStat
		err := rows.Scan(&stat.Name, &stat.MatchCount, &stat.MaxSummonerLevel)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (db *Database) PrintPlayerStats() error {
	// Проверка соединения с базой данных
	if err := db.conn.Ping(); err != nil {
		return fmt.Errorf("connection to the database failed: %w", err)
	}

	rows, err := db.conn.Query(`
		SELECT Name, COUNT(*) as Matches, SummonerLevel
		FROM PlayerMatches
		GROUP BY Name, SummonerLevel
		ORDER BY Matches DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Player Name; Matches Count; Summoner Level")

	for rows.Next() {
		var name string
		var matches int
		var summonerLevel int32

		err := rows.Scan(&name, &matches, &summonerLevel)
		if err != nil {
			return err
		}

		fmt.Printf("%s; %d; %d\n", name, matches, summonerLevel)
	}

	// проверяем на наличие ошибок при итерации
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
