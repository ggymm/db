select *
from user
where id < 10
    and id > 5
   or (id < 100 and id > 50);