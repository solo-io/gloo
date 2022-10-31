// @ts-ignore
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';

export const getGlooInstanceStatus = (
  instance?: GlooInstance.AsObject,
  specificCheckSubfield?: keyof GlooInstance.GlooInstanceSpec.Check.AsObject
): 0 | 1 | 2 | 3 => {
  const check = instance?.spec?.check || {};

  if (!!specificCheckSubfield) {
    const checkSummary = check[specificCheckSubfield];

    if (checkSummary?.errorsList.length) {
      return UpstreamStatus.State.REJECTED;
    }

    if (checkSummary?.warningsList.length) {
      return UpstreamStatus.State.WARNING;
    }
  } else {
    let warningSeen = false;
    Object.values(check).forEach(checkSummary => {
      if (checkSummary === undefined) {
        // This defaults to pending if there is no data returned (errors or otherwise).
        return UpstreamStatus.State.PENDING;
      }

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
  specificCheckSubfield?: keyof GlooInstance.GlooInstanceSpec.Check.AsObject
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
    instanceB?.spec?.cluster === undefined
  ) {
    return 1;
  }
  if (
    instanceA?.metadata?.name === undefined ||
    instanceA?.metadata?.namespace === undefined ||
    instanceA?.spec?.cluster === undefined
  ) {
    return -1;
  }

  const nameDiff = instanceA.metadata.name.localeCompare(
    instanceB.metadata.name
  );
  const namespaceDiff = instanceA.metadata.namespace.localeCompare(
    instanceB.metadata.namespace
  );
  const clusterDiff = instanceA.spec.cluster.localeCompare(
    instanceB.spec.cluster
  );

  return nameDiff || namespaceDiff || clusterDiff;
}
