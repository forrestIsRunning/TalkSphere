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