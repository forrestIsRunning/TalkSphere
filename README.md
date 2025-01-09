
引入自己的脚手架
```bash
git clone git@github.com:Forrest-Tao/My_favorable_scaffold.git backend
```

```bash
git add .
git commit -m "Initial commit"
git push -u origin main
```

mysql环境
```bash
docker-compose -p mysql -f /Users/Zhuanz/go/src/forrest/TalkSphere/backend/deploy/mysql-docker-compose.yaml up -d
```
与数据库交互

```bash
➜  ~ docker ps -a | grep mysql
ab9cdd12c62a   mysql:8.0              "docker-entrypoint.s…"   43 hours ago   Up 43 hours   0.0.0.0:3306->3306/tcp, 33060/tcp

#进入docker容器
➜  ~ docker exec -it ab9cdd12c62a sh

#show databases
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| performance_schema |
| talksphere         |
+--------------------+
3 rows in set (0.01 sec)

#change database
mysql> use talksphere;
Database changed

mysql> select * from users \G
*************************** 1. row ***************************
        id: 1
created_at: 2025-01-03 19:10:27.900
updated_at: 2025-01-03 19:10:27.900
deleted_at: NULL
   user_id: 198927022297088
  username: yansaitao
  password: 3537313430307973742341bdf8249049c40e8e0ce7305e78e5471b
     email: yansaitao@qq.com
    avatar:
       bio:
*************************** 2. row ***************************
        id: 2
created_at: 2025-01-04 15:52:22.188
updated_at: 2025-01-04 15:52:22.188
deleted_at: NULL
   user_id: 481263710375936
  username: alice
  password: 313233bdf8249049c40e8e0ce7305e78e5471b
     email: email@gamil.com
    avatar:
       bio:
2 rows in set (0.00 sec)
```

## timeLine

- 2025.1.3 
  - 项目前端、后端骨架搭建
  - 登陆、注册前后端连调
- 2025.1.4
  - 实现JWT
  - 对接腾讯OSS，实现bucket创建和object的获取、删除、更新操作
  - 实现用户头像上传，bio更新
  - 表结构设计；简化业务
  - fix 登陆后，跳转forum失败问题
- 2025.1.5
  - 表结构设计
  - 业务功能设计
  - 板块模块的CRUD
  - 使用casbin完成 用户RBAC
- 2025.1.8
  - fix JWT
  - 帖子模块
    - 用户发表帖子
    - 根据id获取帖子详情
    - 根据id删除帖子
    - 根据id更新帖子
    - 获取某用户的所有帖子
    - 获取某板块下的所有帖子
- 2025.1.9
  - 交互模块
    - 评论
      - 根评论
      - 子评论
    - 点赞
      - 给贴子点赞
      - 给用户评论点赞
    - 收藏帖子