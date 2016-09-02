proxy "ExtractVariablesFixture" {}

proxy_endpoint "default" {
  http_proxy_connection {
    base_path    = "/v0/variables"
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
      step "extract-vars" {}
    }
  }
}

policy extract_variables "extract-vars" {
  source {
    clear_payload = true
    value         = "response"
  }

  variable_prefix             = "myprefix"
  ignore_unresolved_variables = true

  uri_path {
    pattern {
      ignore_case = true
      value       = "/accounts/{id}"
    }
  }

  query_param "code" {
    pattern {
      ignore_case = true
      value       = "DBN{dbncode}"
    }
  }

  header "Authorization" {
    pattern {
      ignore_case = false
      value       = "Bearer {oauthtoken}"
    }
  }

  form_param "greeting" {
    pattern {
      value = "hello {user}"
    }
  }

  variable "request.content" {
    pattern {
      value = "hello {user}"
    }
  }

  json_payload {
    variable "name" {
      type      = "string"
      json_path = "{example}"
    }
  }

  xml_payload {
    namespace "apigee" {
      value = "http://apigee.com"
    }

    namespace "gmail" {
      value = "http://mai.google.com"
    }

    variable "name" {
      type  = "boolean"
      xpath = "/apigee:test/apigee:example"
    }

    variable "name2" {
      type  = "boolean"
      xpath = "/apigee:test/apigee:example2"
    }
  }
}
