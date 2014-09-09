docker build --rm=true -t pablosan/poser-linux .
docker run -it -v /var/bin:/gopath/bin --name poser-linux pablosan/poser-linux
mkdir -p ../bin
cp /var/bin/poser ../bin/poser
docker rm poser-linux
docker rmi pablosan/poser-linux
