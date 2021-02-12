import React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { useNavigate } from 'react-router';
import { Link } from 'react-router-dom';

const SoloLinkWrapper = styled.div`
  a {
    text-decoration: none;
  }
`;

type SoloLinkLooksProps = {
  displayInline?: boolean;
  stylingOverrides?: string;
};
export const SoloLinkLooks = styled.div<SoloLinkLooksProps>`
  display: ${(props: SoloLinkLooksProps) =>
    props.displayInline ? 'inline' : 'block'};
  color: ${colors.seaBlue};
  cursor: pointer;

  &:hover,
  &:focus {
    color: ${colors.lakeBlue};
  }

  ${(props: SoloLinkLooksProps) => props.stylingOverrides}
`;

export type SimpleLinkProps = {
  displayElement: React.ReactNode;
  link: string;
  inline?: boolean;
  stylingOverrides?: string;
};

export const SoloLink = (props: SimpleLinkProps) => {
  return (
    <SoloLinkWrapper>
      <Link to={props.link}>
        <SoloLinkLooks
          displayInline={props.inline}
          stylingOverrides={props.stylingOverrides}>
          {props.displayElement}
        </SoloLinkLooks>
      </Link>
    </SoloLinkWrapper>
  );
};

// For tables' readability
export const RenderSimpleLink = (props?: SimpleLinkProps) => {
  if (!props) return null;

  return <SoloLink {...props} />;
};
