import styled from '@emotion/styled';
import { useGetGraphqlApiDetails } from 'API/hooks';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import {
  StatusHealth,
  WarningCircle,
} from 'Components/Features/Overview/OverviewBoxSummary';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { colors } from 'Styles/colors';

const StyledBanner = styled.div`
  background: ${colors.flashlightGold};
  border: solid 1px ${colors.darkFebruaryGrey};
  padding: 10px;
  margin-bottom: 20px;
`;

const StyledStatusHealth = styled(StatusHealth)`
  margin-bottom: 0;
  margin-left: 50px;
  display: flex;
  flex-direction: row;
  align-items: center;
`;

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
    <StyledBanner className='grid w-full '>
      <StyledStatusHealth isWarning className=''>
        <div>
          <WarningCircle>
            <WarningExclamation />
          </WarningCircle>
        </div>
        <div>
          <div className='text-xl'>No Resolvers defined</div>
        </div>
      </StyledStatusHealth>
    </StyledBanner>
  );
};

export default GraphqlDefineResolversPrompt;
