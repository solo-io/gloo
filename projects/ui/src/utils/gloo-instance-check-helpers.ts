import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/glooinstance_pb';
import { Issue } from 'Components/Common/NotificationBox';

const getLink = (
  key: string,
  cluster?: string,
  glooInstanceRef?: ObjectRef.AsObject,
  objRef?: ObjectRef.AsObject
): { detailsLink: string; linkTitle: string } => {
  switch (key) {
    case 'gateways':
      return {
        detailsLink: glooInstanceRef
          ? `/gloo-instances/${glooInstanceRef.namespace}/${glooInstanceRef.name}/gloo-admin/gateways`
          : '',
        linkTitle: 'View Gateway Details',
      };

    case 'proxies':
      return {
        detailsLink: glooInstanceRef
          ? `/gloo-instances/${glooInstanceRef.namespace}/${glooInstanceRef.name}/gloo-admin/proxy`
          : '',
        linkTitle: 'View Proxy Details',
      };

    case 'settings':
      return {
        detailsLink: glooInstanceRef
          ? `/gloo-instances/${glooInstanceRef.namespace}/${glooInstanceRef.name}/gloo-admin/settings`
          : '',
        linkTitle: 'View Settings Details',
      };

    case 'virtualServices':
      return {
        detailsLink:
          glooInstanceRef && objRef && cluster
            ? `/gloo-instances/${glooInstanceRef.namespace}/${glooInstanceRef.name}/virtual-services/${cluster}/${objRef.namespace}/${objRef.name}`
            : '',
        linkTitle: 'View Virtual Service Details',
      };

    case 'upstreams':
      return {
        detailsLink: `/gloo-instances/:namespace/:name/upstreams/:upstreamClusterName/:upstreamNamespace/:upstreamName`,
        linkTitle: 'View Upstream Details',
      };

    case 'upstreamGroups':
      return {
        detailsLink: `/gloo-instances/:namespace/:name/upstream-groups/:upstreamGroupClusterName/:upstreamGroupNamespace/:upstreamGroupName`,
        linkTitle: 'View Upstream Group Details',
      };
  }
  return { detailsLink: '', linkTitle: '' };
};

export const getIssues = (
  glooInstance: GlooInstance.AsObject
): {
  errors: Issue[];
  warnings: Issue[];
} => {
  const glooInstanceRef: ObjectRef.AsObject | undefined =
    glooInstance.metadata?.name && glooInstance.metadata?.namespace
      ? {
          name: glooInstance.metadata.name,
          namespace: glooInstance.metadata.namespace,
        }
      : undefined;
  const cluster: string | undefined = glooInstance.spec?.cluster;

  const errors: Issue[] = [];
  const warnings: Issue[] = [];
  const check = glooInstance.spec?.check;
  if (check) {
    for (const [key, checkSummary] of Object.entries(check)) {
      const errorsList = checkSummary?.errorsList;
      if (errorsList?.length) {
        errorsList.forEach(({ message, ref }) => {
          errors.push({
            message,
            ...getLink(key, cluster, glooInstanceRef, ref),
          });
        });
      }
      const warningsList = checkSummary?.warningsList;
      if (warningsList?.length) {
        warningsList.forEach(({ message, ref }) => {
          warnings.push({
            message,
            ...getLink(key, cluster, glooInstanceRef, ref),
          });
        });
      }
    }
  }

  return {
    errors,
    warnings,
  };
};
