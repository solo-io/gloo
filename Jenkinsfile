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
],
volumes: [
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
    secretVolume(secretName: 'soloio-docker-hub', mountPath: '/etc/docker'),
    secretVolume(secretName: 'soloio-github', mountPath: '/etc/github')
]) {

    properties([
        parameters ([
            stringParam(
                defaultValue: '',
                description: 'Commit hash for gloo',
                name: 'GLOO_HASH'),
            stringParam(
                defaultValue: '',
                description: 'Commit hash for gloo-plugins',
                name: 'GLOO_PLUGINS_HASH')
        ])
    ])

    node('gloo-builder') {
        stage('Build thetool') {
            container('golang') {
                echo 'Building thetool...' 
                sh '''
                    curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/release/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep
                    git clone https://github.com/solo-io/thetool.git
                    mkdir -p ${GOPATH}/src/github.com/solo-io/
                    ln -s `pwd`/thetool ${GOPATH}/src/github.com/solo-io/thetool
                    cd ${GOPATH}/src/github.com/solo-io/thetool
                    dep ensure -vendor only
                    CGO_ENABLED=0 go build
                    }
                '''
            }
        }

        stage('Initialize thetool') {
            echo 'Initializing thetool...'
            sh '''
                mkdir thetool-work
                cd thetool-work
                ../thetool init -g param.GLOO_HASH --no-defaults
                ../thetool add -r https://github.com/solo-io/gloo-plugins.git -c param.GLOO_PLUGINS_HASH
                ../thetool build all -d
            '''
        }

        stage('Build Envoy') {
            container('envoy-build') {
                echo 'Building envoy...'
                sh '''
                    pwd
                    ls
                '''
            }
        }

        stage('Build gloo') {
            container('golang') {
                echo 'Building gloo...'
                sh '''
                    pwd
                    ls
                '''
            }
        }
    }
}