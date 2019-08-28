import styled from '@emotion/styled';
import { List } from 'antd';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import React from 'react';
import { colors } from 'Styles';
import { getFunctionList } from 'utils/helpers';

const ExtraInfoContainer = styled.div`
  margin-top: -18px;
`;

const ToggleContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  padding-right: 10px;
`;

const ShowToggle = styled.div`
  color: ${colors.seaBlue};
  font-size: 14px;
  line-height: 30px;
  height: 30px;
  cursor: pointer;
`;

const FunctionsContainer = styled.div`
  position: relative;
  max-height: 200px;
  overflow-y: scroll;
  padding-left: 20px;
  padding-right: 5px;
  background: ${colors.januaryGrey};
`;
const ScrollMirror = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 15px;
  height: 100%;
  border-right: 1px solid ${colors.scrollbarBorderGrey};
  background: ${colors.scrollbarBackgroundGrey};
`;

const ListBlock = styled.div`
  margin: 8px 0 5px;

  .ant-list {
    background: white;
  }
`;

interface Props {
  upstream: Upstream.AsObject;
}

export function ExtraInfo(props: Props) {
  const [showModal, setShowModal] = React.useState(true);

  // TODO: process upstream spec to support all types
  const [functionsList, setFunctionsList] = React.useState<
    { key: string; value: string }[]
  >(() => getFunctionList(props.upstream.upstreamSpec!));

  return (
    <ExtraInfoContainer>
      <ToggleContainer>
        <ShowToggle
          onClick={() => {
            console.log(functionsList);
            if (functionsList.length !== 0) {
              setShowModal(s => !s);
            }
          }}>
          {`${showModal ? 'Hide' : 'Show'} Functions`}
        </ShowToggle>
      </ToggleContainer>
      {showModal && (
        <FunctionsContainer>
          <ScrollMirror />
          <ListBlock>
            <List
              size='small'
              bordered
              dataSource={functionsList}
              locale={{
                emptyText: 'No Functions'
              }}
              renderItem={item => (
                <List.Item
                  style={{
                    padding: '0 5px'
                  }}>
                  <List.Item.Meta title={item.value} />
                </List.Item>
              )}
            />
          </ListBlock>
        </FunctionsContainer>
      )}
    </ExtraInfoContainer>
  );
}
