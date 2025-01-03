
    My_favorable_scaffold
    │  main.go  主程序
    │  
    ├─conf 配置层
    │      config.yaml 配置文件
    │      
    ├─controller  处理层
    │      code.go 
    │      request.go
    │      response.go 封装response
    │      vaildator.go 
    │      
    ├─dao  数据访问层
    │  ├─mysql 
    │  │      mysql.go
    │  │      
    │  └─redis
    │          redis.go
    │          
    ├─logger 日志
    │      logger.go
    │      
    ├─middlewares 中间件
    │      auth.go
    │      ratelimit.go
    │      
    ├─models 用于存放获取请求参数的结构体
    ├─pkg 外部组件
    │  ├─encrypt
    │  │      encrypt.go 加密
    │  │      
    │  ├─jwt
    │  │      jwt.go 生成 jwt和解析 jwt
    │  │      
    │  └─snowflake
    │          snowflake.go 生成ID
    │          
    ├─router
    │      router.go 路由
    │      
    └─setting
            setting.go viper解析配置文件
            
