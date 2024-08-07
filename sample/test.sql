--创建表
create table 'user'
(
    'user_id'      INT64 NOT NULL,
    'username'     VARCHAR        DEFAULT NULL,
    'nickname'     VARCHAR        DEFAULT NULL,
    'email'        VARCHAR        DEFAULT NULL,
    'phone_number' VARCHAR        DEFAULT NULL,
    'account'      VARCHAR        DEFAULT NULL,
    'password'     VARCHAR        DEFAULT NULL,
    'status'       INT32          DEFAULT NULL,
    'extras'       VARCHAR        DEFAULT NULL,
    'create_time'  VARCHAR        DEFAULT NULL,
    'create_id'    INT64 NOT NULL DEFAULT '1',
    'update_time'  VARCHAR        DEFAULT NULL,
    'update_id'    INT32 NOT NULL DEFAULT '1',
    'del_flag'     INT32 NOT NULL DEFAULT '1',
    PRIMARY KEY ('user_id'),
    INDEX          ACCOUNT_INDEX ('account')
);


--查看表
show tables;


--插入测试数据
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('1', '名称1', '昵称1', '邮箱1', '手机号1', '账号1', '密码1', '1', '额外信息1', '创建时间1', '1', '更新时间1', '1', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('2', '名称2', '昵称2', '邮箱2', '手机号2', '账号2', '密码2', '2', '额外信息2', '创建时间2', '2', '更新时间2', '2', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('3', '名称3', '昵称3', '邮箱3', '手机号3', '账号3', '密码3', '3', '额外信息3', '创建时间3', '3', '更新时间3', '3', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('4', '名称4', '昵称4', '邮箱4', '手机号4', '账号4', '密码4', '4', '额外信息4', '创建时间4', '4', '更新时间4', '4', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('5', '名称5', '昵称5', '邮箱5', '手机号5', '账号5', '密码5', '5', '额外信息5', '创建时间5', '5', '更新时间5', '5', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('6', '名称6', '昵称6', '邮箱6', '手机号6', '账号6', '密码6', '6', '额外信息6', '创建时间6', '6', '更新时间6', '6', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('7', '名称7', '昵称7', '邮箱7', '手机号7', '账号7', '密码7', '7', '额外信息7', '创建时间7', '7', '更新时间7', '7', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('8', '名称8', '昵称8', '邮箱8', '手机号8', '账号8', '密码8', '8', '额外信息8', '创建时间8', '8', '更新时间8', '8', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('9', '名称9', '昵称9', '邮箱9', '手机号9', '账号9', '密码9', '9', '额外信息9', '创建时间9', '9', '更新时间9', '9', '1');
insert into user ('user_id', 'username', 'nickname', 'email', 'phone_number', 'account', 'password', 'status', 'extras', 'create_time', 'create_id', 'update_time', 'update_id', 'del_flag')
value ('10', '名称10', '昵称10', '邮箱10', '手机号10', '账号10', '密码10', '10', '额外信息10', '创建时间10', '10', '更新时间10', '10', '1');


--查询表
select * from user;
select * from user where user_id < 5;
select * from user where user_id < 5 and username = "名称1";


--更新表
update user set username = "名称1-修改" where user_id = 1;


--删除表
delete from user where user_id = 8;
