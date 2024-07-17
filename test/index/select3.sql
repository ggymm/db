select *
from user
where id < 50
   or id > 100
   or (id = 60 or id = 70 or id = 80);