proxy "conditional-policy" {}

proxy_endpoint "default" {
  pre_flow {
    response {
      step "Timer" {
        condition = "request.header.responsetime=\"true\""
      }
    }
  }

  http_proxy_connection {
    base_path    = "/timer"
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

policy script "Timer" {
  resource_url = "py://timer"

  content = <<EOF
response_time = flow.getVariable("target.received.start.timestamp") - flow.getVariable("target.sent.start.timestamp");
response.setVariable("header.X-Apigee-target", flow.getVariable("target.url"));
response.setVariable("header.X-Apigee-start-time", flow.getVariable("target.sent.start.time"));
response.setVariable("header.X-Apigee-start-timestamp", flow.getVariable("target.sent.start.timestamp"));
response.setVariable("header.X-Apigee-end-time", flow.getVariable("target.sent.end.time"));
response.setVariable("header.X-Apigee-end-timestamp", flow.getVariable("target.received.start.timestamp"));
response.setVariable("header.X-Apigee-target-responseTime", response_time);
EOF
}
