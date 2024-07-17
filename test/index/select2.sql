select *
from user
where id < 5
    and id > 10
   or (id < 100 and id > 50);