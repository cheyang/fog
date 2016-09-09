# infrastructure of infrastructure

 原则：
 * 轻量级
 * 简单
 * 没有侵入性

 # Fog

 ## 背景和要解决的问题

 随着微服务和serverless架构的广泛使用，基于docker的容器云架构大行其道，swarm，swarmkit， kubernetes 和mesos以及其上的的maraton等各路诸侯都有自己的忠实拥趸，我深深认同容器云将成为应用开发的基础架构。但是如何在各种各样的公有云和私有云姿势正确的部署各种各样的容器云框架，并且做到一键式的完成。所以我希望能够提供一个批量创建一组虚机并且可以灵活的部署容器云的框架，可以方便的在其上做二次开发。所以我称我们的工具为基础架构的基础架构。

 ## 解决的方法和想法

docker machine + golang + ansible + docker

使用docker machine的原因是它提供了一套标准的创建虚机并且用ssh连接到虚机进行工作的API接口，并且有一大群云厂商提供了SPI的实现，并且负责维护，这就可以重用docker社区在IAAS层标准化的工作。并且这些厂商都在docker machine的实现上提供了专业的优化，也是我们可以重用的部分。

使用golang的原因  
1） docker machine的类库是golang实现  
2) 单独的binary容易部署  

使用ansible的原因  
1) ansible agent less的架构非常适合部署配置比较简单的docker集群  
2) offical的playbook支持保证最专业的部署, k8s,mesos对playbook的支持都是最全面的  
3) 强大的社区支持  
4) 更为强大和灵活的可配置型(比如扩展master，etcd)
使用docker

## 价值

* 重用和整合

不重新制造任何轮子，而是把已有的轮子整合起来

1. 资源或生态

docker machine已经得到了众多云大厂的广泛支持，从Amazon，阿里云， 

2. 知识

* 轻量级和无侵入性
