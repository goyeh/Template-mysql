package sql

import (
	"all-news/conf"
	"all-news/lib"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	_ "github.com/goyeh/mysql"
)

type Article struct {
	Uid     int64         `json:"uid"`
	Title   string        `json:"title"`
	Content string        `json:"desc"`
	Author  string        `json:"author"`
	Email   string        `json:"email"`
	Topic   string        `json:"topic"`
	Cat     string        `json:"cat"`
	Link    string        `json:"link"`
	Detail  ArticleDetail `json:"detail"`
	Rating  int64         `json:"rating"`
	Created time.Time     `json:"created"`
}

type ArticleDetail map[string]interface{}

func (a *ArticleDetail) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error marshaling ArticleDetail to JSON: %v", err)
	}
	return jsonBytes, nil
}

func (a *ArticleDetail) Scan(value interface{}) error {
	json.Unmarshal(value.([]byte), a)
	return nil
}

type HomeMadeArticle struct {
	Title   string `json:"title"`
	Content string `json:"desc"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Cat     string `json:"cat"`
	Link    string `json:"link"`
	Image   string `json:"image"`
}

type AccountsStruct struct {
	Email     string
	Otpkey    string
	Plan      string
	Allocated int64
	Used      int64
	End       time.Time
}

var (
	db       *sql.DB
	Accounts map[string]AccountsStruct
)

func InitDB() {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Checking Database Connection:", r)
		}
	}()

	db = openDB()
	lib.Debug("Check DB connection:", db.Ping())

	LoadAccounts()
}

// TODO save accounts to Database
func LoadAccounts() {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("LoadAccounts from Database:", r)
		}
	}()

	rows, err := db.Query("SELECT apikey, email, otpkey, plan, allocated, used, end FROM accounts")
	if lib.CheckErr(err) {
		panic(err)
	}
	defer rows.Close()
	Accounts = make(map[string]AccountsStruct) // reset the map. fresh reload.

	for rows.Next() {
		var d AccountsStruct
		var apikey string
		err = rows.Scan(&apikey, &d.Email, &d.Otpkey, &d.Plan, &d.Allocated, &d.Used, &d.End)
		if lib.CheckErr(err) {
			panic(err)
		}

		Accounts[apikey] = AccountsStruct{
			Email:     d.Email,
			Otpkey:    d.Otpkey,
			Plan:      d.Plan,
			Allocated: d.Allocated,
			Used:      d.Used,
			End:       d.End,
		}
	}

	return
}

/*******************************************************************
 *      _                     _                 _   _      _
 *     (_)                   | |     /\        | | (_)    | |
 *      _ _ __  ___  ___ _ __| |_   /  \   _ __| |_ _  ___| | ___
 *     | | '_ \/ __|/ _ \ '__| __| / /\ \ | '__| __| |/ __| |/ _ \
 *     | | | | \__ \  __/ |  | |_ / ____ \| |  | |_| | (__| |  __/
 *     |_|_| |_|___/\___|_|   \__/_/    \_\_|   \__|_|\___|_|\___|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Inserts new articles into the datase, Rejects duplicates
 * ---------------------------------------------------------------- */
func InsertArticle(topic string, title string, content string, name string, email string, cat string, link string, image string) (id int64, returnErr error) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Error Inserting Articles into Database:", r)
		}
	}()

	id = 0
	now := time.Now().Format("2006-01-02 15:04:05.0000")
	sqlString := `INSERT INTO articles (title,content,author,email,topic,cat,link,detail,created) values(?,?,?,?,?,?,?,?,?);` // ON DUPLICATE KEY UPDATE rating = rating + 1`
	var detail string
	jsonByte, _ := json.Marshal(map[string]interface{}{"img": image})
	detail = string(jsonByte)

	res, err := db.Exec(sqlString, lib.TrimLen(title, 128), lib.NilString(content), lib.NilString(name), lib.NilString(email), lib.NilString(topic), lib.NilString(lib.TrimLen(cat, 512)), lib.NilString(link), lib.NilString(detail), now)
	if err == nil {
		id, _ = res.RowsAffected()
		lastId, _ := res.LastInsertId()
		lib.Info("Insert general Completed...", "id:", lastId, " rows:", id)
	} else {
		if !strings.Contains(err.Error(), "Duplicate") {
			lib.Info("Insert Statement result...", err, "\nDEBUG:", sqlString, res)
		}
		returnErr = err
	}
	return
}

