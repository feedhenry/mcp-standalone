node ("go") {
  sh "mkdir -p src/github.com/feedhenry/"
  dir ("src/github.com/feedhenry/mcp-standalone") {
    checkout scm
    sh "make setup"
    sh "make check"
    sh "make build"
    sh "make build_cli"  
  }
}
