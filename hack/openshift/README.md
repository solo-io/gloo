# running locally on minishift

to update containers locally, you'll want to do the following:

1. build the container with `recompile.sh PROJECTNAME TAG`
2. change `Always` to `IfNotPresent` in the template and re-apply with apply.sh
3. Delete the desired pod if it doesn't autoamtically restart.
