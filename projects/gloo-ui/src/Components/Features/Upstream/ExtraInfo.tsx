import React from 'react';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import styled from '@emotion/styled/macro';
import { List, Button } from 'antd';

const ExtraInfoContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  padding-right: 10px;
`;

const FunctionContainer = styled.div`
  max-height: 250px;
  overflow-y: scroll;
  padding: 5px;
`;
interface Props {
  upstream: Upstream.AsObject;
}

export function ExtraInfo(props: Props) {
  const [showModal, setShowModal] = React.useState(true);
  const { upstreamSpec } = props.upstream;
  const functions =
    (upstreamSpec &&
      upstreamSpec.aws &&
      upstreamSpec.aws.lambdaFunctionsList) ||
    [];

  return (
    <div>
      <ExtraInfoContainer>
        <Button
          type='link'
          style={{ cursor: 'pointer' }}
          onClick={() => setShowModal(s => !s)}
          disabled={functions.length === 0}>
          {`${showModal ? 'Hide' : 'Show'} Functions`}
        </Button>
      </ExtraInfoContainer>
      <FunctionContainer>
        {showModal && (
          <List
            size='small'
            bordered
            dataSource={functions}
            locale={{ emptyText: 'No Functions' }}
            renderItem={item => (
              <List.Item style={{ padding: '0 5px' }}>
                <List.Item.Meta title={item.logicalName} />
              </List.Item>
            )}
          />
        )}
      </FunctionContainer>
    </div>
  );
}
