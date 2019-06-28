import * as React from 'react';
import { StringCardsList } from 'Components/Common/StringCardsList';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { DetailsSectionTitle } from './VirtualServiceDetails';

const DomainsContainer = styled.div`
  background: ${colors.januaryGrey};
  display: flex;
  justify-content: flex-start;
  align-content: center;
  align-items: center;
`;

export const Domains = () => {
  const [domains, setDomains] = React.useState(['solo.io', 'domain.com']);

  const addDomain = (domain: string) => {
    setDomains([...domains, domain]);
  };

  const removeDomain = (removeIndex: number) => {
    let newList = [...domains];
    newList.splice(removeIndex, 1);
    setDomains(newList);
  };
  return (
    <React.Fragment>
      <DetailsSectionTitle>Domains</DetailsSectionTitle>
      <DomainsContainer>
        <StringCardsList
          values={domains}
          valueDeleted={removeDomain}
          createNewPromptText='new domain'
          createNew={addDomain}
        />
      </DomainsContainer>
    </React.Fragment>
  );
};
