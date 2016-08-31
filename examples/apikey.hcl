proxy "ApiKeySample" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "ValidateAPIKey"{}
      step "CheckQuota"{}
      step "Remove-Query-Param"{}
    }
  }

  http_proxy_connection {
    base_path    = "/mocktarget_key"
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

policy verify_api_key "ValidateAPIKey" {
  apikey {
    ref = "request.queryparam.apikey"
  }
}

policy quota "CheckQuota" {
  interval {
    ref = "verifyapikey.ValidateAPIKey.apiproduct.developer.quota.interval"
  }

  time_unit {
    ref = "verifyapikey.ValidateAPIKey.apiproduct.developer.quota.timeunit"
  }

  allow {
    count_ref = "verifyapikey.ValidateAPIKey.apiproduct.developer.quota.limit"
  }

  identifier {
    ref = "request.queryparam.apikey"
  }
}

policy assign_message "Remove-Query-Param" {
  ignore_unresolved_variables = true

  remove {
    query_param "apikey" {}
  }

  assign_to {
    create_new = false
    type       = "request"
  }
}
