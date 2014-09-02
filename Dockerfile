#
# Poser REST API server Dockerfile
#
# Golang/Martini web server running on a 'scratch' image
#
# https://github.com/pablosan/poser
#

FROM scratch

ADD bin/poser /poser

CMD ["/poser"]
