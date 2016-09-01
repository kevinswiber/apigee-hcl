# Apigee HCL

* Write Apigee proxy configuration using [HCL (HashiCorp Configuration Language)](https://github.com/hashicorp/hcl).
* Output valid XML-based proxy bundles for deployment to Apigee.

> The goal of HCL is to build a structured configuration language that is both human and machine friendly for use with command-line tools, but specifically targeted towards DevOps tools, servers, etc.

## Example

```hcl
# hello.hcl

proxy "hello" {
  display_name = "Hello"
}

proxy_endpoint "default" {
  pre_flow {
    request {
      # Handle preflight OPTIONS calls for cross origin requests
      step "add-cors" {
        condition = "request.verb == \"OPTIONS\""
      }
    }
  }

  http_proxy_connection {
    base_path    = "/v1/hello"
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
```

### Generate proxy bundle

`$ apigee-hcl -i hello.hcl -o ./build`

This will generate an Apigee API proxy based on the `hello.hcl` configuration.  
The output will be generated into the `./build` directory.

The bundle can then be deployed using [apigeetool](https://github.com/apigee/apigeetool-node).

## Install

If you have Go v1.6+ installed, simply:

`$ go install github.com/kevinswiber/apigee-hcl`


## Policy support

- [x] Assign Message
- [x] JavaScript
- [ ] Extract Variables
- [x] Raise Fault
- [ ] Service Callout
- [ ] OAuth v2.0
- [x] Verify API Key
- [x] Response Cache
- [ ] XML to JSON
- [x] Spike Arrest
- [x] Quota
- [ ] XSL Transform
- [ ] Basic Authentication
- [ ] Statistics Collector
- [ ] Key Value Map Operations
- [ ] Message Logging
- [ ] Populate Cache
- [ ] Lookup Cache
- [ ] JSON to XML
- [ ] Access Control
- [ ] Java Callout
- [ ] JSON Threat Protection
- [ ] Access Entity
- [ ] SOAP Message Validation
- [ ] Regular Expression Protection
- [ ] Concurrent Rate Limit
- [ ] XML Threat Protection
- [ ] Generate SAML Assertion
- [ ] Invalidate Cache
- [ ] Set OAuth v2.0 Info
- [ ] Get OAuth v2.0 Info
- [ ] Monetization Limits Check
- [ ] OAuth v1.0a
- [ ] Reset Quota
- [x] Python Script

## License

Apache License, v2.0 
