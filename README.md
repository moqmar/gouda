# Gouda.io
The no-worries documentation tool

Gouda (short for **GO** **U**nified **D**ocumentation **A**ggregator) is a static documentation generator that takes simple Markdown files and creates a nice looking static webpage from it.

## Download
Give us a sec, we just started out. Give us a star and put the repo on your watch list to be the first to download gouda once we're ready. 

## Planned features
- Host your own website (free)
- Hosted plan (`<project>.gouda.io`)
    - 2€ per month and project (any team size)
    - 3€ with own domain/deployment location
- Online WYSIWYG editor (with git integration)
- custom CSS, but no custom templates (→ consistency)
- search
    - dynamic index
    - static index
- multi language
- versioning (git)
- GitHub integration (links to issues)

## Typical repository structure
```

+ docs/
  + examples.md
  + developers/
    + index.md
    + api.md
  + gouda.yaml
+ src/
  + ...
+ README.md
```

### gouda.yml
```yaml

title: ...

target: ./html

links:
  "/developers/go":
    godoc: ../src
  "/": ../README.md

before:
  - echo "Generating documentation..."
after:
  - echo "Throw the cheese!!!"

deploy:
  github:
    branch: gh-pages
    cname: docs.example.org
  exampleServer:
    ssh: test@example.org:22
    path: ~/docs
    keyfile: ~/.ssh/zwiebelsuppe
    before:
      - "rm -rf *"
    after:
      - "nginx -s reload"
  exampleServerViaRsync:
    rsync: test@example.org:22
    path: ~/docs
  ftpDeploy:
    ftp: test:[password]@example.org:21
  s3Deploy:
    s3: my-docs-bucket
  zeitNowDeploy:
    now: example-qadmntdfnh.now.sh
  surgeDeploy:
    surge: docs.example.org
```

## The application
### Static mode
```

$ gouda build
$ gouda serve
$ gouda deploy [target]
$ gouda edit # Start dynamic mode without authentication
```

## Dynamic/Server mode
```
$ gouda-dynamic [host] [port]
$ docker run --name gouda -v "$PWD":/data -p 8080:80 gouda/dynamic
```
