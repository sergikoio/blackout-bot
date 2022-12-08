package bot

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"

	"blackout-bot/internal/db"
)

const (
	dbStatus           = "status"
	dbLastSend         = "last-send"
	dbLastScheduleSend = "last-schedule-send"
	dbServerConfig     = "server-config"
)

type botDB struct {
	db *db.DB
}

func newBotDB(db *db.DB) *botDB {
	return &botDB{
		db: db,
	}
}

type statusDB struct {
	*botDB
}

func (db *botDB) statusDB() *statusDB {
	return &statusDB{
		botDB: db,
	}
}

func (db *statusDB) set(status string) error {
	return db.db.Set([]byte(dbStatus), status)
}

func (db *statusDB) get() string {
	var phrase string
	err := db.db.Get([]byte(dbStatus), &phrase)
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return ""
		default:
			log.Error().Err(err).Send()
			return ""
		}
	}

	return phrase
}

type lastSendDB struct {
	*botDB
}

func (db *botDB) lastSendDB() *lastSendDB {
	return &lastSendDB{
		botDB: db,
	}
}

func (db *lastSendDB) set(send time.Time) error {
	return db.db.Set([]byte(dbLastSend), send)
}

func (db *lastSendDB) get() time.Time {
	var send time.Time
	err := db.db.Get([]byte(dbLastSend), &send)
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return time.Time{}
		default:
			log.Error().Err(err).Send()
			return time.Time{}
		}
	}

	return send
}

type lastScheduleSendDB struct {
	*botDB
}

func (db *botDB) lastScheduleSendDB() *lastScheduleSendDB {
	return &lastScheduleSendDB{
		botDB: db,
	}
}

func (db *lastScheduleSendDB) set(send time.Time) error {
	return db.db.Set([]byte(dbLastScheduleSend), send)
}

func (db *lastScheduleSendDB) get() time.Time {
	var send time.Time
	err := db.db.Get([]byte(dbLastScheduleSend), &send)
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return time.Time{}
		default:
			log.Error().Err(err).Send()
			return time.Time{}
		}
	}

	return send
}

type serverConfigDB struct {
	*botDB
}

func (db *botDB) serverConfigDB() *serverConfigDB {
	return &serverConfigDB{
		botDB: db,
	}
}

func (db *serverConfigDB) set(serverConfig ServerConfig) error {
	return db.db.Set([]byte(dbServerConfig), serverConfig)
}

func (db *serverConfigDB) get() ServerConfig {
	var serverConfig ServerConfig
	err := db.db.Get([]byte(dbServerConfig), &serverConfig)
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return ServerConfig{}
		default:
			log.Error().Err(err).Send()
			return ServerConfig{}
		}
	}

	return serverConfig
}
