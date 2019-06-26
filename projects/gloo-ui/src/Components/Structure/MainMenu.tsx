import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import { NavLink } from 'react-router-dom';
import { colors } from 'Styles';
import { ReactComponent as GlooE } from 'assets/GlooEE.svg';
import { ReactComponent as SettingsGear } from 'assets/settings-gear.svg';
import { ReactComponent as HelpBubble } from 'assets/help-icon.svg';
import styled from '@emotion/styled/macro';
import { Popover } from 'antd';

const NavLinkStyles = {
  display: 'inline-block',
  color: 'white',
  textDecoration: 'none',
  fontSize: '18px',
  marginRight: '50px',
  fontWeight: 300
};

const Container = styled.div`
  min-width: 1070px;
  padding: 10px 0;
  height: 55px;
  line-height: 36px;
  background: ${colors.seaBlue};
`;
const InnerContainer = styled.div`
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

const TitleDiv = styled.div`
  position: relative;
  display: inline-block;
  color: ${colors.puddleBlue};
  font-size: 18px;
  margin-right: 50px;
  padding-right: 50px;
  border-right: 1px solid ${colors.lakeBlue};
  cursor: default;

  > svg {
    position: absolute;
    left: -58px;
    width: auto;
    height: 35px;
  }
`;

const activeStyle = {
  borderBottom: `8px solid ${colors.pondBlue}`,
  cursor: 'default'
};
const activeSettingsStyle = {
  cursor: 'default'
};

const HelpHolder = styled.div`
  display: flex;
  height: 36px;
  line-height: 36px;
  align-items: center;
  line-height: 46px;
  float: right;
  margin-right: 10px;
  padding-right: 10px;
  border-right: 1px solid ${colors.lakeBlue};
  cursor: pointer;
`;
const CommLinkCss = css`
  display: block;
  color: white;
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 5px;

  &:hover,
  &:focus {
    color: ${colors.januaryGrey};
  }
`;

const DocumentationLink = styled.a`
  ${CommLinkCss};
`;
const VersionDisplay = styled.div`
  margin-top: 8px;
  border-top: 1px solid white;
  padding-top: 8px;
  font-weight: 300;
`;

export const MainMenu = () => {
  return (
    <Container>
      <InnerContainer>
        <TitleDiv>
          <GlooE /> Gloo Enterprise
        </TitleDiv>
        <NavLink
          style={NavLinkStyles}
          to='/virtualservices/'
          activeStyle={activeStyle}>
          Virtual Services
        </NavLink>
        <NavLink
          style={NavLinkStyles}
          to='/upstreams/'
          activeStyle={activeStyle}>
          Upstreams
        </NavLink>
        <NavLink style={NavLinkStyles} to='/stats/' activeStyle={activeStyle}>
          Stats
        </NavLink>
        <NavLink
          style={{
            ...NavLinkStyles,
            float: 'right',
            fontSize: '33px',
            marginRight: '0',
            display: 'flex',
            height: '36px',
            lineHeight: '36px',
            alignItems: 'center'
          }}
          to='/settings/'
          activeStyle={activeSettingsStyle}>
          <SettingsGear />
        </NavLink>

        <HelpHolder>
          <Popover
            trigger='click'
            mouseLeaveDelay={0.2}
            content={
              <div>
                <DocumentationLink
                  href='https://gloo.solo.io/enterprise/'
                  target='_blank'>
                  Join the Community
                </DocumentationLink>

                <VersionDisplay>Version: {'unknown'}</VersionDisplay>
              </div>
            }>
            <HelpBubble />
          </Popover>
        </HelpHolder>
      </InnerContainer>
    </Container>
  );
};
