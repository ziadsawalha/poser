poser
=====

Poser is modeled after VCR, an HTTP interactions testing tool [originally created in Ruby](https://github.com/vcr/vcr) and also [ported to Python](https://github.com/kevin1024/vcrpy).

So why re-invent the wheel? Having now built several Docker-fied REST API web servers I thought it would be useful to have a very small docker container that could stand-in for my external HTTP dependencies (e.g. auth, downstream http services, whatever).

While scratching that particular itch, I also realized this could fit another use case: REST API mock-ups, a la [apiary.io](http://apiary.io/) and [mockable.io](http://www.mockable.io/).

So Poser aims to provide a small, fast, canned request/response web server that can be used either to eliminate external dependencies for reliable/fast testing or to quickly mock up a web API to test new ideas.

Poser is written in Go so it's pretty snappy. And the Docker container is tiny (~5.3 Megabytes)! This means pulling the image only takes a few seconds and spinning up a new container takes milliseconds.

### Installing Poser

  * Using Go:
    * Make sure you have [go](http://golang.org/) installed (for Mac users I strongly recommend HomeBrew: just `brew install go` or `port install go` if you use macports)
    * run `go get github.com/pablosan/poser`
    * Then run `$GOPATH/bin/poser -scenes path/to/your/hello.yaml -port <some_legal_port>`
  * Build from source:

    ```shell
    git clone git@github.com:pablosan/poser.git
    cd poser
    go install
    $GOPATH/bin/poser -scenes path/to/your/hello.yaml -port <some_legal_port>
    ```

  * Using Docker

    `docker run -d -p 8080:3000 -v /path/to/your/hello.yaml:/var/scenes.json --name poser pablosan/poser -scene /var/scenes.json`

  Then make a call to poser:

    ```shell
    $ curl http://localhost:3000/hello -i
    host:3000/hello -i
    HTTP/1.1 200 OK
    Date: Fri, 05 Sep 2014 03:58:32 GMT
    Content-Length: 26
    Content-Type: text/plain; charset=utf-8

    {"message":"Hello World!"}
    ```

If you're wondering what goes in the scenes file, check out the [examples](examples).

### Building the Docker container

Building the final Docker container is a two step process. The first `docker build` is done from within the `linux-bin` directory. This build relies on the google/golang container and statically compiles the `poser` binary.

Once the binary is built, the second `docker build` command is done from within the project's root directory. This is where the magic happens: using a scratch image means this Docker image starts at a size of zero bytes. Because we statically compiled the binary there are no dependencies. So we can simply add the binary and one other static file (favicon.ico) to the image. This is how `poser` keeps its svelte figure.

Because of this process, __building the Docker image must occur on a Docker or CoreOS host OS__ (you cannot run this from the Mac OS command line, e.g. using boot2docker).

Once you are on your Docker/CoreOS machine (i.e. ssh'ed into a Docker or CoreOS instance), from the project's root directory:

  ```shell
  cd linux-bin
  ./build.sh
  cd ..
  docker build --rm=true --no-cache=true -t pablosan/poser .
  ```

Now you can launch the Docker container without performing a `docker pull`.

### Using Poser

__Command Line Arguments:__

  * `-scenes <path/to/scenes.(json|yaml)>` (Optional): Specifies the scenes file poser should use for read/write. Defaults to `scenes.json`.
  * `-port <some-legal-port>` (Optional): Specifies the HTTP port that should be used. Defaults to `3000`.
  * `-play[=false]` (Optional*): Enable playback mode. This mode requires at least one valid scene to be present in the scenes file. Defaults to `false`.
  * `-record[=false] (Optional*)`: Enable record/provy mode. This mode requires either `baseurl` to be defined at the root of the scenes document, or the `-proxy` parameter. Poser will proxy any requests to the proxy url provided, proxy the response back to the client, and write a new scene to the scenes file. Defaults to `false`.
  * `-proxy <http://some.valid.url.com[:<some-valid-port>]>` (Optional): Specifies the URL to which requests should be proxied. Only applies to `-record` mode. Defaults to "".

_*NOTE: At least one of `-play` or `-record` are required. If both are specified, poser will give precedence to `-play`: i.e. it will only proxy the request when no matching pre-recorded scene is found._
