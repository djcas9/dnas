package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jinzhu/gorm"
	"github.com/karlbunch/tablewriter"
	"github.com/mgutz/ansi"
)

var (
	count  = 0
	dbname = []byte("dnas")
	layout = "01/02/06 03:04pm (MST)"
	lime   = ansi.ColorCode("green+h:black")
	red    = ansi.ColorCode("red")
	green  = ansi.ColorCode("green")
	reset  = ansi.ColorCode("reset")
)

func FlushDatabase(db gorm.DB) {
	db.DropTableIfExists(&Question{})
	db.DropTableIfExists(&Answer{})
	db.DropTableIfExists(&Client{})
}

func DatabaseConnect(options *Options) (db gorm.DB, err error) {
	var connect string
	var dbType string

	if options.Mysql {

		dbType = "mysql"

		connect = options.DbUser + ":" +
			options.DbPassword + "@" + options.DbHost + "/" +
			options.DbDatabase

		if options.DbSsl {
			if options.DbSkipVerify {
				connect = connect + "?tls=skip-verify"
			} else {
				connect = connect + "?tls=true"
			}
		}

	}

	if options.Postgres {

		dbType = "postgres"

		connect = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s",
			options.DbUser, options.DbDatabase, options.DbPassword, options.DbHost, options.DbPort,
		)

		if options.DbSsl {
			if options.DbSkipVerify {
				connect = connect + " sslmode=require"
			} else {
				connect = connect + " sslmode=verify-full"
			}
		} else {
			connect = connect + " sslmode=disable"
		}

	}

	if options.Sqlite3 {
		dbType = "sqlite3"
		connect = options.DbPath
	}

	db, err = gorm.Open(dbType, connect)

	if err != nil {
		return db, err
	}

	// Diable Logger
	if options.Quiet {
		db.LogMode(false)
	} else {
		db.LogMode(options.DatabaseOutput)
	}

	err = db.DB().Ping()

	if err != nil {
		return db, err
	}

	if options.DbFlush {
		FlushDatabase(db)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	if !db.HasTable(&Question{}) {
		db.CreateTable(&Question{})
	}

	if !db.HasTable(&Answer{}) {
		db.CreateTable(&Answer{})
	}

	if !db.HasTable(&Client{}) {
		db.CreateTable(&Client{})
	}

	db.AutoMigrate(&Question{}, &Answer{})

	return db, nil
}

func CreateClient(db gorm.DB, options *Options) (id int64) {

	var c Client

	db.Where(
		&Client{
			Hostname:  options.Hostname,
			Interface: options.Interface,
			MacAddr:   options.InterfaceData.HardwareAddr.String(),
			Ip:        options.Ip,
		},
	).First(&c)

	if c.Id != 0 {
		db.Model(&c).Update(
			&Client{
				LastSeen: time.Now().Unix(),
			},
		)
	} else {
		c.LastSeen = time.Now().Unix()
		c.Hostname = options.Hostname
		c.MacAddr = options.InterfaceData.HardwareAddr.String()
		c.Interface = options.Interface
		c.Ip = options.Ip
		db.Table("clients").Create(&c)
	}

	options.Client = &c

	return c.Id
}

func prettyPrint(message *Question, count int) {

	fmt.Printf("---   %d   ---\n", count)

	fmt.Printf("\nQuestion: \033[0;31;49m%s\033[0m\n\n", message.Question)

	fmt.Printf("\033[0;32;49mAnswers (%d):\033[0m\n\n", len(message.Answers))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"RR", "Name", "Data", "Timestamp"})
	// table.SetBorder(false)

	for _, aa := range message.Answers {

		table.Append([]string{
			aa.Record,
			aa.Name,
			lime + aa.Data + reset,
			aa.CreatedAt.Format(layout),
		})
	}

	table.Render()
	fmt.Printf("\n")
}

func (question *Question) ToDatabase(db gorm.DB, options *Options) (err error) {
	var q Question
	var id int64 = 0

	db.Where(
		&Question{
			SrcIp:    question.SrcIp,
			DstIp:    question.DstIp,
			Question: question.Question,
		},
	).First(&q)

	if q.Id != 0 {

		db.Model(&q).Update(
			&Question{
				UpdatedAt: time.Now().Unix(),
				SeenCount: q.SeenCount + 1,
			},
		)

		id = q.Id

	} else {
		options.Client.QuestionCount = options.Client.QuestionCount + 1

		db.Exec("UPDATE clients SET question_count=? WHERE id=?",
			options.Client.QuestionCount, options.Client.Id)

		var time int64 = time.Now().Unix()

		question.SeenCount = 1

		question.CreatedAt = time
		question.UpdatedAt = time

		db.Table("questions").Create(question)

		id = question.Id
	}

	for _, a := range question.Answers {
		a.ClientId = options.Client.Id
		a.QuestionId = id
		a.CreatedAt = time.Now()

		db.Table("answers").Create(a)
	}

	return nil
}

// ToStdout write data to standard out
func (message *Question) ToStdout(options *Options) {

	count++

	prettyPrint(message, count)

	if options.Hexdump {
		dump, _ := hex.DecodeString(message.Packet)
		fmt.Printf("\033[0;32;49mHexdump:\033[0m\n\n%s\n", hex.Dump(dump))
	}

}

// ToJSON convert Message struct to json
func (message *Question) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
