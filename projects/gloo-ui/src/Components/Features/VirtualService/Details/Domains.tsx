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
  padding: 7px 7px 3px;
`;

interface Props {
  domains: string[];
  domainsChanged: (newDomainsList: string[]) => any;
}
export const Domains: React.FC<Props> = props => {
  const [domains, setDomains] = React.useState(props.domains);

  // need to hook this up to api
  const addDomain = (domain: string) => {
    setDomains([...domains, domain]);
    props.domainsChanged([...domains, domain]);
  };

  const removeDomain = (removeIndex: number) => {
    let newList = [...domains];
    newList.splice(removeIndex, 1);
    setDomains(newList);

    props.domainsChanged(newList);
  };

  const isValidDomain = (domain: string) => {
    if (domain === '*') {
      return true;
    }

    return RegExp(
      '(?:[a-z0-9*](?:[a-z0-9-]{0,61}[a-z0-9])?.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]'
    ).test(domain);
  };

  return (
    <React.Fragment>
      <DetailsSectionTitle>Domains</DetailsSectionTitle>
      <DomainsContainer>
        <StringCardsList
          values={domains}
          valueDeleted={removeDomain}
          createNewPromptText='new.domain.com'
          createNew={addDomain}
          valueIsValid={isValidDomain}
        />
      </DomainsContainer>
    </React.Fragment>
  );
};
