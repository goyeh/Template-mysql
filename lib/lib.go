package lib

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html"
	"io"
	"log"
	"log/syslog"
	"math"
	"math/rand"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/gomail.v2"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var (
	logOut     *syslog.Writer
	DebugLevel string
)

// ----------------------------------------------------------------------------
//      ____                  __                   ____
//     / __ \____ _____  ____/ /___  ____ ___     / __ \____ _____  ____ ____
//    / /_/ / __ `/ __ \/ __  / __ \/ __ `__ \   / /_/ / __ `/ __ \/ __ `/ _ \
//   / _, _/ /_/ / / / / /_/ / /_/ / / / / / /  / _, _/ /_/ / / / / /_/ /  __/
//  /_/ |_|\__,_/_/ /_/\__,_/\____/_/ /_/ /_/  /_/ |_|\__,_/_/ /_/\__, /\___/
//                                                               /____/
func randomRange(min int, max int) (randomRange int) {
	rand.Seed(time.Now().UnixNano())
	randomRange = rand.Intn(max-min) + min
	return
}

/***    _       _ _
 *     (_)     (_) |
 *      _ _ __  _| |_
 *     | | '_ \| | __|
 *     | | | | | | |_
 *     |_|_| |_|_|\__|  * * */
func LogInit(level string, name string) {
	var err error
	DebugLevel = level
	logOut, err = syslog.New(syslog.LOG_INFO, name)
	log.SetOutput(logOut)
	log.Println("Log Init:", logOut, err)
}

/****   _             _    _      _
 *     | |           | |  | |    | |
 *     | | ___   __ _| |__| | ___| |_ __   ___ _ __ ___
 *     | |/ _ \ / _` |  __  |/ _ \ | '_ \ / _ \ '__/ __|
 *     | | (_) | (_| | |  | |  __/ | |_) |  __/ |  \__ \
 *     |_|\___/ \__, |_|  |_|\___|_| .__/ \___|_|  |___/
 *               __/ |             | |
 *              |___/              |_|                        */
func Debug(msg ...interface{}) { logCore("DEBUG", msg...) }
func Info(msg ...interface{})  { logCore("INFO", msg...) }
func Warn(msg ...interface{})  { logCore("WARN", msg...) }
func Error(msg ...interface{}) { logCore("ERROR", msg...) }
func Crit(msg ...interface{})  { logCore("CRIT", msg...) }

/* ------------------------------------- */

/***    _              _____
 *     | |            / ____|
 *     | | ___   __ _| |     ___  _ __ ___
 *     | |/ _ \ / _` | |    / _ \| '__/ _ \
 *     | | (_) | (_| | |___| (_) | | |  __/
 *     |_|\___/ \__, |\_____\___/|_|  \___|
 *               __/ |
 *              |___/                           */
func logCore(level string, msg ...interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error detected logging:", r)
		}
	}()
	var err error
	if strings.Contains(DebugLevel, level) {
		switch level {
		case "INFO":
			err = logOut.Info(fmt.Sprint(msg...))
		case "DEBUG":
			err = logOut.Debug(fmt.Sprint(msg...))
		case "ERROR":
			err = logOut.Err(fmt.Sprint(msg...))
		case "CRIT":
			err = logOut.Crit(fmt.Sprint(msg...))
		}
		if strings.Contains(DebugLevel, "STDOUT") {
			fmt.Println(msg...)
		} else if err != nil {
			panic(err)
		}
	}
}

// -------------------------------------------------
//     ___ _           _     ___
//    / __| |_  ___ __| |__ | __|_ _ _ _ ___ _ _
//   | (__| ' \/ -_) _| / / | _|| '_| '_/ _ \ '_|
//    \___|_||_\___\__|_\_\ |___|_| |_| \___/_|
func CheckErr(err error) (isErr bool) {
	defer func() {
		r := recover()
		if r != nil {
			Error("Error detected:", r)
		}
	}()
	isErr = false
	if err != nil {
		isErr = true
		a, b, c, _ := runtime.Caller(1)
		Error(err, " in ", "Process ID:", a, "In Module:", b, "Line:", c) //Return Error object
	}
	return
}