/*********************************************************************************
 *           _               _     _____            _             _
 *          | |             | |   / ____|          | |           | |
 *       ___| |__   ___  ___| | _| |     ___  _ __ | |_ _ __ ___ | |
 *      / __| '_ \ / _ \/ __| |/ / |    / _ \| '_ \| __| '__/ _ \| |
 *     | (__| | | |  __/ (__|   <| |___| (_) | | | | |_| | | (_) | |
 *      \___|_| |_|\___|\___|_|\_\\_____\___/|_| |_|\__|_|  \___/|_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * get the last article in the database.
 * ----------------------------------------------------------------------------- */
func CheckControl(target string, platform string) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("check control error:", r)
		}
	}()
	sqlStmt := fmt.Sprintf("SELECT count(*) from control where target = '%s';", target)
	if CountSQL(sqlStmt) < 1 {
		sqlStmt := fmt.Sprintf("INSERT INTO control (target,timestamp,lastupdate, platform, note) VALUES ('%[1]s','%[2]v','%[3]v', '%[4]s', '%[5]s');", target, time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05.0000"), time.Now().Format("2006-01-02 15:04:05.0000"), platform, "Created on first Touch")
		retval := RunSQL(sqlStmt)
		if retval < 1 {
			lib.Error("Control Table problem:", sqlStmt)
			panic(retval)
		}
		lib.Info("Control Record Inserted:", retval)
	}
}

/******************_*****************_***_******_***********
 *                | |     /\        | | (_)    | |
 *       __ _  ___| |_   /  \   _ __| |_ _  ___| | ___
 *      / _` |/ _ \ __| / /\ \ | '__| __| |/ __| |/ _ \
 *     | (_| |  __/ |_ / ____ \| |  | |_| | (__| |  __/
 *      \__, |\___|\__/_/    \_\_|   \__|_|\___|_|\___|
 *       __/ |
 *      |___/                                             */
func getArticle(sqlString string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Getting Articles:", r)
		}
	}()

	var aList []Article
	reply = conf.Reply{Code: 206, Msg: "[]"}
	rowCount := 0

	lib.Debug("Getting Data:", sqlString)
	rows, err := db.Query(sqlString)
	lib.CheckErr(err)
	for rows.Next() {
		var a Article
		lib.Debug("Counting", rowCount)
		switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
		case sql.ErrNoRows:
			lib.Error("No rows retrieved", err)
			reply.Code = 204
			reply.Msg = "[]"
		case nil:
			rowCount++
			lib.Debug("Result:", a)
			aList = append(aList, a)
		default:
			reply.Code = 400
			reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
		}
	}
	if rowCount > 0 {
		jsonReply, err := json.MarshalIndent(aList, "", "")
		lib.Debug("Reply:", jsonReply)
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/*********************************************************************************
 *                 _   _               _                 _   _      _
 *                | | | |             | |     /\        | | (_)    | |
 *       __ _  ___| |_| |     __ _ ___| |_   /  \   _ __| |_ _  ___| | ___
 *      / _` |/ _ \ __| |    / _` / __| __| / /\ \ | '__| __| |/ __| |/ _ \
 *     | (_| |  __/ |_| |___| (_| \__ \ |_ / ____ \| |  | |_| | (__| |  __/
 *      \__, |\___|\__|______\__,_|___/\__/_/    \_\_|   \__|_|\___|_|\___|
 *       __/ |
 *      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * get the last article in the database.
 * ----------------------------------------------------------------------------- */
func GetLastSingleJsonArticle(topic string, limit int) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Getting last Articles:", r)
		}
	}()
	sqlArticles := fmt.Sprintf(`SELECT * FROM articles WHERE topic='%[2]s' ORDER BY uid DESC LIMIT %[1]v;`, limit, topic)
	lib.Debug("Getting Data:", sqlArticles)
	reply = getArticle(sqlArticles)
	return
}

/*******************************************************************************
 *                 _   _   _           _                 _   _      _
 *                | | | \ | |         | |     /\        | | (_)    | |
 *       __ _  ___| |_|  \| | _____  _| |_   /  \   _ __| |_ _  ___| | ___
 *      / _` |/ _ \ __| . ` |/ _ \ \/ / __| / /\ \ | '__| __| |/ __| |/ _ \
 *     | (_| |  __/ |_| |\  |  __/>  <| |_ / ____ \| |  | |_| | (__| |  __/
 *      \__, |\___|\__|_| \_|\___/_/\_\\__/_/    \_\_|   \__|_|\___|_|\___|
 *       __/ |
 *      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID from the control
 * ---------------------------------------------------------------------------- */
