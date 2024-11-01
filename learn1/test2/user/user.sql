create table `user` (
      `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
      `user_id` BIGINT(20) not null,
      `username` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
      `password` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
      `email` VARCHAR(64) COLLATE utf8mb4_general_ci,
      `gender` TINYINT NOT NULL DEFAULT(0),
      `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
      `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`),
      UNIQUE KEY `idx_username` (`username`) USING BTREE,
      UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*
*   ---------- 创建表 ----------
*
*   1. 在 name.sql 中放置创建表的语句
* 
*   2. 在数据库中创建表
* 
*   3. 终端输入命令自动生成 crud 的代码块:
*        goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/zero" -table="user" -dir=./model 
*
*   ---------- 更改配置文件 ----------
*
*   4. api/config/config.go 增加一个字段:  Mysql struct { 
*                                             DataSource string 
*                                         }
*
*   5. api/etc/user-api.yaml 增加 mysql 对应的内容: Mysql:
*                                                    DataSource: root:123456@tcp(127.0.0.1:3306)/zero
*
*   6. api/internal/svc 文件的更改：  
*          -1. ServiceContext  结构体  加上相应的字段
*          -2. NewServiceContext 函数  返回上面加入的字段

*   ---------- 业务逻辑编写 ----------
*
*   7. 业务逻辑调用结构体中新增接口的对应的方法
*
*/