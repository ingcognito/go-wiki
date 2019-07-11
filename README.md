# go-wiki
Go-Wiki is a Slack Bot written in Golang that retrieves Wikipedia information.

<img src="docs/helloworld.png" width="500">

## Setup
---

**Dependencies**

- You will need `docker` and `docker-compose` installed. 
https://docs.docker.com/install/

**Configuration**

- Export Environment Variables to be used.
[*If you do not have a Slack Token, please follow the steps in 'How to acquire a Slack Token'* :point_down:](https://github.com/noah-ing/go-wiki#how-to-acquire-a-slack-token)
```
export GOWIKI_SLACK_TOKEN='xoxb-123456789012-123456789012-50rQ8TcaSdUoZDyAbHb0hjqs'
export GOWIKI_DB_CONFIG="user=postgres dbname=bot host=db port=5432 sslmode=disable"
```

**Build**
- Once your environment variables are properly configured, go to your terminal and run
```
docker-compose up
```
This will enable your server to run with a postgres db. You will beable to see the logs.

**Go to your slack channel and click on + Add App**

<img src="docs/addapp.png" width="300">

**Search go-wiki and add to channel**

<img src="docs/gowikiadd.png" width="300">


**Use the command `go wiki` to search a topic you would like to share with your team!** :rocket: 
```
gowiki Toronto
```
go-wiki stores each call to a database for caching so the next time around it will show up without the need to call Wikipedia.
:squirrel:	

## How to acquire a Slack Token
---
**Below are steps to aquire a slack token necessary for this golang application to work.**

1. **Create a Slack App** https://api.slack.com/apps?new_app=1

<img src="docs/createslackapp.png" width="300">


2. **Create a bot user**

<img src="docs/addbotuser.png" width="300">


3. **Install your app to your workplace**

<img src="docs/installappworkplace.png" width="300">


4. **Now go to OAuth & Permissions and copy the Bot User OAuth Access Token**

<img src="docs/oauthaccesstoken.png" width="300">


[:clap: **Congratulations! :confetti_ball: Now you have your Slack Token! Now go back up to Setup!**](https://github.com/noah-ing/go-wiki#setup)

### Further Development
---
I have left the entire struct of a wikipedia page in main.go. Please feel free to call anything that your heart desires from Wikipedia to your Slack Channel :heart:. Just remember to *add the rows to the pages table in the database, update the db query and response!*


### References
---
- Creating struct using https://mholt.github.io/json-to-go/
