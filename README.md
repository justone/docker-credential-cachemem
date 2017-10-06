# docker-credential-cachemem

Cache Docker credentials in memory

# Raw usage

```
cat sample/cred.json | ./docker-credential-cachemem store
echo foo | ./docker-credential-cachemem get
echo foo | ./docker-credential-cachemem erase
./docker-credential-cachemem list
```
