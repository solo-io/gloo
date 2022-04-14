import { useGetGraphqlApiDetails } from 'API/hooks';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import {
  StatusHealth,
  WarningCircle,
} from 'Components/Features/Overview/OverviewBoxSummary';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';

const GraphqlDefineResolversPrompt: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  // gets the schema from the api
  const { data: graphqlApi, error: graphqlApiError } =
    useGetGraphqlApiDetails(apiRef);

  // This decides whether to show the resolver prompt,
  // which is shown above the details+explorer tabs.
  const [showResolverPrompt, setShowResolverPrompt] = React.useState(false);
  React.useEffect(() => {
    if (
      graphqlApi?.spec?.executableSchema?.executor === undefined ||
      graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap
        ?.length === 0
    ) {
      setShowResolverPrompt(true);
    } else {
      setShowResolverPrompt(false);
    }
  }, [
    !!graphqlApi?.spec?.executableSchema?.executor,
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.length,
  ]);

  if (!showResolverPrompt) return null;
  return (
    <div className='grid w-full '>
      <StatusHealth isWarning className=' place-content-center'>
        <div>
          <WarningCircle>
            <WarningExclamation />
          </WarningCircle>
        </div>
        <div>
          <>
            <div className='text-xl '>No Resolvers defined</div>
            <div className='text-lg '>Define resolvers</div>
          </>
        </div>
      </StatusHealth>
    </div>
  );
};

export default GraphqlDefineResolversPrompt;
