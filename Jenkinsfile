node {
  def project = 'weather-179004'
  def appName = 'weather'
  def feSvcName = "${appName}-app"
  def appImageTag = "gcr.io/${project}/${appName}-app:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"
  def apiImageTag = "gcr.io/${project}/${appName}-api:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"

  try {
      stage('Checkout'){

        checkout scm
      }
  }

  stage 'Build front and backend image '{
    sh("cd ./containers/weather-app; docker build -t ${appImageTag} .")
    sh("cd ./containers/weather-api; docker build -t ${apiImageTag} .")
}

  stage 'Run Test'{
    sh("echo Test placeholder")
  }

  stage 'Push to Registry'{
    sh("cd ./containers/weather-app; gcloud docker -- push ${appImageTag}")
    sh("cd ./containers/weather-api; gcloud docker -- push ${appImageTag}")
}

  stage "Deploy Application"
  stage 'Run Go tests'
  switch (env.BRANCH_NAME) {
    // Roll out to canary environment
    case "canary":
        // Change deployed image in canary to the one we just built
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-app:1.0.0#${appImageTag}#' ./k8s/canary/frontend-canary.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-api:1.0.0#${apiImageTag}#' ./k8s/canary/backend-canary.yaml")
        sh("kubectl --namespace=production apply -f k8s/canary/frontend-canary.yaml")
        sh("kubectl --namespace=production apply -f k8s/canary/backend-canary.yaml")
        sh("kubectl --namespace=production apply -f k8s/services/frontend.yaml")
        sh("kubectl --namespace=production apply -f k8s/services/backend.yaml")
        // sh("echo http://`kubectl --namespace=production get service/${feSvcName} --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > ${feSvcName}")
        break
    // Roll out to production
    case "master":
        // Change deployed image in canary to the one we just built
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-app:1.0.0#${appImageTag}#' ./k8s/production/frontend-production.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-api:1.0.0#${apiImageTag}#' ./k8s/production/backend-production.yaml")
        sh("kubectl --namespace=production apply -f k8s/services/frontend.yaml")
        sh("kubectl --namespace=production apply -f k8s/services/backend.yaml")
        sh("kubectl --namespace=production apply -f k8s/production/frontend-production.yaml")
        sh("kubectl --namespace=production apply -f k8s/production/backend-production.yaml")
        // sh("echo http://`kubectl --namespace=production get service/${feSvcName} --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > ${feSvcName}")
        break
    // Roll out a dev environment
    default:
        // Create namespace if it doesn't exist
        sh("kubectl get ns ${env.BRANCH_NAME} || kubectl create ns ${env.BRANCH_NAME}")
        // Don't use public load balancing for development branches
        sh("sed -i.bak 's#LoadBalancer#ClusterIP#' ./k8s/services/frontend.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-app:1.0.0#${appImageTag}#' ./k8s/dev/frontend-dev.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-api:1.0.0#${apiImageTag}#' ./k8s/dev/backend-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/frontend.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/backend.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/frontend-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/backend-dev.yaml")
        echo 'To access your environment run `kubectl proxy`'
        echo "Then access your service via http://localhost:8001/api/v1/proxy/namespaces/${env.BRANCH_NAME}/services/${feSvcName}:80/"
  }

  catch (err) {

        currentBuild.result = "FAILURE"
        echo "The Build failed.."
            // mail body: "project build error is here: ${env.BUILD_URL}" ,
            // from: 'xxxx@yyyy.com',
            // replyTo: 'yyyy@yyyy.com',
            // subject: 'project build failed',
            // to: 'zzzz@yyyyy.com'

        throw err
    }

}




