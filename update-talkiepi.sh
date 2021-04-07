#! /bin/bash

if [[ $EUID -ne 0 ]];
then
    exec sudo /bin/bash "$0" "$@"
fi

if [ ! -d /home/mumble ]; then
    adduser --disabled-password --disabled-login --gecos "" mumble
    usermod -a -G cdrom,audio,video,plugdev,users,dialout,dip,input,gpio mumble
    apt-get install golang libopenal-dev libopus-dev git
fi

sudo -i -u mumble bash << EOF
mkdir ~/gocode
mkdir ~/bin   
export GOPATH=/home/mumble/gocode
export GOBIN=/home/mumble/bin
cd $GOPATH
pwd
if [ ! -d src/github.com/ramonaerne/talkiepi/ ]; then
    go get github.com/ramonaerne/talkiepi
fi
cd $GOPATH/src/github.com/ramonaerne/talkiepi
git pull
go build -o /home/mumble/bin/talkiepi cmd/talkiepi/main.go 
exit
EOF

cp /home/mumble/gocode/src/github.com/ramonaerne/talkiepi/conf/systemd/mumble.service /etc/systemd/system/mumble.service
systemctl enable mumble.service
systemctl restart mumble.service
