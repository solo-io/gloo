import { css } from '@emotion/core';
import { InputNumber, Tooltip } from 'antd';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloTable } from 'Components/Common/SoloTable';
import { Formik, useFormikContext } from 'formik';
import { WeightedDestination } from 'proto/gloo/projects/gloo/api/v1/proxy_pb';
import React, { useState } from 'react';
import { useParams } from 'react-router';
import { upstreamGroupAPI } from 'store/upstreamGroups/api';
import { colors } from 'Styles';
import useSWR, { mutate } from 'swr';
import { ConfigurationToggle } from '../VirtualService/Details/VirtualServiceDetails';

export const WeightInput = ({ ...props }) => {
  let shadow = `box-shadow: 0px 0px 3px ${colors.lakeBlue};
  border: 1px solid ${colors.seaBlue}; color: ${colors.septemberGrey}`;

  const form = useFormikContext<{
    upstreams: WeightedDestination.AsObject[];
  }>();
  const field = form.getFieldProps(props.name);
  let totalWeight = 0;

  form.values.upstreams.forEach(weightedDest => {
    totalWeight += weightedDest.weight;
  });

  if (totalWeight > 100 || totalWeight < 100) {
    shadow = `box-shadow: 0px 0px 3px ${colors.grapefruitOrange};
  border: 1px solid ${colors.grapefruitOrange}; color: ${colors.pumpkinOrange}`;
  }
  if (totalWeight === 100) {
    shadow = `box-shadow: 0px 0px 3px #0DCE93;
  border: 1px solid #0DCE93; color: #0F9D72`;
  }

  return (
    <>
      <InputNumber
        css={css`
          width: 70px;
          ${shadow}
        `}
        defaultValue={field.value}
        value={field.value}
        min={0}
        max={100}
        step={5}
        formatter={value => `${value}%`}
        parser={value => value!.replace('%', '')}
        onChange={value => form.setFieldValue(field.name, value)}
      />
    </>
  );
};

export function getWeightColor(weight: number) {
  let color = colors.pondBlue;
  if (weight <= 10) color = colors.pondBlue;
  if (weight > 10 && weight <= 30) color = '#40A7D5';
  if (weight > 30 && weight <= 50) color = colors.oceanBlue;
  if (weight > 50 && weight <= 80) color = '#07638D';
  if (weight > 80 && weight <= 100) color = '#253E58';

  return color;
}

