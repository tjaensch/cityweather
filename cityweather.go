package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	}
	Sys struct {
		Country string `json:"country"`
	}
	Weather []weather `json:"weather"`
}

func (w weatherData) Fahrenheit() int {
	return int(w.Main.Temp*9/5 - 459.67)
}

type weather struct {
	Description string `json:"description"`
}

var (
	err      error
	res      weatherData
	response *http.Response
	body     []byte
)

func main() {

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.HandleFunc("/", handler)
	http.HandleFunc("/showweather", showweather)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, rootForm)
}

const rootForm = `
  <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="city, temperature, golang, go_appengine">
    <meta name="author" content="Thomas Jaensch">
    <title>City Weather</title>
    <!-- Bootstrap Core CSS -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <!-- Custom CSS -->
    <link href="/css/cityweather.css" rel="stylesheet">
    <!-- Custom Fonts -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.3.0/css/font-awesome.min.css" rel="stylesheet" type="text/css">
    <link href="http://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,700,300italic,400italic,700italic" rel="stylesheet" type="text/css">
    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
    <!-- Header -->
    <header id="top" class="header">
        <div class="text-vertical-center">
            <h1>City Weather</h1>
            <h3>Enter a city name to see the current weather conditions*</h3>
            <p><small>*returns raw data from openweathermap.org API, source code on <a href="https://github.com/tjaensch/cityweather" target="blank">GitHub</a></small></p>
            <br>
            <form action="/showweather" method="post" accept-charset="utf-8">
                    <input type="text" name="city" placeholder="Enter city name" id="city" />
                    <input type="submit" value="Submit" class="btn btn-dark btn-lg"/>
            </form>
        </div>
    </header>
    <!-- Footer -->
    <footer>
        <div class="container">
            <div class="row text-center">
                <p class="text-muted">&copy; Thomas Jaensch 2015</p>
            </div>
        </div>
    </footer>
</body>
</html>
`

var upperTemplate = template.Must(template.New("showweather").Parse(upperTemplateHTML))

func showweather(w http.ResponseWriter, r *http.Request) {
	city_value := r.FormValue("city")

	safe_city_value := url.QueryEscape(city_value)
	apikey := "&APPID=e637873503756b3e4182c1b0e80e8881"
	fullUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s", safe_city_value+apikey)

	response, err = http.Get(fullUrl)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &res)

	tempErr := upperTemplate.Execute(w, res)
	if tempErr != nil {
		http.Error(w, tempErr.Error(), http.StatusInternalServerError)
	}
}

const upperTemplateHTML = ` 
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="city, temperature, golang, go_appengine">
    <meta name="author" content="Thomas Jaensch">
    <title>City Weather</title>
    <!-- Bootstrap Core CSS -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <!-- Custom CSS -->
    <link href="/css/cityweather.css" rel="stylesheet">
    <!-- Custom Fonts -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.3.0/css/font-awesome.min.css" rel="stylesheet" type="text/css">
    <link href="http://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,700,300italic,400italic,700italic" rel="stylesheet" type="text/css">
    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
    <body>
    <!-- Header -->
    <header id="top" class="header">
        <div class="text-vertical-center">
            <h3>{{.Name}} ({{.Sys.Country}}): {{range .Weather}}
                {{.Description}}
            {{end}}, {{.Fahrenheit}} ºF</h3>
            <a href="/" class="btn btn-dark btn-lg">Try Again</a>
        </div>
    </header>
    <!-- Footer -->
    <footer>
        <div class="container">
            <div class="row text-center">
                <p class="text-muted">&copy; Thomas Jaensch 2015</p>
            </div>
        </div>
    </footer>
   </body>
</html>
`
