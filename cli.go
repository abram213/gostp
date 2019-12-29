package gostp

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

// Migrate soft migration for initialized models, data won't be deleted
func Migrate() {
	Db.AutoMigrate(Models...)
}

// DropTables drops all data. Be careful!
func DropTables() {
	Db.DropTableIfExists(Models...)
}

// СheckArguments - checks cli arguments and do stuff
func СheckArguments(port *string) {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "migrate" {
			Migrate()
			fmt.Println("Migration successful")
			os.Exit(0)
		} else if arg == "migrate:refresh" {
			DropTables()
			Migrate()
			fmt.Println("All tables dropped, migration successful")
			os.Exit(0)
		} else if arg == "newuser" {
			GenerateUser()
			os.Exit(0)
		} else if arg == "help" {
			fmt.Println("Available commands. Program will be stoped after execution:")

			color.New(color.FgGreen).Println("---------")
			color.New(color.FgYellow).Print("migrate")
			fmt.Print(" - soft migration\n")

			color.New(color.FgGreen).Println("---------")
			color.New(color.FgYellow).Print("migrate:refresh")
			fmt.Println(" - delete all data and migrate from scratch")

			color.New(color.FgGreen).Println("---------")
			color.New(color.FgYellow).Print("newuser")
			fmt.Println(" - create new user")

			color.New(color.FgGreen).Println("---------")
			fmt.Println("Available arguments. Program will continue running with these arguments:")

			color.New(color.FgGreen).Println("---------")
			color.New(color.FgYellow).Print("-p")
			fmt.Print(" - port and host to listen. For example:")
			color.New(color.FgYellow).Print(" :7777")
			fmt.Print(" for all hosts, or")
			color.New(color.FgYellow).Print(" 127.0.0.1:7777")
			fmt.Print(" only for localhost\n")

			color.New(color.FgGreen).Println("---------")
			color.New(color.FgYellow).Print("-key")
			fmt.Print(" - secret key which will be used to generate tokens. Default is")
			color.New(color.FgYellow).Print(" 1234")
			fmt.Print(".")
			color.New(color.FgRed).Print(" Please, change it.\n")
			fmt.Println("")
			os.Exit(0)
		}
		// If found at least one argument - countinue
		continueCondition := false
		for index, loopArg := range os.Args {
			if loopArg == "-p" {
				if index+1 <= len(os.Args) {
					*port = os.Args[index+1]
					continueCondition = true
				} else {
					fmt.Println("Not enough arguments to start at custom port")
					os.Exit(0)
				}
			}
			if loopArg == "-key" {
				if index+1 <= len(os.Args) {
					Settings.SigningKey = os.Args[index+1]
					continueCondition = true
				} else {
					fmt.Println("Not enough arguments to start with custom secret key")
					os.Exit(0)
				}
			}
		}
		if !continueCondition {
			fmt.Println("Comand not found")
			os.Exit(0)
		}
	}
}
