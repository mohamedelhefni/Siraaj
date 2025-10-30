import { analytics, AnalyticsCore } from './analytics';

export { analytics, AnalyticsCore };
export type { AnalyticsConfig, EventData, SessionInfo, DeviceInfo, Analytics } from './types';

// For UMD builds - set as default export for easier usage
export default analytics;
