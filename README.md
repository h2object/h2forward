# h2forward
http reverse proxy with http api to control the virtual hosts to forward

## Quick Start

#### Install

````
$: go get github.com/h2object/h2forward

````

#### Start

````
$: h2forward -w=/path/to/work -l="0.0.0.0:80" -a="127.0.0.1:9000" -d start
````

#### Stop

````
$: h2forward -w=/path/to/work stop
````

#### Virtual Host

-	Set Virtual Host

````
$: h2forward -a="127.0.0.1:9000" virtualhost set "www.example.com" "http://127.0.0.1:8080"
````

-	Del Virtual Host

````
$: h2forward -a="127.0.0.1:9000" virtualhost del "www.example1.com" "www.example2.com" 
````

-	Get Virtual Host

**Get All***

````
$: h2forward -a="127.0.0.1:9000" virtualhost get 
````

**Get One***

````
$: h2forward -a="127.0.0.1:9000" virtualhost get "www.example.com"
````

## Enjoy It!
