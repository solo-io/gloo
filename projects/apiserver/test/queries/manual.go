package queries

import "strings"

// VirtualServices

const AllVirtualServices = `
{
	allNamespaces{
		name
		virtualServices{
			metadata{
				name
			}
		}
	}
}
`

const CreateVirtualServiceMutation = `
mutation MVS($vs: InputVirtualService!){
  virtualServices{
    create(virtualService: $vs){
      metadata{
        name
      }
    }
  }
}`

const CreateVirtualServiceVariables = `{
  "vs": {
    "metadata": {
      "name": "whoops6",
      "resourceVersion": "",
      "namespace": "mmm"
    },
    "displayName": "yesss"
  }
}`

const basicUpstreamsQueryNice = `{"query":"{
allNamespaces{
  name
    upstreams{
      metadata{
        name
      }}}
} "}`

var BasicUpstreamsQuery = strings.Replace(basicUpstreamsQueryNice, "\n", " ", -1)

const BasicUpstreamsQueryMatch = `{"data":{"allNamespaces":[{"name":`
