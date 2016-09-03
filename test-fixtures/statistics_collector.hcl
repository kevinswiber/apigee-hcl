proxy "StatisticsCollectorFixture" {}

proxy_endpoint "default" {
  http_proxy_connection {
    base_path    = "/v0/stats"
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

  pre_flow {
    response {
      step "publishPurchaseDetails" {}
    }
  }
}

policy statistics_collector "publishPurchaseDetails" {
  statistic "productID" {
    ref   = "product.id"
    type  = "string"
    value = "999999"
  }

  statistic "price" {
    ref   = "product.price"
    type  = "string"
    value = "0"
  }

  statistic "hi" {
    ref  = "hi"
    type = "string"
  }
}