// -------------------------------------------------
//     ___ _           _  DB ___
//    / __| |_  ___ __| |__ | __|_ _ _ _ ___ _ _
//   | (__| ' \/ -_) _| / / | _|| '_| '_/ _ \ '_|
//    \___|_||_\___\__|_\_\ |___|_| |_| \___/_|
func CheckDbErr(err error, db *sql.DB, msg ...interface{}) (isErr bool) {
	defer func() {
		r := recover()
		if r != nil {
			Error("DB Error detected:", r)
		}
	}()
	isErr = false
	if err != nil {
		isErr = true
		a, b, c, _ := runtime.Caller(1)
		Error("DB", db.Ping(), err, "Process ID:", a, "In Module:", b, "Line:", c, msg) //Return the Database Error object
	}
	return
}

//         _
//        | |
//     ___| | ___  ___  ___
//    / __| |/ _ \/ __|/ _ \
//   | (__| | (_) \__ \  __/
//    \___|_|\___/|___/\___|
//  Safe close routine
func DeferClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		Crit(err)
	}
}

/*************************************************
 *            _ _  _____ _        _
 *           (_) |/ ____| |      (_)
 *      _ __  _| | (___ | |_ _ __ _ _ __   __ _
 *     | '_ \| | |\___ \| __| '__| | '_ \ / _` |
 *     | | | | | |____) | |_| |  | | | | | (_| |
 *     |_| |_|_|_|_____/ \__|_|  |_|_| |_|\__, |
 *                                         __/ |
 *                                        |___/
 * * * * * * * * * * * * * * * * * * * * * * * * *
 * prevents a nil string from infecting data
 * --------------------------------------------- */
func NilString(s string) string {
	if len(s) == 0 {
		return ""
	}
	return s
}

/*****************************************************
 *      _        _           _
 *     | |      (_)         | |
 *     | |_ _ __ _ _ __ ___ | |     ___ _ __
 *     | __| '__| | '_ ` _ \| |    / _ \ '_ \
 *     | |_| |  | | | | | | | |___|  __/ | | |
 *      \__|_|  |_|_| |_| |_|______\___|_| |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Trim the length of a string based on the max length
 * -------------------------------------------------- */
func TrimLen(str string, size int) (splited string) {
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	var start, stop int
	for i := 0; i < splitedLength; i += 1 {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited = str[start:stop]
	}
	return splited
}

/****************************************************************************************
 *                                _ _______    _____                        _
 *                               | |__   __|  |  __ \                      (_)
 *      _ __ ___ _ __   ___  _ __| |_ | | ___ | |  | | ___  _ __ ___   __ _ _ _ __
 *     | '__/ _ \ '_ \ / _ \| '__| __|| |/ _ \| |  | |/ _ \| '_ ` _ \ / _` | | '_ \
 *     | | |  __/ |_) | (_) | |  | |_ | | (_) | |__| | (_) | | | | | | (_| | | | | |
 *     |_|  \___| .__/ \___/|_|   \__||_|\___/|_____/ \___/|_| |_| |_|\__,_|_|_| |_|
 *              | |
 *              |_|
 * --------------------------------------------------------------------------------------
 * Function to report to domains, this will require paramters, since we can not predict
 * how each application is managing its configuration control.
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
type Services struct {
	Service  string `json:"service"`
	Dns      string `json:"dns"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Cb       string `json:"callback"`
	Expected string `json:"expected"`
	Deps     string `json:"deps"`
	Meta     string `json:"meta"`
}

/******************************************************************************
 *      _   _               _   _   _               _ _
 *     | | | |             | | | | | |        /\   | | |
 *     | |_| |__  _ __ ___ | |_| |_| | ___   /  \  | | | _____      __
 *     | __| '_ \| '__/ _ \| __| __| |/ _ \ / /\ \ | | |/ _ \ \ /\ / /
 *     | |_| | | | | | (_) | |_| |_| |  __// ____ \| | | (_) \ V  V /
 *      \__|_| |_|_|  \___/ \__|\__|_|\___/_/    \_\_|_|\___/ \_/\_/
 * ----------------------------------------------------------------------------
 * This function will temp store the value in a map and then remove it, it will
 * return true or false if the item is in the map, Now sets delay on second response
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
var throttle = make(map[string]bool)

func ThrottleAllow(ip string, timeout int) (retVal bool) {
	if throttle[ip] {
		Warn("-=Throttle=-To frequent calls from:", ip)
		time.Sleep(time.Duration(timeout) * time.Second) //Random next cycle.
		retVal = true                                    // false will result is receiging to frequent message
	} else {
		throttle[ip] = true
		go func() {
			time.Sleep(time.Duration(timeout) * time.Nanosecond) //Random next cycle.
			delete(throttle, ip)
		}()
		retVal = true
	}
	return
}

/********_____*_***********************_*_**********************************
 *      / ____| |                /\   | (_)
 *     | (___ | |_ __ _ _   _   /  \  | |___   _____
 *      \___ \| __/ _` | | | | / /\ \ | | \ \ / / _ \
 *      ____) | || (_| | |_| |/ ____ \| | |\ V /  __/
 *     |_____/ \__\__,_|\__, /_/    \_\_|_| \_/ \___|
 *                       __/ |
 *                      |___/
 * Call this routine to stay alive, but check for break condition to exit.
 * Paramter is a clean up function to call.                               */
