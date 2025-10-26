// Core Analytics Types
export interface AnalyticsConfig {
  apiUrl?: string;
  projectId?: string;
  autoTrack?: boolean;
  debug?: boolean;
  bufferSize?: number;
  flushInterval?: number;
  timeout?: number;
  maxRetries?: number;
}

export interface EventData {
  event_name: string;
  user_id: string;
  session_id: string;
  session_duration: number;
  url: string;
  referrer: string;
  user_agent: string;
  timestamp: string;
  browser: string;
  os: string;
  device: string;
  project_id: string;
}

export interface SessionInfo {
  sessionId: string;
  userId: string;
}

export interface DeviceInfo {
  browser: string;
  os: string;
  device: string;
}

export interface Analytics {
  track(eventName: string, properties?: Record<string, any>): void;
  pageView(url?: string | null): void;
  identify(userId: string, traits?: Record<string, any>): void;
  trackClick(elementId: string, properties?: Record<string, any>): void;
  trackForm(formId: string, properties?: Record<string, any>): void;
  trackError(error: Error | string, context?: Record<string, any>): void;
  setUserProperties(properties: Record<string, any>): void;
  flush(async?: boolean): Promise<any>;
  destroy(): void;
}

