import styled from '@emotion/styled';
import { InputNumber } from 'antd';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppState } from 'store';
import {
  getSettings,
  updateRefreshRate,
  updateWatchNamespaces
} from 'store/config/actions';
import { colors } from 'Styles';
import { css } from '@emotion/core';
import { TallyContainer } from 'Components/Common/DisplayOnly/TallyInformationDisplay';

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
  const { namespacesList } = useSelector((store: AppState) => store.config);

  const settings = useSelector((state: AppState) => state.config.settings);
  const currentRefreshRate = useSelector(
    (state: AppState) => state.config.settings.refreshRate
  );
  const watchNamespacesList = useSelector(
    (state: AppState) => state.config.settings.watchNamespacesList
  );

  const [availableNS, setAvailableNS] = React.useState(namespacesList);

  const [refreshRateSeconds, setRefreshRateSeconds] = React.useState(
    () => settings.refreshRate!.seconds || 1
  );
  const [refreshRateNanos, setRefreshRateNanos] = React.useState(
    () => settings.refreshRate!.seconds || 0
  );

  React.useEffect(() => {
    if (currentRefreshRate) {
      setRefreshRateSeconds(currentRefreshRate!.seconds);
      setRefreshRateNanos(currentRefreshRate!.nanos);
    }
  }, []);

  React.useEffect(() => {
    if (!settings) {
      dispatch(getSettings());
    }
  }, [settings.refreshRate, settings.watchNamespacesList]);

  React.useEffect(() => {
    setAvailableNS(
      namespacesList.filter(ns => !watchNamespacesList.includes(ns))
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

  const handleUpdateRefreshRate = () => {
    dispatch(
      updateRefreshRate({
        refreshRate: { seconds: refreshRateSeconds, nanos: refreshRateNanos }
      })
    );
    setTimeout(() => dispatch(getSettings()), 300);
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
          updateFn={handleUpdateRefreshRate}
        />
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
