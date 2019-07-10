export const soloConstants = {
  smallBuffer: 18,
  buffer: 20,
  largeBuffer: 23,

  smallRadius: 8,
  radius: 10,
  largeRadius: 16,

  transitionTime: '.3s'
};

/*
export function getHealthStatus(
  
): {
  value: number;
  color: string;
} {
  const installStatus = installState
    ? installState
    : applicationState && applicationState.installationState
    ? applicationState.installationState.status
    : InstallStatus.PENDING_INSTALL;

  if (installStatus === InstallStatus.PENDING_INSTALL) {
    return soloConstants.healthStatus.Pending;
  }

  const healthState =
    applicationState && applicationState.applicationHealth
      ? applicationState.applicationHealth.state
      : HealthState.HEALTHY;

  if (
    installStatus === InstallStatus.INSTALLED &&
    healthState === HealthState.HEALTHY
  ) {
    return soloConstants.healthStatus.Good;
  }

  return soloConstants.healthStatus.Error;
}*/
