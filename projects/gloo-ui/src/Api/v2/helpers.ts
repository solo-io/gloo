import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import {
  BoolValue,
  UInt32Value
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';

export function getBoolVal(obj: BoolValue.AsObject): BoolValue {
  let boolVal = new BoolValue();
  boolVal.setValue(obj.value);
  return boolVal;
}
export function getUInt32Val(obj: UInt32Value.AsObject): UInt32Value {
  let intVal = new UInt32Value();
  intVal.setValue(obj.value);
  return intVal;
}

export function getResourceRef(name: string, namespace: string): ResourceRef {
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  return ref;
}

export function getStatus(name: string, namespace: string): ResourceRef {
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  return ref;
}

export function getDuration(seconds: number, nanos: number): Duration {
  let dur = new Duration();
  dur.setSeconds(seconds);
  dur.setNanos(nanos);
  return dur;
}

////// SET ATTEMPTS /////
export function setBoolVal(
  setFunc: (bv: BoolValue) => void,
  boolValObj?: BoolValue.AsObject
): void {
  if (!!boolValObj) {
    let boolVal = new BoolValue();
    boolVal.setValue(boolValObj.value);
    setFunc(boolVal);
  }
}
export function setUInt32Val(
  setFunc: (iv: UInt32Value) => void,
  uIntValObj?: UInt32Value.AsObject
): void {
  if (uIntValObj !== undefined) {
    let uintVal = new UInt32Value();
    uintVal.setValue(uIntValObj.value);
    setFunc(uintVal);
  }
}

export function setResourceRef(
  setFunc: (rr: ResourceRef) => void,
  resRefObj?: ResourceRef.AsObject
): void {
  if (resRefObj !== undefined) {
    let resRef = new ResourceRef();
    resRef.setName(resRefObj.name);
    resRef.setNamespace(resRefObj.namespace);
    setFunc(resRef);
  }
}
export function setDuration(
  setFunc: (d: Duration) => void,
  durationObj?: Duration.AsObject
): void {
  if (durationObj !== undefined) {
    let duration = new Duration();
    duration.setSeconds(
      Math.trunc(durationObj.seconds || durationObj.nanos / 1000)
    );
    duration.setNanos(
      durationObj.seconds ? durationObj.nanos : durationObj.nanos % 1000
    );
    setFunc(duration);
  }
}

export function setStatus(
  setFunc: (s: Status) => void,
  statusObj?: Status.AsObject
): void {
  if (statusObj !== undefined) {
    let status = new Status();
    status.setReason(statusObj.reason);
    status.setReportedBy(statusObj.reportedBy);
    status.setState(statusObj.state);
    setFunc(status);
  }
}

export function setMetadata(
  setFunc: (s: Metadata) => void,
  metadataObj?: Metadata.AsObject
): void {
  if (!!metadataObj) {
    let metadata = new Metadata();
    metadata.setName(metadataObj.name);
    metadata.setNamespace(metadataObj.namespace);
    metadata.setCluster(metadataObj.cluster);
    metadata.setResourceVersion(metadataObj.resourceVersion);
    setFunc(metadata);
  }
}
