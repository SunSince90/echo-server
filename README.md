# Echo server

A very simple server serving requests at `localhost:8080`.

Use this for testing.

## Usage

Clone this:

```bash
git clone https://github.com/SunSince90/echo-server.git
```

Build:

```bash
make build
```

Run:

```bash
./bin/server
```

Contact it:

```bash
curl localhost:8080/hey
```

Alternatively, if you want a long string to be printed, you may call

```bash
# Set paragraphs to whatever number you desire. Default is 1
call localhost:8080/lorem-ipsum?paragraphs=3
```

## Docker/Kubernetes

Build and push the docker container:

```bash
make docker-build docker-push IMG=<your-repository/image:tag>
```

Deploy it on Kubernetes:

```bash
kubectl create deployment echo-server --image <your-repository/image:tag>
```

Expose it:

```bash
kubectl create service loadbalancer echo-server --tcp=80:8080,8080:80
```

Contact it:

```bash
# Get the load balancer address
kubectl get service echo-server

curl http://<external-ip>/hey
```
