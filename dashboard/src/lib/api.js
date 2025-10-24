// API configuration
const API_BASE_URL = 'http://localhost:8080/api';

/**
 * Fetch analytics stats from the backend
 * @param {string} startDate - Start date in YYYY-MM-DD format
 * @param {string} endDate - End date in YYYY-MM-DD format
 * @param {number} limit - Limit for top results (default 50)
 * @param {Object} filters - Optional filters {source, country, browser, event, project}
 * @returns {Promise<Object>} Analytics stats
 */
export async function fetchStats(startDate, endDate, limit = 50, filters = {}) {
    const params = new URLSearchParams();
    if (startDate) params.append('start', startDate);
    if (endDate) params.append('end', endDate);
    if (limit) params.append('limit', limit.toString());

    // Add filters
    if (filters.source) params.append('source', filters.source);
    if (filters.country) params.append('country', filters.country);
    if (filters.browser) params.append('browser', filters.browser);
    if (filters.event) params.append('event', filters.event);
    if (filters.project) params.append('project', filters.project);

    const response = await fetch(`${API_BASE_URL}/stats?${params}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch stats: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Fetch all events with pagination
 * @param {string} startDate - Start date
 * @param {string} endDate - End date
 * @param {number} limit - Number of events to fetch
 * @param {number} offset - Offset for pagination
 * @returns {Promise<Object>} Events data
 */
export async function fetchEvents(startDate, endDate, limit = 100, offset = 0) {
    const params = new URLSearchParams();
    if (startDate) params.append('start', startDate);
    if (endDate) params.append('end', endDate);
    params.append('limit', limit.toString());
    params.append('offset', offset.toString());

    const response = await fetch(`${API_BASE_URL}/events?${params}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch events: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Fetch online users count
 * @param {number} timeWindowMinutes - Time window in minutes (default 5)
 * @returns {Promise<Object>} Online users data
 */
export async function fetchOnlineUsers(timeWindowMinutes = 5) {
    const params = new URLSearchParams();
    params.append('window', timeWindowMinutes.toString());

    const response = await fetch(`${API_BASE_URL}/online?${params}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch online users: ${response.statusText}`);
    }
    return response.json();
}

/**
 * Fetch list of available projects
 * @returns {Promise<Array<string>>} List of project IDs
 */
export async function fetchProjects() {
    const response = await fetch(`${API_BASE_URL}/projects`);
    if (!response.ok) {
        throw new Error(`Failed to fetch projects: ${response.statusText}`);
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
