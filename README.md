# cmdscp
基于SSH协议的批量文件传输和命令执行工具

$ ./cmdscp -g demo2 -f test.txt -d "/home/t.txt"

  ==================== 139.129.212.98 =======================  
Warning ssh connection failed

  ==================== 139.129.239.41 ======================= 
File sent successfully

  ==================== 139.129.237.40 ======================= 
File sent successfully


$ ./cmdscp -g demo2 -c "id"        

  ==================== 139.129.212.98 =======================  
Warning ssh connection failed

  ==================== 139.129.239.41 ======================= 
uid=0(root) gid=0(root) groups=0(root)

  ==================== 139.129.237.40 ======================= 
uid=0(root) gid=0(root) groups=0(root)
