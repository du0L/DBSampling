#databases sampling

一个简单的MSSQL数据库取样工具

后期会加入MySQL支持，故未重构代码

没有写输出文件，用起来比较麻烦，需要重定向到一个文件

添加了过滤数据条数的功能


使用方法：
```
Usage of databases.exe:
  -IP string
        databases IP address , default is '127.0.0.1' (default "127.0.0.1")
  -Windows_verification
        use Windows verification(true or false), default is true (default true)
  -bypass int
        bypass data count,default is 0
  -password string
        databases password , default is 'password' (default "password")
  -port string
        databases port , default is 1433 (default "1433")
  -username string
        databases username , default is 'sa' (default "sa")
```

eg:

使用机器验证登陆数据库

    databases.exe -IP=192.168.0.9 -Windows_verification=true >out.html

使用用户名密码登陆数据库

    databases.exe -IP 192.168.0.9 -Windows_verification=false -username=sa -password=123456 >out.html

