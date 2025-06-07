## LodBal

A sample implementation how the internals of a load balancer looks like. Inspired from the reverseproxy logic of the [Go standard library](https://golang.org/pkg/net/http/httputil/#ReverseProxy).

PS: This is not a production-ready load balancer, but rather a learning project to understand how load balancing works.

## How to use it
1. Clone the repository:

```code
https://github.com/akp111/lodbal.git
cd loadbal
```

2. Install the dependencies:

```code
go mod tidy
```

3. Start the backend servers:

```code
./start_dummy_servers.sh
```

3. Run the load balancer (in another terminal tab):

```code
go run main.go
```

4. Test the load balancer (in another terminal tab):
```code
curl -v http://localhost:8000
```
