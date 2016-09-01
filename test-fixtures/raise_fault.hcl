proxy "RaiseFaultFixture" {
  display_name = "Raise Fault Fixture"
}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "throw-418" {}
      step "throw-500" {}
    }
  }

  http_proxy_connection {
    base_path    = "/v0/hello"
    virtual_host = ["default", "secure"]
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

policy raise_fault "throw-500" {
  continue_on_error           = false
  enabled                     = true
  display_name                = "Throw 500"
  ignore_unresolved_variables = false
  fault_response {
    set {
      header "brewing" {
        value = "always"
      }
      status_code   = 500
      reason_phrase = "Server Error"
    }
  }
}

policy raise_fault "throw-418" {
  continue_on_error           = false
  enabled                     = true
  display_name                = "Throw 418"
  fault_response {
  remove {
    header      "Accept"    {}
    header      "X-Requested-With"{}
  }

  set {
    header "brewing" {
      value = "always"
    }
    payload = {
      content_type = "application/json"
      variable_prefix = "$"
      variable_suffix = "#"
      value        = "{ \"status\": \"418\", \"reason\": \"I'm a teapot.\" }"
    }
    status_code   = 418
    reason_phrase = "I'm a teapot."
  }
  }
}
