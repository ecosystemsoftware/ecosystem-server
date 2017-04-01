// Copyright 2017 Jonathan Pincas

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Use the forked version of the go-jwt-middlware, not the auth0 version

package ghost

import (
	"net/http"

	"github.com/jpincas/ghost/ghost"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
)

func init() {

	ServeCmd.Flags().String("smtppw", "", "SMTP server password for outgoing mail")
	ServeCmd.Flags().BoolP("demomode", "d", false, "Run server in demo mode")
	ServeCmd.Flags().BoolP("debug", "b", false, "Run server in debug mode")
	ServeCmd.Flags().StringP("secret", "s", "", "Secure secret for signing JWT")
	ServeCmd.Flags().StringP("pgpw", "p", "", "Postgres superuser password")
	ServeCmd.Flags().StringP("configfile", "c", "config", "Name of config file (without extension)")
	ServeCmd.Flags().BoolP("noprompt", "n", false, "Override prompt for confirmation")

	viper.BindPFlags(ServeCmd.Flags())

}

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the ghost server",
	Long:  `Start the ghost API Server`,
	RunE:  serve,
}

func serve(cmd *cobra.Command, args []string) error {

	ghost.App.Setup(viper.GetString("configfile"))
	preServe()
	startServer()
	return nil
}

//ActivatePackages is a hook for activating packages from main
var BeforeServe func()

func preServe() {

	//Setup the email system if required
	if ghost.App.Config.ActivateEmail {
		ghost.App.MailServer.Setup()
	}

	//Check to make sure a secret has been provided
	//No default provided as a security measure, server will exit of nothing provided
	if viper.GetString("secret") == "" {
		ghost.LogFatal("ghost.SERVE", false, "No signing secret provided", nil)
	}

	//Establish a temporary connection as the super user
	dbTemp := ghost.SuperUserDBConfig.ReturnDBConnection("")

	//Generate a random server password, set it and get out
	serverPW := ghost.RandomString(16)
	_, err := dbTemp.Exec(fmt.Sprintf(ghost.SQLToSetServerRolePassword, serverPW))
	if err != nil {
		ghost.LogFatal("ghost.SERVE", false, "Error setting server role password:", err)
	}

	dbTemp.Close()

	//Establish a permanent connection
	ghost.App.DB = ghost.ServerUserDBConfig.ReturnDBConnection(serverPW)

	BeforeServe()

}

func startServer() {

	ghost.Log("ghost.SERVE", true, "Server started on port "+viper.GetString("apiPort"), nil)
	http.ListenAndServe(":"+viper.GetString("apiPort"), ghost.App.Router)

}