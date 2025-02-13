def func_telegram_sendMessage(message, token, chatid) {
    try {
        sh """
            curl -s -X POST https://api.telegram.org/bot${token}/sendMessage \
            -d chat_id=${chatid} -d parse_mode=markdown \
            -d text='${message}'
        """
    } catch(Exception e) {
        currentBuild.result = 'SUCCESS'
    }
}

pipeline {
    agent any
    environment {
        SERVICE = "${GIT_URL.tokenize('/.')[-2]}"
        BRANCH_NAME = "${GIT_BRANCH.replaceFirst('origin/', '')}"
        WORKDIR = '/opt/cicd/'
        TOKEN = credentials('telegram_bot_token')
        CHAT = credentials('telegram_chat_id')
        MESSAGE_BASE = "${env.BRANCH_NAME == "develop" ? "\\[ DEV4 ] *${env.SERVICE}*: " : "\\[ STAGE1 ] *${env.SERVICE}*: "}"
        REGISTRY = credentials('docker_registry_host')
        DEV = credentials('dev_server_ssh_cmd')
        STAGE = credentials('stage1_ssh_cmd')
        TARGET_SRV = "${env.BRANCH_NAME == "develop" ? "${env.DEV}" : "${env.STAGE}"}"
        TARGET_PATH = "${env.BRANCH_NAME == "develop" ? "/opt/touchon/gobin" : "/opt/touchon"}"
        IMG_TAG = "${env.BRANCH_NAME == "develop" ? "develop" : "${env.BRANCH_NAME.replaceFirst('release/', '')}"}"
    }
    stages {
        stage('Notification') {
            steps {
                echo "${env.NV_NAME}"
                script {
                    initMessage = "${env.MESSAGE_BASE}STARTED"
                }
                func_telegram_sendMessage("$initMessage", "${env.TOKEN}", "${env.CHAT}")
            }
        }
        stage('Pull') {
            steps {
                sh """
                  git -C ${env.WORKDIR}${env.SERVICE} checkout ${env.BRANCH_NAME}
                  git -C ${env.WORKDIR}${env.SERVICE} pull
                """
            }
        }
        stage('Build') {
            steps {
                sh """
                    docker buildx build \
                    -t ${env.REGISTRY}/${env.SERVICE}:${env.IMG_TAG} \
                    --platform linux/arm64 \
                    --push \
                    ${env.WORKDIR}${env.SERVICE}
                """
            }
        }
        stage('Publish') {
            steps {
                sh """
                    ssh ${env.TARGET_SRV} << EOF
                    set -e
                    cd ${env.TARGET_PATH}
                    docker compose pull ${env.SERVICE}
                    docker compose up --force-recreate --build -d ${env.SERVICE}
                    docker system prune -af
                    << EOF
                """
            }
        }
    }
    
    post {
        success {
            script {
                gitCommit = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} log -n 1 --pretty=format:'%h'", returnStdout: true)
                gitCommiter = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} show -s --pretty=%an", returnStdout: true)
                gitCommitComment = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} show --pretty=format:'%B' --no-patch -n 1 $gitCommit", returnStdout: true)
                gitCommitComment = gitCommitComment.replace("_", "\\_")
                gitCommitComment = gitCommitComment.replace("*", "\\*")
                gitCommitComment = gitCommitComment.replace("[", "\\[")
                gitCommitComment = gitCommitComment.replace("`", "\\`")
                successMessage = "${env.MESSAGE_BASE}SUCSESS%0ACommit $gitCommit by $gitCommiter$gitCommitComment"
                func_telegram_sendMessage("$successMessage", "${env.TOKEN}", "${env.CHAT}")
            }
        }
        aborted {
            script {
                abortMessage = "${env.MESSAGE_BASE}ABORTED"
                func_telegram_sendMessage("$abortMessage", "${env.TOKEN}", "${env.CHAT}")
            }
        }
        failure {
            script {
                failMessage = "${env.MESSAGE_BASE}FAILURE"
                func_telegram_sendMessage("$failMessage", "${env.TOKEN}", "${env.CHAT}")
            }
        }
    }
}