var Alive bool

func StayAlive(fn func()) {
	Debug("Stay Alive Started")
	Alive = true
	for {
		time.Sleep(time.Second * 10) // Put it here so to handle all retries.
		if !Alive {
			fn()
			break
		}
	}
	Debug("Stay Alive Ended")
}

// ---___--_--_------------___-------_------_-----------
//   | __|(_)| | ___  ___ | __|__ __(_) ___| |_  ___
//   | _| | || |/ -_)|___|| _| \ \ /| |(_-<|  _|(_-<
//   |_|  |_||_|\___|     |___|/_\_\|_|/__/ \__|/__/
func Exists(filePath string) (exists bool) {
	_, err := os.Stat(filePath)
	if err != nil {
		exists = false
	} else {
		exists = true
	}
	return
}

//  ------------_---_---_-----------------_---_----------------_----------
//             | | | | | |               | | | |              | |
//    __ _  ___| |_| |_| | ___  __ _ _ __| |_| |__   ___  __ _| |_
//   / _` |/ _ \ __|  _  |/ _ \/ _` | '__| __| '_ \ / _ \/ _` | __|
//  | (_| |  __/ |_| | | |  __/ (_| | |  | |_| |_) |  __/ (_| | |_
//   \__, |\___|\__\_| |_/\___|\__,_|_|   \__|_.__/ \___|\__,_|\__|
//    __/ |
//   |___/
func NextHeartBeat(heartBeat int) (nextHeartBeat int) {
	nextHeartBeat = heartBeat
	nextHeartBeat = RandomRange(nextHeartBeat/2, nextHeartBeat*2)
	return
}

// ----------------------------------------------------------------------------
//      ____                  __                   ____
//     / __ \____ _____  ____/ /___  ____ ___     / __ \____ _____  ____ ____
//    / /_/ / __ `/ __ \/ __  / __ \/ __ `__ \   / /_/ / __ `/ __ \/ __ `/ _ \
//   / _, _/ /_/ / / / / /_/ / /_/ / / / / / /  / _, _/ /_/ / / / / /_/ /  __/
//  /_/ |_|\__,_/_/ /_/\__,_/\____/_/ /_/ /_/  /_/ |_|\__,_/_/ /_/\__, /\___/
//                                                               /____/
func RandomRange(min int, max int) (randomRange int) {
	rand.Seed(time.Now().UnixNano())
	randomRange = rand.Intn(max-min) + min
	return
}

// ---------------------------------------------------------------------------------
//   _____ _               _     _____                            _   _
//  / ____| |             | |   / ____|                          | | (_)
// | |    | |__   ___  ___| | _| |     ___  _ __  _ __   ___  ___| |_ _  ___  _ __
// | |    | '_ \ / _ \/ __| |/ / |    / _ \| '_ \| '_ \ / _ \/ __| __| |/ _ \| '_ \
// | |____| | | |  __/ (__|   <| |___| (_) | | | | | | |  __/ (__| |_| | (_) | | | |
// \_____|_| |_|\___|\___|_|\_\\_____\___/|_| |_|_| |_|\___|\___|\__|_|\___/|_| |_|
func CheckConnect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		Error("Error connecting to ", host, " port ", port)
		return false
	} else {
		if conn != nil {
			defer conn.Close()
			return true
		} else {
			return false
		}
	}
}

/*************************************************
 *                         _ __  __       _ _
 *                        | |  \/  |     (_) |
 *      ___  ___ _ __   __| | \  / | __ _ _| |
 *     / __|/ _ \ '_ \ / _` | |\/| |/ _` | | |
 *     \__ \  __/ | | | (_| | |  | | (_| | | |
 *     |___/\___|_| |_|\__,_|_|  |_|\__,_|_|_|
 * -----------------------------------------------
 * Sends an email based on config values
 * * * * * * * * * * * * * * * * * * * * * * * * */
