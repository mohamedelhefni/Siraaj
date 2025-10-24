// API configuration
const API_BASE_URL = 'http://localhost:8080/api';

/**
 * Fetch analytics stats from the backend
 * @param {string} startDate - Start date in YYYY-MM-DD format
 * @param {string} endDate - End date in YYYY-MM-DD format
 * @returns {Promise<Object>} Analytics stats
 */
export async function fetchStats(startDate, endDate) {
    const params = new URLSearchParams();
    if (startDate) params.append('start', startDate);
    if (endDate) params.append('end', endDate);

    const response = await fetch(`${API_BASE_URL}/stats?${params}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch stats: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Track a new event
 * @param {Object} event - Event data
 * @returns {Promise<Object>} Response
 */
export async function trackEvent(event) {
    const response = await fetch(`${API_BASE_URL}/track`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(event),
    });
    if (!response.ok) {
        throw new Error(`Failed to track event: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Fetch debug events
 * @returns {Promise<Object>} Debug events
 */
export async function fetchDebugEvents() {
    const response = await fetch(`${API_BASE_URL}/debug/events`);
    if (!response.ok) {
        throw new Error(`Failed to fetch debug events: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Health check
 * @returns {Promise<Object>} Health status
 */
export async function healthCheck() {
    const response = await fetch(`${API_BASE_URL}/health`);
    if (!response.ok) {
        throw new Error(`Failed to check health: ${response.statusText}`);
    }
    return response.json();
}
