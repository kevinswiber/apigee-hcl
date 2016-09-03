proxy "variables" {}

proxy_endpoint "default" {
  flow "default" {
    response {
      step "parse-xml-response"{}
      step "read-variables"{}
      step "xml-to-json"{}
      step "parse-json-response"{}
      step "read-variables"{}
      step "variables"{}
    }
  }

  http_proxy_connection {
    base_path    = "/variables"
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

policy extract_variables "parse-xml-response" {
  ignore_unresolved_variables = true
  variable_prefix             = "mock"

  xml_payload {
    stop_payload_processing = false

    variable "city" {
      type  = "string"
      xpath = "//city"
    }

    variable "state" {
      type  = "string"
      xpath = "//state"
    }
  }
}

policy javascript "read-variables" {
  time_limit   = 200
  resource_url = "jsc://read-variables.js"

  content = <<EOF
var a = context.getVariable("mock.city");
print(a);

var b = context.getVariable("mock.state");
print(b);

var c = context.getVariable("mock.firstName");
print(c);

var d = context.getVariable("mock.lastName");
print(d);
EOF
}

policy xml_to_json "xml-to-json" {
  options {}
}

policy extract_variables "parse-json-response" {
  ignore_unresolved_variables = true
  variable_prefix             = "mock"

  json_payload {
    variable "firstName" {
      json_path = "$.root.firstName"
    }

    variable "lastName" {
      json_path = "$.root.lastName"
    }
  }
}

policy assign_message "variables" {
  assign_to {
    create_new = false
    type       = "response"
  }

  ignore_unresolved_variables = true

  set {
    # Edge flow variables
    header "system.timestamp" {
      value = "{system.timestamp}"
    }

    header "system.time" {
      value = "{system.time}"
    }

    header "organization.name" {
      value = "{organization.name}"
    }

    header "apiproxy.name" {
      value = "{apiproxy.name}"
    }

    header "apiproxy.revision" {
      value = "{apiproxy.revision}"
    }

    header "proxy.basepath" {
      value = "{proxy.basepath}"
    }

    header "proxy.name" {
      value = "{apiproxy.name}"
    }

    header "proxy.pathsuffix" {
      value = "{proxy.pathsuffix}"
    }

    header "message.headers.count" {
      value = "{message.headers.count}"
    }

    header "message.headers.names" {
      value = "{message.headers.names}"
    }

    header "client.ip" {
      value = "{client.ip}"
    }

    header "request.uri" {
      value = "{request.uri}"
    }

    header "request.headers.names" {
      value = "{request.headers.names}"
    }

    header "request.header.user-agent.values" {
      value = "{request.header.user-agent.values}"
    }

    header "request.path" {
      value = "{request.path}"
    }

    header "request.querystring" {
      value = "{request.querystring}"
    }

    header "request.queryparams.names" {
      value = "{request.queryparams.names}"
    }

    header "request.queryparam.w" {
      value = "{request.queryparam.w}"
    }

    header "request.verb" {
      value = "{request.verb}"
    }

    header "target.url" {
      value = "{target.url}"
    }

    header "target.host" {
      value = "{target.host}"
    }

    header "target.ip" {
      value = "{target.ip}"
    }

    # Variables populated by parsing XML response with an ExtractVariables policy
    header "mock.city" {
      value = "{mock.city}"
    }

    header "mock.state" {
      value = "{mock.state}"
    }

    # Variables populated by parsing JSON response with an ExtractVariables policy
    header "mock.firstName" {
      value = "{mock.firstName}"
    }

    header "mock.lastName" {
      value = "{mock.lastName}"
    }
  }
}
