// @ts-ignore
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';

export const getGlooInstanceStatus = (
  instance?: GlooInstance.AsObject,
  specificCheckSubfield?: string
): 0 | 1 | 2 | 3 => {
  const check = instance?.spec?.check || {};

  if (!!specificCheckSubfield) {
    const checkSummary: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject =
      // @ts-ignore getting past ts' inability to deal with [] referencing
      check[specificCheckSubfield];

    if (checkSummary?.errorsList.length) {
      return UpstreamStatus.State.REJECTED;
    }

    if (checkSummary?.warningsList.length) {
      return UpstreamStatus.State.WARNING;
    }
  } else {
    let warningSeen = false;
    Object.keys(check).forEach(key => {
      const checkSummary: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject =
        // @ts-ignore getting past ts' inability to deal with [] referencing
        check[key];

      if (checkSummary.errorsList.length > 0) {
        return UpstreamStatus.State.REJECTED;
      }

      if (checkSummary.warningsList.length > 0) {
        warningSeen = true;
      }
    });

    if (warningSeen) {
      return UpstreamStatus.State.WARNING;
    }
  }

  return UpstreamStatus.State.ACCEPTED;
};

export const getGlooInstanceListStatus = (
  glooInstances: GlooInstance.AsObject[] | undefined,
  specificCheckSubfield?: string
): 0 | 1 | 2 | 3 => {
  if (!glooInstances) {
    return UpstreamStatus.State.PENDING;
  }

  let warningSeen = false;
  glooInstances.forEach(instance => {
    const status = getGlooInstanceStatus(instance, specificCheckSubfield);
    if (status === UpstreamStatus.State.REJECTED) {
      return UpstreamStatus.State.REJECTED;
    }
    if (status === UpstreamStatus.State.WARNING) {
      warningSeen = true;
    }
  });
  if (warningSeen) {
    return UpstreamStatus.State.WARNING;
  }

  return UpstreamStatus.State.ACCEPTED;
};

export function sortGlooInstances(
  instanceA?: GlooInstance.AsObject,
  instanceB?: GlooInstance.AsObject
) {
  if (
    instanceB?.metadata?.name === undefined ||
    instanceB?.metadata?.namespace === undefined ||
    instanceB.metadata.clusterName === undefined
  ) {
    return 1;
  }
  if (
    instanceA?.metadata?.name === undefined ||
    instanceA?.metadata?.namespace === undefined ||
    instanceA.metadata.clusterName === undefined
  ) {
    return -1;
  }

  const nameDiff = instanceA.metadata.name.localeCompare(
    instanceB.metadata.name
  );
  const namespaceDiff = instanceA.metadata.namespace.localeCompare(
    instanceB.metadata.namespace
  );
  const clusterDiff = instanceA.metadata.clusterName.localeCompare(
    instanceB.metadata.clusterName
  );

  return nameDiff || namespaceDiff || clusterDiff;
}
