# cmdscp
基于SSH协议的批量文件传输和命令执行工具

$ ./cmdscp -g demo2 -f test.txt -d "/home" -p /root/pwdfile
``` Bash
  ==================== 191.168.1.1 =======================  
Warning ssh connection failed

  ==================== 191.168.1.2 ======================= 
File sent successfully

  ==================== 191.168.1.3 ======================= 
File sent successfully
```

$ ./cmdscp -g demo2 -c "id"      
``` Bash
  ==================== 191.168.1.1 =======================  
Warning ssh connection failed

  ==================== 191.168.1.2 ======================= 
uid=0(root) gid=0(root) groups=0(root)

  ==================== 191.168.1.3 ======================= 
uid=0(root) gid=0(root) groups=0(root)
```