export const UpstreamGroupDetails = () => {
  const [editMode, setEditMode] = useState(false);
  const [showConfig, setShowConfig] = useState(false);
  let { upstreamgroupname, upstreamgroupnamespace } = useParams();
  const { data: upstreamGroup, error } = useSWR(
    !!upstreamgroupname && !!upstreamgroupnamespace
      ? ['getUpstreamGroup', upstreamgroupname, upstreamgroupnamespace]
      : null,
    (key: string, upstreamgroupname: string, upstreamgroupnamespace: string) =>
      upstreamGroupAPI.getUpstreamGroup({
        ref: {
          name: upstreamgroupname!,
          namespace: upstreamgroupnamespace!
        }
      })
  );

  if (error) {
    return <div>error</div>;
  }
  if (!upstreamGroup) {
    return <div>Loading...</div>;
  }

  function handleSaveYamlEdit(editedYaml: string) {
    mutate(
      ['getUpstreamGroup', upstreamgroupname, upstreamgroupnamespace],
      upstreamGroupAPI.updateUpstreamGroupYaml({
        editedYamlData: {
          ref: { name: upstreamgroupname!, namespace: upstreamgroupnamespace! },
          editedYaml: editedYaml
        }
      })
    );
  }

  const handleUpdateWeight = (values: {
    upstreams: WeightedDestination.AsObject[];
  }) => {
    mutate(
      ['getUpstreamGroup', upstreamgroupname, upstreamgroupnamespace],
      upstreamGroupAPI.updateUpstreamGroup({
        ...upstreamGroup,
        upstreamGroup: {
          ...upstreamGroup!.upstreamGroup,
          destinationsList: values.upstreams,
          metadata: {
            ...upstreamGroup?.upstreamGroup?.metadata!
          }
        }
      })
    );

    setEditMode(false);
  };

  if (upstreamGroup?.upstreamGroup?.destinationsList === undefined) {
    return <div>Loading</div>;
  }

  return (
    <>
      <Breadcrumb />
      <SectionCard
        cardName={upstreamgroupname!}
        logoIcon={<UpstreamGroupIcon />}
        health={upstreamGroup?.upstreamGroup?.status?.state}
        headerSecondaryInformation={[
          { title: 'namespace', value: upstreamgroupnamespace! }
        ]}>
        <div
          css={css`
            display: flex;
            flex-direction: row;
            justify-content: flex-end;
          `}>
          <ConfigurationToggle onClick={() => setShowConfig(s => !s)}>
            {showConfig ? `Hide ` : `View `}
            YAML Configuration
          </ConfigurationToggle>

          <FileDownloadLink
            fileContent={upstreamGroup?.raw?.content!}
            fileName={upstreamGroup?.raw?.fileName!}
          />
        </div>
        <div>
          {showConfig && (
            <ConfigDisplayer
              content={upstreamGroup?.raw?.content || ''}
              asEditor
              yamlError={undefined}
              saveEdits={handleSaveYamlEdit}
            />
          )}
        </div>
        <p
          css={css`
            font-size: 18px;
            color: ${colors.novemberGrey};
          `}>
          Upstreams
        </p>
        <Formik
          initialValues={{
            upstreams: upstreamGroup?.upstreamGroup?.destinationsList || []
          }}
          onSubmit={handleUpdateWeight}>
          {({ values, handleSubmit, errors }) => {
            return (
              <>
                <SoloTable
                  dataSource={
                    upstreamGroup?.upstreamGroup?.destinationsList || []
                  }
                  columns={[
                    {
                      title: 'Weight',
                      dataIndex: 'weight',
                      render: (
                        weight: number,
                        resource: any,
                        index: number
                      ) => {
                        return (
                          <>
                            {editMode ? (
                              <>
                                <WeightInput
                                  name={`upstreams[${index}].weight`}
                                />
                              </>
                            ) : (
                              <Tooltip
                                placement='right'
                                title={'Click to edit'}>
                                <div
                                  onClick={() => setEditMode(true)}
                                  css={css`
                                    text-align: center;
                                    padding: 2px 3px;
                                    cursor: pointer;
                                    width: 60px;
                                    color: white;
                                    border-radius: 5px;
                                    background: ${getWeightColor(weight)};
                                  `}>
                                  {weight}%
                                </div>
                              </Tooltip>
                            )}
                          </>
                        );
                      }
                    },
                    { title: 'Name', dataIndex: 'destination.upstream.name' },
                    {
                      title: 'Namespace',
                      dataIndex: 'destination.upstream.namespace'
                    }
                  ]}
                />
                <div
                  css={css`
                    display: grid;
                    grid-template-columns: 150px 150px;
                    grid-gap: 15px;
                    margin-top: 10px;
                  `}>
                  <SoloButton onClick={handleSubmit} text='Update Settings' />
                  <SoloButton
                    css={css`
                      background-color: ${colors.juneGrey};
                      &:disabled {
                        background: ${colors.juneGrey};
                      }
                      &:hover {
                        background: ${colors.juneGrey};
                      }
                      &:focus,
                      &:active {
                        background: ${colors.juneGrey};
                      }
                    `}
                    onClick={() => setEditMode(false)}
                    text='Cancel'
                    disabled={!editMode}
                  />
                </div>
              </>
            );
          }}
        </Formik>
      </SectionCard>
    </>
  );
};
