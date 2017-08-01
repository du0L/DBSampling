# DBSampling

一个简单的MSSQL数据库取样工具

后期会加入MySQL支持，故未重构代码

没有写输出文件，用起来比较麻烦，需要重定向到一个文件

使用方法：
databases.exe -help
  -IP string
    	databases IP address , default is '127.0.0.1' (default "127.0.0.1")
  -Windows_verification
    	use Windows verification(true or false), default is true (default true)
  -password string
    	databases password , default is 'password' (default "password")
  -username string
    	databases username , default is 'sa' (default "sa")
exit status 2

eg:
使用机器验证登陆数据库
databases.exe -IP=192.168.0.9 -Windows_verification=true >out.html
使用用户名密码登陆数据库
databases.exe -IP 192.168.0.9 -Windows_verification=flase -uesranme=sa -password=123456 >out.htlm