func GetNextArticle(callerId string, platform string, topic string) (message string, Behind int) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Get next article:", r)
		}
	}()
	var a Article
	CheckControl(callerId, platform)
	sqlCount := fmt.Sprintf(`SELECT count(*) FROM articles WHERE topic = '%[2]s' AND created > (SELECT timestamp FROM control WHERE target = '%[1]s') ;`, callerId, topic)
	rowCount := db.QueryRow(sqlCount)
	switch err := rowCount.Scan(&Behind); err {
	case sql.ErrNoRows:
		lib.Warn("No rows.", err)
	case nil: // No errors, and has rows

		sqlArticle := fmt.Sprintf(`SELECT * FROM articles WHERE topic = '%[2]s' AND created > (SELECT timestamp FROM control WHERE target = '%[1]s') ORDER BY created LIMIT 1 ;`, callerId, topic)
		row := db.QueryRow(sqlArticle)
		switch err := row.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
		case sql.ErrNoRows:
			lib.Warn("No rows, adding First record.", err)
		case nil:
			sqlStatement := fmt.Sprintf("UPDATE control SET timestamp = '%[1]s' WHERE target = '%[2]s';", a.Created.Format("2006-01-02 15:04:05.0000"), callerId)
			lib.Debug("Update Control:", sqlStatement)
			RunSQL(sqlStatement)
			if conf.NEWSDETAIL {
				message = fmt.Sprintf("*%[1]s*\n _%[2]s_ [%[3]s]", a.Title, a.Content, a.Link, a.Rating)
			} else {
				message = fmt.Sprintf("*%[1]s*\n [%[2]s]", a.Title, a.Link)
			}
			lib.Debug(message)
		default: //This should never happen, but left in just in case.
			lib.Debug("Panic?:", a)
			panic(err)
		}

	default:
		lib.Debug("Panic?:", a)
		panic(err)
	}

	return
}

/*******************************************************************************
 *                 _   _   _           _                 _   _      _      ____        _  __
 *                | | | \ | |         | |     /\        | | (_)    | |    |  _ \      | |/ /
 *       __ _  ___| |_|  \| | _____  _| |_   /  \   _ __| |_ _  ___| | ___| |_) |_   _| ' / ___ _   _
 *      / _` |/ _ \ __| . ` |/ _ \ \/ / __| / /\ \ | '__| __| |/ __| |/ _ \  _ <| | | |  < / _ \ | | |
 *     | (_| |  __/ |_| |\  |  __/>  <| |_ / ____ \| |  | |_| | (__| |  __/ |_) | |_| | . \  __/ |_| |
 *      \__, |\___|\__|_| \_|\___/_/\_\\__/_/    \_\_|   \__|_|\___|_|\___|____/ \__, |_|\_\___|\__, |
 *       __/ |                                                                    __/ |          __/ |
 *      |___/                                                                    |___/          |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID from the control
 * ---------------------------------------------------------------------------- */
func GetNextArticleByKeyWords(callerId string, keyword string, platform string, topic string) (message string) {
	var a Article
	CheckControl(callerId, platform)
	keyword = strings.ReplaceAll(keyword, " ", "|")
	sqlArticle := fmt.Sprintf(`SELECT * FROM articles WHERE topic = '%[2]s' created > (SELECT timestamp FROM control WHERE target = '%[1]s') ORDER BY created LIMIT 1 ;`, callerId, topic)
	if len(keyword) > 3 {
		sqlArticle = fmt.Sprintf(`SELECT * FROM articles WHERE topic = '%[3]s' title rlike '%[2]s' AND created > (SELECT timestamp FROM control WHERE target = '%[1]s') ORDER BY created LIMIT 1 ;`, callerId, keyword, topic)
	}
	lib.Debug("Collecting:", sqlArticle, "Message:", message)
	row := db.QueryRow(sqlArticle)
	switch err := row.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
	case sql.ErrNoRows:
		lib.Warn("No rows, adding First record.", err)
	case nil:
		sqlStatement := fmt.Sprintf("UPDATE control SET timestamp = '%[1]s' WHERE target = '%[2]s';", a.Created.Format("2006-01-02 15:04:05.0000"), callerId)
		lib.Debug("Update Control:", sqlStatement)
		RunSQL(sqlStatement)
		if conf.NEWSDETAIL {
			message = fmt.Sprintf("*%[1]s*\n _%[2]s_ [%[3]s]", a.Title, a.Content, a.Link)
		} else {
			message = fmt.Sprintf("*%[1]s*\n [%[2]s]", a.Title, a.Link)
		}
		lib.Debug(message)
	default:
		lib.Debug("Panic?:", a)
		panic(err)
	}
	return
}

