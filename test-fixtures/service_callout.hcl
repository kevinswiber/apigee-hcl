proxy "ServiceCalloutFixture" {}

proxy_endpoint "default" {
  pre_flow {
    request {
      step "service-callout" {}
    }
  }

  http_proxy_connection {
    base_path    = "/v0/service-callout"
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

policy service_callout "service-callout" {
  display_name = "Inline Request Message"

  request {
    variable = "authenticationRequest"

    set {
      query_param "address" {
        value = "{request.queryparam.postalcode}"
      }

      query_param "region" {
        value = "{request.queryparam.country}"
      }

      query_param "sensor" {
        value = "false"
      }
    }
  }

  response = "GeoCodingResponse"

  timeout = 30000

  http_target_connection {
    url = "http://maps.googleapis.com/maps/api/geocode/json"
  }
}
