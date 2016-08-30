proxy "base64encoder" {}

proxy_endpoint "default" {
  flow "UsernamePasswordAuthentication" {
    request {
      step "EncodeAuthHeader"{}
      step "EchoVariables"{}
    }
  }

  http_proxy_connection {
    base_path    = "/base64encoder"
    virtual_host = ["secure"]
  }

  route_rule "default" {}
}

target_endpoint "default" {
  http_target_connection {
    url = "http://weather.yahooapis.com"
  }
}

policy javascript "EncodeAuthHeader" {
  display_name = "EncodeAuthHeader"
  time_limit   = 200

  include_url = [
    "jsc://core-min.js",
    "jsc://enc-utf16-min.js",
    "jsc://enc-base64-min.js",
  ]

  resource_url = "jsc://encodeAuthHeader.js"
}

policy assign_message "EchoVariables" {
  assign_to {
    create_new = false
    type       = "response"
  }

  ignore_unresolved_variables = true

  set {
    header "X-Encoded-Credentials" {
      value = "{encodedAuthHeader}"
    }
  }
}
