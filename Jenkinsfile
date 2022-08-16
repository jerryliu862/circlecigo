pipeline {
    agent any
    
    environment {
        CI = 'true'
        RED = '#FF0000'
        YELLOW = '#FFFF00'
        GREEN = '#008000'
    }

    stages {
        stage('Pre-Build') {
            steps {
                echo "Build Started"
                slackSend(color: '#008000', message: "Build Started: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})", channel: "#17live_wso_cicd")
            }
        }
        
        stage('Copy Folder') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'development') {    
                        echo 'dev copy folder'   
                        sh 'cp -R ./* /home/foritech/Projects/17Live_WSO/17Live_WSO_BE/'
                    } 
                }
            }
        }
        
        stage('Remove Docker Containers') {
            steps {
                script {
                    if (env.BRANCH_NAME == "development") {
                        sh 'sudo docker rm -f 17live-wso-be-container-dev || true'
                    }
                    
                    if (env.BRANCH_NAME == "stage") {
                        sh 'sudo docker rm -f 17live-wso-be-container-stg || true'
                    }                    
                }
            }
        }

        stage('Remove Docker Image') {
            steps {
                script {
                    if (env.BRANCH_NAME == "development") { 
                        sh 'sudo docker image rm 17live-wso-be-image-dev || true'
                    }
                    if (env.BRANCH_NAME == "stage") { 
                        sh 'sudo docker image rm 17live-wso-be-image-stg || true'
                    }
                }
            }
        }

        stage('List Docker Containers') {
            steps {
                script {
                    if (env.BRANCH_NAME == "development") { 
                        sh 'sudo docker ps --all'
                    }
                    if (env.BRANCH_NAME == "stage") { 
                        sh 'sudo docker ps --all'
                    }
                }    
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    if (env.BRANCH_NAME == "development") { 
                        sh 'sudo docker build . -t "17live-wso-be-image-dev"'
                    }
                    if (env.BRANCH_NAME == "stage") { 
                        sh 'sudo docker build . -t "17live-wso-be-image-stg"'
                    }
                }    
            }
        }

        stage('Build Docker Container') {
            steps {
                script {
                    if (env.BRANCH_NAME == "development") { 
                        sh 'sudo docker run -d --name 17live-wso-be-container-dev -p 4102:8080 17live-wso-be-image-dev'
                    }
                    if (env.BRANCH_NAME == "stage") { 
                        sh 'sudo docker run -d --name 17live-wso-be-container-stg -p 4112:8080 17live-wso-be-image-stg'
                    }
                }    
            }
        }
    }

    post { 
        success {
            echo "Build Success"
            slackSend(color: '#008000', message: "Build Succeed: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})", channel: "#17live_wso_cicd")
        }

        failure { 
            echo "Build Failed"
            slackSend(color: '#FF0000', message: "Build Failed: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})", channel: "#17live_wso_cicd")
        }
    }
}
