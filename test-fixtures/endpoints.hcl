proxy "helloworld" {
  display_name = "Hello World"
  description  = "Get started with your first API proxy."
}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "check-quota" {}

      # Handle preflight OPTIONS calls for cross origin requests
      step "add-cors" {
        condition = "request.verb == \"OPTIONS\""
      }
    }
  }

  flow "TokenEndpoint" {
    condition = "proxy.pathsuffix MatchesPath \"/accesstoken\""

    request {
      step "GenerateAccessToken" {}
    }
  }

  post_flow {
    response {
      step "FakePolicy" {}
    }
  }

  fault_rule "FaultRule1" {
    step      "FakeFaultRuleStep"{}
    condition = "proxy.pathsuffix MatchesPath \"/accesstoken\""
  }

  fault_rule "FaultRule2" {
    step      "FakeFaultRuleStep"{}
    step      "FakeFaultRuleStep2"{}
    condition = "proxy.pathsuffix MatchesPath \"/accesstoken\""
  }

  default_fault_rule "default" {
    step           "FakeFaultRuleStep"{}
    condition      = "proxy.pathsuffix MatchesPath \"/accesstoken\""
    always_enforce = true
  }

  http_proxy_connection {
    base_path    = "/v0/hello"
    virtual_host = ["default", "secure"]

    properties {
      "allow.http10"              = true
      "request.streaming.enabled" = true
    }
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

    load_balancer {
      algorithm = "weighted"

      server "server1" {
        weight = 1
      }

      server "server2" {
        weight      = 2
        is_fallback = true
      }
    }

    properties {
      testprop = 123.45
    }
  }

  local_target_connection {
    path           = "/v1/streetcarts"
    proxy_endpoint = "yolo"
    api_proxy      = "hey"
  }

  script_target {
    resource_url = "node://server.js"

    environment_variables {
      NAME = "VALUE"
    }

    arguments = ["arg1", "arg2"]
  }

  ssl_info {
    enabled             = true
    client_auth_enabled = false
    key_store           = "ref://keystoreref"
    key_alias           = "myKeyAlias"
    ciphers             = ["TLS_RSA_WITH_3DES_EDE_CBC_SHA", "TLS_RSA_WITH_DES_CBC_SHA"]
    protocols           = ["TLSv1", "TLSv1.2"]
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
  }

  assign_to {
    create_new = false
    transport  = "http"
    type       = "response"
  }
}

policy quota "check-quota" {
  async             = false
  continue_on_error = false
  enabled           = true
  type              = "calendar"
  display_name      = "Check Quota"

  allow {
    count     = 5
    count_ref = "request.header.allowed_quota"
  }

  interval {
    ref   = "request.header.quota_count"
    value = 1
  }

  distributed = false
  synchronous = false

  time_unit {
    ref   = "request.header.quota_timeout"
    value = "minute"
  }

  start_time = "2016-3-31 00:00:00"

  asynchronous_configuration {
    sync_interval_in_seconds = 20
    sync_message_count       = 5
  }
}
