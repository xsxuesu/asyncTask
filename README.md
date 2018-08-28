         Author Fanxu(746439274@qq.com)
         
# 程序编译说明
    
    #编译
     进入install/bin/目录
     go build http.go
    #如果有配置文件修改
    则将config/online.json 复制到 install/env/online.json
    
    
# 安装启动说明
    #进入服务器 服务器版本要求 Centos7 64位 及以上 Ubuntu16 64位 及以上
    #第一步：下载服务接口程序
    sudo git clone http://42.51.64.22/backend/install.git /home/baas
    #第二步：修改文件权限
    cd /home/baas
    sudo chmod a+x install.sh
    sudo chmod a+x start.sh
    #第三步：安装baas所需要的基础环境（根据网速不同安装速度有所不同）
    sudo sh install.sh
    #********检测是否安装成功(重要)********
    docker images | grep latest
    #查看是否有以下images tag为latest
   ![work](install.png)
        
    #如果缺少image则重新执行 sudo sh install.sh
    
    #启动BaaS服务器接口服务（以上环境检测安装完毕则可以运行程序 必须要用sudo权限启动）
    #启动接口服务
    cd /home/baas
    sudo ./start.sh up
    #停止接口服务
    cd /home/baas
    sudo ./start.sh down
    #重新启动接口服务
    cd /home/baas
    sudo ./start.sh restart
    #检测程序运行日志
    tail -f /tmp/http.go
    
    #检测程序是否运行中
    ps -ef | grep "env/online.json" | grep "http"
    
# 框架说明
## 框架 包含http 服务 rpc 服务 异步队列服务

    启用http 服务 则 go run http.go

    启用rpc 服务 则 go run rpc.go

    启用队列服务 在入口 加入
      queue.Init()
      go queue.Queue.Consumer()
  
## 文件夹说明

    config 配置文件
    
    control 控制逻辑代码文件 既可以写 http 接收 也可以写 rpc 逻辑
    
    helpers 自定义 工具类 函数库
    
    model 数据库 字段映射 或者 一些 interface 复杂数据结构类型
    
    queue 异步队列处理
    
    route 路由 针对http 规则 http://xxxx/user/index  对应执行 control/user/UserIndex 方法
    
    rpc control 则需要继承 control.RpcServer 客户端调用 "RpcServer.UserIndex"
    
    
    
<!-- #第三步：下载安装文件
    curl -L http://42.51.64.22/backend/asyncTask/blob/master/install/install.zip?raw=true > install.zip 
    #第四步：安装解压软件
    (centos)sudo yum install -y unzip zip 
    (ubuntu)sudo apt-get install -y unzip zip -->
    #第五步：解压安装文件
    unzip install.zip
    sudo chmod a+x install.sh
    #第六步：安装baas所需要的基础环境（根据网速不同安装速度有所不同）
    sudo sh install.sh
    #********检测是否安装成功(重要)********
    docker images | grep latest
    #查看是否有以下images tag为latest