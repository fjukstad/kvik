# OpenCPU interface 
This is a simple go interface to OpenCPU. See the [Dockerfile](Dockerfile) for
the docker container and [kompute.go](kompute.go) for the interface itself. 

# Example 
See the [examples](examples/) folder for an example on how to use it. 

# Docker container 
```
docker run -t -p 8004:8004 -v /Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/kompute/home/:/home --name opencpu TAG
```
