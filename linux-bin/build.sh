docker build --rm=true -t pablosan/poser-linux .
docker run -it -v /var/bin:/gopath/bin --name poser-linux pablosan/poser-linux
cp /var/bin/poser ../bin/linux/poser
docker rm poser-linux
docker rmi pablosan/poser-linux
