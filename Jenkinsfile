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
        SERVICE = 'touchon-server'
        WORKDIR = '/opt/cicd_v2/'
        TOKEN = credentials('telegram_bot_token')
        CHAT = credentials('telegram_chat_id')
        MESSAGE_BASE = "\\[ DEV4 ] *${env.SERVICE}*: "
        REGISTRY = credentials('docker_registry_host')
        DEV_SRV = credentials('dev_server_ssh_cmd')
    }
    stages {
        stage('Notification') {
            steps {
                // script {
                //     initMessage = "${env.MESSAGE_BASE}STARTED"
                // }
                // func_telegram_sendMessage("$initMessage", "${env.TOKEN}", "${env.CHAT}")
                echo 'Pulling...' + env.GIT_BRANCH
                echo GIT_URL.tokenize('/.')[-2]
                println scm.branches
                // sh 'printenv'
            }
        }
        // stage('Pull') {
        //     steps {
        //         sh """
        //           git -C ${env.WORKDIR}${env.SERVICE} checkout develop
        //           git -C ${env.WORKDIR}${env.SERVICE} pull
        //         """
        //     }
        // }
        // stage('Build') {
        //     steps {
        //         sh """
        //             docker buildx build \
        //             -t ${env.REGISTRY}/${env.SERVICE}:develop \
        //             --platform linux/arm64 \
        //             --push \
        //             ${env.WORKDIR}${env.SERVICE}
        //         """
        //     }
        // }
        // stage('Publish') {
        //     steps {
        //         sh """
        //             ssh ${env.DEV_SRV} << EOF
        //             set -e
        //             cd /opt/touchon/gobin
        //             docker-compose pull ${env.SERVICE}
        //             docker create --name temp_${env.SERVICE}_container ${env.REGISTRY}/${env.SERVICE}:develop
        //             docker cp temp_${env.SERVICE}_container:/opt/service/db.sqlite ./${env.SERVICE}/db.sqlite
        //             docker rm -f temp_${env.SERVICE}_container
        //             docker-compose up --force-recreate --build -d ${env.SERVICE}
        //             docker system prune -af
        //             << EOF
        //         """
        //     }
        // }
    }
    
    // post {
    //     success {
    //         script {
    //             gitCommit = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} log -n 1 --pretty=format:'%h'", returnStdout: true)
    //             gitCommiter = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} show -s --pretty=%an", returnStdout: true)
    //             gitCommitComment = sh (script: "git -C ${env.WORKDIR}${env.SERVICE} show --pretty=format:'%B' --no-patch -n 1 $gitCommit", returnStdout: true)
    //             successMessage = "${env.MESSAGE_BASE}SUCSESS%0ACommit $gitCommit by $gitCommiter$gitCommitComment"
    //             // func_telegram_sendMessage("$successMessage", "${env.TOKEN}", "${env.CHAT}")
    //         }
    //     }
    //     aborted {
    //         script {
    //             abortMessage = "${env.MESSAGE_BASE}ABORTED"
    //             // func_telegram_sendMessage("$abortMessage", "${env.TOKEN}", "${env.CHAT}")
    //         }
    //     }
    //     failure {
    //         script {
    //             failMessage = "${env.MESSAGE_BASE}FAILURE"
    //             // func_telegram_sendMessage("$failMessage", "${env.TOKEN}", "${env.CHAT}")
    //         }
    //     }
    // }
}
