import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';
import {
  GetSettingsRequest,
  UpdateSettingsRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { useGetSettings, useUpdateSettings } from 'Api';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { NamespacesContext } from 'GlooIApp';
import { InputNumber } from 'antd';

interface Props {}

const UpdateRefreshRateText = styled.div`
  color: ${colors.seaBlue};
  text-decoration: underline;
  &:hover {
    color: ${colors.oceanBlue};
  }
`;
const RefreshRateTitle = styled.div`
  color: black;
  font-weight: 600;
`;

// TODO: style input element to better match spec
const StyledInput = styled(InputNumber)`
  width: 60px;
  height: 22px;
  border: none;
`;

interface RefreshRateProps {
  updateFn: () => void;
  seconds: number;
  nanos: number;
  setSeconds: React.Dispatch<React.SetStateAction<number>>;
  setNanos: React.Dispatch<React.SetStateAction<number>>;
}
const RefreshRate: React.FC<RefreshRateProps> = props => {
  const { updateFn, seconds, nanos, setSeconds, setNanos } = props;

  return (
    <div>
      <div style={{ display: 'flex', flexDirection: 'row' }}>
        <RefreshRateTitle>{`Refresh Rate:`}</RefreshRateTitle>
        <StyledInput
          size='small'
          style={{ width: '60px', border: 'none', height: '22px' }}
          value={seconds}
          defaultValue={seconds}
          onChange={seconds => setSeconds(seconds!)}
          formatter={seconds => `${seconds}s`}
          parser={seconds => seconds!.replace('s', '')}
        />
        :
        <StyledInput
          size='small'
          style={{ width: '60px', border: 'none', height: '22px' }}
          value={nanos}
          defaultValue={nanos}
          onChange={nanos => setNanos(nanos!)}
          formatter={nanos => `${nanos}ns`}
          parser={nanos => nanos!.replace('ns', '')}
        />
        <div onClick={updateFn} style={{ cursor: 'pointer' }}>
          <UpdateRefreshRateText>{'Change'}</UpdateRefreshRateText>
        </div>
      </div>
    </div>
  );
};

// TODO: consolidate update functions to avoid repetition
export const WatchedNamespacesPage = (props: Props) => {
  const allNamespaces = React.useContext(NamespacesContext);
  const [availableNS, setAvailableNS] = React.useState(
    allNamespaces.namespacesList
  );

  let req = new GetSettingsRequest();
  const { data, loading, error } = useGetSettings(req);

  const [refreshRateSeconds, setRefreshRateSeconds] = React.useState(() =>
    data ? data.settings!.refreshRate!.seconds : 1
  );
  const [refreshRateNanos, setRefreshRateNanos] = React.useState(() =>
    data ? data.settings!.refreshRate!.nanos : 0
  );

  const updateRequest = React.useRef(new UpdateSettingsRequest());
  const { refetch: makeRequest } = useUpdateSettings(null);

  const [watchedNamespacesList, setWatchedNamespacesList] = React.useState<
    string[]
  >([]);

  React.useEffect(() => {
    if (data && data.settings) {
      setWatchedNamespacesList(data.settings.watchNamespacesList);
      setRefreshRateSeconds(data.settings!.refreshRate!.seconds);
      setRefreshRateNanos(data.settings!.refreshRate!.nanos);
      const resourceRef = new ResourceRef();
      const { metadata } = data.settings!;
      resourceRef.setName(metadata!.name);
      resourceRef.setNamespace(metadata!.namespace);
      updateRequest.current.setRef(resourceRef);
    }
  }, [loading]);

  React.useEffect(() => {
    setAvailableNS(
      allNamespaces.namespacesList.filter(
        ns => !watchedNamespacesList.includes(ns)
      )
    );
  }, [watchedNamespacesList]);

  if (!data || loading) {
    return <div>Loading...</div>;
  }

  const addNamespace = (newNamespace: string) => {
    const newArray = [...watchedNamespacesList];
    newArray.push(newNamespace);
    setWatchedNamespacesList(newArray);

    updateRequest.current.setWatchNamespacesList(newArray);
    makeRequest(updateRequest.current);
  };

  const removeNamespace = (removeIndex: number) => {
    let newList = [...watchedNamespacesList];
    newList.splice(removeIndex, 1);
    setWatchedNamespacesList(newList);

    updateRequest.current.setWatchNamespacesList(newList);
    makeRequest(updateRequest.current);
  };

  const updateRefreshRate = () => {
    const duration = new Duration();
    duration.setSeconds(refreshRateSeconds);
    duration.setNanos(refreshRateNanos);

    updateRequest.current.setRefreshRate(duration);
    makeRequest(updateRequest.current);
  };

  return (
    <SectionCard
      cardName={'Watched Namespaces'}
      secondaryComponent={
        <RefreshRate
          seconds={refreshRateSeconds}
          nanos={refreshRateNanos}
          setSeconds={setRefreshRateSeconds}
          setNanos={setRefreshRateNanos}
          updateFn={updateRefreshRate}
        />
      }
      logoIcon={<RelatedCircles />}>
      <StringCardsList
        values={watchedNamespacesList}
        valueDeleted={removeNamespace}
        createNew={addNamespace}
        createNewPromptText={'watch new namespace'}
        asTypeahead
        presetOptions={availableNS}
      />
    </SectionCard>
  );
};
