#
# Poser REST API server Dockerfile
#
# Golang/Martini web server running on a 'scratch' image
#
# https://github.com/pablosan/poser
#

FROM google/golang

# Expose /gopath/bin so we can get the built binary
VOLUME /gopath/bin

CMD CGO_ENABLED=0 go get -a -ldflags '-s' github.com/pablosan/poser
