# go-wiki
Go-Wiki is a Slack Bot written in Golang that retrieves Wikipedia information

##How to Use

Export Environment Variables to be used.
```
export GOWIKI_SLACK_TOKEN='xoxb-123456789012-123456789012-50rQ8TcaSdUoZDyAbHb0hjqs'
export GOWIKI_DB_CONFIG="user=postgres dbname=bot host=db port=5432 sslmode=disable"
```

Once your environment variables are exported, run
```
docker-compose up
```

### References
- Creating struct using https://mholt.github.io/json-to-go/
- Following standards https://github.com/golang-standards/project-layout