/*****************************************************************************************************
 *                 _   _   _           _       _                               _   _      _
 *                | | | \ | |         | |     | |                   /\        | | (_)    | |
 *       __ _  ___| |_|  \| | _____  _| |_    | |___  ___  _ __    /  \   _ __| |_ _  ___| | ___
 *      / _` |/ _ \ __| . ` |/ _ \ \/ / __|   | / __|/ _ \| '_ \  / /\ \ | '__| __| |/ __| |/ _ \
 *     | (_| |  __/ |_| |\  |  __/>  <| || |__| \__ \ (_) | | | |/ ____ \| |  | |_| | (__| |  __/
 *      \__, |\___|\__|_| \_|\___/_/\_\\__\____/|___/\___/|_| |_/_/    \_\_|   \__|_|\___|_|\___|
 *       __/ |
 *      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------ */
func GetLatestJsonArticle(callerId string, limit int, next bool, platform string, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Can not get the next Article:", r)
		}
	}()
	var aList []Article
	CheckControl(callerId, platform)
	rowCount := 0
	sqlArticles := ""
	if next {
		sqlArticles = fmt.Sprintf(`SELECT * FROM articles WHERE created > (SELECT timestamp FROM control WHERE target = '%[1]s' and topic = '%[3]s') ORDER BY created LIMIT %[2]v ;`, callerId, limit, topic)
	} else {
		sqlArticles = fmt.Sprintf(`SELECT * FROM articles WHERE created < (SELECT timestamp FROM control WHERE target = '%[1]s' and topic = '%[3]s') ORDER BY created DESC LIMIT %[2]v ;`, callerId, limit, topic)
	}
	rows, err := db.Query(sqlArticles)
	lib.CheckErr(err)
	for rows.Next() {
		var a Article
		lib.Debug("Counting", rowCount)
		switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
		case sql.ErrNoRows:
			lib.Error("No rows retrieved", err)
			reply.Code = 204
			reply.Msg = "[]"
		case nil:
			rowCount++
			lib.Debug("Result:", a)
			aList = append(aList, a)
		default:
			reply.Code = 400
			reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
		}
	}
	if rowCount < 1 {
		reply.Code = 206
		reply.Msg = "[]"
	} else {
		lib.Debug("Rows:", rowCount, aList)
		sqlStatement := fmt.Sprintf("UPDATE control SET timestamp = '%[1]s' WHERE target = '%[2]s';", aList[0].Created.Format("2006-01-02 15:04:05.0000"), callerId)
		lib.Debug("Update Control:", sqlStatement)
		RunSQL(sqlStatement)
		jsonReply, err := json.MarshalIndent(aList, "", "")
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/*****************************************************************************************************
 *                 _       _                 ____        _____    _
 *                | |     | |               |  _ \      |_   _|  | |
 *       __ _  ___| |_    | |___  ___  _ __ | |_) |_   _  | |  __| |
 *      / _` |/ _ \ __|   | / __|/ _ \| '_ \|  _ <| | | | | | / _` |
 *     | (_| |  __/ || |__| \__ \ (_) | | | | |_) | |_| |_| || (_| |
 *      \__, |\___|\__\____/|___/\___/|_| |_|____/ \__, |_____\__,_|
 *       __/ |                                      __/ |
 *      |___/                                      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------ */
func GetJsonByArticleId(articleId string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve Article:", r, articleId)
		}
	}()
	var aList []Article
	rowCount := 0
	sqlArticles := ""
	sqlArticles = fmt.Sprintf(`SELECT * FROM articles WHERE uid = %[1]s ;`, articleId)
	rows, err := db.Query(sqlArticles)
	lib.CheckErr(err)
	for rows.Next() {
		var a Article
		lib.Debug("Counting", rowCount)
		switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
		case sql.ErrNoRows:
			lib.Error("No rows retrieved", err)
			reply.Code = 204
			reply.Msg = "[]"
		case nil:
			rowCount++
			lib.Debug("Result:", a)
			aList = append(aList, a)
		default:
			reply.Code = 400
			reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
		}
	}
	if rowCount < 1 {
		reply.Code = 206
		reply.Msg = "[]"
	} else {
		lib.Debug("Rows:", rowCount, aList)
		jsonReply, err := json.MarshalIndent(aList, "", "")
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/***********************************************************************************************************************
 *                 _   _   _           _       _                 ____                      _   _      _      _____    _
 *                | | | \ | |         | |     | |               |  _ \          /\        | | (_)    | |    |_   _|  | |
 *       __ _  ___| |_|  \| | _____  _| |_    | |___  ___  _ __ | |_) |_   _   /  \   _ __| |_ _  ___| | ___  | |  __| |
 *      / _` |/ _ \ __| . ` |/ _ \ \/ / __|   | / __|/ _ \| '_ \|  _ <| | | | / /\ \ | '__| __| |/ __| |/ _ \ | | / _` |
 *     | (_| |  __/ |_| |\  |  __/>  <| || |__| \__ \ (_) | | | | |_) | |_| |/ ____ \| |  | |_| | (__| |  __/_| || (_| |
 *      \__, |\___|\__|_| \_|\___/_/\_\\__\____/|___/\___/|_| |_|____/ \__, /_/    \_\_|   \__|_|\___|_|\___|_____\__,_|
 *       __/ |                                                          __/ |
 *      |___/                                                          |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------------------------- */
func GetNextJsonByArticleId(articleId string, limit int, filter string, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve Article:", r, articleId)
		}
	}()
	var a Article
	var aList []Article
	rowCount := 0
	articleUid, err := strconv.Atoi(articleId)
	if lib.CheckErr(err) {
		panic(err)
	} else {
		sqlArticles := fmt.Sprintf(`SELECT * FROM articles WHERE uid > %[1]v and topic = $[3]s ORDER BY uid ASC LIMIT %[2]v ;`, articleUid, limit, topic)
		if len(filter) > 0 {
			filter = strings.ReplaceAll(filter, ",", "|")
			sqlArticles = fmt.Sprintf(`SELECT * FROM articles WHERE uid > %[1]v AND topic = '%[4]s' AND content rlike '%[3]s' ORDER BY uid ASC LIMIT %[2]v ;`, articleUid, limit, filter, topic)
		}
		rows, err := db.Query(sqlArticles)
		lib.CheckErr(err)
		for rows.Next() {
			lib.Debug("Counting", rowCount)
			switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
			case sql.ErrNoRows:
				lib.Error("No rows retrieved", err)
				reply.Code = 204
				reply.Msg = "[]"
			case nil:
				rowCount++
				lib.Debug("Result:", a)
				aList = append(aList, a)
			default:
				reply.Code = 400
				reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
			}
		}
	}
	if rowCount < 1 {
		reply.Code = 206
		reply.Msg = "[]"
	} else {
		lib.Debug("Rows:", rowCount, aList)
		jsonReply, err := json.MarshalIndent(aList, "", "")
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/************************************************************************************************************************
 *                 _   _____                    _                 ____                      _   _      _      _____    _
 *                | | |  __ \                  | |               |  _ \          /\        | | (_)    | |    |_   _|  | |
 *       __ _  ___| |_| |__) | __ _____   __   | |___  ___  _ __ | |_) |_   _   /  \   _ __| |_ _  ___| | ___  | |  __| |
 *      / _` |/ _ \ __|  ___/ '__/ _ \ \ / /   | / __|/ _ \| '_ \|  _ <| | | | / /\ \ | '__| __| |/ __| |/ _ \ | | / _` |
 *     | (_| |  __/ |_| |   | | |  __/\ V / |__| \__ \ (_) | | | | |_) | |_| |/ ____ \| |  | |_| | (__| |  __/_| || (_| |
 *      \__, |\___|\__|_|   |_|  \___| \_/ \____/|___/\___/|_| |_|____/ \__, /_/    \_\_|   \__|_|\___|_|\___|_____\__,_|
 *       __/ |                                                           __/ |
 *      |___/                                                           |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------------------------- */
func GetPrevJsonByArticleId(articleId string, limit int, filter string, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve Article:", r, articleId)
		}
	}()
	var a Article
	var aList []Article
	rowCount := 0
	articleUid, err := strconv.Atoi(articleId)
	if lib.CheckErr(err) {
		panic(err)
	} else {
		sqlArticles := fmt.Sprintf(`SELECT * FROM articles WHERE uid < %[1]v AND topic = '%[3]s' ORDER BY uid DESC LIMIT %[2]v ;`, articleUid, limit, topic)
		if len(filter) > 0 {
			filter = strings.ReplaceAll(filter, ",", "|")
			sqlArticles = fmt.Sprintf(`SELECT * FROM articles WHERE uid < %[1]v AND topic = '%[4]s' AND content rlike '%[3]s' ORDER BY uid DESC LIMIT %[2]v ;`, articleUid, limit, filter, topic)
		}
		rows, err := db.Query(sqlArticles)
		lib.CheckErr(err)
		for rows.Next() {
			lib.Debug("Counting", rowCount)
			switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
			case sql.ErrNoRows:
				lib.Error("No rows retrieved", err)
				reply.Code = 204
				reply.Msg = "[]"
			case nil:
				rowCount++
				lib.Debug("Result:", a)
				aList = append(aList, a)
			default:
				reply.Code = 400
				reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
			}
		}
	}

	if rowCount < 1 {
		reply.Code = 206
		reply.Msg = "[]"
	} else {
		lib.Debug("Rows:", rowCount, aList)
		jsonReply, err := json.MarshalIndent(aList, "", "")
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/*****************************************************************************************************
 *                 _       _                 ____     _______ _ _   _
 *                | |     | |               |  _ \   |__   __(_) | | |
 *       __ _  ___| |_    | |___  ___  _ __ | |_) |_   _| |   _| |_| | ___
 *      / _` |/ _ \ __|   | / __|/ _ \| '_ \|  _ <| | | | |  | | __| |/ _ \
 *     | (_| |  __/ || |__| \__ \ (_) | | | | |_) | |_| | |  | | |_| |  __/
 *      \__, |\___|\__\____/|___/\___/|_| |_|____/ \__, |_|  |_|\__|_|\___|
 *       __/ |                                      __/ |
 *      |___/                                      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------ */
func GetJsonByTitle(search string, limit int, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve Article through search:", r, search)
		}
	}()
	var aList []Article
	safeStr, err := url.QueryUnescape(search)
	if !lib.CheckErr(err) {
		searchStr := strings.ReplaceAll(safeStr, " ", "|")
		rowCount := 0
		sqlArticles := fmt.Sprintf(`SELECT * FROM articles WHERE topic = '%[3]s' title rlike '%[1]s' ORDER BY created DESC LIMIT %[2]v ;`, searchStr, limit, topic)
		lib.Debug("Search String:", sqlArticles)
		rows, err := db.Query(sqlArticles)
		lib.CheckErr(err)
		for rows.Next() {
			var a Article
			lib.Debug("Counting", rowCount)
			switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
			case sql.ErrNoRows:
				lib.Error("No rows retrieved", err)
				reply.Code = 204
				reply.Msg = `[]`
			case nil:
				rowCount++
				lib.Debug("Result:", a)
				aList = append(aList, a)
			default:
				reply.Code = 400
				reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
			}
		}
		if rowCount < 1 {
			reply.Code = 206
			reply.Msg = "[]"
		} else {
			lib.Debug("Rows:", rowCount, aList)
			jsonReply, err := json.MarshalIndent(aList, "", "")
			if !lib.CheckErr(err) {
				reply.Code = 200
				reply.Msg = string(jsonReply)
			}
		}
	} else {
		reply.Msg = "Unable to decode the search string"
		reply.Code = 400
	}

	return
}

