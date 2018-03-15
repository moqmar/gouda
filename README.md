# Gouda.io
> The no-worries documentation tool

![](https://static.mo-mar.de/gouda.png)

---

Gouda (short for **Go** **U**nified **D**ocumentation **A**ggregator) is a static documentation generator that's using Markdown, is extendable using plugins and can even import existing documentation from JSDoc, GoDoc and more, and give it an uniform look.  
Besides the static generation, there will be a dynamic mode with a feature set similar to a lot of wiki software, which will also be available on a hosted plan.

## Download
Wait a moment, we just started out. Give us a star and put the repo on your watch list to be the first to use gouda once we're ready!  
Binaries and documentation on installation and usage will be available with version 0.1. You can see on the [roadmap](ROADMAP.md) how far we already are.

## Planned features
- **Online WYSIWYG editor** (with git integration)
- **Fast and advanced search**; lightweight search for statically generated documentation
- **Multi-Language**
- **Fully automated index and sitemap generation**, with pagination ("next"/"previous" page) for book-style documentation
- **Versioning** (using git commits/tags)
- **PDF, ePub and MHTML generation**
- **Deployment to GitHub Pages, Surge, SSH, SFTP, rsync, FTP and many more**
- **Custom CSS and templates**
- Lots of awesome Markdown plugins that make the creation of beautiful documentation a breeze
  - **Interactive REST toolbox** - prepare request templates, and let users modify and send them. Will also support various methods for authentication and shared fields.
  - **LaTeX Math**
  - **Graphing/Flow Charts/Sequence Diagrams/...** (with a graphical editor in the WYSIWYG editor)
  - **Links to GitHub**, GitLab, Gogs, or Taiga artifacts (autocompletion in editor)
  - **Hashtags** (autocompletion in editor)
  - **Global bibliography** with citation/references (autocompletion in editor)

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

<details><summary>example 'gouda.yml'</summary>

```yaml

title: My awesome program

target: ./html

links:
  "/developers/go":
    godoc: ../src
  "/": ../README.md

exclude:
  - /excluded-folder/
  - /you-can-use-gitignore-syntax.txt

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
</details>

## The application
### Static mode
```
$ gouda build
$ gouda serve
$ gouda deploy [target=*]
$ gouda edit # Start dynamic mode without authentication
```

### Dynamic/Server mode
```
$ gouda dynamic [host=::] [port=8080]
$ docker run --name gouda -v "$PWD":/data -p 8080:80 gouda/dynamic
```

## License

    Gouda - the no-worries documentation tool
    Copyright (C) 2018  Moritz Marquardt & Frederik Kammel

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