type Mail struct {
	HOST    string
	PORT    int
	FROM    string
	PASS    string
	TO      string
	SUBJECT string
	BODY    string
	FILES   []string
}

func SendMail(mail Mail) {
	defer func() {
		r := recover()
		if r != nil {
			Error("Sending Email problem:", r)
		}
	}()
	addresses := strings.Split(mail.TO, ";")
	mess := gomail.NewMessage()
	mess.SetHeader("From", mail.FROM)       // Set E-Mail sender
	mess.SetHeader("To", addresses...)      // Set E-Mail receivers
	mess.SetHeader("Subject", mail.SUBJECT) // Set E-Mail subject
	mess.SetBody("text/html", mail.BODY)    // Set E-Mail body. plain text or html with text/html
	for i := range mail.FILES {
		mess.Attach(mail.FILES[i])
	}
	logCore("DEBUG", "Sending email")
	del := gomail.NewDialer(mail.HOST, mail.PORT, mail.FROM, mail.PASS) // Settings for SMTP server
	del.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	logCore("DEBUG", "Sent email:", mess)
	res := del.DialAndSend(mess)
	logCore("DEBUG", "Result:", res)
}

/*****************************************************************
 *      _  __ _______ _                ______ _
 *     (_)/ _|__   __| |              |  ____| |
 *      _| |_   | |  | |__   ___ _ __ | |__  | |___  ___
 *     | |  _|  | |  | '_ \ / _ \ '_ \|  __| | / __|/ _ \
 *     | | |    | |  | | | |  __/ | | | |____| \__ \  __/
 *     |_|_|    |_|  |_| |_|\___|_| |_|______|_|___/\___|
 * --------------------------------------------------------------
 * IfThenElse evaluates a condition, if true returns the first
 * parameter otherwise the second
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

/********************************************************************************************
 *                    _       _     ______ _ _                         _ _____
 *                   | |     | |   |  ____(_) |        /\             | |  __ \
 *     __      ____ _| |_ ___| |__ | |__   _| | ___   /  \   _ __   __| | |__) |   _ _ __
 *     \ \ /\ / / _` | __/ __| '_ \|  __| | | |/ _ \ / /\ \ | '_ \ / _` |  _  / | | | '_ \
 *      \ V  V / (_| | || (__| | | | |    | | |  __// ____ \| | | | (_| | | \ \ |_| | | | |
 *       \_/\_/ \__,_|\__\___|_| |_|_|    |_|_|\___/_/    \_\_| |_|\__,_|_|  \_\__,_|_| |_|
 * ------------------------------------------------------------------------------------------
 * Watches for a file change and if it change returns nil
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
func WatchFileAndRun(filePath string, fn func()) {
	defer func() {
		r := recover()
		if r != nil {
			Error("Watching file:", r)
		}
	}()
	fn()
	initialStat, err := os.Stat(filePath)
	CheckErr(err)
	for {
		stat, err := os.Stat(filePath)
		CheckErr(err)
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			logCore("INFO", "File Changed:", stat.Name())
			fn()
			initialStat, err = os.Stat(filePath)
			if CheckErr(err) {
				panic(err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

/******************************************************************
 *      _    _ _______ __  __ _      _____ _
 *     | |  | |__   __|  \/  | |    / ____| |
 *     | |__| |  | |  | \  / | |   | |    | | ___  __ _ _ __
 *     |  __  |  | |  | |\/| | |   | |    | |/ _ \/ _` | '_ \
 *     | |  | |  | |  | |  | | |___| |____| |  __/ (_| | | | |
 *     |_|  |_|  |_|  |_|  |_|______\_____|_|\___|\__,_|_| |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Remove all HTML tags and character codes.
 * ------------------------------------------------------------- */
func HTMLClean(htmlString string) (pureText string) {
	defer func() {
		r := recover()
		if r != nil {
			logCore("INFO", "Error Stripping HTML:", r)
		}
	}()
	p := bluemonday.StrictPolicy()
	cleanText := p.Sanitize(htmlString)
	pureText = html.UnescapeString(cleanText)
	return
}

func GetImagePath(html string) string {
	re := regexp.MustCompile(`\bhttps?:[^)''"]+\.(?:jpg|jpeg|gif|png)\?{0,1}[a-zA-Z0-9\=]*`)
	str := string(re.Find([]byte(html)))
	return str
}
