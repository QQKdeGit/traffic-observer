# Traffic Observer

Traffic Observer is a cloud native traffic real-time monitoring system. Including traffic data preprocessing, feature extraction, malicious traffic detection, detection result management, visualization, user management and other functions.

See the [中文文档](/docs/readme_cn.md) for Chinese readme.



## Architecture overview

<img alt="Architecture" src="/images/architecture.svg" width="50%">

Look down from above Traffic Observer, the whole project can be divided into four parts: HTTP Proxy, Traffic Detector, Prometheus and Grafana. They work as follows:

1. HTTP Proxy: Receive an HTTP request to visit a web, extract traffic information, and send it as POST request parameters to Traffic Detector
2. Traffic Detector: Whenever a POST request from the HTTP Proxy is received, Traffic Detector parse the parameters into an array of HTTP request, analyze the traffic for malicious traffic, and return the result
3. [Prometheus](https://prometheus.io): Send a request to the HTTP Proxy, obtain the analyzed traffic information, and store it in the time-series database
4. [Grafana](https://grafana.com): Get data from Prometheus, provide visualization, user management and other functions



## Technology stack introduction

- `Go` : Develop an HTTP Proxy with Go
- `FastAPI` : Equipped with Traffic Detector
- `Prometheus` : Provide time-series database
- `Grafana` : Provide visualization, user management function
- `Docker` : Container solution

The system architecture fully considers the demands of the cloud native era, so Traffic Observer is designed to be as **loosely coupled, scalable, lightweight and portable** as possible.



On the whole, the system mainly adopts three concepts of microservice, containerization and DevOps, which are as follows:

#### Microservice

In order to be loosely coupled and scalable, the system is basically doomed to a microservice architecture. The traditional monolithic architecture can no longer meet the needs of the cloud native era. The loose coupling between microservices and the support of modern container engines make it have good horizontal scalability.

#### DevOps

In order to speed up the construction, testing and release of Traffic Observer, GitHub Actions is used to build a DevOps Pipeline, the core of which is CI/CD. Once the repository receives the commit from the branch master, the pipeline is triggered, followed by static code checking and project building. After the build is complete, the HTTP Proxy and the Traffic Detector will be packaged into Docker images and published to Docker Hub so that we can deploy it in the phase of CI.

#### Containerization

Containerization is also the foundation of the cloud native era. It is lightweight and portable. It packages code and all required components so that application can run consistently in any environment and on any infrastructure. It mainly uses the capability provided by the Linux Kernel to implement resource isolation through `namespace`, implement resource restrictions through `cgroups`, and implement Copy on Write through `UnionFS`. If traditional deployment solutions are used, applications are difficult to migrate, and developers have to spend a lot of time adapting to the new environment. The microservice and DevOps mentioned above simply won't happen. It can be said that containerization is an inevitable choice for cloud native applications. No containerization, no cloud native.



On the single modules of Traffic Observer, the functions of each module and the reasons for choosing this technology are as follows:

#### HTTP Proxy

The reason why Go is chosen as the development language of the HTTP Proxy is that Go has execution performance like C and near-perfect compilation speed. At the same time, its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, which almost perfectly meets the needs of the HTTP Proxy in Traffic Observer to send requests to Traffic Detector and process Prometheus requests at the same time.

#### Traffic Detector

Traffic Observer uses FastAPI to equip machine model to analyze malicious traffic. FastAPI has extremely high performance comparable to Go and is one of the fastest Python Frameworks available. And it minimizes code duplication, reduces the time of debugging and is designed to be easy to use and learn.

#### Prometheus

Prometheus, a Cloud Native Computing Foundation project, is a systems and service monitoring system. It is able to collect and store its metrics as time series data, i.e.  metrics information is stored with the timestamp at which it was recorded, alongside optional key-value pairs called labels. This feature undoubtedly fits well with Traffic Observer's demands for storing traffic information.

#### Grafana

Grafana is an open source platform for visualizing large measurement data, which provides a powerful and elegant way to create, share, explore data. The Grafana and Prometheus open source projects are de facto standards for observability, with wide grassroots adoption.