# cityweather

Small app written in Go with Bootstrap CSS that returns city weather data from the openweathermap.org API.

*With working Go installation:*
CD into application folder and run "go run cityweather.go"

*OR*

*Run application as a Docker container:*

	CD into application folder
	Run "docker build -t cityweather ."
	Run "docker-machine ip" to find IP address your docker daemon is running on
	Run "docker run --publish 6060:8080 --name cityweather --rm cityweather"
	Open http://YOUR-DOCKER-DAEMON-IP:6060/ in a web browser and you should see something like this:

![cityweather.png](https://github.com/tjaensch/cityweather/blob/master/cityweather.png)