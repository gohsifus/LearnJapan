update users
set login    = 'gohsifus',
    password = 'qawsed345rf',
    email    = 'gohsifus@gmail.com'
where id = 49;

delete
from users
where id <> 49;


