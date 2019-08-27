import * as React from 'react';
import { StringCardsList } from 'Components/Common/StringCardsList';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { isEqual } from 'lodash';

const DomainsContainer = styled.div`
  background: ${colors.januaryGrey};
  display: flex;
  justify-content: flex-start;
  align-content: center;
  align-items: center;
  padding: 7px 7px 3px;
  border-radius: ${soloConstants.smallRadius}px;
`;

interface Props {
  domains: string[];
  domainsChanged: (newDomainsList: string[]) => any;
}

function equivalentProps(
  oldProps: Readonly<Props>,
  nextProps: Readonly<Props>
): boolean {
  return isEqual(oldProps.domains, nextProps.domains);
}

export const Domains: React.FC<Props> = React.memo(props => {
  const [domains, setDomains] = React.useState(props.domains);

  React.useEffect(() => {
    if (!isEqual(props.domains, domains)) {
      setDomains(props.domains);
    }
  }, [props.domains]);

  // need to hook this up to api
  const addDomain = (domain: string) => {
    props.domainsChanged([...domains, domain]);
    setDomains([...domains, domain]);
  };

  const removeDomain = (removeIndex: number) => {
    let newList = [...domains];
    newList.splice(removeIndex, 1);

    props.domainsChanged(newList);
    setDomains(newList);
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
}, equivalentProps);
