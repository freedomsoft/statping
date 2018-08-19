// Statup
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statup
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package handlers

import (
	"github.com/hunterlong/statup/core"
	"github.com/hunterlong/statup/types"
	"github.com/hunterlong/statup/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SetupHandler(w http.ResponseWriter, r *http.Request) {
	if core.CoreApp.Services != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	w.WriteHeader(http.StatusOK)
	port := 5432
	if os.Getenv("DB_CONN") == "mysql" {
		port = 3306
	}
	var data interface{}
	if os.Getenv("DB_CONN") != "" {
		data = &types.DbConfig{
			DbConn:      os.Getenv("DB_CONN"),
			DbHost:      os.Getenv("DB_HOST"),
			DbUser:      os.Getenv("DB_USER"),
			DbPass:      os.Getenv("DB_PASS"),
			DbData:      os.Getenv("DB_DATABASE"),
			DbPort:      port,
			Project:     os.Getenv("NAME"),
			Description: os.Getenv("DESCRIPTION"),
			Email:       "",
			Username:    "admin",
			Password:    "",
		}
	}
	ExecuteResponse(w, r, "setup.html", data)
}

func ProcessSetupHandler(w http.ResponseWriter, r *http.Request) {
	if core.CoreApp.Services != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	dbHost := r.PostForm.Get("db_host")
	dbUser := r.PostForm.Get("db_user")
	dbPass := r.PostForm.Get("db_password")
	dbDatabase := r.PostForm.Get("db_database")
	dbConn := r.PostForm.Get("db_connection")
	dbPort, _ := strconv.Atoi(r.PostForm.Get("db_port"))
	project := r.PostForm.Get("project")
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	sample := r.PostForm.Get("sample_data")
	description := r.PostForm.Get("description")
	domain := r.PostForm.Get("domain")
	email := r.PostForm.Get("email")

	config := &core.DbConfig{&types.DbConfig{
		DbConn:      dbConn,
		DbHost:      dbHost,
		DbUser:      dbUser,
		DbPass:      dbPass,
		DbData:      dbDatabase,
		DbPort:      dbPort,
		Project:     project,
		Description: description,
		Domain:      domain,
		Username:    username,
		Password:    password,
		Email:       email,
		Error:       nil,
		Location:    ".",
	}}

	err := config.Save()
	if err != nil {
		utils.Log(4, err)
	}

	if err != nil {
		utils.Log(3, err)
		config.Error = err
		SetupResponseError(w, r, config)
		return
	}

	core.Configs, err = core.LoadConfig()
	if err != nil {
		utils.Log(3, err)
		config.Error = err
		SetupResponseError(w, r, config)
		return
	}

	err = core.DbConnection(core.Configs.Connection, false, ".")
	if err != nil {
		utils.Log(3, err)
		core.DeleteConfig()
		config.Error = err
		SetupResponseError(w, r, config)
		return
	}

	admin := &types.User{
		Username: config.Username,
		Password: config.Password,
		Email:    config.Email,
		Admin:    true,
	}
	core.CreateUser(admin)

	if sample == "on" {
		core.LoadSampleData()
	}

	core.InitApp()
	resetCookies()
	time.Sleep(2 * time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SetupResponseError(w http.ResponseWriter, r *http.Request, a interface{}) {
	ExecuteResponse(w, r, "setup.html", a)
}
