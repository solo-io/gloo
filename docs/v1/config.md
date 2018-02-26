<a name="top"/>

## Contents
  - [Config](#v1.Config)



<a name="config"/>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Config"/>

### Config
[]()Config is a top-level config object. It is used internally by gloo as a container for the entire user config.


```yaml
upstreams: [{Upstream}]
virtual_hosts: [{VirtualHost}]

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| upstreams | [Upstream](upstream.md#v1.Upstream) | repeated | The list of all upstreams defined by the user. |
| virtual_hosts | [VirtualHost](virtualhost.md#v1.VirtualHost) | repeated | the list of all virtual hosts defined by the user. |





 

 

 

