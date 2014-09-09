poser
=====

Poser is modeled after VCR, an HTTP interactions testing tool [originally created in Ruby](https://github.com/vcr/vcr) and also [ported to Python](https://github.com/kevin1024/vcrpy).

So why re-invent the wheel? Having now built several Docker-fied REST API web servers I thought it would be useful to have a very small docker container that could stand-in for my external HTTP dependencies (e.g. auth, downstream http services, whatever).

While scratching that particular itch, I also realized this could fit another use case: REST API mock-ups, a la [apiary.io](http://apiary.io/) and [mockable.io](http://www.mockable.io/).

So Poser aims to provide a small, fast, canned request/response web server that can be used either to eliminate external dependencies for reliable/fast testing or to quickly mock up a web API to test new ideas.

Poser is written in Go so it's pretty snappy. And the Docker container is tiny (<5 Megabytes)! This means pulling the image only takes a few seconds and spinning up a new container takes milliseconds.

### Installing Poser

  * Using Go:
    * Make sure you have [go](http://golang.org/) installed (for Mac users I strongly recommend HomeBrew: just `brew install go` or `port install go` if you use macports)
    * run `go get github.com/pablosan/poser`
    * Then run `$GOPATH/bin/poser -scenes path/to/your/hello.yaml -port <some_legal_port>`
  * Build from source:
    * `git clone git@github.com:pablosan/poser.git`
    * `cd poser`
    * `go install`
    * `$GOPATH/bin/poser -scenes path/to/your/hello.yaml -port <some_legal_port>`
  * Using Docker
    * `docker run -d -p 8080:3000 -v /path/to/your/hello.yaml:/var/scenes.json --name poser pablosan/poser -scene /var/scenes.json`

  Then make a call to poser:

    ```
    $ curl http://localhost:3000/hello -i
    host:3000/hello -i
    HTTP/1.1 200 OK
    Date: Fri, 05 Sep 2014 03:58:32 GMT
    Content-Length: 26
    Content-Type: text/plain; charset=utf-8

    {"message":"Hello World!"}
    ```

If you're wondering what goes in the scenes file, check out the [examples](examples).
