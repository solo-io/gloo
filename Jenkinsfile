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
                defaultValue: 'cf08737718cf62bf597f88aa2068c6f6b28b9992',
                description: 'Commit hash for gloo',
                name: 'GLOO_HASH'),
            stringParam(
                defaultValue: '282a844ea3ed2527f5044408c9c98bc7ee027cd2',
                description: 'Commit hash for gloo-plugins',
                name: 'GLOO_PLUGINS_HASH')
        ])
    ])

    node('gloo-builder') {
        stage('Build thetool') {
            container('golang') {
                echo 'Building thetool...' 
                sh '''
                    go get -u github.com/solo-io/thetool
                    mkdir thetool-work
                    cp ${GOPATH}/bin/thetool thetool-work/.
                '''
            }
        }

        stage('Initialize thetool') {
            echo 'Initializing thetool...'
            sh '''
                cd thetool-work
                ls -l
                ./thetool init -g $GLOO_HASH --no-defaults
                ./thetool add -r https://github.com/solo-io/gloo-plugins.git -c $GLOO_PLUGINS_HASH
            '''
        }

        stage('Build Envoy') {
            container('envoy-build') {
                echo 'Building envoy...'
                sh '''
                    pwd
                    ls
                    cd thetool-work
                    ls
                    ./thetool build envoy -d
                '''
            }
        }

        stage('Build gloo') {
            container('golang') {
                echo 'Building gloo...'
                sh '''
                    pwd
                    ls
                    cd thetool-work
                    ls
                    ./thetool build gloo -d
                '''
            }
        }
    }
}