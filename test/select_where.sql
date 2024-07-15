select *
from user
where extras < 3
  and (id < 2 or id = 1 or id = 3 or (id = 1 and id = 2))
  and (username = "名称1" or nickname = "昵称2" or email = "邮箱3" or email = "邮箱4");