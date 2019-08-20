import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, TableActionCircle } from 'Styles';
import { ReactComponent as DocumentSVG } from 'assets/document.svg';
import { ReactComponent as TableDownloadIcon } from 'assets/download.svg';

const Link = styled.div`
  display: flex;
  align-items: center;
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

const DocumentIcon = styled(DocumentSVG)`
  margin-right: 5px;
`;

const ActionCircle = styled<'div', { gridArea?: string }>('div')`
  ${props => !!props.gridArea && `grid-area: ${props.gridArea}`};
`;

interface Props {
  fileContent: string;
  fileName: string;
}

export const FileDownloadLink = (props: Props) => {
  const downloadYaml = (): void => {
    const templElement = document.createElement('a');
    const file = new Blob([props.fileContent], { type: 'text/plain' });
    templElement.href = URL.createObjectURL(file);
    templElement.download = props.fileName;
    document.body.appendChild(templElement); // Required for this to work in FireFox
    templElement.click();
  };

  return (
    <Link onClick={() => downloadYaml()}>
      <DocumentIcon /> {props.fileName}
    </Link>
  );
};

interface CircleProps extends Props {
  gridArea?: string;
}

export const FileDownloadActionCircle = (props: CircleProps) => {
  const downloadYaml = (): void => {
    const templElement = document.createElement('a');
    const file = new Blob([props.fileContent], { type: 'text/plain' });
    templElement.href = URL.createObjectURL(file);
    templElement.download = props.fileName;
    document.body.appendChild(templElement); // Required for this to work in FireFox
    templElement.click();
  };

  return (
    <ActionCircle gridArea={props.gridArea}>
      <TableActionCircle onClick={() => downloadYaml()}>
        <TableDownloadIcon />
      </TableActionCircle>
    </ActionCircle>
  );
};
