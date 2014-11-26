# How to run Kvik

To make it easy to get Kvik running on your local system, we're made a Docker
image with Kvik already installed! To get that running on our system: 

- [Install Docker](https://docs.docker.com/installation)
- Pull down the [Docker repository](https://registry.hub.docker.com/u/fjukstad/kvik/) by running
```
docker pull fjukstad/kvik
```
- Start the Docker container. 
```
docker run -p 80:8000 -p 8080:8080 -p 8888:8888 -t fjukstad/kvik
```
- Kvik now runs on port 80, [go have a look](http://localhost)


## Building from a Dockerfile
You could also build the Kvik container from a Dockerfile. From this directory,
run

```
docker build -t USERNAME/kvik .
```

and then

```
docker run -p 80:8000 -p 8080:8080 -p 8888:8888 -t USERNAME/kvik
```

## Making changes to Kvik 
If you want to change out the data analysis in Kvik, make your changes to
[data-engine.r](data-engine.r),  uncomment [the line](https://github.com/fjukstad/kvik/blob/187c60d60216203318538f946275d438203cba74/src/docker/Dockerfile#L99) that reads 

```
# COPY data-engine.r /root/kvik/src/src/github.com/fjukstad/dataengine/
```

in the Dockerfile and build the image from the Dockerfile. 

# I can't make it work! 
Let me know by submitting an [issue!](https://github.com/fjukstad/kvik/issues) 
