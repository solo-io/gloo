import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors, healthConstants, soloConstants } from 'Styles';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { useGetEnvoyList } from 'Api/v2/useEnvoyClientV2';
import { EnvoyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/envoy_pb';
import { SectionCard } from 'Components/Common/SectionCard';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { JsonDisplayer } from 'Components/Common/DisplayOnly/JsonDisplayer';

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

export const Envoy = (props: Props) => {
  const [envoysOpen, setEnvoysOpen] = React.useState<boolean[]>([]);

  const { data, loading, error, setNewVariables } = useGetEnvoyList();
  const [allEnvoys, setAllEnvoys] = React.useState<EnvoyDetails.AsObject[]>([]);

  React.useEffect(() => {
    if (!!data) {
      setAllEnvoys(data.toObject().envoyDetailsList);
      setEnvoysOpen(data.toObject().envoyDetailsList.map(e => false));
    }
  }, [loading]);

  if (!data || (!data && loading)) {
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
            health={healthConstants.Good.value}
            healthMessage={'Envoy Status'}>
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
                  <JsonDisplayer content={envoy.raw.content} />}
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
