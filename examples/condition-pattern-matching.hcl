proxy "matchtest" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "SomePolicy" {
        condition = "(proxy.pathsuffix Matches \"/cat\")"
      }
    }
  }

  post_flow {
    response {
      step "SetResponse" {}
    }
  }

  http_proxy_connection {
    base_path    = "/matchtest"
    virtual_host = ["default", "secure"]
  }

  route_rule "default" {}
}

policy assign_message "SomePolicy" {
  display_name = "SomePolicy"

  assign_variable {
    name  = "condition.status"
    value = "Condition Succeeded"
  }

  ignore_unresolved_variables = true

  assign_to {
    create_new = false
    type       = "request"
  }
}

policy javascript "SetResponse" {
  time_limit   = 200
  display_name = "SetResponse"

  resource_url = "jsc://setresponse.js"

  content = <<EOF
var conditionStatus = context.getVariable('condition.status');
if (conditionStatus === null | conditionStatus === '') {
    context.setVariable("response.content", "Condition Failed for proxy.pathsuffix: " + context.getVariable("proxy.pathsuffix"));
} else {
    context.setVariable("response.content", conditionStatus + " for proxy.pathsuffix: " + context.getVariable("proxy.pathsuffix"));
}
EOF
}
