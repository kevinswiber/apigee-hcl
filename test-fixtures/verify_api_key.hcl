proxy "VerifyAPIKeyfixture" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "verify-apikey" {}
    }
  }

  http_proxy_connection {
    base_path    = "/v0/hello"
    virtual_host = ["default", "secure"]
  }

  route_rule "preflight" {
    condition = "request.verb == \"OPTIONS\""
  }

  route_rule "default" {
    target_endpoint = "default"
  }
}

target_endpoint "default" {
  http_target_connection {
    url = "http://mocktarget.apigee.net"
  }
}

policy verify_api_key "verify-apikey" {
  async             = false
  continue_on_error = false
  enabled           = true
  display_name      = "Verify APIKey"

  apikey {
    ref   = "request.header.apikey"
  }
}
