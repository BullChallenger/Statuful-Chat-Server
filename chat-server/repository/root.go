package repository

import (
	"chat-server/config"
	"chat-server/repository/kafka"
	"chat-server/types/schema"
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

func (repository *Repository) ServerSet(ip string, available bool) error {
	_, err := repository.db.Exec("INSERT INTO serverInfo(`ip`, `available`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `available` = VALUES(`available`)", ip, available)
	return err
}

func (repository *Repository) InsertChatting(user, message, roomName string) error {
	_, err := repository.db.Exec("INSERT INTO chatting.chat(room, name, message) VALUES (?, ?, ?)", roomName, user, message)
	return err
}

func (repository *Repository) ReadChatList(roomName string) ([]*schema.Chat, error) {
	qs := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` ASC LIMIT 10"})

	if cursor, err := repository.db.Query(qs, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		var result []*schema.Chat

		for cursor.Next() {
			c := new(schema.Chat)
			if err = cursor.Scan(&c.ID, &c.Room, &c.Name, &c.Message, &c.When); err != nil {
				return nil, err
			} else {
				result = append(result, c)
			}
		}

		if len(result) == 0 {
			return []*schema.Chat{}, nil
		} else {
			return result, nil
		}
	}
}

func (repository *Repository) RoomList() ([]*schema.Room, error) {
	qs := query([]string{"SELECT * FROM", room})

	if cursor, err := repository.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		var result []*schema.Room

		for cursor.Next() {
			r := new(schema.Room)
			if err = cursor.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt); err != nil {
				return nil, err
			} else {
				result = append(result, r)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
	}
}

func (repository *Repository) MakeRoom(name string) error {
	_, err := repository.db.Exec("INSERT INTO chatting.room(name) VALUES(?)", name)
	return err
}

func (repository *Repository) Room(name string) (*schema.Room, error) {
	r := new(schema.Room)
	qs := query([]string{"SELECT * FROM", room, "WHERE NAME = ?"})
	err := repository.db.QueryRow(qs, name).Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt)

	if err := noResult(err); err != nil {
		return nil, err
	} else {
		return nil, nil
	}

	return r, err
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}

func noResult(err error) error {
	if strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil
	} else {
		return err
	}
}
