export function FormatEUR(cents: number | null | undefined): string {
    if (cents == null) return '—';
    return new Intl.NumberFormat('nl-NL', {
        style: 'currency',
        currency: 'EUR'
    }).format(cents / 100);
}

export function FormatPct(ratio: number): string {
    const pct = ratio * 100;
    const sign = pct > 0 ? '+' : '';
    return `${sign}${pct.toFixed(1)}%`;
}

export function ArbitrageClass(ratio: number): string {
    if (ratio >= 0.1) return 'arb-great';
    if (ratio >= 0.03) return 'arb-good';
    if (ratio >= 0) return 'arb-ok';
    return 'arb-bad';
}