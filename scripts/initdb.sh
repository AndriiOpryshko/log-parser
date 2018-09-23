#!/bin/bash

docker-compose up -d

echo "sleeps 10 second"
sleep 10


ContainerName=scripts_logserv_mongodb_1
User=logger
Pass=eiko1Aeraiceechi

docker exec -i $ContainerName mongo admin <<EOF
use admin
var user = {
  "user" : "$User",
  "pwd" : "$Pass",
  roles : ["userAdminAnyDatabase", "dbAdminAnyDatabase", "readWriteAnyDatabase"]
}
db.createUser(user);
exit
EOF

echo "add user $User with pas $Pass for auth db admin"
