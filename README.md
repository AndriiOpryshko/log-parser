# log-parser

##  Preparation
- Command: ./script/initdb.sh

- Description: It runs mongo db by docker-compose and inits auth with creds 
from Vars User and Pas in shell script.

## Running 
- Command: ./script/run.sh

- Description: It rund mongo db by docker-compose, gets libs and runs
 parser.

## Config file
- Path: ./config.yml

- Description: consists of mongo config and config of parsing files

## Test log files
- Path: ./testlogs

- Descriptio: file 1.log formated by first format and 2.log formated by
 second format

