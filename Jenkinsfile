pipeline {
    agent any

    tools { go '1.20' }

    stages {
        stage('Build') {
            steps {
                // Build the application
                sh "mkdir -p output"
                sh "go build -o output ./cmd/wikilinkui"
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
                sh "aws s3 cp ./output/wikilinkui s3://cloudschoolproject-buildartifacts/wikilinkui_v_$BUILD_NUMBER"
            }
        }
    }
}