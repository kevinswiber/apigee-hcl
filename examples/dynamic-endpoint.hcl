proxy "dynamic-endpoint" {}

proxy_endpoint "default" {
  http_proxy_connection {
    base_path    = "/dynamic-endpoint"
    virtual_host = ["default"]
  }

  route_rule "fbroute" {
    condition       = "request.queryparam.routeTo = \"fb\""
    target_endpoint = "facebook"
  }

  route_rule "twroute" {
    condition       = "request.queryparam.routeTo = \"tw\""
    target_endpoint = "twitter"
  }

  route_rule "default" {
    target_endpoint = "twitter"
  }
}

target_endpoint "facebook" {
  http_target_connection {
    url = "https://api.facebook.com"
  }
}

target_endpoint "twitter" {
  http_target_connection {
    url = "https://dev.twitter.com/rest/public"
  }
}
