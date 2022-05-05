import lodash from 'lodash';
import React, { createContext, useContext, useEffect, useState } from 'react';

const defaultAppSettings = {
  keyboardHandler: 'keyboard',
};
const LOCAL_STORAGE_KEY = 'app-settings';
/**
 * Gets the app settings from local storage.
 * Uses defaultAppSettings as a fallback.
 */
const getCachedAppSettingsOrDefault = () => {
  const cachedAppSettings = localStorage.getItem(LOCAL_STORAGE_KEY);
  return (
    cachedAppSettings === null
      ? defaultAppSettings
      : JSON.parse(cachedAppSettings)
  ) as AppSettings;
};

//
// TYPES
//
type AppSettings = typeof defaultAppSettings;
interface IAppSettingsContext {
  appSettings: AppSettings;
  onAppSettingsChange(newSettings: AppSettings): void;
}

//
// CONTEXT
//
const AppSettingsContext = createContext({} as IAppSettingsContext);

//
// PROVIDER
//
export const AppSettingsProvider: React.FC = props => {
  const [appSettings, setAppSettings] = useState<AppSettings>(
    getCachedAppSettingsOrDefault()
  );
  const onAppSettingsChange = (newSettings: AppSettings) => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(newSettings));
    setAppSettings(newSettings);
  };

  // Checks localStorage on-page-focus.
  // This allows multiple windows/tabs to be open and share state.
  useEffect(() => {
    const checkAppSettings = () => {
      const newAppSettings = getCachedAppSettingsOrDefault();
      if (!lodash.isEqual(newAppSettings, appSettings))
        onAppSettingsChange(newAppSettings);
    };
    window.addEventListener('focus', checkAppSettings);
    return () => {
      window.removeEventListener('focus', checkAppSettings);
    };
  }, [appSettings, setAppSettings]);

  return (
    <AppSettingsContext.Provider value={{ appSettings, onAppSettingsChange }}>
      {props.children}
    </AppSettingsContext.Provider>
  );
};

//
// HOOK
//
export const useAppSettings = () => useContext(AppSettingsContext);
