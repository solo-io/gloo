import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { SoloTable, RenderStatus } from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as ClustersIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as ClusterKubernetesIcon } from 'assets/cluster-instance-icon.svg';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as SmallGreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as ArrowToggle } from 'assets/arrow-toggle.svg';
import { colors } from 'Styles/colors';
import { useListClusterDetails } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { SoloLinkLooks } from 'Components/Common/SoloLink';
import { RegisterClusterModal } from './RegisterClusterModal';
import { DataError } from 'Components/Common/DataError';
import { sortGlooInstances } from 'utils/gloo-instance-helpers';

const ClustersContainer = styled.div`
  position: relative;
`;

const EmptyClustersArea = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  height: 220px;

  svg {
    margin-right: 25px;
  }
`;

const EmptyTitle = styled.div`
  font-weight: 500;
  font-size: 20px;
  line-height: 26px;
`;

const EmptyDescription = styled.div`
  color: ${colors.septemberGrey};
  font-size: 18px;
  line-height: 26px;
`;

const RegisterCTAHolder = styled.div`
  position: absolute;
  right: 0;
  top: -36px;

  display: flex;
  align-items: center;
  font-weight: 500;
  cursor: pointer;

  svg {
    width: 22px;
    fill: ${colors.forestGreen};
    margin-right: 8px;
  }
`;

type ClusterRowColorerProps = { clusterRows: number[] };
const ClusterRowColorer = styled.div<ClusterRowColorerProps>`
  ${(props: ClusterRowColorerProps) => `
    table {
      thead tr {
        background: white !important;
      }

      tbody {
        ${props.clusterRows.map(
          rowInd => `
          tr:nth-of-type(${rowInd}) {
            background: ${colors.februaryGrey};
            z-index: 2;
          }`
        )}
      }
    }
  `}
`;

type RealClusterNameProps = {
  isOpen?: boolean;
};
const RealClusterName = styled.div`
  position: relative;
  display: flex;
  align-items: center;
  font-weight: 500;
  color: ${colors.novemberGrey};

  padding-left: 15px;
  ${(props: RealClusterNameProps) =>
    props.isOpen === undefined ? '' : 'cursor: pointer;'}

  svg {
    width: 26px;
    height: 26px;
    padding: 3px;
    margin-right: 5px;

    background: white;
    border-radius: 100%;
  }

  svg.arrow {
    position: absolute;
    left: -5px;
    top: 9px;
    width: 10px;
    height: auto;
    padding: 0;
    background: transparent;

    ${(props: RealClusterNameProps) =>
      props.isOpen ? `transform: rotate(0deg);` : `transform: rotate(180deg);`}
    transition: transform .3s;

    * {
      fill: ${colors.septemberGrey};
    }
  }
`;

const GlooInstanceName = styled.div`
  position: relative;
  display: flex;
  padding-left: 20px;

  &:before {
    content: '';
    position: absolute;
    border-left: 1px dotted ${colors.marchGrey};
    border-bottom: 1px dotted ${colors.marchGrey};
    left: 8px;
    top: -20px;
    width: 10px;
    height: 29px;
  }
`;
const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;
  background: ${colors.februaryGrey};
  min-width: 20px;
  width: 20px;
  height: 20px;
  border-radius: 100%;
  margin-right: 5px;

  svg {
    margin-left: 5px;
    width: 12px;
  }
`;

type ClusterTableFields = {
  key: string;
  name: {
    naming: SimpleLinkProps | string;
    icon?: React.ReactNode;
    rowInd?: number;
  };
  namespace?: string;
  version?: string;
  proxyReplicas?: number;
  virtualServices?: number;
};

