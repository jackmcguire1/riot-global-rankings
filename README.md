# riot-global-rankings
> This repository is for the AWS x League Of Legends Hackathon - Global Power Rankings - LoL Esports Data

[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/riot-global-rankings)](https://goreportcard.com/report/github.com/jackmcguire1/riot-global-rankings)
[![codecov](https://codecov.io/gh/jackmcguire1/riot-global-rankings/graph/badge.svg?token=CWCV1WAKMM)](https://codecov.io/gh/jackmcguire1/riot-global-rankings)
# About
- This repository contains the API for the Global Power Rankings Hackathon sponsored by AWS x Riot


- Team rankings are based on accumulated wins and losses of tournament matches where the type is 'BestOf' and mode is 'Classic'.

# Devpost Project Submission
> https://devpost.com/software/riot-rankings

# Requirements
- MongoDB ATLAS cluster - AWS - US-EAST-1
- mongodb/brew/mongodb-database-tools
- Python3
- GO 1.21+
- AWS SAM [lambda, API Gateway, etc]
- Docker Desktop 4.24+
- Docker Compose 2.22.0+

# SETUP
- setup an MongoDB Atlas cluster https://www.mongodb.com/cloud/atlas/register
- > go mod download

## CLI tools
>  brew install mongodb/brew/mongodb-database-tools

## Swagger
 Copy the contents of ./swagger.yml to [Swagger Editor](https://editor-next.swagger.io/) and begin invoking the Serverless API endpoints.
> https://app.swaggerhub.com/apis/jackmcguire1/Riot-Global-Rankings/1.0.1
## Run With Docker
> docker-compose up -d

> docker compose watch

## Import data to Mongo

1. run ```python3 ./tools/get_riot_files.py```
2. create a database on your MongoDB Atlas cluster called ```riot```
3. > cd ./esports-data
4. >mongoimport --uri 'mongodb+srv://{ADDRESSS}' --collection tournaments --type json --file tournaments.json --jsonArray
5. > mongoimport --uri 'mongodb+srv://{ADDRESSS}' --collection mappings --type json --file mapping_data.json --jsonArray
6. > mongoimport --uri 'mongodb+srv://{ADDRESSS}' --collection leagues --type json --file leagues.json --jsonArray7.
7. > mongoimport --uri 'mongodb+srv://{ADDRESSS}' --collection teams --type json --file teams.json --jsonArray
8. > mongoimport --uri 'mongodb+srv://{ADDRESSS}' --collection players --type json --file players.json --jsonArray

# Relevant links
- [Devpost Project Submission](https://devpost.com/software/riot-rankings)
- [Dev Post Hackathon](https://lolglobalpowerrankings.devpost.com/)
- [League Esports Data Guide](https://docs.google.com/document/d/1wFRehKMJkkRR5zyjEZyaVL9H3ZbhP7_wP0FBE5ID40c/edit)
- [League Esports API GUIDE ](https://docs.google.com/document/d/1Klodp4YqE6bIOES026ecmNb_jS5IOntRqLv5EmDAXyc/edit)