# News-Trawler-API
The cunction of this application is
1. Collect News feeds RSS from remote news sources.
1. Have an API for internal application presentation
1. Sources:RSS Feeds: from https://blog.feedspot.com/world_news_rss_feeds/#rightModal
1. Health. https://www.sciencedaily.com/rss/top/health.xml
1. Science feeds. https://www.sciencedaily.com/newsfeeds.htm
1. https://techxplore.com/feeds/


the API will have the following functions
1. Get the latest news based on the clients last news ID
1. Get the News on a specific day
1. Get the news based on a range of days, with Date from and date to.
1. Get the news based on range from today with number of days.
1. Get the news based on Tags in combination of Date Range or last received.

Formats.
News is a Get and Post request
Result is JSON object containing
 1. Header
 1. Detail/Article
 1. Image links (If Any)
 1. Source link to the original Article.

Target functions.
1. Post to Telegram
1. Post to Social
1. Requests from Connect
1. Requests from Exchange, based on the token Tags


## Posting News to Telegram
curl -s -X POST https://api.telegram.org/bot{token}/sendMessage -d chat_id={id} -d text='{message}'

token=1204200932:AAFR-Rr_kSzqSR4XnpcslTtVc0ddSRL1z_U

Chat IDs as follows:
To get the following: curl -s -X POST https://api.telegram.org/bot1204200932:AAFR-Rr_kSzqSR4XnpcslTtVc0ddSRL1z_U/getUpdates
Its Own Chat id=77612747
developer Channel ID=-308117337
Management Channel ID=-317581054
News Channel:-1001407700871
Link:https://telegram.me/AENSmartbot

Example:
curl -s -X POST https://api.telegram.org/bot1204200932:AAFR-Rr_kSzqSR4XnpcslTtVc0ddSRL1z_U/sendMessage -d chat_id=-308117337 -d text='{message}'
  
curl -s -X POST https://api.telegram.org/bot1204200932:AAFR-Rr_kSzqSR4XnpcslTtVc0ddSRL1z_U/sendMessage -d chat_id=77612747 -d text='{message}'

curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"chat_id": "@secret_news_general", "text": "General News channel, Dedicated to the trending news."}' \
  https://api.telegram.org/bot6279371821:AAF7DfAWC66UvTRP5GsYyWOJbgUscYsQYVw/sendMessage


## Paramaters. 
url:8002
"/"         gets the next article using the IP as a key
"/12abc"    Gets the next article using your defined key
"/12abc/5"  Gets the next 5 articles using your defined key

## Global Reserverd Keys:
-"last"      Gets the last article Global
-"last/2"    Gets last 2 articles

## Works with ID or default /# is optional, and limits the number of articles, default is 1
"prev/#" or "/{id}/prev/#"      Gets the previous article for your ID
"next/#" or "/{id}/next/#"      Same as above for next article using your IP

## Adding Limit:
"prev/{limit}" or "/{id}/prev{limit}"      Gets the previous article for your ID limited by the limit count
"next/{limit}" or "/{id}/next{limit}"      Same as above for next article using your IP limited by the limit count

## Searching for Articles returns json
"/find/Some text to search/2"              Find some text, with limit to how many articles to retreive
"/tag/some tags/5"                         Finds articles and searches the Tags/Categories
"/content/some text to search/10"          Seach the article content for some text, return a limit of pages

# Posting new articles
Duplicates will be rejected, based on the title.
Payload in raw json format
```
{ 
    "title" : "Mandatory, The Title displayed in Bold",
    "desc"  : "The Text of the article, Summary is good",
    "author": "Who wrote the article",
    "email" : "Email of the writer",
    "cat"   : "Tags for this article, Example, Bitcoin, Token, etc",
    "link"  : "Link to the Article, if it has one online."
}
```


# API Reference

The API format is based on REST, which means the parameters are part of the request. the format of this request is as a GET format, the response will be in JSON format.

Every article has an ID. it is up to the client to maintain their own ID. when a request is made without an ID they will receive the first article of the day. and its article number. the client then will send the number of the received article, and in return will receive the next one in sequence. The client should keep a receord of this sequence number for the next request. these requests should be spaced out, and not more than 1 request per minute, it is suggested that a 15 minut gap as a minimum is maintained.

The get your first article:

news.aenx.org  or aenx.org/news

To retreive the next Article sequence, use the last article number and the next in sequence will be returned.
For specific features they are listed as such:

`fetch/#` Will return the specific article by the ID  
`/article#/limit#/tags` Will return # number of articles matching the search term   
Example:  
```news.aenx.org/100/4``` will return the next 4 articles after the article number 100   

