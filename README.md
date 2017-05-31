# Chef Restaurant
Chef Restaurant is a software helping the deployment of multiple Chef components on a centralized manner. It's meant to be typically run on a CI/CD agent.

It looks for the content of the last commit and based on that, it detects if the file that has been modified was a Chef role, environment, cookbook and uploads it on the Chef server.

It also bumps the version of Chef cookbooks to enforce versioning on them.

Chef Restaurant needs the following configuration settings to work properly :
- a file `config/config.json` needs to be configured with the following items :
  - a Slack hook URL where he can push the Slack notifications of what is being done;
  - a Slack channel where the notifications will be sent to;
  - a GitHub organization;
  - a GitHub repo.

Chef Restaurant also assumes the following softwares are already installed :
- Chef Knife with a configuration pointing to a Chef server
- Git

## How can I use it ?

Simply run the Chef Restaurant executable with no parameter on the repo where your Chef Git repo is.

```
./chef-restaurant
```

If you want to run an early commit using Chef Restaurant, you can use the following parameter to do so.

```
./chef-restaurant -commit <commitID>
```
