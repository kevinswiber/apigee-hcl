proxy "helloworld" {
  display_name = "Hello World"
}

proxy_endpoint "default" {
  pre_flow {
    request {
      # Handle preflight OPTIONS calls for cross origin requests
      step "add-cors" {
        condition = "request.verb == \"OPTIONS\""
      }
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

policy assign_message "add-cors" {
  continue_on_error           = false
  enabled                     = true
  ignore_unresolved_variables = true
  display_name                = "Add CORS"

  add {
    header "Access-Control-Allow-Origin" {
      value = "{request.header.origin}"
    }

    header "Access-Control-Allow-Headers" {
      value = "origin, x-requested-with, accept"
    }

    header "Access-Control-Max-Age" {
      value = "3628800"
    }

    header "Access-Control-Allow-Methods" {
      value = "GET, PUT, POST, DELETE"
    }

    form_param "formyp" {
      value = "formypvalue"
    }

    query_param "queryyp" {
      value = "queryypvalue"
    }
  }

  copy {
    source        = "Hello"
    header        "Accept"      {}
    header        "X-Requested-With"{}
    payload       = true
    version       = true
    path          = true
    verb          = true
    status_code   = true
    reason_phrase = true
  }

  remove {
    header      "Accept"    {}
    header      "X-Requested-With"{}
    form_param  "hello"     {}
    query_param "qhello"    {}
    payload     = true
  }

  set {
    header "Access-Control-Allow-Origin" {
      value = "{request.header.origin}"
    }

    header "Access-Control-Allow-Headers" {
      value = "origin, x-requested-with, accept"
    }

    header "Access-Control-Max-Age" {
      value = "3628800"
    }

    header "Access-Control-Allow-Methods" {
      value = "GET, PUT, POST, DELETE"
    }

    form_param "formyp" {
      value = "formypvalue"
    }

    query_param "queryyp" {
      value = "queryypvalue"
    }

    payload = {
      content_type = "application/json"
      value        = "{ \"hello\": \"world\" }"
    }

    verb          = "PUT"
    path          = "/hello"
    status_code   = 200
    version       = "3"
    reason_phrase = "OK"
  }

  assign_to {
    create_new = false
    transport  = "http"
    type       = "response"
    value      = "partnerresponse"
  }
}
