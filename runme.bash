#!/usr/bin/env bash

# postgres is running, right?
! pgrep postgres$ >/dev/null && echo "postgres is not running!" && exit 1

# did we already do this?
tmux list-sessions | grep '^lncddb-web' && echo "tmux attach -t lncddb-web" && exit 0

# we dont have our guys running somewhere else, right?
pgrep postgrest && echo "already running postgrest!" && exit 1
pgrep ./web$ && echo "already running web interface?" && exit 1

## start tmux container
tmux new-session -s lncddb-web -d;

## run postgrest: database api 

tmux new-window -t lncddb-web -n postgrest -d 'postgrest -j abcd -m 300 postgres://postgres@localhost:5432/lncddb -a postgres'
# if permissions were working in the elm code
#tmux new-window -t lncddb-web -n postgrest -d 'postgrest -j abcd -m 300 postgres://postgres@localhost:5432/lncddb -a web'

## run elm code

tmux new-window -t lncddb-web -n auth -d './web -secret abcd -dbrole lncd'

## attach to container
tmux attach -t lncddb-web

