# starcloud

CDN for Starboard assets.

Current version only serves `starboard-notebook@0.13.2`, `starboard-notebook@0.14.0`, `starboard-notebook@0.14.1`.

## Deploying

Install `flyctl`, on a Mac you can use homebrew:
```shell
brew install superfly/tap/flyctl
```

Alternatively use this on Linux or Mac:

```shell
curl -L https://fly.io/install.sh | sh
```

Then log in:
```shell
flyctl auth login
```

And then deploy :) Easy as that.
```shell
flyctl deploy
```