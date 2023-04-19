package route

import (
	"all-news/conf"
	"all-news/lib"
	"all-news/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

var (
	limiter *rate.Limiter
)

/***************************************
 *      _____       _ _    ____
 *     |_   _|     (_) |  / /\ \
 *       | |  _ __  _| |_| |  | |
 *       | | | '_ \| | __| |  | |
 *      _| |_| | | | | |_| |  | |
 *     |_____|_| |_|_|\__| |  | |
 *                        \_\/_/   *  */
func Init() {
	limiter = rate.NewLimiter(rate.Limit(2), 5) //TODO add these values to the config
	router := fasthttprouter.New()
	server := &fasthttp.Server{
		Name:               "MyStaticServer",
		Handler:            router.Handler,
		ReadTimeout:        5 * time.Second,
		WriteTimeout:       10 * time.Second,
		IdleTimeout:        30 * time.Second,
		MaxRequestBodySize: 1 * 1024 * 1024, // 1 MB
	}

	router.GET("/", webserver)                       // Gets the next Article, using IP as a control
	router.ServeFiles("/static/*filepath", "static") // Gets the next Article, using IP as a control
	router.GET("/api/V1", newsHandler)               // Gets the next Article, using IP as a control

	lib.Debug("Connection:", ":"+conf.PORT)
	lib.CheckErr(server.ListenAndServe(":" + conf.PORT)) //This holds live.
	panic("Abnormal exit:" + "Possible Port conflict?" + conf.PORT)
}

/**********************_*********************************************
 *                    | |
 *     __      __ ___ | |__   ___   ___  _ __ __   __ ___  _ __
 *     \ \ /\ / // _ \| '_ \ / __| / _ \| '__|\ \ / // _ \| '__|
 *      \ V  V /|  __/| |_) |\__ \|  __/| |    \ V /|  __/| |
 *       \_/\_/  \___||_.__/ |___/ \___||_|     \_/  \___||_|      */
func webserver(ctx *fasthttp.RequestCtx) {
	if !limiter.Allow() { // return a 429 Too Many Requests error if the limit is exceeded
		ctx.Error("Too Many Requests", fasthttp.StatusTooManyRequests)
		return
	}
	// serve files from the "html" folder
	switch string(ctx.Path()) {
	case "/":
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SendFile("./static/index.html")
	default:
		// return a 404 Not Found error for unknown paths
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

/********************************************************************
 *                              _    _                 _ _
 *                             | |  | |               | | |
 *      _ __   _____      _____| |__| | __ _ _ __   __| | | ___ _ __
 *     | '_ \ / _ \ \ /\ / / __|  __  |/ _` | '_ \ / _` | |/ _ \ '__|
 *     | | | |  __/\ V  V /\__ \ |  | | (_| | | | | (_| | |  __/ |
 *     |_| |_|\___| \_/\_/ |___/_|  |_|\__,_|_| |_|\__,_|_|\___|_|    */
func newsHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Warn("newsHandler problem:", r)
		}
	}()
	if !limiter.Allow() {
		ctx.Error("Too many requests", fasthttp.StatusTooManyRequests)
		return
	}
	//TO-DO add count for the key used, as this is 1 request.
	accessKey := string(ctx.QueryArgs().Peek("access_key"))
	if account, ok := sql.Accounts[accessKey]; ok {
		if account.Used > account.Allocated {
			ctx.Error("Usage Limit Reached", fasthttp.StatusUnauthorized)
			return
		} else if account.End.Before(time.Now()) {
			ctx.Error("Account period expired", fasthttp.StatusUnauthorized)
			return
		} else { // All good, update values.
			account.Used++
			sql.Accounts[accessKey] = account
		}
	} else {
		ctx.Error("Access key is invalid", fasthttp.StatusUnauthorized)
		return
	}

	keywords := string(ctx.QueryArgs().Peek("keywords"))
	date := string(ctx.QueryArgs().Peek("date"))
	categories := string(ctx.QueryArgs().Peek("categories"))
	sources := string(ctx.QueryArgs().Peek("sources"))
	limit := string(ctx.QueryArgs().Peek("limit"))
	offset := string(ctx.QueryArgs().Peek("offset"))
	sort := string(ctx.QueryArgs().Peek("sort"))

	lib.Debug("Query:", accessKey, ":", keywords, ":", date, ":", categories, ":", sources, ":", limit, ":", offset, ":", sort, ";")
	// TODO: Handle the actual news API request here
	// using the accessKey, keywords, and countries values

	ctx.SetStatusCode(fasthttp.StatusOK)
}

