node {
  def project = 'weather-179004'
  def appName = 'weather'
  def feSvcName = "${appName}-app"
  def uiImageTag = "gcr.io/${project}/${appName}-ui:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"
  def apiImageTag = "gcr.io/${project}/${appName}-api:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"
  def proxyImageTag = "gcr.io/${project}/${appName}-proxyi:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"

  checkout scm

  stage 'Build ui and api containers '
    sh("cd ./containers/weather-ui; docker build -t ${uiImageTag} .")
    sh("cd ./containers/weather-api; docker build -t ${apiImageTag} .")

  stage 'Run Test'
    sh("echo Test placeholder")

  stage 'Push Containers to Registry'
    sh("cd ./containers/weather-ui; gcloud docker -- push ${uiImageTag}")
    sh("cd ./containers/weather-api; gcloud docker -- push ${apiImageTag}")

  stage 'Deploy '

    // Roll out a dev environment
        // Create namespace if it doesn't exist
        sh("kubectl get ns ${env.BRANCH_NAME} || kubectl create ns ${env.BRANCH_NAME}")
        // Don't use public load balancing for development branches
        sh("sed -i.bak 's#LoadBalancer#ClusterIP#' ./k8s/services/frontend.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-ui:1.0.0#${uiImageTag}#' ./k8s/dev/frontend-dev.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-api:1.0.0#${apiImageTag}#' ./k8s/dev/backend-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/frontend.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/backend.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/frontend-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/backend-dev.yaml")
        echo 'To access your environment run `kubectl proxy`'
        echo "Then access your service via http://localhost:8001/api/v1/proxy/namespaces/${env.BRANCH_NAME}/services/${feSvcName}:80/"
}
