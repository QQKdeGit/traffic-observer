# Traffic Observer

Traffic Observer 是一个[云原生](https://github.com/cncf/toc/blob/main/DEFINITION.md)的网络流量实时监控系统。包含网络流量数据包预处理功能、特征提取功能、恶意流量检测、检测结果管理功能、检测结果图表展示功能、用户管理功能等。



## 架构概述

<p align="center">
  <img alt="Architecture" src="/docs/images/architecture-cn.svg" width="50%">
</p>

鸟瞰 Traffic Observer ，整个项目可以分为四个部分：HTTP 代理、流量检测器、Prometheus 以及 Grafana 。其功能分别为：

1. HTTP 代理：接收想访问某网页的 HTTP 请求，提取流量信息，并将其作为参数发送一个 POST 请求给流量检测器
2. 流量检测器：每当收到来自 HTTP 代理发送的 POST 请求时，就将参数解析成 HTTP 请求的数组，分析该网络流量是否为恶意流量，并返回结果
3. [Prometheus](https://prometheus.io)：向 HTTP 代理发送请求，获取已分析好的流量数据，并存入时序数据库（time-series database）中
4. [Grafana](https://grafana.com)：从 Prometheus 中获取数据，并提供可视化、用户管理等功能



## 技术栈介绍

- `Go` ：使用 Go 开发 HTTP 代理
- `FastAPI` ：搭载流量检测器
- `Prometheus` ：提供时序数据库功能
- `Grafana` ：提供可视化、用户管理功能
- `Docker` ：容器解决方案

系统架构充分考虑到了云原生时代的需求，于是 Traffic Observer 尽量设计成**低耦合（loosely coupled）、可扩展（scalable）、轻量（lightweight）、可移植（portable）**的。



从整体来看，整个系统主要采用了**微服务、容器化和 DevOps **三个概念，下面逐一介绍：

#### 微服务

要做到低耦合和可扩展就基本注定了我们需要采用微服务架构，传统的单体架构已经无法满足云原生时代的需求。微服务模块之间的低耦合度、现代容器引擎的支持使得它们能很好地被横向扩展。

#### DevOps

为了加速 Traffic Observer 的构建、测试和发布，我利用了 GitHub Actions 搭建 DevOps Pipeline ，它的核心是持续集成与持续交付（CI/CD）。仓库一旦收到 master 分支的 commit ，便会触发 pipeline ，随后进行代码的静态检查和项目构建。构建完成后会将 HTTP 代理和流量检测器打包成 Docker 镜像发布至 Docker Hub ，以便我们在持续集成阶段部署。

#### 容器化

容器化也是云原生时代的基础，它轻量、可移植，它将软件代码和所需的所有组件打包在一起，这样应用才能够在任何环境和任何基础架构上一致地运行。它主要使用 Linux Kernel 提供的能力，通过 `namespace` 实现了资源隔离，通过 `cgroups` 实现了资源限制，通过 `UnionFS` 实现了 Copy on Write 的文件操作。如果使用传统的部署方案，应用难以移动，开发者需要花费大量的时间适配新的环境，以上所说的微服务与 DevOps 将根本无法推行。可以说，容器化是云原生应用的必然选择，没有容器化就没有云原生。



从单个模块来看，每个模块分别负责的功能以及为什么选用此技术完成的理由如下：

#### HTTP 代理

之所以选用 Go 作为 HTTP 代理的开发语言，就是因为 Go 拥有近乎 C 的执行性能和几乎完美的编译速度。同时 Go 的并发机制使编写充分利用多核和联网机器的程序变得容易，这几乎完美符合了 Traffic Observer 中 HTTP 代理既要给流量检测器发送请求同时又要处理 Prometheus 请求的功能所需。

#### 流量检测器

Traffic Observer 使用了 FastAPI 框架来搭载用于分析恶意流量的机器模型。FastAPI 拥有着比肩 Go 的极高性能，是目前最快的 Python web 框架之一。同时代码极其简短，易于使用和学习并且 bug 更少。

#### Prometheus

作为 Cloud Native Computing Foundation 项目，Prometheus 是一个系统和服务的监控系统。它能收集指标并存储为时间序列数据（即指标信息与记录时的时间戳，以及被称为标签的可选键值对一起存储），这一点无疑完美匹配于 Traffic Observer 存储流量信息的需求。

#### Grafana

Grafana 是用于可视化大型测量数据的开源程序，它提供了强大且优雅的方式去创建、共享、浏览数据。同时它和 Prometheus 是事实上的可观察性标准，并得到了广泛的基层采用。

