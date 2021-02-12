import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { FailoverSchemeSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover_pb';
import { ReactComponent as ArrowDown } from 'assets/arrow-toggle.svg';
import { ReactComponent as ClusterIcon } from 'assets/cluster-instance-icon.svg';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as WeightIcon } from 'assets/weight-balance-icon.svg';
import { colors } from 'Styles/colors';
import { CardWhiteSubsection } from 'Components/Common/Card';
import { CircleIconHolder, IconHolder } from 'Styles/StyledComponents/icons';
import UpstreamFailoverGroupTable from './UpstreamFailoverGroupTable';

const QuickStats = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;

const QuickStatContainer = styled.div`
  display: flex;
  align-items: center;
  margin-right: 45px;
`;

const PriorityValue = styled.div`
  margin-left: 15px;
  font-weight: bold;
  font-size: 14px;
`;

const QuickStatValue = styled.div`
  margin-left: 15px;
  font-weight: bold;
  font-size: 17px;
`;

const QuickStatTitle = styled.div`
  font-size: 14px;
  margin-left: 4px;
`;

const Divider = styled.div`
  margin: 0 30px;
  width: 1px;
  height: 42px;
  background: ${colors.marchGrey};
`;

const PriorityNumber = styled.div`
  font-size: 20px;
  color: white;
`;

const getPriorityColor = (priority: number): string => {
  switch (priority) {
    case 1:
      return colors.planeOfWaterBlue;
    case 2:
      return colors.oceanBlue;
    case 3:
      return colors.lakeBlue;
    default:
      return colors.puddleBlue;
  }
};

type ArrowIconProps = {
  expanded: boolean;
};
const ArrowIconHolder = styled(IconHolder)<ArrowIconProps>`
  ${(props: ArrowIconProps) => props.expanded && `transform: rotate(180deg);`}
`;

type Props = {
  priority: number;
  group: FailoverSchemeSpec.FailoverEndpoints.AsObject;
};

const UpstreamFailoverGroup = ({ priority, group }: Props) => {
  const [isExpanded, setExpanded] = useState(false);

  const numClusters = group.priorityGroupList?.length ?? 0;
  const numUpstreams = group.priorityGroupList?.reduce(
    (sum, group) => sum + group.upstreamsList?.length ?? 0,
    0
  );
  const isWeighted = group.priorityGroupList?.some(
    group => group.localityWeight !== undefined
  );

  const toggleExpanded = () => {
    // if user is selecting text, don't toggle
    if (getSelection()?.toString().length) {
      return;
    }
    setExpanded(exp => !exp);
  };

  return (
    <CardWhiteSubsection>
      <QuickStats onClick={toggleExpanded}>
        <QuickStatContainer>
          <CircleIconHolder backgroundColor={getPriorityColor(priority)}>
            <PriorityNumber>{priority}</PriorityNumber>
          </CircleIconHolder>
          <PriorityValue>Priority: {priority}</PriorityValue>
        </QuickStatContainer>
        <Divider />
        <QuickStatContainer>
          <CircleIconHolder
            backgroundColor={colors.seaBlue}
            iconColor={{ strokeNotFill: true, color: 'white' }}>
            <UpstreamIcon />
          </CircleIconHolder>
          <QuickStatValue>{numUpstreams}</QuickStatValue>
          <QuickStatTitle>{`Upstream ${
            numUpstreams === 1 ? 'Endpoint' : 'Endpoints'
          }`}</QuickStatTitle>
        </QuickStatContainer>
        <QuickStatContainer>
          <IconHolder width={33} applyColor={{ color: colors.seaBlue }}>
            <ClusterIcon />
          </IconHolder>
          <QuickStatValue>{numClusters}</QuickStatValue>
          <QuickStatTitle>
            {numClusters === 1 ? 'Cluster' : 'Clusters'}
          </QuickStatTitle>
        </QuickStatContainer>
        <QuickStatContainer>
          <CircleIconHolder
            backgroundColor={isWeighted ? colors.seaBlue : colors.juneGrey}
            iconColor={{ color: 'white' }}
            iconSize='24px'>
            <WeightIcon />
          </CircleIconHolder>
          <QuickStatValue>{isWeighted ? 'Varied' : 'Equal'}</QuickStatValue>
          <QuickStatTitle>Weights</QuickStatTitle>
        </QuickStatContainer>
        <ArrowIconHolder
          expanded={isExpanded}
          width={20}
          applyColor={{ color: colors.septemberGrey }}>
          <ArrowDown />
        </ArrowIconHolder>
      </QuickStats>
      {isExpanded && (
        <UpstreamFailoverGroupTable group={group} isWeighted={isWeighted} />
      )}
    </CardWhiteSubsection>
  );
};

export default UpstreamFailoverGroup;
