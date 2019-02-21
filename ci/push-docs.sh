#!/usr/bin/env bash

# Requires $tag shell variable and $GITHUB_TOKEN environment variable

set -e
xargs=$(which gxargs || which xargs)

# Validate settings.
[ "$TRACE" ] && set -x

CONFIG=$@

for line in $CONFIG; do
  eval "export ${line}"
done

github_token_no_spaces=$(echo $GITHUB_TOKEN | tr -d '[:space:]')
branch="docs-gloo-$tag"

set +x
echo "Cloning solo-docs repo"
git clone https://soloio-bot:$github_token_no_spaces@github.com/solo-io/solo-docs.git
[ "$TRACE" ] && set -x

git config --global user.name "soloio-bot"
(cd solo-docs && git checkout -b $branch)

# Gloo
if [ -d "solo-docs/gloo/docs/v1/github.com/solo-io/gloo" ]; then
	rm -r solo-docs/gloo/docs/v1/github.com/solo-io/gloo
fi
cp -r docs/v1/github.com/solo-io/gloo solo-docs/gloo/docs/v1/github.com/solo-io/gloo

# Solo Kit
if [ -d "solo-docs/gloo/docs/v1/github.com/solo-io/solo-kit" ]; then
	rm -r solo-docs/gloo/docs/v1/github.com/solo-io/solo-kit
fi
cp -r docs/v1/github.com/solo-io/solo-kit solo-docs/gloo/docs/v1/github.com/solo-io/solo-kit

# Gogoproto
if [ -d "solo-docs/gloo/docs/v1/gogoproto" ]; then
	rm -r solo-docs/gloo/docs/v1/gogoproto
fi
cp -r docs/v1/gogoproto solo-docs/gloo/docs/v1/gogoproto

# Google
if [ -d "solo-docs/gloo/docs/v1/google" ]; then
	rm -r solo-docs/gloo/docs/v1/google
fi
cp -r docs/v1/google solo-docs/gloo/docs/v1/google

(cd solo-docs && git add .)

if [[ $( (cd solo-docs && git status --porcelain) | wc -l) -eq 0 ]]; then
  echo "No changes to gloo docs, exiting."
  rm -rf solo-docs
  exit 0;
fi

(cd solo-docs && git commit -m "Add docs for tag $tag")
(cd solo-docs && git push --set-upstream origin $branch)

curl -v -H "Authorization: token $github_token_no_spaces" -H "Content-Type:application/json" -X POST https://api.github.com/repos/solo-io/solo-docs/pulls -d \
'{"title":"Update docs for gloo '"$tag"'", "body": "Update docs for gloo '"$tag"'", "head": "'"$branch"'", "base": "master"}'

rm -rf solo-docs