Example search you can enter the tag terms on the end. delimited by Comma.   
```
news.aenx.org/100/4/btc,bitcoin
```
will return the next 5 articles with the words btc or bitcoin within them.   
Example output woudld look like the following:
``` 
[
    {
        "uid": 12,
        "title": "Organic growth? Bitcoin SV activity up 761%!a(MISSING)head of BSV conference ",
        "desc": "Bitcoin SV users appear to have got very excited ahead of the conference. ",
        "author": "Cointelegraph By Andrew Fenton",
        "email": "",
        "cat": "Bitcoin SV Craig Wright",
        "link": "https://cointelegraph.com/news/organic-growth-bitcoin-sv-activity-up-761-ahead-of-bsv-conference",
        "detail": "",
        "created": "2020-10-08T13:39:22Z"
    },
    {
        "uid": 16,
        "title": "Bitcoin vs. USD: why only a weaker dollar will push BTC above $20,000",
        "desc": "Investors should keep an eye on the tight inverse correlation between the strength of the U.S. dollar and Bitcoin.",
        "author": "Cointelegraph By Michaël van de Poppe",
        "email": "",
        "cat": "Bitcoin Altcoin Markets Market Analysis Bitcoin Price Ethereum Price",
        "link": "https://cointelegraph.com/news/bitcoin-vs-usd-why-only-a-weaker-dollar-will-push-btc-above-20-000",
        "detail": "",
        "created": "2020-10-08T13:39:22Z"
    },
    {
        "uid": 18,
        "title": "Should Bitcoin traders be worried about lower highs ever since $12K?",
        "desc": "Bitcoin price technical analysis shows some key levels that traders should watch this week as BTC remains range-bound below $11,000.",
        "author": "Cointelegraph By Michaël van de Poppe",
        "email": "",
        "cat": "Altcoin Markets Price analysis Bitcoin price Ethereum price",
        "link": "https://cointelegraph.com/news/should-bitcoin-traders-be-worried-about-lower-highs-ever-since-12k",
        "detail": "",
        "created": "2020-10-08T13:39:22Z"
    },
    {
        "uid": 42,
        "title": "Ethereum Soars Over 125%!S(MISSING)ince March: What to Expect Now?",
        "desc": "Almost all markets across the world have been in turmoil owing to the economic uncertainties brought about by the coronavirus pandemic, and in that regard, the crypto market has been no different. However, one of the major cryptocurrencies to have made a remarkable recovery since hitting its lowest levels in March is Ethereum (ETH), and it is important to take a closer look at it. In this regard, it should be noted that ETH is the second-biggest cryptocurrency in the world, and its recovery might have an impact on the wider crypto market.\nMajor Triggers\nThe ... \t\t\t\n\t\t\t\t\tRead The Full Article On CryptoCurrencyNews.com\n\n\t\t\t\t\t\n\t\t\t\t\t\t\n\t\t\t\t\t\n\t\t\t\t\t\n\t\t\t\t\t\t \n\t\t\t\t\t\n\t\t\t\t\tGet latest cryptocurrency news on bitcoin, ethereum, initial coin offerings, ICOs, ethereum and all other cryptocurrencies. Learn How to trade on cryptocurrency exchanges.\n\t\t\t\t\tAll content provided by Crypto Currency News is subject to our Terms Of Use and Disclaimer.\n\t\t\t\t",
        "author": "Ankit Singhania",
        "email": "",
        "cat": "Ethereum News cryptocurrency ethereum digital currency editorial eth price ethereum ethereum bitcoin ethereum coin ethereum coin price ethereum mining ethereum news ethereum news today ethereum price chart ethereum value ethereum vs bitcoin financial news mining ethereum prconnect where to buy ethereum",
        "link": "https://cryptocurrencynews.com/ethereum-soars-expectations-04-28-20/?utm_campaign=rss__ccn\u0026utm_source=rss\u0026utm_medium=rss",
        "detail": "",
        "created": "2020-10-08T13:39:26Z"
    }
]
```


Sample Pipe script for jenkins
```
pipeline {
    agent any
    tools {
        go 'Golang-1.16'
    }
    environment {
        GO114MODULE = 'on'
        BRANCH = 'master'
        GITPATH = 'git@gitlab.aensmart.net:aenprojects/bots/news-trawler-api.git'
        TARGETSERVER01 = 'jenkins@babylon5.aenxchange.com'
        TARGETPATH = '/var/exchdb/news-bot'
    }
    stages {
        stage('CleanWorkspace') {
            steps {
                cleanWs()
            }
        }
        stage('Checkout Repositories') {
            steps {
                git branch: "${BRANCH}", credentialsId: "Jenkins4Gitlab", url: "${GITPATH}"
                sh 'pwd'
                sh 'ls -la'
            }
        }
        stage('Build') {
            steps {
                sh './deploy.sh -j build -v ${BUILD_NUMBER}'
            }
        }
        //stage('Dabatase') {
        //    steps {
        //        sh './deploy.sh -j database -v ${BUILD_NUMBER} -s jenkins@babylon5.aenxchange.com -f /var/exchdb/news-bot'
        //        findText alsoCheckConsoleOutput: true, regexp: 'command not found'
        //    }
        //}
        stage('Deploy') {
            steps {
                script {
                    sh 'ls -l'
                    sh 'chmod a+x ./deploy.sh'
                    sh './deploy.sh -j deploy -v ${BUILD_NUMBER} -s jenkins@babylon5.aenxchange.com -f /var/exchdb/news-bot -c newscentre'
                }
            }
        }
    }
    post {
        success {
            script{
                sh "curl -s -X POST https://api.telegram.org/bot846438395:AAE7SJaVfFWyO1sbgk8nQNkbHQoROI6oQGo/sendMessage -d chat_id=${telegram} -d text='[PROD Success] News-Trawler-Bot' "
            }
        }
        failure {
            script{
                sh "curl -s -X POST https://api.telegram.org/bot846438395:AAE7SJaVfFWyO1sbgk8nQNkbHQoROI6oQGo/sendMessage -d chat_id=${telegram} -d text='[PROD Failed] News-Trawler-Bot' "
            }
        }
    }
}
```


