## Distribution
Summary of distribution process as it exists right now

#### Release
Distribution is kicked off during release builds of solo-projects. The `distribution` make target currently only calls 
a go script located in the distribution directory. This script has 2 main purposes:
1) consolidate all necessary release files and zip them 
2) save all of the local data to google storage

There are three main things being saved into google storage:
1) folder containing all uncompressed release assets
2) a folder with a randomly generated name containing a zipped version of the release assets
3) an index file which tracks the release version with the random folder name it is saved to

The release file located in the randomly generated folder is publicly accessible so it is important that no inference 
can be drawn from one release to the next in terms of location 

#### Testing

The entire distribution script is located in a single go script and can therefore can be run locally for testing/development purposes. 
Since the process interacts directly with cloud storage, to test the remote features gcloud credentials have to be present locally.
The authentication flow is explained [here](https://cloud.google.com/docs/authentication/production). Once the `GOOGLE_APPLICATION_CREDENTIALS` 
is set the program will attempt the remote features. The program only takes 1 argument and that is the version to be saved.

```bash
go run $GOPATH/src/github.com/solo-io/solo-projects/install/distribution <version>
```
If run like this, the helm templates have to be available beforehand, if run from the make target directly, it will take care of that.
```bash
make distribution TAGGED_VERSION=v<version>
```

##### Notes
In order to properly test the environment like the one in the cloud build run `make clean` to clear the directory of all generated files