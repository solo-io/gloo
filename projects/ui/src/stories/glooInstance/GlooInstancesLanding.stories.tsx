import { expect } from '@storybook/jest';
import { ComponentMeta, ComponentStory } from '@storybook/react';
import { waitFor, within } from '@storybook/testing-library';
import {
  useListGlooInstances,
  useListUpstreams,
  useListVirtualServices,
} from 'API/hooks';
import { ListVirtualServicesResponse } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { ListUpstreamsResponse } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { DiProvider } from 'react-magnetic-di';
import { MemoryRouter } from 'react-router';
import {
  createGlooCheck,
  createGlooInstanceObj,
  createObjSpec,
} from 'stories/mocks/generators';
import { createSwrInjectable } from 'stories/mocks/story-helpers';
import { GlooInstancesLanding } from '../../Components/Features/GlooInstance/GlooInstancesLanding';

export default {
  title: 'GlooInstance / GlooInstancesLanding',
  component: GlooInstancesLanding,
} as ComponentMeta<typeof GlooInstancesLanding>;

const Template: ComponentStory<typeof GlooInstancesLanding> = args => {
  const useListGlooInstancesDi = createSwrInjectable(useListGlooInstances, [
    createGlooInstanceObj({
      spec: createObjSpec({
        check: createGlooCheck({
          // `undefined` gloo instance checks shouldn't break this page.
          // If they do, the page won't render and the storybook test will fail.
          matchableHttpGateways: undefined,
        }),
      }),
    }),
  ]);
  const useListUpstreamsDi = createSwrInjectable(useListUpstreams, {
    upstreamsList: [],
    total: 0,
  } as ListUpstreamsResponse.AsObject);
  const useListVirtualServicesDi = createSwrInjectable(useListVirtualServices, {
    virtualServicesList: [],
    total: 0,
  } as ListVirtualServicesResponse.AsObject);
  return (
    <DiProvider
      use={[
        useListGlooInstancesDi,
        useListUpstreamsDi,
        useListVirtualServicesDi,
      ]}>
      <MemoryRouter>
        <GlooInstancesLanding />
      </MemoryRouter>
    </DiProvider>
  );
};

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstancesLanding>;

Primary.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  await waitFor(async () => {
    const landingPage = canvas.getByTestId('gloo-instances-landing');
    expect(landingPage).not.toBeNull();
  });
};
