package conf

import (
	"[app name]/lib"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	ConfigFilePath    string
	DEBUG             string
	VERSION           string
	AppName           string
	LOGDIR            string
	PORT              string
	HEARTBEAT         int
	RSSFILES          string
	ATOMFILES         string
	JSONFILES         string
	CHATID            string
	BOTID             string
	POSTDELAY         int
	MYSQL_USER        string
	MYSQL_PASS        string
	MYSQL_DB          string
	MYSQL_PORT        string
	MYSQL_HOST        string
	POSTAGE           int
	DISCORDTOKEN      string
	NEWSDETAIL        bool
	NEWSLIMIT         int
	MONITORAPI        string
	MYDNS             string
	THROTTLE          int
	UPDATE_NEWSDETAIL bool
)

/***********************************
 *                     _
 *                    | |
 *      _ __ ___ _ __ | |_   _
 *     | '__/ _ \ '_ \| | | | |
 *     | | |  __/ |_) | | |_| |
 *     |_|  \___| .__/|_|\__, |
 *              | |       __/ |
 *              |_|      |___/
 * * * * * * * * * * * * * * * * * *
 * Define the reply structure
 * ------------------------------- */
type Reply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var RestartSequece bool

func init() { // Autoloaded on run
	RestartSequece = true
	flag.StringVar(&ConfigFilePath, "c", "app.conf", "config file path")
	flag.Parse()
	lib.CheckErr(godotenv.Load(ConfigFilePath))

	LoadConfig()
	watchFile := ConfigFilePath
	go lib.WatchFileAndRun(watchFile, LoadConfig) //This trick allows the config to be reloaded on edit.
}

/*******_*********************_**_____*************__*_***********
 *     | |                   | |/ ____|           / _(_)
 *     | |     ___   __ _  __| | |     ___  _ __ | |_ _  __ _
 *     | |    / _ \ / _` |/ _` | |    / _ \| '_ \|  _| |/ _` |
 *     | |___| (_) | (_| | (_| | |___| (_) | | | | | | | (_| |
 *     |______\___/ \__,_|\__,_|\_____\___/|_| |_|_| |_|\__, |
 *                                                       __/ |
 *                                                      |___/   */
func LoadConfig() {
	defer func() {
		r := recover()
		if r != nil {
			log.Print("Possible app.conf error:", r)
		}
	}()

	DEBUG = getEnv("DEBUG", "ERROR DEBUG INFO")
	VERSION = getEnv("VERSION", "0.2.1")
	AppName = getEnv("APPNAME", filepath.Base(os.Args[0]))
	LOGDIR = getEnv("LOGDIR", ".")
	PORT = getEnv("PORT", "7451")
	HEARTBEAT = getEnvAsInt("HEARTBEAT", 60)
	RSSFILES = getEnv("RSSFILES", "./rss")
	ATOMFILES = getEnv("ATOMFILES", "./atom")
	JSONFILES = getEnv("JSONFILES", "./json")
	CHATID = getEnv("CHATID", "77612747")
	BOTID = getEnv("BOTID", "1204200932:AAFR-Rr_kSzqSR4XnpcslTtVc0ddSRL1z_U")
	POSTDELAY = getEnvAsInt("POSTDELAY", 5)
	MYSQL_USER = getEnv("MYSQL_USER", "news")
	MYSQL_PASS = getEnv("MYSQL_PASS", "NewsMe101")
	MYSQL_DB = getEnv("MYSQL_DB", "news")
	MYSQL_PORT = getEnv("MYSQL_PORT", "3306")
	MYSQL_HOST = getEnv("MYSQL_HOST", "localhost")
	POSTAGE = getEnvAsInt("POSTAGE", 362)
	DISCORDTOKEN = getEnv("DISCORDTOKEN", "NzczODIzNTE2NjA0MTA0NzA0.X6O1Tw.pMjvmtozxsv2K7FAj69tHVSTdoc")
	NEWSDETAIL = getEnvAsBool("NEWSDETAIL", true)
	NEWSLIMIT = getEnvAsInt("NEWSLIMIT", 100)
	MONITORAPI = "http://" + getEnv("MONITORAPI", "domains.aenxchange.com:7440")
	MYDNS = getEnv("MYDNS", "news.aensmart.com")
	THROTTLE = getEnvAsInt("THROTTLE", 10)
	UPDATE_NEWSDETAIL = getEnvAsBool("UPDATE_NEWSDETAIL", false)

	lib.LogInit(DEBUG, AppName) //Global debug levels.
	lib.Info("Logfile:", AppName)
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	defer func() {
		r := recover()
		if r != nil {
			log.Print("Get Environment String:", key, r)
		}
	}()

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Get Environment Bool:", name, r)
		}
	}()
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Get Environment Int:", name, r)
		}
	}()
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Get Environment Duration:", name, r)
		}
	}()

	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		retVal := time.Duration(value)
		return retVal
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt64(name string, defaultVal int64) int64 {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Get Environment Int64:", name, r)
		}
	}()
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return defaultVal
}

//  ----------------------------------------------------------------
//              _   _   _                 _   _                _
//             | | | | | |               | | | |              | |
//    __ _  ___| |_| |_| | ___  __ _ _ __| |_| |__   ___  __ _| |_
//   / _` |/ _ \ __|  _  |/ _ \/ _` | '__| __| '_ \ / _ \/ _` | __|
//  | (_| |  __/ |_| | | |  __/ (_| | |  | |_| |_) |  __/ (_| | |_
//   \__, |\___|\__\_| |_/\___|\__,_|_|   \__|_.__/ \___|\__,_|\__|
//    __/ |
//   |___/
func nextHeartBeat() (nextHeartBeat int) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Next Heartbeat:", nextHeartBeat, r)
		}
	}()

	nextHeartBeat = HEARTBEAT
	if strings.Contains(DEBUG, "DEBUG") {
		nextHeartBeat = lib.RandomRange(nextHeartBeat/6, nextHeartBeat/2)
	} else {
		nextHeartBeat = lib.RandomRange(nextHeartBeat/2, nextHeartBeat*2)
	}
	return
}
