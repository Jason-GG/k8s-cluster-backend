# k8s-cluster-backend

[![GitHub stars](https://img.shields.io/github/stars/Jason-GG/k8s-cluster-backend.svg)](https://github.com/Jason-GG/k8s-cluster-backend/stargazers)
&nbsp;
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](/LICENSE)
&nbsp;
[![GoDocWidget]][GoDocReference]

[GoDocWidget]: https://godoc.org/k8s.io/client-go?status.svg
[GoDocReference]:https://godoc.org/k8s.io/client-go 

## Getting Started

Download links:

SSH clone URL: git@github.com:Jason-GG/k8s-cluster-backend.git

HTTPS clone URL: https://github.com/Jason-GG/k8s-cluster-backend.git


## Description
the project is based on client-go, fiber, gorm to make a backend platform to manage k8s cluster.
please use the project: [vue3-k8s-dashboard-webssh]

[vue3-k8s-dashboard-webssh]:https://github.com/Jason-GG/vue3-k8s-dashboard-webssh.git

### Prerequisites
initialization the project with the command:
```
go get
```

### Deployment

```
make build
```

### deploy for the linux machine

```
make build-linux
```
### clean

```
make clean
```
## Important

Before launching this project, ensure you've completed the following configurations:

1, "config/prod" File: This file should represent your cluster's .kube/config file.
<br/>

2, "gaget/sql_info.json" File: Fill this file with the required database details. For the user table, manually add an account to enable the API to proceed with the authentication process.