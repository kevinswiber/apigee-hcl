proxy "SpikeArrestfixture" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "spike-arrest" {}
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

policy spike_arrest "spike-arrest" {
  async             = false
  continue_on_error = false
  enabled           = true
  display_name      = "Spike Arrest"

  identifier {
    ref   = "request.header.apikey"
  }

  message_weight {
    ref = "request.header.weight"
  }

  rate {
    value = "10pm"
  }
}
