pipeline {
    agent { dockerfile { args '-u root:root' } }

    stages {
        stage('Linting') {
            steps {
                echo 'Linting..'
                sh 'golangci-lint run --timeout=5m'
            }
        }
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'go build -buildvcs=false'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sh 'go test ./...'
            }
        }
    }
}