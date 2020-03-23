import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';
import { css } from '@emotion/core';
import React from 'react';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as DevPortalIcon } from 'assets/developer-portal-icon.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';

import { SectionCard } from 'Components/Common/SectionCard';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { ReactComponent as PortalIcon } from 'assets/portal-icon.svg';
import { Container, Header, HealthScoreContainer } from '../Admin/AdminLanding';
import { healthConstants } from 'Styles';
import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { SwaggerExplorer } from './SwaggerExplorer';
import { ErrorBoundary } from '../Errors/ErrorBoundary';

export const DevPortalOverview = () => {
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <Container>
        <Header>
          <div>
            <div className='text-2xl '>Developer Portal</div>
            <div className='text-lg text-gray-700'>
              Advanced Administration for your Gloo Configuration
            </div>
          </div>
          <HealthScoreContainer health={healthConstants.Good.value}>
            <span className='text-blue-500'>
              <DevPortalIcon className='w-32 h-32 fill-current' />
            </span>
          </HealthScoreContainer>
        </Header>
        <div className='grid grid-cols-3 px-4 py-5 sm:p-6'>
          <StatusTile
            exploreMoreLink={{
              prompt: 'View Portals',
              link: '/dev-portal/portals'
            }}
            healthStatus={1}
            titleText='Portals'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'
            titleIcon={
              <span className='text-blue-500'>
                <PortalIcon className='w-8 h-8 fill-current ' />
              </span>
            }>
            <div className='grid grid-cols-2 gap-4 '>
              <TallyInformationDisplay
                tallyCount={2}
                tallyDescription={`portals`}
                color='blue'
              />
              <TallyInformationDisplay
                tallyCount={3}
                tallyDescription={`published APIs`}
                color='blue'
              />
            </div>
          </StatusTile>
          <StatusTile
            exploreMoreLink={{
              prompt: 'View APIs',
              link: '/dev-portal/apis'
            }}
            healthStatus={1}
            titleText='APIs'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'
            titleIcon={
              <span className='text-blue-500'>
                <CodeIcon className='w-8 h-8 fill-current' />
              </span>
            }>
            <>
              {/* <div className='grid grid-cols-2 gap-4 '>
                <TallyInformationDisplay
                  tallyCount={5}
                  tallyDescription={'APIs'}
                  color='blue'
                />
                <TallyInformationDisplay
                  tallyCount={25}
                  tallyDescription={'Endpoints'}
                  color='blue'
                />
              </div> */}
              <TallyInformationDisplay
                tallyDescription={`You have no APIs detected. Get started by creating an API!`}
                color='blue'
              />
            </>
          </StatusTile>

          <StatusTile
            exploreMoreLink={{
              prompt: 'View Users & Groups',
              link: '/dev-portal/users'
            }}
            healthStatus={1}
            titleText='Users & Groups'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'
            titleIcon={
              <span className='text-blue-500'>
                <UserIcon className='w-8 h-8 fill-current ' />
              </span>
            }>
            <div className='grid grid-cols-2 gap-4 '>
              <TallyInformationDisplay
                tallyCount={2}
                tallyDescription='Users'
                color='blue'
              />
              <TallyInformationDisplay
                tallyCount={3}
                tallyDescription='Groups'
                color='blue'
              />
            </div>
          </StatusTile>
        </div>
        <div>
          <div className='text-2xl text-gray-900 '>
            Developer Portal Documentation
          </div>
        </div>
        <div className='grid grid-cols-3 px-4 py-5 sm:p-6'>
          <StatusTile
            exploreMoreLink={{ prompt: 'View Documentation', link: '' }}
            titleText='Creating a Portal'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'></StatusTile>
          <StatusTile
            exploreMoreLink={{ prompt: 'View Documentation', link: '' }}
            titleText='Generating a Gloo API'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'></StatusTile>

          <StatusTile
            exploreMoreLink={{ prompt: 'View Documentation', link: '' }}
            titleText='Granting User Access'
            description='Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.'></StatusTile>
        </div>
      </Container>
    </ErrorBoundary>
  );
};
