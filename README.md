
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

