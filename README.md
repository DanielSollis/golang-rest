## Installation

To install the command line interface you can run the following command from the project's root directory:

```
$ go install pingthings/pingcli
```

## Command Line Interface

The CLI handles both serving and issuing requests to the server.

While serving you can use SIGINT (ctrl + c) to gracefully shutdown the server.

## Example requests

With the server up you can issue requests to it either using the CLI or through
standard means like Curl, postman, etc.

You can list all servers in the database:

```
$ curl http://localhost:8080/allsensors
```

query a specific sensor by name:

```
$ curl http://localhost:8080/sensor/C1MAG
```

insert a sensor by providing JSON data:

```
$ curl http://localhost:8080/sensor \
 --include \
 --header "Content-Type: application/json" \
 --request "POST" \
 --data '{"name": "C2MAG", "location": {"latitude": 38.4,"longitude": 26.9},
"tags": {"name": "C2MAG","unit": "amps", "ingress": "brazil", "distiller": "foo"}}'
```

update a sensor already in the database:

```
$ curl http://localhost:8080/sensor/C2MAG \
 --include \
 --header "Content-Type: application/json" \
 --request "PUT" \
 --data '{"name": "C2MAG", "location": {"latitude": 50.3,"longitude": 26.9},
"tags": {"name": "C2MAG","unit": "amps", "ingress": "brazil", "distiller": "foo"}}'
```

or perform a health check:

```
$ curl http://localhost:8080/health
```
