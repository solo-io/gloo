<a name="top"/>

## Contents
  - [Status](#v1.Status)

  - [Status.State](#v1.Status.State)


<a name="status"/>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Status"/>

### Status
Status indicates whether a config resource (currently only [virtualhosts](TODO) and [upstreams](TODO)) has been (in)validated by gloo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [Status.State](#v1.Status.State) |  | State is the enum indicating the state of the resource |
| reason | [string](#string) |  | Reason is a description of the error for Rejected resources. If the resource is pending or accepted, this field will be empty |





 


<a name="v1.Status.State"/>

### Status.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| Pending | 0 | Pending status indicates the resource has not yet been validated |
| Accepted | 1 | Accepted indicates the resource has been validated |
| Rejected | 2 | Rejected indicates an invalid configuration by the user Rejected resources may be propagated to the xDS server depending on their severity |


 

 

