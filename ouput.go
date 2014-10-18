package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func MysqlConnect(options *Options) (db gorm.DB, err error) {

	var connect string

	connect = options.MysqlUser + ":" +
		options.MysqlPassword + "@" + options.MysqlHost + "/" +
		options.MysqlDatabase

	// if options.MysqlTLS {
	// if options.MysqlSkipVerify {
	// connect += "?tls=skip-verify"
	// } else {
	// connect += "?tls=true"
	// }
	// }

	db, err = gorm.Open("mysql", connect)

	if err != nil {
		return db, err
	}

	// Diable Logger
	db.LogMode(false)

	// defer db.Close()

	err = db.DB().Ping()

	if err != nil {
		return db, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.CreateTable(&Question{})
	db.CreateTable(&Answer{})

	db.AutoMigrate(&Question{}, &Answer{})

	return db, nil
}

func prettyPrint(message *Question, count int) {

	fmt.Printf("---   %d   ---\n", count)

	fmt.Printf("\nQuestion: \033[0;31;49m%s\033[0m\n\n", message.Question)

	fmt.Printf("\033[0;32;49mAnswers (%d):\033[0m\n\n", len(message.Answers))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"RR", "Name", "Data", "Last Seen", "Active"})
	// table.SetBorder(false)

	for i, aa := range message.Answers {

		active := green + "Yes" + reset

		if !message.Answers[i].Active {
			active = red + "No" + reset
		}

		table.Append([]string{
			aa.Record,
			aa.Name,
			lime + aa.Data + reset,
			aa.UpdatedAt.Format(layout),
			active,
		})
	}

	table.Render()
	fmt.Printf("\n")
}

func (question *Question) ToMysql(db gorm.DB, options *Options) (err error) {

	// db.Create(question)

	db.FirstOrCreate(
		question,
		Question{SrcIp: question.SrcIp, DstIp: question.DstIp, Question: question.Question},
	)

	// _, err = db.Exec("INSERT INTO questions (question, packet, src_ip, dst_ip, timestamp, protocol) VALUES (?, ?, ?, ?, ?, ?);",
	// message.Question, message.Packet, message.SrcIP, message.DstIP, message.Timestamp, message.Protocol)

	// if err != nil {
	// fmt.Println(err.Error())
	// }

	// for _, aa := range message.Answers {

	// _, err = db.Exec("INSERT INTO answers (question_id, name, record, data, created_at, updated_at, active) VALUES (?, ?, ?, ?, ?, ?, ?);",
	// insertCount, aa.Name, aa.Record, aa.Data, aa.CreatedAt, aa.UpdatedAt, aa.Active)

	// if err != nil {
	// fmt.Println(err.Error())
	// }
	// }

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
