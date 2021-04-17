### 简介

```shell script
author: 郭飞
```

### 背景

> 服务器上经常有用户不规范地运行一些程序，导致服务器资源掉底。Glimit利用cgroup来限制用户的资源利用率。

### 原理

* 此程序运行在本机。
* 用户登录时，访问此程序的api，需要用户shell的pid和用户名作为GET参数。(如何访问？此命令将写入/etc/bashrc文件)
* 服务端收到请求后，将创建cgroup组，限制此pid和此pid的所有子进程。

### 依赖

* make
* crontab
* golang
* cgroup

安装cgroup

```shell script
#centos7
yum -y install libcgroup libcgroup-tools

#centos6
yum -y install libcgroup
sudo service cgconfig start     #启动cgroup服务
sudo chkconfig cgconfig on      #开机自启
```

### 编译

```shell script
make    #解压后进入项目目录，执行make命令
```

### 启动

* 启动前应修改配置文件，如centos6的cgroup目录为/cgroup
* 编译完成后，请将dlimit目录移动到适当的目录，并启动：```sudo nohup ./dlimit-server >> /dev/null &```

### 程序对系统的变动

* 监听1216端口
* 在/etc/bashrc中添加dlimit函数，用户登录时将执行此函数。
* 往/etc/crontab中添加计划任务，将过期的cgroup组进行清理。
