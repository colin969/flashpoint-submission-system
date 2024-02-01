# flashpoint-submission-system

Submission management system for https://flashpointarchive.org/ or something like that. Now with all kinds of functionality duct-taped on the sides.

![spinny spin](static/crystal-spin.webm)

<img src="static/opal.png" alt="drawing" width="200"/>

## How to run this thing

- it's using discord for user auth, so you need a discord app with oauth config
- set up a discord bot to read user roles (FYI roles are hardcoded for Flashpoint discord server), roles are used for
  permission inside the system
- set up a discord bot to post notifications, can be the same bot as the previous one
- start a mysql instance, `make db` will do the work for you if you're a fan docker-compose
- start a curation validator server https://github.com/FlashpointProject/Curation-Validation-Bot (make command available
  in this repo)
- start an archive indexer if you want to upload stuff to
  flashfreeze https://github.com/Dri0m/recursive-archive-indexer (make command available in this repo)
- fill in all the stuff in .env (which is complex and needs more description here, yea)
- start the thing using `go run ./main/*.go`

it looks something like this

![submit page](github/ss2.png)

and this

![submissions page](github/ss3.png)

and this

![submission page](github/ss4.png)

and also this

![profile page](github/ss1.png)
