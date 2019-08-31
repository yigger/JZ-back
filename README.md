## 简介
此仓库是 **洁账** 的 go 服务端，是基于 [Echo](https://github.com/labstack/echo) 框架开发的，由于本人也是第一次接触 go 语言，所以如果各位看官觉得代码写得不好或者有待优化的地方欢迎提 PR 指导指导。

## 小程序端地址
https://github.com/yigger/jiezhang

## 温馨提示：本项目处于开发中状态，暂不可同时配合小程序端使用。

#### 目录
- [环境配置](#环境配置)
- [运行](#运行)
- [目录结构](#目录结构)
- [热加载](#热加载)

## 环境配置
```
1. 源码下载
go get github.com/yigger/JZ-back

2. Mysql 配置文件
cp conf/conf.d/database.json.example database.json
编辑 database.json，填写你本地环境的参数

3. Database 迁移
go run migrate/run.go

3. Redis 配置文件
cp conf/conf.d/redis.json.example redis.json

4. 常规配置文件
cp conf/conf.d/jz.json.example jz.json
自行准备好 app_id 和 app_secret，此参数请登录微信小程序进行获取

```

## 运行
```
go run main.go

然后打开小程序开发工具，查看 network，检查接口连接情况
```

## 目录结构
```sh
conf            配置项
controller      控制器层（C）
log             日子类
middleware      中间件
model           模型（M）
service         服务层（业务逻辑封装）
public          静态资源
tests           测试类
utils           工具类
```

## 热加载
https://github.com/cosmtrek/air

## 联系
yigger#163.com