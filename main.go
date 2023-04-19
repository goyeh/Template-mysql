package main

import (
	"all-news/lib"
	"all-news/route"
	"all-news/rss"
	"all-news/sql"
)

/*
https://medium.com/@vosmith08/scraping-the-web-in-go-b85ac37c14df
Done: Read folder for RSS files, Change them to json format. Add Category, that identifies the Table.
Done: Add paramters in the RSS list for Category, Target Media Telegram etc, and keys.
{	category:"",
	media:[{"Title":"Telegram", "KEY":"2o387230812308"},{"TITLE":"FaceBook","KEY":"09809098098098"}],
	rss:[	"http:\\ksajhdsajhlsj",
			"http:\\otherurl"    ]
}
*/

/***********************_**************
 *                     (_)
 *      _ __ ___   __ _ _ _ __
 *     | '_ ` _ \ / _` | | '_ \
 *     | | | | | | (_| | | | | |
 *     |_| |_| |_|\__,_|_|_| |_|
 * ---------------------------------- */
func main() {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Application Failure:", r)
		}
	}()
	lib.Info("Initialize Database")
	sql.InitDB()
	lib.Info("Initilize Stay Alive")
	go rss.StayAlive()
	lib.Info("Initilize Feed inport")
	rss.Init() //Initialise and load json files for the feed
	go rss.LoadFeeds() // Get the RSS feeds
	lib.Info("Initilize Posting to Channels")
	go rss.Post() // Post to various channels etc.
	lib.Info("Initilize Post Clean up")
	go rss.Cleanposts() // Remove old posts.
	//go openDiscord()

	route.Init()

}
