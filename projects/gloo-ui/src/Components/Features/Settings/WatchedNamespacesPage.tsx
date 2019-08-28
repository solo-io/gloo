import styled from '@emotion/styled';
import { InputNumber } from 'antd';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppState } from 'store';
import { getSettings, updateSettings } from 'store/config/actions';
import { colors } from 'Styles';
import { css } from '@emotion/core';

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
      <div
        css={css`
          display: flex;
          flex-direction: row;
        `}>
        <RefreshRateTitle>{`Refresh Rate:`}</RefreshRateTitle>
        <InputNumber
          size='small'
          value={seconds}
          defaultValue={seconds}
          onChange={seconds => setSeconds(seconds!)}
          formatter={seconds => `${seconds}s`}
          parser={seconds => seconds!.replace('s', '')}
          css={css`
            width: 60px;
            height: 22px;
            border: none;
          `}
        />
        :
        <InputNumber
          size='small'
          value={nanos}
          defaultValue={nanos}
          onChange={nanos => setNanos(nanos!)}
          formatter={nanos => `${nanos}ns`}
          parser={nanos => nanos!.replace('ns', '')}
          css={css`
            width: 60px;
            height: 22px;
            border: none;
          `}
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