export const Clusters = () => {
  const [registerModalOpen, setRegisterModalOpen] = useState(false);
  const [tableData, setTableData] = useState<ClusterTableFields[]>([]);
  const [clusterRows, setClusterRows] = useState<number[]>([]); //Tracks which rows of table data are clusters
  const [showingClusterRows, setShowingClusterRows] = useState<number[]>();

  const { data: clusterDetails, error: cError } = useListClusterDetails();

  useEffect(() => {
    if (clusterDetails) {
      let newTableData: ClusterTableFields[] = [];
      let newClusterRows: number[] = [];
      let newShowingClusterRows: number[] = [];

      clusterDetails
        .sort((cDetsA, cDetsB) => cDetsA.cluster.localeCompare(cDetsB.cluster))
        .forEach((clusterDetail, ind) => {
          newTableData.push({
            key: clusterDetail.cluster,
            name: {
              naming: clusterDetail.cluster,
              icon: undefined,
              rowInd: ind,
            },
          });
          newClusterRows.push(newTableData.length);

          if (!showingClusterRows || showingClusterRows.includes(ind)) {
            newShowingClusterRows.push(ind);

            clusterDetail.glooInstancesList
              .sort((giA, giB) => sortGlooInstances(giA, giB))
              .forEach(glooInstance => {
                newTableData.push({
                  key: glooInstance.metadata?.uid ?? '' + clusterDetail.cluster,
                  name: {
                    naming: !!glooInstance.metadata
                      ? {
                          displayElement: glooInstance.metadata?.name ?? '',
                          link: glooInstance.metadata
                            ? `/gloo-instances/${glooInstance.metadata.namespace}/${glooInstance.metadata.name}/`
                            : '',
                        }
                      : 'No metadata - ohno',
                    icon: <GlooIcon />,
                  },
                  namespace: glooInstance.metadata?.namespace ?? '',
                  version: glooInstance.spec?.controlPlane?.version ?? '',
                  proxyReplicas: glooInstance.spec?.check?.proxies?.total ?? 0,
                  virtualServices:
                    glooInstance.spec?.check?.virtualServices?.total ?? 0,
                });
              });
          }
        });

      setTableData(newTableData);
      setClusterRows(newClusterRows);
      setShowingClusterRows(newShowingClusterRows);
    } else {
      setTableData([]);
    }
  }, [clusterDetails, showingClusterRows?.length]);

  if (!!cError) {
    return <DataError error={cError} />;
  } else if (!clusterDetails) {
    return <Loading message={`Retrieving cluster information...`} />;
  }

  const renderClusterName = (data: {
    naming: SimpleLinkProps | string;
    icon?: React.ReactNode;
    rowInd?: number;
  }) => {
    if (typeof data.naming === 'string') {
      return (
        <RealClusterName
          onClick={() =>
            setShowingClusterRows(oldRows => {
              console.log(oldRows, data);
              if (!oldRows) {
                return [data.rowInd!];
              } else if (oldRows.includes(data.rowInd!)) {
                return oldRows.filter(rInd => rInd !== data.rowInd);
              }

              return [...oldRows, data.rowInd!];
            })
          }
          isOpen={
            showingClusterRows && showingClusterRows.includes(data.rowInd!)
          }>
          <ArrowToggle className='arrow' />
          <ClusterKubernetesIcon />
          {data.naming}
        </RealClusterName>
      );
    } else {
      return (
        <GlooInstanceName>
          <GlooIconHolder>
            <GlooIcon />
          </GlooIconHolder>

          {RenderSimpleLink(data.naming)}
        </GlooInstanceName>
      );
    }
  };

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 275,
      render: renderClusterName,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Version',
      dataIndex: 'version',
    },
    {
      title: 'Proxy Replicas',
      dataIndex: 'proxyReplicas',
    },
    {
      title: 'Virtual Services',
      dataIndex: 'virtualServices',
    },
  ];

  return (
    <ClustersContainer>
      {!tableData.length ? (
        <SectionCard cardName=''>
          <EmptyClustersArea>
            <IconHolder width={90}>
              <ClustersIcon />
            </IconHolder>
            <div>
              <EmptyTitle>
                There are no Registered Clusters to display.
              </EmptyTitle>
              <EmptyDescription>
                <SoloLinkLooks
                  displayInline={true}
                  onClick={() => setRegisterModalOpen(true)}>
                  Register a Cluster
                </SoloLinkLooks>{' '}
                amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor
              </EmptyDescription>
            </div>
          </EmptyClustersArea>
        </SectionCard>
      ) : (
        <>
          <RegisterCTAHolder onClick={() => setRegisterModalOpen(true)}>
            <SmallGreenPlus /> Register a Cluster
          </RegisterCTAHolder>
          <SectionCard
            cardName={'Clusters'}
            logoIcon={
              <IconHolder width={20}>
                <ClustersIcon />
              </IconHolder>
            }
            noPadding={true}>
            <ClusterRowColorer clusterRows={clusterRows}>
              <SoloTable
                columns={columns}
                dataSource={tableData}
                removePaging
                removeShadows
                curved={true}
              />
            </ClusterRowColorer>
          </SectionCard>
        </>
      )}
      <RegisterClusterModal
        modalOpen={registerModalOpen}
        onClose={() => setRegisterModalOpen(false)}
      />
    </ClustersContainer>
  );
};
