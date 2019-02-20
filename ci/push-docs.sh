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
branch="docs-glooe-$tag"

git clone https://github.com/solo-io/gloo-docs.git
git config --global user.name "soloio-bot"
(cd gloo-docs && git checkout -b $branch)

if [ -d "gloo-docs/docs/v1/github.com/solo-io/solo-projects" ]; then
	rm -r gloo-docs/docs/v1/github.com/solo-io/solo-projects
fi
cp -r projects/gloo/doc/docs/v1/github.com/solo-io/solo-projects gloo-docs/docs/v1/github.com/solo-io/solo-projects

rm gloo-docs/docs/cli/glooctl*
cp projects/gloo/doc/docs/cli/glooctl* gloo-docs/docs/cli/

(cd gloo-docs && git add .)

if [[ $( (cd gloo-docs && git status --porcelain) | wc -l) -eq 0 ]]; then
  echo "No changes to solo-projects docs, exiting."
  rm -rf gloo-docs
  exit 0;
fi

(cd gloo-docs && git commit -m "Add docs for tag $tag")
(cd gloo-docs && git push --set-upstream origin $branch)

curl -v -H "Authorization: token $github_token_no_spaces" -H "Content-Type:application/json" -X POST https://api.github.com/repos/solo-io/gloo-docs/pulls -d \
'{"title":"Update docs for glooe '"$tag"'", "body": "Update docs for glooe '"$tag"'", "head": "'"$branch"'", "base": "master"}'

rm -rf gloo-docs
