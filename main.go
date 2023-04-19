package main

import ( // after go mod init [app name] add the app name to these modules
	"[app name]/conf"
	"[app name]/lib"
	"[app name]/route"
	"[app name]/sql"
)

/*
	Welcome to the Template for starting
*/

/***********************_**************
 *                     (_)
 *      _ __ ___   __ _ _ _ __
 *     | '_ ` _ \ / _` | | '_ \
 *     | | | | | | (_| | | | | |
 *     |_| |_| |_|\__,_|_|_| |_|
 * ---------------------------------- */
func main() {
	defer func() { //This is a normal trap to prevent app crashing out. no point in main, but higher functions works well.
		r := recover()
		if r != nil {
			lib.Error("Application Failure:", r)
		}
	}()
	lib.Info("Initialize Database")
	sql.InitDB()
	lib.Info("Initilize Posting to Channels")

	route.Init()
}
