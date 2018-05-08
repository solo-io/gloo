<a name="top"></a>

## Contents
  - [Role](#v1.Role)



<a name="role"></a>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Role"></a>

### Role
A Role is a container for a set of Virtual Services that will be used to generate a single proxy config
to be applied to one or more Envoy nodes. The Role is best understood as an in-mesh application&#39;s localized view
of the rest of the mesh.
Each domain for each Virtual Service contained in a Role cannot appear more than once, or the Role
will be invalid.


```yaml
name: string
virtual_services: [string]
status: (read only)
metadata: {Metadata}

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | string |  | Name of the role. Envoy nodes will be assigned a config corresponding with role they are assigned. Envoy instances must specify the role they belong to when they register to Gloo. Currently this is done by specifying the name of the role as a prefix to the Envoy&#39;s Node ID which can be specified with the `--service-node` flag, or in the Envoy instance&#39;s bootstrap config. Names must be unique and follow the following syntax rules: One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| virtual_services | string | repeated | the list of names of the virtual services this role encapsulates. |
| status | [Status](status.md#v1.Status) |  | Status indicates the validation status of the role resource. Status is read-only by clients, and set by gloo during validation |
| metadata | [Metadata](metadata.md#v1.Metadata) |  | Metadata contains the resource metadata for the role |





 

 

 

