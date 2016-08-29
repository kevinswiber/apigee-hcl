proxy "quotafixture" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "check-quota" {}
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

policy quota "check-quota" {
  async             = false
  continue_on_error = false
  enabled           = true
  type              = "calendar"
  display_name      = "Check Quota"

  allow {
    count     = 5
    count_ref = "request.header.allowed_quota"

    class {
      ref = "request.queryparam.time_variable"

      allow {
        class = "peak_time"
        count = 5000
      }

      allow {
        class = "off_peak_time"
        count = 1000
      }
    }
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

  identifier {
    ref = "verifyapikey.verify-api-key.client_id"
  }

  message_weight {
    ref = "request.header.weight"
  }
}
