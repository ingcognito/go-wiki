# go-wiki
Go-Wiki is a Slack Bot written in Golang that retrieves Wikipedia information

## How to Use
Export Environment Variables to be used.
```
export GOWIKI_SLACK_TOKEN='xoxb-123456789012-123456789012-50rQ8TcaSdUoZDyAbHb0hjqs'
export GOWIKI_DB_CONFIG="user=postgres dbname=bot host=db port=5432 sslmode=disable"
```
If you do not have a Slack Token, please follow the steps in 'How to acquire a Slack Token'

Once your environment variables are exported, run
```
docker-compose up
```

- Go to your slack channel and click on + Add App

![Go to your slack channel and click on + Add App](docs/addapp.png "Add App")

- Search go-wiki and add to channel

![Search go-wiki and add to channel](docs/gowikiadd.png "Add to channel")

- Search a topic you would like to share with your team!
```
gowiki Toronto
```


### How to acquire a Slack Token
-Create a Slack App https://api.slack.com/apps?new_app=1

![Create a Slack App](docs/createslackapp.png "Create Slack App")

-Create a bot user

![Create a bot user](docs/addbotuser.png "Add Bot User")

-Install your app to your workplace

![Install your app to your workplace](docs/installappworkplace.png "Install App Work Place")

-Now go to OAuth & Permissions and copy the Bot User OAuth Access Token

![Now go to OAuth & Permissions and copy the Bot User OAuth Access Token](docs/oauthaccesstoken.png "OAuth Access Token")

Now you have your Slack Token!

### References
- Creating struct using https://mholt.github.io/json-to-go/
- Following standards https://github.com/golang-standards/project-layout