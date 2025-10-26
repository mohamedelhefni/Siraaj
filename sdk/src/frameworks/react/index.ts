import { useState, useEffect, useCallback, useRef } from 'react';
import { analytics } from '../../core/analytics';
import type { AnalyticsConfig } from '../../core/types';

export { analytics };

export interface UseAnalyticsReturn {
    track: (eventName: string, properties?: Record<string, any>) => void;
    pageView: (url?: string | null) => void;
    identify: (userId: string, traits?: Record<string, any>) => void;
    trackClick: (elementId: string, properties?: Record<string, any>) => void;
    trackForm: (formId: string, properties?: Record<string, any>) => void;
    trackError: (error: Error | string, context?: Record<string, any>) => void;
}

export function useAnalytics(): UseAnalyticsReturn {
    const track = useCallback((eventName: string, properties?: Record<string, any>) => {
        analytics.track(eventName, properties);
    }, []);

    const pageView = useCallback((url?: string | null) => {
        analytics.pageView(url);
    }, []);

    const identify = useCallback((userId: string, traits?: Record<string, any>) => {
        analytics.identify(userId, traits);
    }, []);

    const trackClick = useCallback((elementId: string, properties?: Record<string, any>) => {
        analytics.trackClick(elementId, properties);
    }, []);

    const trackForm = useCallback((formId: string, properties?: Record<string, any>) => {
        analytics.trackForm(formId, properties);
    }, []);

    const trackError = useCallback((error: Error | string, context?: Record<string, any>) => {
        analytics.trackError(error, context);
    }, []);

    return {
        track,
        pageView,
        identify,
        trackClick,
        trackForm,
        trackError,
    };
}

export interface AnalyticsProviderProps {
    config: AnalyticsConfig;
    children: React.ReactNode;
}

export function AnalyticsProvider({ config, children }: AnalyticsProviderProps) {
    const initialized = useRef(false);

    useEffect(() => {
        if (!initialized.current) {
            // Re-initialize analytics with new config
            Object.assign((analytics as any).config, config);
            initialized.current = true;
        }
    }, [config]);

    return children;
}

export function usePageTracking() {
    const { pageView } = useAnalytics();

    useEffect(() => {
        pageView();
    }, [pageView]);
}
