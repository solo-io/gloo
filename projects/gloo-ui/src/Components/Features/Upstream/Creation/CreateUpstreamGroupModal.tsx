import { css } from '@emotion/core';
import styled from '@emotion/styled';
import { Dialog } from '@reach/dialog';
import '@reach/dialog/styles.css';
import { Divider, Transfer } from 'antd';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloTable } from 'Components/Common/SoloTable';
import { Formik } from 'formik';
import { UpstreamGroup } from 'proto/gloo/projects/gloo/api/v1/proxy_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { upstreamAPI } from 'store/upstreams/api';
import { colors } from 'Styles';
import useSWR, { mutate, trigger } from 'swr';
import { getIconFromSpec } from 'utils/helpers';
import * as yup from 'yup';
import { WeightInput } from '../UpstreamGroupDetails';
import { upstreamGroupAPI } from 'store/upstreamGroups/api';
interface Props {}
interface Values {
  name: string;
  namespace: string;
  upstreams: { name: string; namespace: string; weight: number }[];
}

const StyledGreenPlus = styled(GreenPlus)`
  cursor: pointer;
  margin-right: 7px;
  .a {
    fill: ${colors.forestGreen};
  }
`;

const ModalContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-content: center;
`;
const Legend = styled.div`
  background-color: ${colors.januaryGrey};
  padding: 13px 13px 13px 10px;
  margin-bottom: 23px;
`;

const ModalTrigger = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 14px;
`;

const validationSchema = yup.object().shape({
  name: yup.string().required(),
  namespace: yup.string().min(2),
  upstreams: yup.array().of(
    yup.object().shape({
      weight: yup.number().max(100)
    })
  )
});

