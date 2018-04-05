#!/usr/bin/env groovy
podTemplate(label: 'gloo-builder',
containers: [
    containerTemplate(
        name: 'golang',
        image: 'golang:1.10',
        ttyEnabled: true,
        command: 'cat'),
    containerTemplate(
        name: 'docker',
        image: 'docker:17.12',
        ttyEnabled: true,
        command: 'cat'),
    containerTemplate(
        name: 'envoy-build',
        image: 'envoyproxy/envoy-build-ubuntu@sha256:6153d9787cb894c2dd6b17a1539eaeba88ae15d79f66f63eec0f4713436d74f0',
        ttyEnabled: true,
        command: 'cat'),
    containerTemplate(
        name: 'helm',
        image: 'devth/helm:2.8.1',
        ttyEnabled: true,
        command: 'cat'),
],
envVars: [
    envVar(key: 'DOCKER_CONFIG', value: '/etc/docker')
],
volumes: [
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
    secretVolume(secretName: 'soloio-docker-hub', mountPath: '/etc/docker'),
    secretVolume(secretName: 'soloio-github', mountPath: '/etc/github')
]) {

    properties([
        parameters ([
            stringParam(
                defaultValue: 'cf08737718cf62bf597f88aa2068c6f6b28b9992',
                description: 'Commit hash for gloo',
                name: 'GLOO_HASH'),
            stringParam(
                defaultValue: '282a844ea3ed2527f5044408c9c98bc7ee027cd2',
                description: 'Commit hash for gloo-plugins',
                name: 'GLOO_PLUGINS_HASH'),
            stringParam(
                defaultValue: 'v0.1.6',
                description: 'Version used for Docker images',
                name: 'GLOO_VERSION'),
            booleanParam(
                defaultValue: false,
                description: 'Update gloo-install',
                name: 'UPDATE_INSTALL')
        ])
    ])

    node('gloo-builder') {
        stage('thetool') {
            container('golang') {
                echo 'Initialize thetool...' 
                sh '''
                    go get -u github.com/solo-io/thetool
                    CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' github.com/solo-io/thetool
                    mkdir thetool-work
                    cd thetool-work
                    cp ${GOPATH}/bin/thetool .
                    ./thetool init -g $GLOO_HASH --no-defaults
                    ./thetool add -r https://github.com/solo-io/gloo-plugins.git -c $GLOO_PLUGINS_HASH
                    ./thetool build envoy -d -v --publish=false
                    ./thetool build gloo -d -v --publish=false
                '''
            }
        }

        stage('Build Envoy') {
            container('envoy-build') {
                echo 'Building envoy...'
                sh '''
                    cd thetool-work/repositories
                    ln -s `pwd` /repositories
                    cd ../envoy
                    ln -s `pwd` /source
                    cd /source
                    ./build-envoy.sh
                '''
            }
        }

        stage('Build gloo') {
            container('golang') {
                echo 'Building gloo...'
                sh '''
                    mkdir -p $GOPATH/pkg/dep
                    cd thetool-work
                    ln -s `pwd` /gloo
                    cd /gloo
                    ./build-gloo.sh
                '''
            }
        }

        stage ('Build k8s service discovery') {
            container('golang') {
                echo 'Building k8s service discovery...'
                sh '''
                    mkdir -p $GOPATH/pkg/dep
                    cd thetool-work
                    ./thetool build gloo-k8s-service-discovery -d
                    if [ -e "/code" ]; then
                      echo "already linked?"
                    else
                      ln -s `pwd` /code
                    fi
                    cd /code
                    ./build-gloo-k8s-service-discovery.sh
                    cd repositories/gloo-k8s-service-discovery
                    git log -n 1 --pretty=format:%h > hash.tmp
                '''
            }
        }
        
        stage ('Build function discovery') {
            container('golang') {
                echo 'Building function discovery...'
                sh '''
                    mkdir -p $GOPATH/pkg/dep
                    cd thetool-work
                    ./thetool build gloo-function-discovery -d
                    if [ -e "/code" ]; then
                      echo "already linked?"
                    else
                      ln -s `pwd` /code
                    fi
                    cd /code
                    ./build-gloo-function-discovery.sh
                    cd repositories/gloo-function-discovery
                    git log -n 1 --pretty=format:%h > hash.tmp
                '''
            }
        }
        
        stage ('Build ingress controller') {
            container('golang') {
                echo 'Building ingress controller...'
                sh '''
                    mkdir -p $GOPATH/pkg/dep
                    cd thetool-work
                    ./thetool build gloo-ingress-controller -d
                    if [ -e "/code" ]; then
                      echo "already linked?"
                    else
                      ln -s `pwd` /code
                    fi
                    cd /code
                    ./build-gloo-ingress-controller.sh
                    cd repositories/gloo-ingress-controller
                    git log -n 1 --pretty=format:%h > hash.tmp
                '''
            }
        }
        
        
        stage('Test Envoy and Gloo') {
            echo 'Testing Envoy and Gloo'
            container('docker') {
                sh '''
                    if [ ! -f `pwd`/thetool-work/envoy/envoy-out/envoy ]; then
                        CID=$(docker run -d  soloio/envoy:v0.1.2 /bin/bash -c exit)
                        docker cp $CID:/usr/local/bin/envoy `pwd`/thetool-work/envoy/envoy-out/envoy 
                        docker rm $CID
                    fi
                '''
            }
            
            container('golang') {
                sh '''
                    export GLOO_BINARY=`pwd`/thetool-work/gloo-out/gloo
                    export ENVOY_BINARY=`pwd`/thetool-work/envoy/envoy-out/envoy

                    go get -u github.com/golang/dep/cmd/dep
                    cd $GOPATH/src
                    mkdir -p github.com/solo-io
                    cd github.com/solo-io
                    git clone https://github.com/solo-io/gloo-testing
                    cd gloo-testing
                    dep ensure -vendor-only
                    cd local_e2e
                    go test
                '''
            }
        }
        
        stage('Publish Docker Images') {
            echo 'Publishing Docker images'
            container('docker') {
                sh '''
                    cd thetool-work
                    export BASE=`pwd`
                    if [ -n "$GLOO_VERSION" ]; then
                      export TAG=$GLOO_VERSION-$BUILD_NUMBER
                    else
                      export TAG=v0.1.6-$BUILD_NUMBER
                    fi
                    cd envoy/envoy-out
                    docker build -t soloio/envoy:$TAG .
                    docker tag soloio/envoy:$TAG soloio/envoy:latest
                    docker push soloio/envoy:$TAG
                    docker push soloio/envoy:latest
                    
                    cd $BASE
                    cd gloo-out
                    cp ../repositories/gloo/Dockerfile .
                    docker build -t soloio/gloo:$TAG .
                    docker tag soloio/gloo:$TAG soloio/gloo:latest
                    docker push soloio/gloo:$TAG
                    docker push soloio/gloo:latest
                    
                    cd $BASE
                    cd repositories/gloo-function-discovery
                    export FTAG=`cat hash.tmp`
                    cd ../../gloo-function-discovery-out
                    cp ../repositories/gloo-function-discovery/Dockerfile .
                    docker build -t soloio/gloo-function-discovery:$TAG .
                    docker tag soloio/gloo-function-discovery:$TAG soloio/gloo-function-discovery:latest
                    docker tag soloio/gloo-function-discovery:$TAG soloio/gloo-function-discovery:$FTAG
                    docker push soloio/gloo-function-discovery:$TAG
                    docker push soloio/gloo-function-discovery:latest
                    docker push soloio/gloo-function-discovery:$FTAG
                    
                    cd $BASE
                    cd repositories/gloo-ingress-controller
                    export FTAG=`cat hash.tmp`
                    cd ../../gloo-ingress-controller-out
                    cp ../repositories/gloo-ingress-controller/Dockerfile .
                    docker build -t soloio/gloo-ingress-controller:$TAG .
                    docker tag soloio/gloo-ingress-controller:$TAG soloio/gloo-ingress-controller:latest
                    docker tag soloio/gloo-ingress-controller:$TAG soloio/gloo-ingress-controller:$FTAG
                    docker push soloio/gloo-ingress-controller:$TAG
                    docker push soloio/gloo-ingress-controller:latest
                    docker push soloio/gloo-ingress-controller:$FTAG
                    
                    cd $BASE
                    cd repositories/gloo-k8s-service-discovery
                    export FTAG=`cat hash.tmp`
                    cd ../../gloo-k8s-service-discovery-out
                    cp ../repositories/gloo-k8s-service-discovery/Dockerfile .
                    docker build -t soloio/gloo-k8s-service-discovery:$TAG .
                    docker tag soloio/gloo-k8s-service-discovery:$TAG soloio/gloo-k8s-service-discovery:latest
                    docker tag soloio/gloo-k8s-service-discovery:$TAG soloio/gloo-k8s-service-discovery:$FTAG
                    docker push soloio/gloo-k8s-service-discovery:$TAG
                    docker push soloio/gloo-k8s-service-discovery:latest
                    docker push soloio/gloo-k8s-service-discovery:$FTAG
                '''
            }
        }
        
        stage('Generate install and publish') {
            echo 'Generate files for install and update gloo-install'
            container('helm') {
                sh '''
                    cd thetool-work
                    if [ -n "$GLOO_VERSION" ]; then
                      export TAG=$GLOO_VERSION-$BUILD_NUMBER
                    else
                      export TAG=v0.1.6-$BUILD_NUMBER
                    fi
                    ./thetool deploy k8s --generate-install -t $TAG 
                    cat install.yaml
                    ./thetool deploy k8s -d -t $TAG
                    cat gloo-chart.yaml
                '''
            }
            
            container('golang') {
                sh '''
                    if [ -n "$GLOO_VERSION" ]; then
                      export TAG=$GLOO_VERSION-$BUILD_NUMBER
                    else
                      export TAG=v0.1.6-$BUILD_NUMBER
                    fi
                    cd thetool-work
                    cp /etc/github/id_rsa $PWD 
                    chmod 400 $PWD/id_rsa 
                    export GIT_SSH_COMMAND="ssh -i $PWD/id_rsa -o 'StrictHostKeyChecking no'" 
                    git clone git@github.com:solo-io/gloo-install.git
                    cd gloo-install
                    git checkout -b $TAG
                    cp ../install.yaml kube/install.yaml
                    cp ../gloo-chart.yaml helm/gloo/values.yaml
                    git diff | cat
                    git add kube/install.yaml
                    git add helm/gloo/values.yaml
                    git config --global user.email "bot@soloio.com"
                    git config --global user.name "Solo Buildbot"
                    git commit -m "Jenkins: updated for $TAG"
                    if [ "$UPDATE_INSTALL" = "true" ]; then
                      git push -u origin $TAG:$TAG
                    fi
                    rm ../id_rsa
                ''' 
            }
        }
    }
}