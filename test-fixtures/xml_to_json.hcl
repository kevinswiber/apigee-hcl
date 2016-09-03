proxy "xmltojsonfixture" {}

proxy_endpoint "default" {
  http_proxy_connection {
    base_path    = "/v1/xmltojson"
    virtual_host = ["default", "secure"]
  }

  route_rule "default" {
    target_endpoint = "default"
  }
}

target_endpoint "default" {
  pre_flow {
    response {
      step "xmltojson1" {}
    }
  }

  http_target_connection {
    url = "http://mocktarget.apigee.net"
  }
}

policy xml_to_json "xmltojson1" {
  display_name = "XML to JSON 1"
  source       = "response"

  options {
    output_variable             = "response"
    recognize_number            = true
    recognize_boolean           = true
    recognize_null              = true
    null_value                  = "NULL"
    namespace_block_name        = "#namespaces"
    default_namespace_node_name = "&"
    namespace_separator         = "***"
    text_always_as_property     = true
    text_node_name              = "TEXT"
    attribute_block_name        = "FOOBLOCK"
    attribute_prefix            = "BAR_"
    output_prefix               = "PREFIX_"
    output_suffix               = "_SUFFIX"
    strip_levels                = 2

    treat_as_array {
      path {
        unwrap = true
        value  = "teachers/teacher/studentnames/name"
      }
    }
  }

  # format = "yahoo"
}
