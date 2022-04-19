import * as React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { ReactComponent as ErrorIcon } from 'assets/big-unsuccessful-x.svg';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { SoloTable } from 'Components/Common/SoloTable';
import * as styles from './UpdateApiModal.style';
import { graphqlConfigApi } from 'API/graphql';
import { GraphQLInspectorDiffOutput } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/diff_pb';

export interface ReactDiffProps {
  originalSchemaString: string;
  newSchemaString: string;
  warningMessage: string;
  setWarningMessage: (value: string) => any;
  validateApi: (schemaString: string) => Promise<any>;
}

const StyledContainer = styled.div<{ level: number }>`
  padding-top: 20px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  color: ${props => {
    if (props.level === GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING) {
      return colors.errorRed;
    } else if (
      props.level === GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS
    ) {
      return colors.errorRed;
    }
    return 'black';
  }};
`;

const StyledIcon = styled(IconHolder)`
  & svg * {
    stroke: white;
    stroke-width: 8px;
  }
`;

const StyledTable = styled(SoloTable)`
  max-height: 350px;
  overflow-y: scroll;
`;

export const GraphqlDiff = (props: ReactDiffProps) => {
  const {
    originalSchemaString,
    newSchemaString,
    setWarningMessage,
    warningMessage,
    validateApi,
  } = props;

  let columns: any = [
    {
      title: 'Type',
      dataIndex: 'level',
      render: (level: number) => {
        let formattedLevel = '';
        if (
          level === GraphQLInspectorDiffOutput.CriticalityLevel.NON_BREAKING
        ) {
          formattedLevel = 'Non breaking';
        } else if (
          level === GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS
        ) {
          formattedLevel = 'Dangerous';
        } else if (
          level === GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING
        ) {
          formattedLevel = 'Breaking';
        }
        return (
          <>
            <StyledContainer level={level}>
              {formattedLevel}
              {(level ===
                GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS ||
                level ===
                  GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING) && (
                <StyledIcon
                  applyColor={{
                    color: colors.errorRed,
                  }}>
                  <ErrorIcon />
                </StyledIcon>
              )}
            </StyledContainer>
          </>
        );
      },
    },
    {
      title: 'Reason',
      dataIndex: 'reason',
      render: (reason: { level: number; message: string }) => {
        return (
          <StyledContainer level={reason.level}>
            {reason.message}
          </StyledContainer>
        );
      },
    },
  ];
  // State
  const [changes, setChanges] = React.useState<
    GraphQLInspectorDiffOutput.Change.AsObject[]
  >([]);

  React.useEffect(() => {
    try {
      validateApi(newSchemaString)
        .then(() => {
          return graphqlConfigApi.getSchemaDiff(
            originalSchemaString,
            newSchemaString
          );
        })
        .then((values: GraphQLInspectorDiffOutput.Change.AsObject[]) => {
          setWarningMessage('');
          const changes = values
            .map((v, idx) => {
              (v as any).level = v.criticality!.level;
              (v as any).key = v.criticality!.level + idx;
              (v as any).reason = {
                message: v.message,
                level: v.criticality!.level,
              };
              return v;
            })
            .sort((a, b) => {
              // Sort by dangerous, then breaking, then non-breaking.

              const aLevel = a.criticality!.level;
              const bLevel = b.criticality!.level;
              if (aLevel === bLevel) {
                return 0;
              }

              if (
                aLevel === GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS
              ) {
                return -1;
              } else if (
                bLevel === GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS
              ) {
                return 1;
              } else if (
                aLevel === GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING
              ) {
                return -1;
              } else if (
                bLevel === GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING
              ) {
                return 1;
              }
              return 0;
            });
          setChanges(changes);
        })
        .catch(err => {
          console.error(err);
          throw err;
        });
    } catch (err: any) {
      console.error(err);
      if (err?.message) {
        setWarningMessage(err.message);
      } else if (typeof err === 'string') {
        setWarningMessage(err);
      }
    }
  }, [originalSchemaString, newSchemaString]);

  return (
    <div>
      {Boolean(warningMessage) && (
        <styles.StyledWarning className='p-2 text-orange-400 border border-orange-400 mb-5'>
          {warningMessage}
        </styles.StyledWarning>
      )}
      <StyledTable
        columns={columns}
        dataSource={changes}
        removePaging
        removeShadows
        curved={true}
      />
    </div>
  );
};
