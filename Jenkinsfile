pipeline {
    agent any

    tools { go '1.20' }

    stages {
        stage('Build') {
            steps {
                // Build the application
                sh "go build ./cmd/ui"
            }
        }
        stage('Test') {
            steps {
                // Test the application
                sh "go test ./..."
            }
        }
        stage('Push') {
            steps {
                // Push to S3
                sh "aws s3 cp ./ui s3://cloudschoolproject-buildartifacts/ui"
            }
        }
    }
}