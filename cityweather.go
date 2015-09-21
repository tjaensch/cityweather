package cityweather

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/showweather", showweather)
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
    <link href="/stylesheets/cityweather.css" rel="stylesheet">
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
            <h3>Enter a city name to see the current weather data*</h3>
            <p><small>*returns raw data from openweathermap.org API</small></p>
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
	fullUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s", safe_city_value)

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	resp, err := client.Get(fullUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Read the content into a byte array
	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		panic(dataReadErr)
	}

	//res := make(map[string]interface{})
	var res weatherData

	json.Unmarshal(body, &res)

	tempErr := upperTemplate.Execute(w, res)
	if tempErr != nil {
		http.Error(w, tempErr.Error(), http.StatusInternalServerError)
	}
}

type weatherData struct {
	Name    string    `json:"name"`
	Weather []weather `json:"weather"`
}

type weather struct {
	Description string `json:"description"`
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
    <link href="/stylesheets/cityweather.css" rel="stylesheet">
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
            <h3>{{.Name}}: {{range .Weather}}
                {{.Description}}
            {{end}}</h3>
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
