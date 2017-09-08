node {
  def project = 'weather-179004'
  def appName = 'weather'
  def feSvcName = "${appName}-app"
  def uiImageTag = "gcr.io/${project}/${appName}-ui:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"
  def apiImageTag = "gcr.io/${project}/${appName}-api:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"
  def proxyImageTag = "gcr.io/${project}/${appName}-proxy:${env.BRANCH_NAME}.${env.BUILD_NUMBER}"

  checkout scm

  stage 'Building docker containers '
    sh("sed -i.bak 's#Microservice Docker Container#Docker Container Version ${env.BRANCH_NAME}.${env.BUILD_NUMBER} Deployed with Kubernetes#' ./containers/weather-ui/index.html")

    sh("cd ./containers/weather-ui; docker build -t ${uiImageTag} .")
    sh("cd ./containers/weather-api; docker build -t ${apiImageTag} .")
    sh("cd ./containers/weather-proxy; docker build -t ${proxyImageTag} .")

  stage 'Run Test placeholder'
    sh("echo Test placeholder")

  stage 'Push Containers to GCP Registry'
    sh("cd ./containers/weather-ui; gcloud docker -- push ${uiImageTag}")
    sh("cd ./containers/weather-api; gcloud docker -- push ${apiImageTag}")
    sh("cd ./containers/weather-proxy; gcloud docker -- push ${proxyImageTag}")

  stage 'Deploy via Kube apply'

    // Roll out a dev environment
        // Create namespace if it doesn't exist
        sh("kubectl get ns ${env.BRANCH_NAME} || kubectl create ns ${env.BRANCH_NAME}")
        // Don't use public load balancing for development branches
        // ** For now in Dev use Loadbalancer sh("sed -i.bak 's#LoadBalancer#ClusterIP#' ./k8s/services/weather-proxy-srv.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-ui:1.0.0#${uiImageTag}#' ./k8s/dev/weather-ui-dev.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-api:1.0.0#${apiImageTag}#' ./k8s/dev/weather-api-dev.yaml")
        sh("sed -i.bak 's#gcr.io/cloud-solutions-images/weather-proxy:1.0.0#${proxyImageTag}#' ./k8s/dev/weather-proxy-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/weather-ui-srv.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/weather-api-srv.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/weather-proxy-srv.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/weather-proxy-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/weather-ui-dev.yaml")
        sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/dev/weather-api-dev.yaml")
        echo 'To access your environment run `kubectl proxy`'
        echo "Then access your service via http://localhost:8001/api/v1/proxy/namespaces/${env.BRANCH_NAME}/services/${feSvcName}:80/"
}
