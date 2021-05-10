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

export GOPATH=/home/mumble/gocode
export GOBIN=/home/mumble/bin

sudo -u mumble bash << EOF
mkdir -p $GOPATH
mkdir -p $GOBIN
echo "path: $GOPATH"
cd $GOPATH
pwd
if [ ! -d $GOPATH/src/github.com/ramonaerne/talkiepi/ ]; then
    echo "go get github.com/ramonaerne/talkiepi"
fi
cd $GOPATH/src/github.com/ramonaerne/talkiepi
git pull
go build -o $GOBIN/talkiepi cmd/talkiepi/main.go
exit
EOF

if [ ! -f /boot/mumble_env.sh ]; then
	cat << EOF
No /boot/mumble_env.sh exists, placing default file from repo.
Make sure to adapt to your configuration!
EOF
	cp /home/mumble/gocode/src/github.com/ramonaerne/talkiepi/conf/boot/mumble_env.sh /boot/
fi

cp /home/mumble/gocode/src/github.com/ramonaerne/talkiepi/conf/systemd/mumble.service /etc/systemd/system/mumble.service
systemctl enable mumble.service
systemctl restart mumble.service
