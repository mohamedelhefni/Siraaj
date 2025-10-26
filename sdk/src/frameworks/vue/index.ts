import { onMounted } from 'vue';
import { analytics } from '../../core/analytics';
import type { AnalyticsConfig } from '../../core/types';

export { analytics };

export function useAnalytics() {
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

export const AnalyticsPlugin = {
    install(app: any, config: AnalyticsConfig) {
        // Re-initialize analytics with config
        Object.assign((analytics as any).config, config);

        app.config.globalProperties.$analytics = analytics;
        app.provide('analytics', analytics);
    },
};

export function usePageTracking() {
    const { pageView } = useAnalytics();

    onMounted(() => {
        pageView();
    });
}
