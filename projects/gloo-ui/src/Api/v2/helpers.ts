import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';

export function getResourceRef(name: string, namespace: string): ResourceRef {
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  return ref;
}
