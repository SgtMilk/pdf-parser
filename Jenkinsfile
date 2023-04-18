pipeline {
    agent { dockerfile true }

    stages {
        stage('Linting') {
            steps {
                echo 'Linting..'
                sh 'golangci-lint run'
            }
        }
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sh 'go test ./...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}