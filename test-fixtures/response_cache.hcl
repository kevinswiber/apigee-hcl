proxy "ResponseCachefixture" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "response-cache" {}
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
  pre_flow {
    response {
      step "response-cache" {}
    }
  }
  http_target_connection {
    url = "http://mocktarget.apigee.net"
  }
}

policy response_cache "response-cache" {
  async             = false
  continue_on_error = false
  enabled           = true
  display_name      = "Response Cache"
  lookup_timeout = 45
  cache_resource = "default"
  exclude_error_response = true
  cache_key {
    prefix = "UserToken"
    key_fragment {
      value = "apiAccessToken"
    }
    key_fragment {
      ref = "request.uri"
    }
  }
  scope = "Exclusive"
  expiry_settings {
    timeout_in_sec {
      value = 3600
    }
    expiry_date {
      ref = "expiryDate"
    }
    time_of_day {
      value = "14:30:00"
    }
  }
  skip_cache_lookup = "request.header.skipcache = true"
  skip_cache_population = "request.header.skipcache = true"
  use_accept_header = false
  use_response_cache_headers = false
}
