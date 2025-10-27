/**
 * Format large numbers into compact notation (5M, 100K, etc.)
 * @param {number} num - The number to format
 * @param {number} decimals - Number of decimal places (default: 1)
 * @returns {string} Formatted number string
 */
export function formatCompactNumber(num, decimals = 1) {
	if (num === null || num === undefined || isNaN(num)) {
		return '0';
	}

	const absNum = Math.abs(num);

	// Less than 1000, show as-is
	if (absNum < 1000) {
		return num.toString();
	}

	const suffixes = [
		{ value: 1e9, symbol: 'B' },  // Billion
		{ value: 1e6, symbol: 'M' },  // Million
		{ value: 1e3, symbol: 'K' }   // Thousand
	];

	for (const { value, symbol } of suffixes) {
		if (absNum >= value) {
			const formatted = (num / value).toFixed(decimals);
			// Remove trailing zeros and decimal point if not needed
			return formatted.replace(/\.0+$/, '') + symbol;
		}
	}

	return num.toString();
}

/**
 * Format number with thousands separators
 * @param {number} num - The number to format
 * @returns {string} Formatted number with commas
 */
export function formatNumberWithCommas(num) {
	if (num === null || num === undefined || isNaN(num)) {
		return '0';
	}
	return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

/**
 * Format percentage
 * @param {number} value - The value
 * @param {number} total - The total
 * @param {number} decimals - Number of decimal places (default: 1)
 * @returns {string} Formatted percentage
 */
export function formatPercentage(value, total, decimals = 1) {
	if (!total || total === 0) return '0%';
	const percentage = (value / total) * 100;
	return `${percentage.toFixed(decimals)}%`;
}

/**
 * Format bytes to human readable format
 * @param {number} bytes - The number of bytes
 * @param {number} decimals - Number of decimal places (default: 2)
 * @returns {string} Formatted byte string
 */
export function formatBytes(bytes, decimals = 2) {
	if (bytes === 0) return '0 Bytes';

	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

	const i = Math.floor(Math.log(bytes) / Math.log(k));

	return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/**
 * Format duration in milliseconds to human readable format
 * @param {number} ms - Duration in milliseconds
 * @returns {string} Formatted duration
 */
export function formatDuration(ms) {
	if (ms < 1000) return `${ms}ms`;
	if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
	if (ms < 3600000) return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`;
	return `${Math.floor(ms / 3600000)}h ${Math.floor((ms % 3600000) / 60000)}m`;
}

/**
 * Format number with automatic unit selection based on context
 * @param {number} num - The number to format
 * @param {string} context - Context: 'compact', 'full', or 'auto' (default)
 * @returns {string} Formatted number
 */
export function formatNumber(num, context = 'auto') {
	if (context === 'compact') {
		return formatCompactNumber(num);
	}
	if (context === 'full') {
		return formatNumberWithCommas(num);
	}
	// Auto: use compact for large numbers, full for small
	return Math.abs(num) >= 10000 ? formatCompactNumber(num) : formatNumberWithCommas(num);
}

/**
 * Format change/delta with +/- sign and optional percentage
 * @param {number} current - Current value
 * @param {number} previous - Previous value
 * @param {boolean} asPercentage - Show as percentage (default: false)
 * @returns {string} Formatted change
 */
export function formatChange(current, previous, asPercentage = false) {
	if (!previous || previous === 0) return '+0';
	
	const diff = current - previous;
	const sign = diff >= 0 ? '+' : '';
	
	if (asPercentage) {
		const percentChange = ((diff / previous) * 100).toFixed(1);
		return `${sign}${percentChange}%`;
	}
	
	return `${sign}${formatCompactNumber(diff)}`;
}
