import React from 'react';
import { TabPanelProps, Tab, Tabs, TabList, TabPanel } from '@reach/tabs';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';

type SelectableTabProps = {
  isSelected?: boolean;
};

const FolderTabHolder = styled.div<SelectableTabProps>`
  button {
    margin-right: 1px;
    padding: 10px 16px 6px;
    border-radius: 10px 10px 0 0;
    border: 1px solid ${colors.marchGrey};
    background: ${colors.februaryGrey};
    color: ${colors.septemberGrey};
    font-size: 18px;
    line-height: 26px;
    text-align: center;
    cursor: pointer;

    ${(props: SelectableTabProps) =>
      props.isSelected
        ? `
    border-bottom: 1px solid white;
    color: ${colors.seaBlue};
    cursor: default;
    z-index: 10;
    background: white;
    font-weight: 500;
    `
        : ``};

    &:focus {
      outline: none;
    }
  }
`;

export const FolderTab = (
  props: {
    disabled?: boolean | undefined;
    isSelected?: boolean | undefined;
  } & TabPanelProps
) => {
  const { children, isSelected, ...rest } = props;
  return (
    <FolderTabHolder isSelected={isSelected}>
      <Tab {...rest}>{children}</Tab>
    </FolderTabHolder>
  );
};

export const FolderTabContent = styled.div`
  padding: 23px 20px 25px;
  border-radius: 10px;
  border: 1px solid ${colors.marchGrey};
  margin-top: -1px;
`;

export const FolderTabList = styled(TabList)`
  display: flex;
  margin-left: 30px;
`;

export const StyledTabs = styled(Tabs)`
  display: grid;
  grid-template-columns: 200px 1fr;
`;

export const StyledTabPanel = styled(TabPanel)`
  &:focus {
    outline: none;
  }
`;
