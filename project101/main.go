package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	l "project101/modules/login"
	db "project101/src/database/init"

	//s "project101/src/srv/chat/server"
	d "project101/modules/chatDialer"

	tea "github.com/charmbracelet/bubbletea"
)
var (
	database = make(chan *sql.DB,1)
)
func init(){	
	
	log.SetFlags(log.Llongfile)//flags that point to exact error line for debugging
	db, err := db.Config()//initialises the database
	   if err != nil{
			log.Fatalf("%v", err)
		}
		fmt.Println("Initialised database...")
	database <- db //retrieves the *sql.Db object
	//initialise the chat server
	cmd := exec.Command("/home/iambronnix/project101/src/srv/chat/server/./main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil{
		log.Fatalf("%v",err)
	}

}

type model struct{
	modules []string //modules provided
	cursor int      //which module cursor is pointing at
	selected map[int]struct{} //which modules selected
}

func initialModel() model{
	return model{
		modules: []string{"login", "exit"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd{
	return nil
}

func main(){
	loginErr, userName := l.Login(<-database)
	 if loginErr != nil{
			log.Fatalf("%v", loginErr)
				}
		fmt.Println(userName)
		d.ClientDialer()
			
}