/*****************************************************************************************************
 *                 _       _                 ____         _____            _             _
 *                | |     | |               |  _ \       / ____|          | |           | |
 *       __ _  ___| |_    | |___  ___  _ __ | |_) |_   _| |     ___  _ __ | |_ ___ _ __ | |_
 *      / _` |/ _ \ __|   | / __|/ _ \| '_ \|  _ <| | | | |    / _ \| '_ \| __/ _ \ '_ \| __|
 *     | (_| |  __/ || |__| \__ \ (_) | | | | |_) | |_| | |___| (_) | | | | ||  __/ | | | |_
 *      \__, |\___|\__\____/|___/\___/|_| |_|____/ \__, |\_____\___/|_| |_|\__\___|_| |_|\__|
 *       __/ |                                      __/ |
 *      |___/                                      |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------ */
func GetJsonByContent(searchStr string, limit int, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve Article through Content search:", r, searchStr)
		}
	}()
	var aList []Article
	safeStr, err := url.QueryUnescape(searchStr)
	if !lib.CheckErr(err) {
		searchStr := strings.ReplaceAll(safeStr, " ", "|")
		rowCount := 0
		sqlArticles := fmt.Sprintf(`SELECT * FROM topic WHERE topic = '%[3]s' cat rlike '%[1]s' ORDER BY created DESC LIMIT %[2]v ;`, searchStr, limit, topic)
		lib.Debug("Search String:", sqlArticles)
		rows, err := db.Query(sqlArticles)
		lib.CheckErr(err)
		for rows.Next() {
			var a Article
			lib.Debug("Counting", rowCount)
			switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
			case sql.ErrNoRows:
				lib.Error("No rows retrieved", err)
				reply.Code = 204
				reply.Msg = "[]"
			case nil:
				rowCount++
				lib.Debug("Result:", a)
				aList = append(aList, a)
			default:
				reply.Code = 400
				reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
			}
		}
		if rowCount < 1 {
			reply.Code = 206
			reply.Msg = "[]"
		} else {
			lib.Debug("Rows:", rowCount, aList)
			jsonReply, err := json.MarshalIndent(aList, "", "")
			if !lib.CheckErr(err) {
				reply.Code = 200
				reply.Msg = string(jsonReply)
			}
		}
	} else {
		reply.Msg = "Unable to decode the search string"
		reply.Code = 400
	}

	return
}

