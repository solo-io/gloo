import styled from '@emotion/styled';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import {
  CheckboxFilterProps,
  ListingFilter,
  RadioFilterProps,
  StringFilterProps,
  TypeFilterProps
} from 'Components/Common/ListingFilter';
import * as React from 'react';
import { Route, RouteComponentProps, Switch } from 'react-router';
import { Envoy } from './Envoy';
import { Gateways } from './Gateways';
import { Proxys } from './Proxy';

const PageChoiceFilter: TypeFilterProps = {
  id: 'pageChoice',
  options: [
    {
      displayName: 'Gateways'
    },
    {
      displayName: 'Proxy'
    },
    {
      displayName: 'Envoy'
    }
  ],
  choice: 'Gateways'
};

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

interface Props extends RouteComponentProps<{ sublocation: string }> {}

export const AdminHub = (props: Props) => {
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);

  const locationChoice =
    props.match.params.sublocation.charAt(0).toUpperCase() +
    props.match.params.sublocation.slice(1);

  const pageChanged = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    const newChoice = types.find(type => type.id === 'pageChoice')!.choice!;
    props.history.replace({
      pathname: `/admin/${newChoice.toLowerCase()}`
    });
  };

  const listDisplay = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ): React.ReactNode => {
    return (
      <React.Fragment>
        <Switch>
          <Route path='/admin/gateways/' component={Gateways} />
          <Route path='/admin/proxy/' component={Proxys} />
          <Route path='/admin/envoy/' component={Envoy} />
        </Switch>
      </React.Fragment>
    );
  };

  return (
    <div>
      <Heading>
        <Breadcrumb />
      </Heading>
      <ListingFilter
        types={[{ ...PageChoiceFilter, choice: locationChoice }]}
        filterFunction={listDisplay}
        onChange={pageChanged}
      />
    </div>
  );
};
