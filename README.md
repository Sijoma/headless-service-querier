# headless-service-querier

This project has two components that can be started with flags. A "server" and a "querier".

The server exposes a worker count through a endpoint. 
The queriers job is then to lookup all IPs of the backing pods and gather the worker count from all the pods.
