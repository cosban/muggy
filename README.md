# muggy
>A terribly written irc bot written in golang

Created for the sake of learning and memes

## Installation and Running
>Prerequisite: go must be installed, you must also know how to use it because
I sure as hell don't

To install, simply run the following

    $ go get github.com/cosban/muggy
    $ cd $GOPATH/src/github.com/cosban/muggy && go install

Once this has been done, you may run it within any directory that you also have
placed the config.ini into. You may find an example config within configs/. To
Run, perform the following command:

    $ muggy

This is excellent for running inside screen sessions and whatnot. The preferred
method to run is in the following manner:

    #!/bin/bash
		if [ "$(pidof muggy)" ]
		then
		    killall muggy
		fi
		muggy > log 2> err &

A running example can be found within #muggy on [freenode](irc.freenode.net) which will always be running the latest version of muggy.

### TODO: All the things