# API Documentation
Sample here:https://mediastack.com/documentation
## API Documentation
The mediastack API was built to provide a powerful, scalable yet easy-to-use REST API interface delivering worldwide live and historical news data in handy JSON format. The API comes with a single news HTTP GET endpoint along with a series of parameters and options you can use to narrow down your news data results. Among other options, you can filter by dates and timeframes, countries, langauges, sources and search keywords.

To get started right away, you can either jump to our 3-Step Quickstart Guide using the button above or scroll down and learn how to authenticate with the mediastack API. If you do not have an account yet, please make sure to get a free API key now to start testing the API and retrieve your first news article.


## Getting Started
API Authentication
To make an API request, you will need an API access key and authenticate with the API by attaching the access_key GET parameter to the base URL and set it to your unique access key. Find below an example request.

Example API Request:

```
Sign Up to Run API Requesthttps://api.mediastack.com/v1/news
    ? access_key = YOUR_ACCESS_KEY
```

A unique API access key is generated for each mediastack account, and it usually never changes. There can only ever be one API access key at a time per account. If you need to re-generate your key, you can do so by logging in to your account dashboard.

256-bit HTTPS EncryptionAvailable on: Standard Plan and higher
To connect to the mediastack API using 256-bit HTTPS (SSL) encryption, you will need to be subscribed to the Standard Plan or higher. If you are on the Free Plan, please note that you will be limited to HTTP connections.

Example API Request:
```
https://api.mediastack.com/v1
```
API Errors
If your API request was unsuccessful, you will receive a JSON error in the format outlined below, carrying code, message and context response objects in order to communicate the type of error that occurred and details to go with it. The API will also return HTTP status codes in accordance with the type of API response sent.

Below you will find an example API error that occurs if an unknown value is set for the API's date parameter.

Example Error:
```
{
   "error": {
      "code": "validation_error",
      "message": "Validation error",
      "context": {
         "date": [
            "NO_SUCH_CHOICE_ERROR"
         ]
      }
   }
}
```

# Comon API Layer

| Type                       | Description                                                                                    |
| --------------------------| ---------------------------------------------------------------------------------------------- |
| invalid_access_key         | An invalid API access key was supplied.                                                       |
| missing_access_key         | No API access key was supplied.                                                               |
| inactive_user              | The given user account is inactive.                                                           |
| https_access_restricted    | HTTPS access is not supported on the current subscription plan.                               |
| function_access_restricted | The given API endpoint is not supported on the current subscription plan.                     |
| invalid_api_function       | The given API endpoint does not exist.                                                         |
| 404_not_found              | Resource not found.                                                                           |
| usage_limit_reached        | The given user account has reached its monthly allowed request volume.                         |
| rate_limit_reached         | The given user account has reached the rate limit.                                             |
| internal_error             | An internal error occurred. business in a specific [MARKET].                                   |


# API Features

Live NewsAvailable on: All plans
The full set of available real-time news articles can be accessed using a simple API request to the mediastack API's news endpoint. Below you will find an example API request as well as a series of optional parameters you can use to filter your news results.

Delayed news on Free Plan: Please note that account subscribed to the Free Plan will receive live news only with a 30-minute delay. To lift this limitation and get news in real-time, please sign up or upgrade to the Standard Plan or higher.

Sample API request
```
https://api.mediastack.com/v1/news?access_key=My_ACCESS_KEY&keywords=tennis&countries=us,gb,de
```

## HTTP GET Request Parameters:



# VPN Method

```
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	vpnConfigPath := "/path/to/vpn/config/file.ovpn"

	err := connectToVPN(vpnConfigPath)
	if err != nil {
		fmt.Printf("Failed to connect to VPN: %v\n", err)
		return
	}

	// The VPN connection has been established successfully
	// Now, make an HTTP request through the VPN connection

	// Parse the RSS feed URL
	rssURL, err := url.Parse("https://techxplore.com/rss-feed/breaking/machine-learning-ai-news/")
	if err != nil {
		fmt.Printf("Failed to parse RSS feed URL: %v\n", err)
		return
	}

	// Create an HTTP client that uses the VPN connection
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	}
	client := &http.Client{
		Transport: transport,
	}

	// Make an HTTP GET request to the RSS feed URL
	req, err := http.NewRequest("GET", rssURL.String(), nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to make HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// TODO: Process the RSS feed data
}
```