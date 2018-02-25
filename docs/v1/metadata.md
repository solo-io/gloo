<a name="top"/>

## Contents
  - [Metadata](#v1.Metadata)



<a name="metadata"/>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Metadata"/>

### Metadata
[]()Metadata contains general properties of config resources useful to clients and the gloo control plane for purposes of versioning, annotating, and namespacing resources.


```yaml
resource_version: string
namespace: string
annotations: map<string,string>

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_version | [string](#string) |  | ResourceVersion keeps track of the resource version of a config resource. This mechanism is used by [gloo-storage](TODO) to ensure safety with concurrent writes/updates to a resource in storage. |
| namespace | [string](#string) |  | Namespace is used for the namespacing of resources. Currently unused by gloo internally. |
| annotations | [map&lt;string,string&gt;](#map&lt;string,string&gt;) |  | Annotations allow clients to tag resources for special use cases. gloo ignores annotations but preserved them on read/write from/to storage. |





 

 

 

