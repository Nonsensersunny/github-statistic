package model

import (
	"fmt"
	"github_statistics/pkg/db/sqlite"
	"time"
)

type Developer struct {
	Index     int64     `json:"index"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Project   string    `json:"project"`
	EventType string    `json:"event_type"`
	Action    string    `json:"action"`
	EventID   int64     `json:"event_id"`
}

func (d Developer) TableInfo() string {
	return `CREATE TABLE IF NOT EXISTS developer (
		'index' INTEGER PRIMARY KEY,
		id INTEGER,
		name VARCHAR(255),
		created_at DATE,
		project VARCHAR(255),
		event_type VARCHAR(64),
		action VARCHAR(64),
		event_id INT
	)`
}

func (d Developer) Insert() error {
	query := fmt.Sprintf("INSERT INTO developer(id, name, created_at, project, event_type, action, event_id) VALUES (%d, '%s', '%v', '%s', '%s', '%s', %d)", d.Id, d.Name, d.CreatedAt, d.Project, d.EventType, d.Action, d.EventID)
	return sqlite.Query(query)
}