// -----------------------------------------------------------------------------
//
//	_____                            _     _    _                 _ _
//
// |  __ \                          | |   | |  | |               | | |
// | |__) |___  __ _ _   _  ___  ___| |_  | |__| | __ _ _ __   __| | | ___ _ __
// |  _  // _ \/ _` | | | |/ _ \/ __| __| |  __  |/ _` | '_ \ / _` | |/ _ \ '__|
// | | \ \  __/ (_| | |_| |  __/\__ \ |_  | |  | | (_| | | | | (_| | |  __/ |
// |_|  \_\___|\__, |\__,_|\___||___/\__| |_|  |_|\__,_|_| |_|\__,_|_|\___|_|
//
//	| |
//	|_|
func requestDefault(ctx *fasthttp.RequestCtx, key string, sub string, kid string, subkid string) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Request Handler Failed:", r, key, sub, kid)
		}
	}()
	if lib.ThrottleAllow(ctx.RemoteIP().String(), conf.THROTTLE) {
		switch key {
		case "test", "help":
			_, _ = fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
			_, _ = fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
			_, _ = fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
			_, _ = fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
			_, _ = fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
			_, _ = fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
			_, _ = fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
			_, _ = fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
			_, _ = fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
			_, _ = fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
			_, _ = fmt.Fprintf(ctx, "Body is %q\n\n", ctx.PostBody())

		case "last": // Just return the last article
			lib.Info("Getting Last for :", sub)
			response(ctx, sql.GetLastJsonByFilter(setLimit(sub), kid, "general"))

		case "fetch": // Get article based on its UID
			lib.Info("Getting Article # :", sub)
			_, err := strconv.Atoi(sub)
			if lib.CheckErr(err) { // No articles ID
				_, _ = fmt.Fprintf(ctx, "Missing Article ID")
			} else {
				response(ctx, sql.GetJsonByArticleId(sub))
			}
		case "find", "tags": // Get limit number of items based on search
			lib.Info("Find by string :", sub)
			response(ctx, sql.GetJsonByTitle(sub, setLimit(kid), "general"))

		case "content": // Get limit number of items based on tags
			lib.Info("Find by string :", sub)
			response(ctx, sql.GetJsonByContent(sub, setLimit(kid), "general"))

		case "next": // Get limit number of items based on tags
			lib.Info("Doing Next :", key, sub, kid, subkid)
			response(ctx, sql.GetNextJsonByArticleId(sub, setLimit(kid), subkid, "general"))

		case "prev": // Get limit number of items based on tags
			lib.Info("Doing Previous with String :", sub, kid, subkid)
			response(ctx, sql.GetPrevJsonByArticleId(sub, setLimit(kid), subkid, "general"))

		default: // If the key is not a function, then its a article number so get the next one.
			lib.Info("Doing Defaults :", key, sub, kid)
			response(ctx, sql.GetNextJsonByArticleId(key, setLimit(sub), kid, "general"))
		}
	} else {
		response(ctx, conf.Reply{Code: 400, Msg: "Too Frequent"})
	}
}

// Use the caller IP as the Article control
func requestFetchIp(ctx *fasthttp.RequestCtx) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Warn("No article count:", r)
		}
	}()
	idStr := ctx.RemoteIP()
	response(ctx, sql.GetLatestJsonArticle(idStr.String(), 10, true, "API", "general"))
}

/****************************************************
 *               _   _      _           _ _
 *              | | | |    (_)         (_) |
 *      ___  ___| |_| |     _ _ __ ___  _| |_
 *     / __|/ _ \ __| |    | | '_ ` _ \| | __|
 *     \__ \  __/ |_| |____| | | | | | | | |_
 *     |___/\___|\__|______|_|_| |_| |_|_|\__|
 * --------------------------------------------------
 * Sets the limit to that of the config if it is above
 * * * * * * * * * * * * * * * * * * * * * * * * * */
func setLimit(value string) (limit int) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Warn("Limit set error:", r)
		}
	}()
	limit, err := strconv.Atoi(value)
	if lib.CheckErr(err) {
		limit = 1
	} // If there is an error, then set the default to 1
	if limit > conf.NEWSLIMIT {
		limit = conf.NEWSLIMIT
	}

	return limit
}

/*****************************************************
 *      _ __ ___  ___ _ __   ___  _ __  ___  ___
 *     | '__/ _ \/ __| '_ \ / _ \| '_ \/ __|/ _ \
 *     | | |  __/\__ \ |_) | (_) | | | \__ \  __/
 *     |_|  \___||___/ .__/ \___/|_| |_|___/\___|
 *                   | |
 *                   |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Package a proper response to the caller.
 * ------------------------------------------------- */
func response(ctx *fasthttp.RequestCtx, reply conf.Reply) {
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Responding to request:", r)
		}
	}()

	if reply.Code < 300 {
		_, _ = fmt.Fprint(ctx, reply.Msg)
		ctx.SetStatusCode(reply.Code)
	} else {
		ctx.Error(reply.Msg, reply.Code)
		ctx.Response.ImmediateHeaderFlush = true
	}
}
