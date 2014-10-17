package main

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/gopass"

	"github.com/sevlyar/go-daemon"
	"github.com/visionmedia/go-flags"

	"os/user"
	"strconv"
	"syscall"
)

func chuser(username string) (uid, gid int) {
	usr, err := user.Lookup(username)

	if err != nil {
		fmt.Printf("failed to find user %q: %s", username, err)
	}

	uid, err = strconv.Atoi(usr.Uid)

	if err != nil {
		fmt.Printf("bad user ID %q: %s", usr.Uid, err)
	}

	gid, err = strconv.Atoi(usr.Gid)

	if err != nil {
		fmt.Printf("bad group ID %q: %s", usr.Gid, err)
	}

	if err := syscall.Setgid(gid); err != nil {
		fmt.Printf("setgid(%d): %s", gid, err)
	}

	if err := syscall.Setuid(uid); err != nil {
		fmt.Printf("setuid(%d): %s", uid, err)
	}

	return uid, gid
}

// Options cli command options
type Options struct {
	Interface string `short:"i" long:"interface" description:"Interface to monitor" value-name:"eth0"`
	Port      int    `short:"p" long:"port" description:"DNS port" default:"53" value-name:"53"`
	Daemon    bool   `short:"D" long:"daemon" description:"Run DNAS in daemon mode"`
	Write     string `short:"w" long:"write" description:"Write JSON output to log file" value-name:"FILE"`
	User      string `short:"u" long:"user" description:"Drop privileges to this user" value-name:"USER"`
	Hexdump   bool   `short:"H" long:"hexdump" description:"Show hexdump of DNS packet"`

	Mysql         bool   `short:"m" long:"mysql" description:"Enable Mysql Output Support" value-name:"PASSWORD"`
	MysqlUser     string `long:"mysql-user" description:"Mysql User" value-name:"root"`
	MysqlPassword string `long:"mysql-password" description:"Mysql Password" value-name:"PASSWORD"`
	MysqlDatabase string `long:"mysql-database" description:"Mysql Database" value-name:"dnas"`
	MysqlHost     string `long:"mysql-host" description:"Mysql Host" value-name:"127.0.0.1"`
	MysqlPort     string `long:"mysql-port" description:"Mysql Port" value-name:"3306"`

	Version bool `short:"v" long:"version" description:"Show version information"`
}

func printUsage(p *flags.Parser) {
	fmt.Printf("\n  %s (%s) - %s\n",
		Name,
		Version,
		Description,
	)

	p.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("%s - %s - Version: %s\n",
		Name,
		Description,
		Version,
	)

	os.Exit(1)
}

// CLIRun start DNAS and process all command-line options
func CLIRun(f func(options *Options)) {

	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		printUsage(parser)
	}

	if options.Version {
		printVersion()
	}

	if options.Mysql {
		if options.MysqlPassword == "" {
			password, err := gopass.GetPass("MySQL Password: ")

			if err != nil {
				panic(err)
			}

			options.MysqlPassword = password
		}
	}

	if options.MysqlHost != "" {
		options.MysqlHost = "tcp(" + options.MysqlHost + ":"
		if options.MysqlPort == "" {
			options.MysqlPort = "3306"
		}

		options.MysqlHost = options.MysqlHost + options.MysqlPort + ")"
	}

	if options.MysqlUser == "" {
		options.MysqlUser = "root"
	}

	if options.MysqlDatabase == "" {
		options.MysqlDatabase = "dnas"
	}

	if options.Interface == "" {
		printUsage(parser)
	}

	if options.Daemon {

		cntxt := &daemon.Context{
			PidFileName: "dnas.pid",
			PidFilePerm: 0644,
			LogFileName: "dnas.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()

		if err != nil {
			log.Fatalln(err)
		}

		if d != nil {
			return
		}

		defer cntxt.Release()

		go f(options)

		err = daemon.ServeSignals()

		if err != nil {
			log.Println("Error:", err)
		}
	} else {
		f(options)
	}
}
