import * as React from 'react';
import { diff, Change, CriticalityLevel } from '@graphql-inspector/core';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { ReactComponent as ErrorIcon } from 'assets/big-unsuccessful-x.svg';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { SoloTable } from 'Components/Common/SoloTable';
import { makeExecutableSchema } from '@graphql-tools/schema';
import * as styles from './UpdateApiModal.style';
import upperFirst from 'lodash/upperFirst';

export interface ReactDiffProps {
  originalSchemaString: string;
  newSchemaString: string;
  warningMessage: string;
  setWarningMessage: (value: string) => any;
  validateApi: (schemaString: string) => Promise<any>;
}

const StyledContainer = styled.div<{ level: CriticalityLevel }>`
  padding-top: 20px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  color: ${props => {
    if (props.level === 'BREAKING') {
      return colors.errorRed;
    } else if (props.level === 'DANGEROUS') {
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
      render: (level: CriticalityLevel) => {
        const formattedLevel = upperFirst(level.toLowerCase()).replace(
          /_/g,
          ' '
        );
        return (
          <>
            <StyledContainer level={level}>
              {formattedLevel}
              {(level === 'BREAKING' || level === 'DANGEROUS') && (
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
      render: (reason: { level: CriticalityLevel; message: string }) => {
        return (
          <StyledContainer level={reason.level}>
            {reason.message}
          </StyledContainer>
        );
      },
    },
  ];
  // State
  const [changes, setChanges] = React.useState<Change[]>([]);

  React.useEffect(() => {
    try {
      const validatePromise = validateApi(newSchemaString)
        .then(() => {
          const original = makeExecutableSchema({
            typeDefs: originalSchemaString,
          });
          const newExec = makeExecutableSchema({
            typeDefs: newSchemaString,
          });
          return diff(original, newExec);
        })
        .then((values: Change[]) => {
          setWarningMessage('');
          const changes = values
            .map((v, idx) => {
              (v as any).level = v.criticality.level;
              (v as any).key = v.criticality.level + idx;
              (v as any).reason = {
                message: v.message,
                level: v.criticality.level,
              };
              return v;
            })
            .sort((a, b) => {
              // Sort by dangerous, then breaking, then non-breaking.

              const aLevel = a.criticality.level;
              const bLevel = b.criticality.level;
              if (aLevel === bLevel) {
                return 0;
              }
              if (aLevel === 'DANGEROUS') {
                return -1;
              } else if (bLevel === 'DANGEROUS') {
                return 1;
              } else if (aLevel === 'BREAKING') {
                return -1;
              } else if (bLevel === 'BREAKING') {
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
