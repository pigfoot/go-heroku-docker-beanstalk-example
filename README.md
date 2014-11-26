## Initialization

Initialize project (godep + github + heroku)

### Prerequisite
Since godep needs specific folder structure, please be aware of following steps.
The keys are 1) godep needs project to be under "src" folder 2) project folder needs version control (like .git)
Reference: http://golang.org/doc/code.html

Your workspace maybe like (run godep save github.com/example_user/example_proj)

bin/
    streak                         # command executable
    todo                           # command executable
    example_proj                   # command executable (your project)
pkg/
    linux_amd64/
        code.google.com/p/goauth2/
            oauth.a                # package object
        github.com/nf/todo/
            task.a                 # package object
src/
    code.google.com/p/goauth2/
        .hg/                       # mercurial repository metadata
        oauth/
            oauth.go               # package source
            oauth_test.go          # test source
    github.com/nf/
        streak/
            .git/                  # git repository metadata
            oauth.go               # command source
            streak.go              # command source
        todo/
            .git/                  # git repository metadata
            task/
                task.go            # package source
            todo.go                # command source
    github.com/example_user/
        example_proj/
            .git/                  # git repository metadata (edit here)
            main.go                # command source
            main_test.go           # command source

to make life simple, we will also use gopm here.

### Setup folder structure && github (already create a project on github)
Github is optional, but it's easy if you would like to opensource it.
It's also ok if using private repository on github or bitbucket.

```bash
$ PROJ="go-heroku-docker-beanstalk-example"
$ REPO="https://github.com/pigfoot/${PROJ}"
$ go get -u github.com/tools/godep
$ go get -u github.com/gpmgo/gopm
$ mkdir -p ${PROJ}-workspace/src && cd ${PROJ}-workspace/src && git clone ${REPO} && cd ${PROJ}
$ export GOPATH=${PWD}/../../
$ export PATH="${GOPATH}/bin:${PATH}"
```

### Create template web service
```bash
$ cat << EOF > server.go
package main

import (
    "io"
    "io/ioutil"
    "net/http"

    "github.com/zenazn/goji"
)

func main() {
    goji.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "pong\n")
    })

    goji.Get("/time", func(w http.ResponseWriter, r *http.Request) {
        res, err := http.Get("http://localhost:8001/timegen")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        } else {
            defer res.Body.Close()
            cnt, _ := ioutil.ReadAll(res.Body)
            io.WriteString(w, string(cnt))
        }
    })

    // Listen and server on :8000 unless "PORT" environment variable is set
    goji.Serve()
}
EOF
$ echo "web: PORT=\$PORT ${PROJ}" > Procfile
$ gofmt -w server.go
$ godep get github.com/zenazn/goji
$ godep save ${PROJ}

### Create template worker

$ mkdir -p ${PROJ}-worker
$ cat << EOF > ${PROJ}-worker/main.go
package main

import (
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/zenazn/goji"
)

var (
    curTime time.Time
)

func main() {
    ticker := time.NewTicker(time.Second * 10)
    go func() {
        for t := range ticker.C {
            curTime = t
        }
    }()

    goji.Get("/timegen", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, fmt.Sprintf("%s\n", curTime))
    })

    // Listen and server on :8000 unless "PORT" environment variable is set
    goji.Serve()
}
EOF
$ echo "worker: PORT=8001 ${PROJ}-worker" >> Procfile
$ gofmt -w ${PROJ}-worker/main.go
$ godep save ${PROJ} ${PROJ}/${PROJ}-worker
$ git add .
$ git commit -m "Create template service (web+worker)"
```

### Build and deploy to heroku
```bash
$ heroku apps:create -b https://github.com/kr/heroku-buildpack-go.git ${PROJ}
$ git push heroku master
$ curl http://${PROJ}.herokuapp.com/ping
$ curl http://${PROJ}.herokuapp.com/time
```

#### Troubleshooting
```bash
$ heroku logs --tail
$ heroku run bash
```
