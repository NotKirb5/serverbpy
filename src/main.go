package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"golang.org/x/crypto/ssh"
	tea "github.com/charmbracelet/bubbletea"
)


var selectedServers []string

type Data struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	IP string `json:"ip"`
	USERNAME string `json:"username"`
}


type model struct {
	servers []string
	cursor int
	selected map[int]struct{}
}

type monitorModel struct {
	servers []string
}

func initserver() {
	content, err := ioutil.ReadFile("./ip.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	fmt.Println(string(content))
	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal: ", err)
	}
	for _, srv := range payload.Servers {
		fmt.Printf("IP: %s, User: %s\n", srv.IP, srv.USERNAME)
		fmt.Println("What is the password for ",string(srv.USERNAME))
		var pass string
		fmt.Scan(&pass)

		config := &ssh.ClientConfig{
			User: srv.USERNAME,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		addr := fmt.Sprintf("%s:22",srv.IP)
		fmt.Println("Connecting to:",addr)

		client, err := ssh.Dial("tcp", addr,config)
		if err != nil {
			log.Printf("Failed to connect to %s: %v\n", srv.IP,err)
			continue
		}
		defer client.Close()

		
		cmds := []string{"hostname","whoami","lsb_release -a","uname -r","uptime","free -h","sudo docker ps","df -h"}

		for _, cmd := range cmds {
			session, err := client.NewSession()
			if err != nil {
				log.Fatal(err)
			}
			output,err := session.CombinedOutput(cmd)
			session.Close()
			if err != nil {
				log.Printf("Failed to run command %s: %v\n",cmd,err)
				continue
			}
			fmt.Printf(string(output))
		}
	}
}


func main() {

	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Theres an error: %v", err)
		os.Exit(1)
	}

	

}

func getusers() ([]string){

	content, err := ioutil.ReadFile("./ip.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal: ", err)
	}

	var svrs []string
	for _, server := range payload.Servers {
		svrs = append(svrs, server.USERNAME)
	}
	return svrs
}

func initalModel() model {


	return model{
		servers: getusers(),
		selected: make(map[int]struct{}),
	}
}


func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.servers)-1 {
                m.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
				case "s":

        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m model) View() string {
    // The header
    s := "What should monitor?\n\n"

    // Iterate over our choices
    for i, server := range m.servers {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x"// selected!

        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, server)
    }

    // The footer
    s += "\nPress q to quit.\n Press s to start \n"

    // Send the UI for rendering
    return s
}


func MonitorModel() {
	//svrs := selectedServers
	fmt.Println(selectedServers)
}


