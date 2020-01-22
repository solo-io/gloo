import { css } from '@emotion/core';
import styled from '@emotion/styled';
import { InputNumber } from 'antd';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';
import { TallyContainer } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import {
  getSettings,
  updateRefreshRate,
  updateWatchNamespaces
} from 'store/config/actions';
import { configAPI } from 'store/config/api';
import { colors } from 'Styles';
import useSWR from 'swr';

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
  seconds: number;
  nanos: number;
}
const RefreshRate: React.FC<RefreshRateProps> = props => {
  const dispatch = useDispatch();
  const { seconds, nanos } = props;
  const [currentSeconds, setCurrentSeconds] = React.useState(seconds);
  const [currentNanos, setCurrentNanos] = React.useState(nanos);

  const handleUpdateRefreshRate = () => {
    dispatch(
      updateRefreshRate({
        refreshRate: { seconds: currentSeconds, nanos: currentNanos }
      })
    );
    setTimeout(() => dispatch(getSettings()), 300);
  };

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
          value={currentSeconds}
          defaultValue={currentSeconds}
          onChange={seconds => setCurrentSeconds(seconds!)}
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
          value={currentNanos}
          defaultValue={currentNanos}
          onChange={nanos => setCurrentNanos(nanos!)}
          formatter={nanos => `${nanos}ns`}
          parser={nanos => nanos!.replace('ns', '')}
          css={css`
            width: 60px;
            height: 22px;
            border: none;
          `}
        />
        <div onClick={handleUpdateRefreshRate} style={{ cursor: 'pointer' }}>
          <UpdateRefreshRateText>{'Change'}</UpdateRefreshRateText>
        </div>
      </div>
    </div>
  );
};

// TODO: consolidate update functions to avoid repetition
export const WatchedNamespacesPage = (props: Props) => {
  const dispatch = useDispatch();
  const { data: namespacesList, error: listNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  const { data: settings, error: settingsError } = useSWR(
    'getSettings',
    configAPI.getSettings
  );

  const currentRefreshRate = settings?.refreshRate;

  const watchNamespacesList = settings?.watchNamespacesList;

  const [availableNS, setAvailableNS] = React.useState(namespacesList);

  const [refreshRateSeconds, setRefreshRateSeconds] = React.useState(
    settings && settings.refreshRate ? settings.refreshRate!.seconds : 1
  );

  const [refreshRateNanos, setRefreshRateNanos] = React.useState(
    settings && settings.refreshRate ? settings.refreshRate!.nanos : 0
  );

  React.useEffect(() => {
    if (currentRefreshRate) {
      setRefreshRateSeconds(currentRefreshRate!.seconds);
      setRefreshRateNanos(currentRefreshRate!.nanos);
    }
  }, [currentRefreshRate]);

  React.useEffect(() => {
    if (!settings) {
      dispatch(getSettings());
    }
  }, [settings?.refreshRate, settings?.watchNamespacesList]);

  React.useEffect(() => {
    setAvailableNS(
      namespacesList?.filter(ns => !watchNamespacesList?.includes(ns))
    );
  }, [watchNamespacesList]);

  if (!watchNamespacesList) {
    return <div>Loading...</div>;
  }

  const addNamespace = (newNamespace: string) => {
    const newArray = [...watchNamespacesList];
    newArray.push(newNamespace);
    dispatch(updateWatchNamespaces({ watchNamespacesList: newArray }));
  };

  const removeNamespace = (removeIndex: number) => {
    if (watchNamespacesList.length > 1) {
      let newList = [...watchNamespacesList];
      newList.splice(removeIndex, 1);
      // setWatchedNamespacesList(newList);
      dispatch(updateWatchNamespaces({ watchNamespacesList: newList }));
    } else {
      dispatch(updateWatchNamespaces({ watchNamespacesList: [] }));
    }

    setTimeout(() => dispatch(getSettings()), 300);
  };

  return (
    <SectionCard
      cardName={'Watched Namespaces'}
      secondaryComponent={
        <RefreshRate seconds={refreshRateSeconds} nanos={refreshRateNanos} />
      }
      logoIcon={<RelatedCircles />}>
      {watchNamespacesList.length === 0 && (
        <TallyContainer
          color='blue'
          style={{ display: 'flex', justifyContent: 'center' }}>
          Currently watching all namespaces
        </TallyContainer>
      )}
      <StringCardsList
        values={watchNamespacesList}
        valueDeleted={removeNamespace}
        createNew={addNamespace}
        createNewPromptText={'watch namespace'}
        asTypeahead
        presetOptions={availableNS}
      />
    </SectionCard>
  );
};
