/** Format a 0–1 confidence as percentage string. */
export function FormatConfidence(c: number | null | undefined): string {
	if (c == null) return '—';
	return `${Math.round(c * 100)}%`;
}

/** CSS class for confidence badge. */
export function ConfidenceClass(c: number): string {
	if (c >= 0.9) return 'conf-high';
	if (c >= 0.7) return 'conf-mid';
	return 'conf-low';
}
