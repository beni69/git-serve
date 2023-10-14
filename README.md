# git-serve

serve the contents of your git repository over HTTP (most likely a static website),
with the correct mime types.

this is different from just putting a static server on the repo folder because you
can request different versions with the url

`/@rev/index.html` -> `rev` can be a branch or tag name, or a commit hash

## configuration

configuration is done through environment variables.

- PORT: http server port (default: 8080)
- REPO: path to your git repo (default: /src)
