import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors, healthConstants, soloConstants } from 'Styles';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { EnvoyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/envoy_pb';
import { SectionCard } from 'Components/Common/SectionCard';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { JsonDisplayer } from 'Components/Common/DisplayOnly/JsonDisplayer';
import { AppState } from 'store';
import { useSelector } from 'react-redux';
import { Status } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import { TallyContainer } from 'Components/Common/DisplayOnly/TallyInformationDisplay';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 18px;
  color: ${colors.novemberGrey};
`;

const EnvoyLogoFullSize = styled(EnvoyLogo)`
  width: 33px !important;
  max-height: none !important;
`;

const ExpandableSection = styled<'div', { isExpanded: boolean }>('div')`
  max-height: ${props => (props.isExpanded ? '1000px' : '0px')};
  overflow: hidden;
  transition: max-height ${soloConstants.transitionTime};
  color: ${colors.septemberGrey};
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

interface Props {}
export const getEnvoyHealth = (code: Status.Code): number => {
  switch (code) {
    case Status.Code.ERROR:
      return healthConstants.Error.value;
    case Status.Code.OK:
      return healthConstants.Good.value;
    default:
      return healthConstants.Pending.value;
  }
};
export const Envoy = (props: Props) => {
  const envoysList = useSelector(
    (state: AppState) => state.envoy.envoyDetailsList
  );
  const [allEnvoys, setAllEnvoys] = React.useState<EnvoyDetails.AsObject[]>([]);

  React.useEffect(() => {
    if (!!envoysList.length) {
      setAllEnvoys(envoysList);
    }
  }, [envoysList.length]);

  const [envoysOpen, setEnvoysOpen] = React.useState<boolean[]>([]);

  React.useEffect(() => {
    if (!!envoysList.length) {
      setAllEnvoys(envoysList);
      setEnvoysOpen(envoysList.map(e => false));
    }
  }, [envoysList.length]);

  if (!envoysList.length) {
    return <div>Loading...</div>;
  }

  if (!allEnvoys.length) {
    return <div>You have no Envoy configurations.</div>;
  }

  const toggleExpansion = (indexToggled: number) => {
    setEnvoysOpen(
      envoysOpen.map((isOpen, ind) => {
        if (ind !== indexToggled) {
          return false;
        }

        return !isOpen;
      })
    );
  };

  return (
    <React.Fragment>
      {allEnvoys.map((envoy, ind) => {
        return (
          <SectionCard
            key={envoy.name + ind}
            cardName={envoy.name}
            logoIcon={<EnvoyLogoFullSize />}
            headerSecondaryInformation={[]}
            health={getEnvoyHealth(envoy!.status!.code!)}
            healthMessage={'Envoy Status'}>
            {envoy!.status!.message !== '' && (
              <TallyContainer color='orange'>
                {envoy!.status!.message!}
              </TallyContainer>
            )}
            <InsideHeader>
              <div>Code Log (Read Only)</div>{' '}
              {!!envoy.raw && (
                <FileDownloadLink
                  fileName={envoy.raw.fileName}
                  fileContent={envoy.raw.content}
                />
              )}
            </InsideHeader>
            {!!envoy.raw && (
              <React.Fragment>
                <ExpandableSection isExpanded={envoysOpen[ind]}>
                  {' '}
                  <JsonDisplayer content={envoy.raw.content} />
                </ExpandableSection>
                <Link onClick={() => toggleExpansion(ind)}>
                  {envoysOpen[ind] ? 'Hide' : 'View'} Settings
                </Link>
              </React.Fragment>
            )}
          </SectionCard>
        );
      })}
    </React.Fragment>
  );
};
