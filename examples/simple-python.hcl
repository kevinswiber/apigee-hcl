proxy "python" {}

proxy_endpoint "default" {
  pre_flow {
    response {
      step "pythonscriptpolicy" {}
    }
  }

  http_proxy_connection {
    base_path    = "/python"
    virtual_host = ["default"]
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

policy script "pythonscriptpolicy" {
  resource_url = "py://setHeader.py"

  content = <<EOF
response.setVariable("header.X-Apigee-Demo-target", flow.getVariable("target.url"));
print 'Reached the script & assigned header variable' 
EOF
}
