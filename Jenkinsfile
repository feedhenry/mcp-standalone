node ("go") {
  sh "mkdir -p src/github.com/feedhenry/mcp-standalone"
  withEnv(["GOPATH=${env.WORKSPACE}/","PATH=${env.PATH}:${env.WORKSPACE}/bin"]) {
    dir ("src/github.com/feedhenry/mcp-standalone") {
      stage("Checkout") {
        checkout scm
      }
      
      stage ("Setup") {
        sh "make setup"
      }
      
      stage("Check") {
        sh "make check"
      }
      
      stage ("Build") {
        sh "make build"
      }
      
      stage ("Build cli") {
        sh "make build_cli"  
      }
    }
  }
}

