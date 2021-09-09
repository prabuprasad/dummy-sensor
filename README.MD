# Test Sensor

The dummy application consists of a Go client (called test-sensor) that provides data from various functions (rectangle / sine).
Once started, the sensor regularly sends measured values to an event broker (pubsub+ from Solace).
A second application called dummy app consists of a Spring-Boot consumer that receives data from Solace via an event listener and stores it in an InfluxDB.
The data stored in the InfluxDB is visualised in Grafana.

```Create docker image
    docker build -t testsensor:0.0.1 .
```
```run testsensor
    docker run -it --rm  --env-file docker-env-settings testsensor:0.0.1
```



