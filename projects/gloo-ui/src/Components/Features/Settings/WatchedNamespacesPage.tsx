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
import { InputNumber } from 'antd';
import { useSelector, useDispatch } from 'react-redux';
import { AppState } from 'store';
import { getSettings, updateSettings } from 'store/config/actions';

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
  const dispatch = useDispatch();
  const { namespacesList, settings } = useSelector(
    (store: AppState) => store.config
  );

  const [availableNS, setAvailableNS] = React.useState(namespacesList);

  React.useEffect(() => {
    if (!settings) {
      dispatch(getSettings());
    }
  }, [settings]);

  const [refreshRateSeconds, setRefreshRateSeconds] = React.useState(() =>
    settings && settings.refreshRate ? settings!.refreshRate!.seconds : 1
  );
  const [refreshRateNanos, setRefreshRateNanos] = React.useState(() =>
    settings && settings.refreshRate ? settings!.refreshRate!.nanos : 0
  );

  const [watchedNamespacesList, setWatchedNamespacesList] = React.useState<
    string[]
  >(settings.watchNamespacesList);

  React.useEffect(() => {
    if (settings && settings.refreshRate) {
      const { watchNamespacesList, refreshRate } = settings!;
      setWatchedNamespacesList(watchNamespacesList);
      setRefreshRateSeconds(refreshRate!.seconds);
      setRefreshRateNanos(settings!.refreshRate!.nanos);
    }
  }, [settings]);

  React.useEffect(() => {
    setAvailableNS(
      namespacesList.filter(ns => !watchedNamespacesList.includes(ns))
    );
  }, [watchedNamespacesList]);

  if (!watchedNamespacesList) {
    return <div>Loading...</div>;
  }

  const addNamespace = (newNamespace: string) => {
    const newArray = [...watchedNamespacesList];
    newArray.push(newNamespace);
    setWatchedNamespacesList(newArray);
    dispatch(updateSettings({ ...settings, watchNamespacesList: newArray }));
  };

  const removeNamespace = (removeIndex: number) => {
    if (watchedNamespacesList.length > 1) {
      let newList = [...watchedNamespacesList];
      newList.splice(removeIndex, 1);
      // setWatchedNamespacesList(newList);
      dispatch(updateSettings({ ...settings, watchNamespacesList: newList }));
    } else {
      dispatch(updateSettings({ ...settings, watchNamespacesList: [] }));
    }

    setTimeout(() => dispatch(getSettings()), 300);
  };

  const updateRefreshRate = () => {
    const duration = new Duration();
    duration.setSeconds(refreshRateSeconds);
    duration.setNanos(refreshRateNanos);
    dispatch(
      updateSettings({
        ...settings,
        refreshRate: { seconds: refreshRateSeconds, nanos: refreshRateNanos }
      })
    );
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
        createNewPromptText={'watch namespace'}
        asTypeahead
        presetOptions={availableNS}
      />
    </SectionCard>
  );
};
