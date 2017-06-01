ssh -i ../keys/vultr root@prod.wherebaryou.com "cd /home/go/src/github.com/mgerb/wbu-server && sqlite3 database.db '.backup database.bak'"
scp -i ../keys/vultr root@prod.wherebaryou.com:/home/go/src/github.com/mgerb/wbu-server/database.bak ./database.bak;

git add --all;

git commit -m "New backup: $(date)";
