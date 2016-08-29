proxy "helloworld" {
  display_name = "helloworld"
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
  display_name                = "Add CORS"
  ignore_unresolved_variables = true

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
    type       = "response"
  }
}

policy quota "check-quota" {
  display_name = "Check Quota"
  type         = "calendar"

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