/*****************************************************************************************************
 *                 _   _               _       _                 ____        ______ _ _ _
 *                | | | |             | |     | |               |  _ \      |  ____(_) | |
 *       __ _  ___| |_| |     __ _ ___| |_    | |___  ___  _ __ | |_) |_   _| |__   _| | |_ ___ _ __
 *      / _` |/ _ \ __| |    / _` / __| __|   | / __|/ _ \| '_ \|  _ <| | | |  __| | | | __/ _ \ '__|
 *     | (_| |  __/ |_| |___| (_| \__ \ || |__| \__ \ (_) | | | | |_) | |_| | |    | | | ||  __/ |
 *      \__, |\___|\__|______\__,_|___/\__\____/|___/\___/|_| |_|____/ \__, |_|    |_|_|\__\___|_|
 *       __/ |                                                          __/ |
 *      |___/                                                          |___/
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Gets the next article depending on the last Article ID and return as a json object
 * ------------------------------------------------------------------------------------------------ */
func GetLastJsonByFilter(limit int, filter string, topic string) (reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve filtered Articles:", r, filter)
		}
	}()
	var aList []Article
	rowCount := 0
	sqlArticles := fmt.Sprintf(`SELECT * FROM %[2]s ORDER BY uid DESC LIMIT %[1]v ;`, limit, topic)
	if len(filter) > 0 {
		filter = strings.ReplaceAll(filter, ",", "|")
		sqlArticles = fmt.Sprintf(`SELECT * FROM article WHERE topic = '%[3]s' content rlike '%[2]s' ORDER BY uid DESC LIMIT %[1]v ;`, limit, filter, topic)
	}
	rows, err := db.Query(sqlArticles)
	lib.CheckErr(err)
	for rows.Next() {
		var a Article
		lib.Debug("Counting", rowCount)
		switch err := rows.Scan(&a.Uid, &a.Title, &a.Content, &a.Author, &a.Email, &a.Topic, &a.Cat, &a.Link, &a.Detail, &a.Rating, &a.Created); err {
		case sql.ErrNoRows:
			lib.Error("No rows retrieved", err)
			reply.Code = 204
			reply.Msg = "[]"
		case nil:
			rowCount++
			lib.Debug("Result:", a)
			aList = append(aList, a)
		default:
			reply.Code = 400
			reply.Msg = fmt.Sprintf(`Error: %v`, err.Error())
		}
	}
	if rowCount < 1 {
		reply.Code = 200
		reply.Msg = "[]"
	} else {
		lib.Debug("Rows:", rowCount, aList)
		jsonReply, err := json.MarshalIndent(aList, "", "")
		if !lib.CheckErr(err) {
			reply.Code = 200
			reply.Msg = string(jsonReply)
		}
	}
	return
}