export const CreateUpstreamGroupModal = (props: Props) => {
  const dispatch = useDispatch();
  const [currentView, setCurrentView] = React.useState<'general' | 'weight'>(
    'general'
  );
  const [showModal, setShowModal] = React.useState(false);
  const [selectedUpstreams, setSelectedUpstreams] = React.useState<
    { name: string; namespace: string; weight: number }[]
  >([]);
  const [targetKeys, setTargetKeys] = React.useState<string[]>([]);
  const { data: upstreamsList, error } = useSWR(
    'listUpstreams',
    upstreamAPI.listUpstreams
  );
  if (!upstreamsList) {
    return <div>Loading...</div>;
  }
  const open = () => setShowModal(true);
  const close = () => setShowModal(false);

  function handleCreateUpstreamGroup(values: Values) {
    const { name, namespace, upstreams } = values;
    const newUpstreamGroup = new UpstreamGroup().toObject();
    mutate(
      !!name && !!namespace ? ['getUpstreamGroup', name, namespace] : null,
      upstreamGroupAPI.createUpstreamGroup({
        upstreamGroup: {
          ...newUpstreamGroup,
          metadata: {
            ...newUpstreamGroup.metadata!,
            name,
            namespace
          },
          destinationsList: values.upstreams.map(us => {
            return {
              weight: us.weight,
              destination: {
                upstream: {
                  name: us.name!,
                  namespace: us.namespace!
                }
              }
            };
          })
        }
      })
    );
    setShowModal(false);
  }
  function formatUpstreamData() {
    return upstreamsList?.map(upstreamDetail => {
      return {
        ...upstreamDetail,
        key: `${upstreamDetail?.upstream?.metadata?.name}::${upstreamDetail?.upstream?.metadata?.namespace}`,
        title: `${upstreamDetail?.upstream?.metadata?.name!}`,
        upstream: upstreamDetail.upstream!
      };
    });
  }

  return (
    <ModalContainer data-testid='create-upstream-modal'>
      <ModalTrigger onClick={open}>
        <>
          <StyledGreenPlus />
          Create Upstream Group
        </>
        <Divider type='vertical' style={{ height: '1.5em' }} />
      </ModalTrigger>
      <Dialog
        aria-label='Create a new Upstream Group'
        css={css`
          width: 750px;
          min-height: 520px;
          border-radius: 10px;
          display: grid;
          padding: 0;
          grid-template-columns: 160px 1fr;
        `}
        isOpen={showModal}
        onDismiss={close}>
        <Formik
          initialValues={{
            name: '',
            namespace: 'gloo-system',
            upstreams: selectedUpstreams
          }}
          onSubmit={handleCreateUpstreamGroup}
          validationSchema={validationSchema}>
          {({ values, handleSubmit, errors, setFieldValue }) => (
            <>
              <div
                css={css`
                  background: ${colors.seaBlue};
                  border-radius: 10px 0 0 10px;
                  font-size: 18px;
                  padding: 10px 0 10px 10px;
                  color: white;
                `}>
                <div
                  onClick={() => setCurrentView('general')}
                  css={css`
                    background-color: ${currentView === 'general'
                      ? 'hsl(199, 65%, 53%)'
                      : ''};
                    border-right: ${currentView === 'general'
                      ? `5px solid ${colors.pondBlue}`
                      : ''};
                  `}>
                  General
                </div>
                <div
                  onClick={() => setCurrentView('weight')}
                  css={css`
                    background-color: ${currentView === 'weight'
                      ? 'hsl(199, 65%, 53%)'
                      : ''};

                    border-right: ${currentView === 'weight'
                      ? `5px solid ${colors.pondBlue}`
                      : ''};
                  `}>
                  Weight
                </div>
              </div>
              <div
                css={css`
                  padding: 25px;
                  display: flex;
                  flex-direction: column;
                  justify-content: space-between;
                `}>
                {currentView === 'general' ? (
                  <>
                    <div
                      css={css`
                        font-size: 22px;
                        color: ${colors.novemberGrey};
                      `}>
                      New Upstream Group
                    </div>
                    <div>
                      <SoloFormInput title='Name' name='name' />
                    </div>
                    <div
                      css={css`
                        align-self: center;
                      `}>
                      <Transfer
                        css={css`
                          .ant-transfer-list {
                            width: 250px;
                            height: 300px;
                          }
                        `}
                        dataSource={formatUpstreamData()}
                        showSearch
                        showSelectAll={false}
                        locale={{
                          itemUnit: 'upstream',
                          itemsUnit: 'upstreams'
                        }}
                        targetKeys={targetKeys}
                        onChange={(
                          targetKeys: string[],
                          direction: string,
                          moveKeys: string[]
                        ) => {
                          setTargetKeys(targetKeys);

                          let selected = upstreamsList
                            ?.filter(usD =>
                              targetKeys.includes(
                                `${usD.upstream?.metadata?.name}::${usD.upstream?.metadata?.namespace}`
                              )
                            )
                            .map(usD => {
                              return {
                                name: usD?.upstream?.metadata?.name!,
                                namespace: usD?.upstream?.metadata?.namespace!,
                                weight: Math.floor(100 / targetKeys.length)
                              };
                            });
                          //account for odd numbers
                          if (
                            Math.floor(100 / targetKeys.length) *
                              targetKeys.length !==
                              100 &&
                            targetKeys.length > 1
                          ) {
                            selected[0].weight +=
                              100 -
                              Math.floor(100 / targetKeys.length) *
                                targetKeys.length;
                          }
                          setSelectedUpstreams(selected!);
                          setFieldValue('upstreams', selected);
                        }}
                        render={item => (
                          <>
                            {getIconFromSpec(item.upstream)}
                            {item.title!}
                          </>
                        )}
                      />
                    </div>
                  </>
                ) : (
                  <>
                    <div
                      css={css`
                        font-size: 22px;
                        color: ${colors.novemberGrey};
                      `}>
                      Upstream Weights
                    </div>
                    <div
                      css={css`
                        padding: 5px;
                        background-color: ${colors.januaryGrey};
                      `}>
                      Weights specify the amount of traffic sent to each of your
                      upstreams in the group. Total weight cannot exceed 100%.
                    </div>
                    <div>Upstreams in {values.name || 'group'}</div>

                    <SoloTable
                      dataSource={selectedUpstreams || []}
                      columns={[
                        {
                          title: 'Weight',
                          dataIndex: 'weight',
                          render: (weight: number, record: any, index: any) => {
                            return (
                              <>
                                <WeightInput
                                  step={1}
                                  name={`upstreams[${index}].weight`}
                                />
                              </>
                            );
                          }
                        },
                        { title: 'Name', dataIndex: 'name' },
                        {
                          title: 'Namespace',
                          dataIndex: 'namespace'
                        }
                      ]}
                    />
                  </>
                )}

                <div
                  css={css`
                    display: flex;
                    justify-content: space-between;
                  `}>
                  <div
                    onClick={close}
                    css={css`
                      color: ${colors.seaBlue};
                      cursor: pointer;
                    `}>
                    Cancel
                  </div>

                  {currentView === 'weight' && (
                    <SoloButton
                      text={'Back'}
                      onClick={() => setCurrentView('general')}
                    />
                  )}
                  {currentView === 'general' && (
                    <SoloButton
                      text={'Next Step'}
                      disabled={errors.name !== undefined}
                      onClick={() => setCurrentView('weight')}
                    />
                  )}
                  {currentView === 'weight' && (
                    <SoloButton text={'Create Group'} onClick={handleSubmit} />
                  )}
                </div>
              </div>
            </>
          )}
        </Formik>
      </Dialog>
    </ModalContainer>
  );
};
