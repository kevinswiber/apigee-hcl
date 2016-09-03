proxy "learn-edge-service-callout-2" {}

proxy_endpoint "default" {
  fault_rule "InvalidApiKey" {
    condition = "(fault.Name matches \"InvalidApiKey\")"
    step      "AM-InvalidApiKeyMessage"{}
  }

  fault_rule "SC-Error" {
    condition = "(fault.name Matches \"ExecutionFailed\")"

    step "RF-RaiseCustomFault" {
      condition = "(mockresponse.status.code = \"404\")"
    }
  }

  pre_flow {
    request {
      step "VerifyApiKey"{}
      step "AM-BuildRequest"{}
      step "SC-GetMockResponse"{}
    }
  }

  http_proxy_connection {
    base_path    = "/v1/learn_edge"
    virtual_host = ["default", "secure"]
  }

  route_rule "noroute" {}
}

policy assign_message "AM-BuildRequest" {
  set {
    path = "/make-error"
  }

  assign_variable {
    name  = "target.copy.pathsuffix"
    value = "false"
  }

  ignore_unresolved_variables = false

  assign_to {
    type       = "request"
    create_new = true
  }
}

policy assign_message "AM-InvalidApiKeyMessage" {
  set {
    payload {
      content_type = "application/json"
      value        = "\\{\"error\": \\{\"message\":\"{fault.name}\", \"detail\":\"Please provide valid API key in the apikey query parameter.\"}}"
    }

    status_code   = 400
    reason_phrase = "Bad Request"
  }

  ignore_unresolved_variables = true
}

policy raise_fault "RF-RaiseCustomFault" {
  fault_response {
    set {
      payload {
        content_type = "application/json"
        value        = "\\{\"error\": \\{\"message\":\"Page Not Found\", \"details\":\"Hello from Learn Edge! This is a custom message..\"}}"
      }

      status_code   = 404
      reason_phrase = "Not Found"
    }

    ignore_unresolved_variables = true
  }
}

policy service_callout "SC-GetMockResponse" {
  request {
    clear_payload               = false
    variable                    = "myrequest"
    ignore_unresolved_variables = false
  }

  response = "mockresponse"

  http_target_connection {
    url = "http://mocktarget.apigee.net"
  }

  timeout = 30000
}

policy verify_api_key "VerifyApiKey" {
  apikey {
    ref = "request.queryparam.apikey"
  }
}
