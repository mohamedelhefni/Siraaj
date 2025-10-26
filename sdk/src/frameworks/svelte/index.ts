import { onMount } from 'svelte';
import { analytics } from '../../core/analytics';
import type { AnalyticsConfig } from '../../core/types';

export { analytics };

export function createAnalytics() {
    const track = (eventName: string, properties?: Record<string, any>) => {
        analytics.track(eventName, properties);
    };

    const pageView = (url?: string | null) => {
        analytics.pageView(url);
    };

    const identify = (userId: string, traits?: Record<string, any>) => {
        analytics.identify(userId, traits);
    };

    const trackClick = (elementId: string, properties?: Record<string, any>) => {
        analytics.trackClick(elementId, properties);
    };

    const trackForm = (formId: string, properties?: Record<string, any>) => {
        analytics.trackForm(formId, properties);
    };

    const trackError = (error: Error | string, context?: Record<string, any>) => {
        analytics.trackError(error, context);
    };

    return {
        track,
        pageView,
        identify,
        trackClick,
        trackForm,
        trackError,
    };
}

export function initAnalytics(config: AnalyticsConfig) {
    Object.assign((analytics as any).config, config);
}

export function usePageTracking() {
    onMount(() => {
        analytics.pageView();
    });
}
