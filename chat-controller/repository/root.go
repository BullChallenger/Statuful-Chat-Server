package repository

import (
	"chat-controller/config"
	"chat-controller/repository/kafka"
	"chat-controller/types/table"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Repository struct {
	config *config.Config
	db     *sql.DB
	Kafka  *kafka.Kafka
}

const (
	room       = "chatting.room"
	chat       = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)

func NewRepository(config *config.Config) (*Repository, error) {
	repository := &Repository{config: config}
	var err error

	if repository.db, err = sql.Open(config.DB.Database, config.DB.URL); err != nil {
		return nil, err
	} else if repository.Kafka, err = kafka.NewKafka(config); err != nil {
		return nil, err
	} else {
		return repository, nil
	}
}

func (repository *Repository) ReadAvailableServerInfo() ([]*table.ServerInfo, error) {
	qs := query([]string{"SELECT * FROM", serverInfo, "WHERE available = 1"})
	if cursor, err := repository.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		var result []*table.ServerInfo

		for cursor.Next() {
			d := new(table.ServerInfo)
			if err = cursor.Scan(&d.IP, &d.Available); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*table.ServerInfo{}, nil
		} else {
			return result, nil
		}
	}
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
