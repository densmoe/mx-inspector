package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
)

type Model struct {
	MendixVersion   string `json:"MendixVersion"`
	Constants       []Constant
	ScheduledEvents []ScheduledEvent
	ProjectSettings ProjectSettings
}

type Constant struct {
	Name         string `json:"Name"`
	DataType     string `json:""`
	DefaultValue string `json:"DefaultValue"`
}

type ScheduledEvent struct {
	Name string `json:"Name"`
}

type ProjectSettings struct {
	Configurations []Configuration
}

type Configuration struct {
	Name string `json:"Name"`
}

func Load(path string) Model {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	model := parse(db)
	// u, err := json.MarshalIndent(model, "", "    ")
	// fmt.Printf(string(u))
	return model
}

func parse(db *sql.DB) Model {
	model := Model{}

	err := db.QueryRow("SELECT _ProductVersion FROM _Metadata").Scan(&model.MendixVersion)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT Contents FROM Unit WHERE ContainmentName = "Documents"`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var raw bson.Raw
		err = rows.Scan(&raw)
		if err != nil {
			log.Fatal(err)
		}

		err = raw.Validate()
		if err != nil {
			log.Fatal(err)
		}
		val := raw.Lookup("$Type")
		result := val.String()

		if result == `"Constants$Constant"` {
			constant := Constant{}
			constant.Name = raw.Lookup("Name").String()
			constant.DefaultValue = raw.Lookup("DefaultValue").String()
			model.Constants = append(model.Constants, constant)
		} else if result == `"ScheduledEvents$ScheduledEvent"` {
			scheduledEvent := ScheduledEvent{}
			scheduledEvent.Name = raw.Lookup("Name").String()
			model.ScheduledEvents = append(model.ScheduledEvents, scheduledEvent)
		}
	}

	rows, err = db.Query(`SELECT Contents FROM Unit WHERE ContainmentName = "ProjectDocuments"`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var raw bson.Raw
		err = rows.Scan(&raw)
		if err != nil {
			log.Fatal(err)
		}

		err = raw.Validate()
		if err != nil {
			log.Fatal(err)
		}
		val := raw.Lookup("$Type")
		result := val.String()

		if result == `"Settings$ProjectSettings"` {
			projectSettings := ProjectSettings{}
			model.ProjectSettings = projectSettings
			// fmt.Println(jsonPrettyPrint(string(raw.Lookup("Settings", "3", "Configurations").String())))
		}
	}

	return model
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}
