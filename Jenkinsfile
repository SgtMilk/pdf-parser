pipeline {
    agent {
        docker { image 'node:16.13.1-alpine' }
    }

    stages {
        stage('Test') {
            steps {
                sh 'pwd'
                sh 'node --version'
            }
        }
    }

    // stages {
    //     stage('Linting') {
    //         steps {
    //             echo 'Linting..'
    //             sh 'pwd'
    //             sh 'echo "$ENV"'
    //             sh 'golangci-lint run'
    //         }
    //     }
    //     stage('Build') {
    //         steps {
    //             echo 'Building..'
    //             sh 'go build'
    //         }
    //     }
    //     stage('Test') {
    //         steps {
    //             echo 'Testing..'
    //             sh 'go test ./...'
    //         }
    //     }
    // }
}