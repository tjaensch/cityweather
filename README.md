# cityweather

Small app written in Go with Bootstrap CSS that returns city weather data from the openweathermap.org API.

*With working Go installation:*
CD into application folder, set `OPENWEATHERMAP_API_KEY`, and run `go run cityweather.go`

Example: `OPENWEATHERMAP_API_KEY=your_api_key go run cityweather.go`

*OR*

*Run application as a Docker container:*

	CD into application folder
	Run "docker build -t cityweather ."
	Run "docker run --publish 6060:8080 --env OPENWEATHERMAP_API_KEY=your_api_key --name cityweather --rm cityweather"
	Open http://localhost:6060/ in a web browser and you should see something like this:

![cityweather.png](https://github.com/tjaensch/cityweather/blob/master/cityweather.png)