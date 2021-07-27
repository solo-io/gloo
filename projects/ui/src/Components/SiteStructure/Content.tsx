import React from 'react';
import { Routes, Route } from 'react-router';
import { ErrorBoundary } from '../Common/ErrorBoundary';
import styled from '@emotion/styled';
import { OverviewLanding } from 'Components/Features/Overview/OverviewLanding';
import { GlooInstancesLanding } from 'Components/Features/GlooInstance/GlooInstancesLanding';
import { VirtualServicesLanding } from 'Components/Features/VirtualService/VirtualServicesLanding';
import { UpstreamsLanding } from 'Components/Features/Upstream/UpstreamsLanding';
import { Breadcrumb } from './Breadcrumb';
import { GlooInstancesDetails } from 'Components/Features/GlooInstance/GlooInstanceDetails';
import { UpstreamDetails } from 'Components/Features/Upstream/UpstreamDetails';
import { AdminLanding } from 'Components/Features/Admin/AdminLanding';
import { GlooInstanceAdministration } from 'Components/Features/GlooInstance/Admin/GlooInstanceAdministration';
import { GlooAdminInnerPagesWrapper } from 'Components/Features/GlooInstance/Admin/GlooAdminInnerPagesWrapper';
import { AdminInnerPagesWrapper } from 'Components/Features/Admin/AdminInnerPagesWrapper';
import { VirtualServiceDetails } from 'Components/Features/VirtualService/VirtualServiceDetails';
import { Clusters } from 'Components/Features/Admin/Clusters';
import { UpstreamGroupDetails } from 'Components/Features/Upstream/UpstreamGroupDetails';
import { WasmLanding } from 'Components/Features/Wasm/WasmLanding';
import { useIsGlooFedEnabled } from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';

const ScrollContainer = styled.div`
  max-height: 100%;
  overflow: auto;
`;

const Container = styled.div`
  padding: 20px 0;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

export const Content = () => {
  const {
    data: glooFedCheckResponse,
    error: glooFedCheckError,
  } = useIsGlooFedEnabled();
  if (!!glooFedCheckError) {
    return <DataError error={glooFedCheckError} />;
  } else if (!glooFedCheckResponse) {
    return <Loading />;
  }

  return (
    <ScrollContainer>
      <Breadcrumb />
      <Container>
        {glooFedCheckResponse.enabled ? (
          <Routes>
            <Route
              path='/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information to get started.</div>
                  }>
                  <OverviewLanding />
                </ErrorBoundary>
              }
            />

            <Route
              path='/gloo-instances/:namespace/:name/gloo-admin/:adminPage'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Gloo Instances.</div>
                  }>
                  <GlooAdminInnerPagesWrapper />
                </ErrorBoundary>
              }
            />
            <Route
              path='/gloo-instances/:namespace/:name/gloo-admin/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Gloo Instances.</div>
                  }>
                  <GlooInstanceAdministration />
                </ErrorBoundary>
              }
            />

            <Route
              path='/gloo-instances/:namespace/:name/virtual-services/:virtualserviceclustername/:virtualservicenamespace/:virtualservicename'
              element={
                <ErrorBoundary
                  fallback={
                    <div>
                      Unable to pull information on Virtual Service Details.
                    </div>
                  }>
                  <VirtualServiceDetails />
                </ErrorBoundary>
              }
            />
            <Route
              path='/gloo-instances/:namespace/:name/upstreams/:upstreamClusterName/:upstreamNamespace/:upstreamName'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Upstream Details.</div>
                  }>
                  <UpstreamDetails />
                </ErrorBoundary>
              }
            />
            <Route
              path='/gloo-instances/:namespace/:name/upstream-groups/:upstreamGroupClusterName/:upstreamGroupNamespace/:upstreamGroupName'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Upstream Details.</div>
                  }>
                  <UpstreamGroupDetails />
                </ErrorBoundary>
              }
            />
            <Route
              path='/gloo-instances/:namespace/:name/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Gloo Instances.</div>
                  }>
                  <GlooInstancesDetails />
                </ErrorBoundary>
              }
            />
            <Route
              path='/gloo-instances/*'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Gloo Instances.</div>
                  }>
                  <GlooInstancesLanding />
                </ErrorBoundary>
              }
            />

            <Route
              path='/virtual-services/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Virtual Services.</div>
                  }>
                  <VirtualServicesLanding />
                </ErrorBoundary>
              }
            />

            <Route
              path='/upstreams/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Upstreams.</div>
                  }>
                  <UpstreamsLanding />
                </ErrorBoundary>
              }
            />

            <Route
              path='/wasm-filters/:filterName/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Virtual Services.</div>
                  }>
                  <WasmLanding />
                </ErrorBoundary>
              }
            />
            <Route
              path='/wasm-filters/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull information on Virtual Services.</div>
                  }>
                  <WasmLanding />
                </ErrorBoundary>
              }
            />

            <Route
              path='/admin/clusters'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull administrative information.</div>
                  }>
                  <Clusters />
                </ErrorBoundary>
              }
            />
            <Route
              path='/admin/federated-resources/:adminPage/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull administrative information.</div>
                  }>
                  <AdminInnerPagesWrapper />
                </ErrorBoundary>
              }
            />
            <Route
              path='/admin/federated-resources/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull administrative information.</div>
                  }>
                  <AdminInnerPagesWrapper />
                </ErrorBoundary>
              }
            />
            <Route
              path='/admin/'
              element={
                <ErrorBoundary
                  fallback={
                    <div>Unable to pull administrative information.</div>
                  }>
                  <AdminLanding />
                </ErrorBoundary>
              }
            />
          </Routes>
        ) : (
          <h1>Coming soon...</h1>
        )}
      </Container>
    </ScrollContainer>
  );
};
