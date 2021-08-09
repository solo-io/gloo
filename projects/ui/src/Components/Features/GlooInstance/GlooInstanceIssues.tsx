import React from 'react';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { getIssues } from 'utils/gloo-instance-check-helpers';
import {
  NotificationBox,
  NotificationType,
} from 'Components/Common/NotificationBox';

type Props = {
  glooInstance: GlooInstance.AsObject;
};

export const GlooInstanceIssues = ({ glooInstance }: Props) => {
  const { errors, warnings } = getIssues(glooInstance);
  if (!errors.length && !warnings.length) {
    return null;
  }
  return (
    <div>
      {errors.length ? (
        <NotificationBox
          type={NotificationType.ERROR}
          issues={errors}
          multipleIssuesMessage='There are multiple errors on this instance.'
        />
      ) : null}
      {warnings.length ? (
        <NotificationBox
          type={NotificationType.WARNING}
          issues={warnings}
          multipleIssuesMessage='There are multiple warnings on this instance.'
        />
      ) : null}
    </div>
  );
};