/***************************************************************************************
 *                 _   _      _     _    ____   __ _______                   _
 *                | | | |    (_)   | |  / __ \ / _|__   __|                 | |
 *       __ _  ___| |_| |     _ ___| |_| |  | | |_   | | __ _ _ __ __ _  ___| |_ ___
 *      / _` |/ _ \ __| |    | / __| __| |  | |  _|  | |/ _` | '__/ _` |/ _ \ __/ __|
 *     | (_| |  __/ |_| |____| \__ \ |_| |__| | |    | | (_| | | | (_| |  __/ |_\__ \
 *      \__, |\___|\__|______|_|___/\__|\____/|_|    |_|\__,_|_|  \__, |\___|\__|___/
 *       __/ |                                                     __/ |
 *      |___/                                                     |___/
 * -------------------------------------------------------------------------------------
 * List all targets that match the platform and return array
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
var thisTarget string

// TODO check this function
func GetListOfTargets(platform string) (aList []string) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Unable to retrieve List of Targets:", r, platform)
		}
	}()
	sqlString := fmt.Sprintf(`SELECT target FROM control WHERE platform = '%[1]s' AND live = 1 ;`, platform)
	rows, err := db.Query(sqlString)
	lib.CheckErr(err)
	lib.Debug("SQL String:", sqlString)
	for rows.Next() {
		switch err := rows.Scan(&thisTarget); err {
		case sql.ErrNoRows:
			lib.Info("No rows retrieved", err)
		case nil:
			lib.Debug("Result:", thisTarget)
			aList = append(aList, thisTarget)
		default:
			lib.Error("Nothing Retreived:", err)
		}
	}
	return
}

/****************************************************************************
 *                _  _____ _        _                            _
 *               | |/ ____| |      | |                          | |
 *      ___  __ _| | (___ | |_ __ _| |_ ___ _ __ ___   ___ _ __ | |_
 *     / __|/ _` | |\___ \| __/ _` | __/ _ \ '_ ` _ \ / _ \ '_ \| __|
 *     \__ \ (_| | |____) | || (_| | ||  __/ | | | | |  __/ | | | |_
 *     |___/\__, |_|_____/ \__\__,_|\__\___|_| |_| |_|\___|_| |_|\__|
 *             | |
 *             |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Run an SQL statement such as an insert.
 * ------------------------------------------------------------------------ */
