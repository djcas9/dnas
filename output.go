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

func DatabaseConnect(options *Options) (db gorm.DB, err error) {

	var connect string

	connect = options.DbUser + ":" +
		options.DbPassword + "@" + options.DbHost + "/" +
		options.DbDatabase

	if options.DbTls {
		if options.DbSkipVerify {
			connect = connect + "?tls=skip-verify"
		} else {
			connect = connect + "?tls=true"
		}
	}

	db, err = gorm.Open("mysql", connect)

	if err != nil {
		return db, err
	}

	// Diable Logger
	if options.Quiet {
		db.LogMode(false)
	} else {
		db.LogMode(options.DatabaseOutput)
	}

	// defer db.Close()

	err = db.DB().Ping()

	if err != nil {
		return db, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.DropTableIfExists(&Question{})
	db.DropTableIfExists(&Answer{})

	if !db.HasTable(&Question{}) {
		db.CreateTable(&Question{})
	}

	if !db.HasTable(&Answer{}) {
		db.CreateTable(&Answer{})
	}

	db.AutoMigrate(&Question{}, &Answer{})

	return db, nil
}

func prettyPrint(message *Question, count int) {

	fmt.Printf("---   %d   ---\n", count)

	fmt.Printf("\nQuestion: \033[0;31;49m%s\033[0m\n\n", message.Question)

	fmt.Printf("\033[0;32;49mAnswers (%d):\033[0m\n\n", len(message.Answers))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"RR", "Name", "Data", "Last Seen"})
	// table.SetBorder(false)

	for _, aa := range message.Answers {

		table.Append([]string{
			aa.Record,
			aa.Name,
			lime + aa.Data + reset,
			aa.UpdatedAt.Format(layout),
		})
	}

	table.Render()
	fmt.Printf("\n")
}

func (question *Question) ToDatabase(db gorm.DB, options *Options) (err error) {
	var q Question

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

	} else {
		question.SeenCount = 1
		question.CreatedAt = time.Now().Unix()
		db.Table("questions").Create(question)
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
