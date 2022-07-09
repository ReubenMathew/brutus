/*
Copyright © 2022 REUBEN NINAN <intelphysic@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const url = "https://brutalist.report/api/v1"

type statusMsg int

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

type model struct {
  body string
  err error
}

func (m model) Init() tea.Cmd {
  return getAllArticles
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case statusMsg:
		m.body = string(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
  return "Thank you for using Brutus"
}

func getAllArticles() tea.Msg {
  c := &http.Client{
    Timeout: 10 * time.Second,
  }
  res, err := c.Get(url + "/topics")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  var rawJsonResp []interface{}
  err = json.NewDecoder(res.Body).Decode(&rawJsonResp)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%s", rawJsonResp[0])
  return rawJsonResp
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "brutus",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
          p := tea.NewProgram(model{})
          if err := p.Start(); err != nil {
            fmt.Println(err)
          }
        },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.brutus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


