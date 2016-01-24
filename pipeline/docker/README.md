# How to

Here's how you run the thing: 

First get a data container up and running. The workers share the /tmp dir
so that they can share results from OpenCPU. Go to the `fileserver` dir and run

```
docker build -t data . 
docker create -v /tmp --name data data /bin/true
```

Then go back to this dir and run

```
docker-compose build
docker-compose scale worker=NUMWORKERS
docker-compose up
```