func RunSQL(sqlStatement string) (id int64) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Warn("Creating Schema:", r)
		}
	}()
	id = -1
	stmt, err := db.Prepare(sqlStatement)
	if lib.CheckErr(err) {
		lib.Warn("Unable to Prepare:", err, sqlStatement)
	} else {
		res, err := stmt.Exec()
		if lib.CheckErr(err) {
			lib.Warn("Unable to Execute:", err, sqlStatement)
		} else {
			id, err = res.RowsAffected()
			lib.CheckErr(err)
			lib.Info("Execution Complete:")
		}
	}
	return
}

/*********countSQL*******************************************************************
 *                            _    _____  ____  _
 *                           | |  / ____|/ __ \| |
 *       ___ ___  _   _ _ __ | |_| (___ | |  | | |
 *      / __/ _ \| | | | '_ \| __|\___ \| |  | | |
 *     | (_| (_) | |_| | | | | |_ ____) | |__| | |____
 *      \___\___/ \__,_|_| |_|\__|_____/ \___\_\______|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Count records based on the SQL request.
 * ------------------------------------------------------------------------ */
func CountSQL(sqlStatement string) (id int64) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Warn("Counting Rows:", r)
		}
	}()
	err := db.QueryRow(sqlStatement).Scan(&id)
	switch {
	case err != nil:
		lib.Debug(err, sqlStatement)
		panic(err)
	default:
		lib.Debug("Total rows:", id)
	}
	return
}

/************************************************************************************
 *       ____                         __  __        _____       _ _____  ____
 *      / __ \                       |  \/  |      / ____|     | |  __ \|  _ \
 *     | |  | | ___  _ __   ___ _ __ | \  / |_   _| (___   __ _| | |  | | |_) |
 *     | |  | |/ _ \| '_ \ / _ \ '_ \| |\/| | | | |\___ \ / _` | | |  | |  _ <
 *     | |__| | (_) | |_) |  __/ | | | |  | | |_| |____) | (_| | | |__| | |_) |
 *      \____/ \___/| .__/ \___|_| |_|_|  |_|\__, |_____/ \__, |_|_____/|____/
 *                  | |                       __/ |          | |
 *                  |_|                      |___/           |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Open Database connection
 * --------------------------------------------------------------------------------*/
func openDB() (dbSource *sql.DB) {
	defer func() {
		r := recover()
		if r != nil {
			log.Print("Possible error:", dbSource, r)
		}
	}()
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.MYSQL_USER, conf.MYSQL_PASS, conf.MYSQL_HOST, conf.MYSQL_PORT, conf.MYSQL_DB)
	dbSource, err := sql.Open("mysql", connectStr)
	if lib.CheckErr(err) {
		lib.Error(dbSource.Ping())
	}
	return
}
