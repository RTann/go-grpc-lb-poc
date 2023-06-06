## gRPC load balancing POC

The purpose of this Proof of Concept (POC) is to evaluate how a gRPC client, with "round-robin" load balancing configuration, handles scaling operations in a Kubernetes environment. I want to understand how the client reacts to changes in the number of server pods (i.e., when they are scaled up and down), and how effective the client-side load balancing is.

This evaluation is particularly to decide if your distributed gRPC service will dynamically adjust to load and what are the availability implications. In particular how gRPC, with its built-in load balancing policies and service discovery via DNS, provides flexibility and efficiency in this regard.

The metrics defined to evaluate the server and client interactions:

1. Request distribution: These metrics give us visibility into the balance of load among the servers. In a well-functioning round-robin setup, we would expect an equal distribution of requests.
2. Latency and performance: These metrics allow us to identify potential bottlenecks or overworked servers.  Currently in this POC, they are mostly meaningless since latency is based on a constant sleep time.
3. Error handling (Client and Server Error Rates): A skewed error rate to show potential issues.
4. Adaptability (Time to Detect Changes, Time to Server Discovery): These metrics evaluate how quickly the system adapts to changes, such as the addition or removal of server pods.
5. Server utilization and health (CPU/Memory Utilization, Server Inactivity Age): These metrics help identify servers that might be under or over-utilized.

Dashboards that would be particularly useful:

- Request Distribution Dashboard: A histogram showing the distribution of requests among the servers.

- Performance Dashboard: A time-series graph showing latencies and response times.

- Error Rates Dashboard: A time-series graph showing error rates.

- Adaptability Dashboard: A graph showing the time to detect changes and time to server discovery.

- Server Utilization Dashboard: A graph showing CPU/memory utilization and server inactivity age.

## Test Scenarios

1. Scaling up the server pods. Here, expect to see the request distribution gradually adjust to include the new servers, and the time to detect changes should correspond to how quickly the new servers are included in the load balancing. Similarly, if you scale down the server pods, you would expect them to gradually disappear from the request distribution, and their server inactivity age to increase.
2. (Not currently supported) Induce some artificial load or errors on specific server pods. You would expect to see an increase in request latencies and error rates, possibly a change in request distribution if the client tries to avoid the troubled servers.

**This is a POC.** In real-world scenarios there more complexity, maybe more metrics and adjustments to the load balancing setup.
