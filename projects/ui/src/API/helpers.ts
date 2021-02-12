import { ObjectMeta } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb';
import {
  ObjectRef,
  ClusterObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export const host = `${
  process.env.NODE_ENV === 'production'
    ? window.location.origin
    : 'http://localhost:8090'
}`;

export function getObjectRefClassFromRefObj(
  requestMeta: ObjectRef.AsObject
): ObjectRef {
  let ref = new ObjectRef();
  ref.setName(requestMeta.name);
  ref.setNamespace(requestMeta.namespace);
  return ref;
}

export function getClusterRefClassFromClusterRefObj(
  requestMeta: ClusterObjectRef.AsObject
): ClusterObjectRef {
  let ref = new ClusterObjectRef();
  ref.setName(requestMeta.name);
  ref.setNamespace(requestMeta.namespace);
  ref.setClusterName(requestMeta.clusterName);
  return ref;
}

export function setObjectMeta(
  setFunc: (rr: ObjectMeta) => void,
  resRefObj?: ObjectMeta.AsObject
): void {
  if (resRefObj !== undefined) {
    let resRef = new ObjectMeta();
    resRef.setName(resRefObj.name);
    resRef.setNamespace(resRefObj.namespace);
    setFunc(resRef);
  }
}

export function objectMetasAreEqual(
  resRef1?:
    | ObjectRef.AsObject
    | ObjectMeta.AsObject
    | ClusterObjectRef.AsObject,
  resRef2?:
    | ObjectRef.AsObject
    | ObjectMeta.AsObject
    | ClusterObjectRef.AsObject,
  onlyCompareRefs?: boolean
) {
  if (!resRef1 || !resRef2) {
    return false;
  }

  if (!onlyCompareRefs) {
    if (
      // @ts-ignore
      (resRef1.clusterName !== undefined || // @ts-ignore
        resRef2.clusterName !== undefined) &&
      // @ts-ignore
      resRef1.clusterName !== resRef2.clusterName
    ) {
      return false;
    }
  }

  return (
    resRef1.name === resRef2.name && resRef1.namespace === resRef2.namespace
  );
}
