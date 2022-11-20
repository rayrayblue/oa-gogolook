# oa-gogolook

## Description
This service is build by clean architecture concept, 
therefore you can replace any layer with your own implementation. 
For example, we use in-memory storage for storing data, 
but you can replace it with any database you want. 
You just need to implement in repository layer using TaskRepository interface.


## How to run
Please follow the instructions below to run the service.

### Configuration
Only one environment variable (SERVER_ADDRESS=0.0.0.0:8888) is required to run the service.
app.env file contains the environment variable, which locates in configs folder.


### build image

```bash
docker build -t oa-gogolook .
```

### Run service from container

```bash
docker run -d -p 8888:8888 --name oa-gogolook oa-gogolook
```

### Testing from source code
If you want to run test from source code, please run the following command.

```bash
make test
```

### Run service from source code
If you want to run service from source code, please run the following command.

```bash
make server
```