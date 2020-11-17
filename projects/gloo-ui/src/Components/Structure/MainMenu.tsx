import { css } from '@emotion/core';
import styled from '@emotion/styled';
import { Popover } from 'antd';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
import { ReactComponent as GlooE } from 'assets/GlooEE.svg';
import { ReactComponent as HelpBubble } from 'assets/help-icon.svg';
import { ReactComponent as SettingsGear } from 'assets/settings-gear.svg';
import * as React from 'react';
import { NavLink } from 'react-router-dom';
import { colors } from 'Styles';
import useSWR from 'swr';
import { configAPI } from '../../store/config/api';
import glooEdge from 'assets/gloo-edge.png';
import glooEdgeE from 'assets/gloo-edge-e.png';

const NavLinkStyles = {
  display: 'inline-block',
  color: 'white',
  textDecoration: 'none',
  fontSize: '18px',
  marginRight: '50px',
  fontWeight: 300
};

const Container = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-width: 1070px;
  height: 55px;
  background: ${colors.seaBlue};
`;
const InnerContainer = styled.div`
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

const TitleDiv = styled.div`
  display: flex;
  align-items: center;
  width: auto;
  color: ${colors.puddleBlue};
  font-size: 18px;
  padding-right: 50px;
  border-right: 1px solid ${colors.lakeBlue};
  cursor: default;

  > svg {
    position: absolute;
    left: 0;
    width: auto;
    height: 35px;
  }
  .settings-gear-a {
    stroke: white;
  }
`;

const activeStyle = {
  borderBottom: `8px solid ${colors.pondBlue}`,
  cursor: 'default',
  fontWeight: 500
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
  const { data: version, error: versionError } = useSWR(
    'getVersion',
    configAPI.getVersion,
    { refreshInterval: 0 }
  );
  const { data: licenseData, error: licenseError } = useSWR(
    'hasValidLicense',
    configAPI.getIsLicenseValid,
    { refreshInterval: 0 }
  );

  const hasValidLicense = licenseData?.isLicenseValid;
  return (
    <div
      className='relative flex items-center justify-between max-w-full pt-2 px-28'
      style={{ backgroundColor: colors.seaBlue }}>
      <div className='flex items-center justify-between min-w-full'>
        <div className='flex items-center justify-start mt-2 justify-items-start '>
          <div className='flex items-center w-auto pb-3 mr-4'>
            {hasValidLicense ? (
              <>
                <img
                  src={glooEdgeE}
                  alt='Gloo Edge Enterprise'
                  className='object-cover h-12 pb-1 w-30'
                />
                <div
                  style={{
                    color: colors.lakeBlue,
                    backgroundColor: colors.lakeBlue
                  }}
                  className='w-px mx-6 h-9'></div>
              </>
            ) : (
              <>
                <img
                  src={glooEdge}
                  alt='Gloo Edge'
                  className='object-cover w-64 h-12 pb-1'
                />

                <div
                  style={{
                    color: colors.lakeBlue,
                    backgroundColor: colors.lakeBlue
                  }}
                  className='w-px mx-6 h-9'></div>
              </>
            )}
          </div>
          <NavLink
            data-testid='overview-navlink'
            style={NavLinkStyles}
            to='/overview'
            activeStyle={activeStyle}>
            Overview
          </NavLink>
          <NavLink
            data-testid='virtual-services-navlink'
            style={NavLinkStyles}
            to='/virtualservices/'
            activeStyle={activeStyle}>
            Virtual Services
          </NavLink>
          <NavLink
            data-testid='upstreams-navlink'
            style={NavLinkStyles}
            to='/upstreams/'
            activeStyle={activeStyle}>
            Upstreams
          </NavLink>
          <NavLink
            data-testid='wasm-navlink'
            style={NavLinkStyles}
            to='/wasm/'
            activeStyle={activeStyle}>
            Wasm
          </NavLink>
        </div>
        <div>
          <NavLink
            data-testid='settings-navlink'
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
            to='/admin/'
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
                    href='https://slack.solo.io/'
                    target='_blank'>
                    Join the Community
                  </DocumentationLink>

                  <VersionDisplay>
                    {hasValidLicense ? 'Version: ' : 'UI Version: '}
                    {version !== undefined ? version : 'unknown'}
                  </VersionDisplay>
                </div>
              }>
              <HelpBubble />
            </Popover>
          </HelpHolder>
        </div>
      </div>
    </div>
  );
};
