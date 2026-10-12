<script lang="ts">
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { GetListings } from '../../bindings/swail/mpservice';
	import { FormatEUR, FormatPct, ArbitrageClass } from '../lib/util';
	import type { Listing } from '../../bindings/swail/models';

	let listings = $state<Listing[]>([]);
	let loading = $state(false);
	let sortBy = $state<'arbitrage' | 'price' | 'confidence' | 'favorites'>('arbitrage');
	let filterMinArbitrage = $state(0);

	function sortListings() {
		const filtered = listings.filter((l) => l.Product.ArbitrageRatio >= filterMinArbitrage);
		return filtered.sort((a, b) => {
			switch (sortBy) {
				case 'arbitrage':
					return b.Product.ArbitrageRatio - a.Product.ArbitrageRatio;
				case 'price':
					return (a.Product.AskingPrice ?? 0) - (b.Product.AskingPrice ?? 0);
				case 'confidence':
					return b.Product.Confidence - a.Product.Confidence;
				case 'favorites':
					return b.Product.FavoriteCount - a.Product.FavoriteCount;
				default:
					return 0;
			}
		});
	}

	function openListing(url: string) {
		// Wails typically uses runtime.BrowserOpenURL
		window.open(url, '_blank');
	}

	async function refreshListings() {
		loading = true;
		try {
			listings = (await GetListings()) ?? [];
		} catch (err) {
			console.error('Failed to load listings:', err);
			listings = [];
		} finally {
			loading = false;
		}
	}

	// Widen the shared .container only while this page is active
	onMount(() => {
		refreshListings();
		const el = document.querySelector('.container');
		if (el) el.classList.add('wide');
		return () => {
			if (el) el.classList.remove('wide');
		};
	});
</script>

<button class="back-btn" onclick={() => push('/')}>← Back</button>

<div class="page-header">
	<h1 class="title">
		<span class="title-accent">Marktplaats</span>
	</h1>
	<p class="subtitle">Live deals ranked by arbitrage opportunity</p>
</div>

<div class="content-layout">
	<!-- Controls -->
	<div class="controls">
		<div class="sort-group">
			<span class="control-label">Sort</span>
			<div class="pill-row">
				<button
					class="pill"
					class:active={sortBy === 'arbitrage'}
					onclick={() => (sortBy = 'arbitrage')}
				>
					Arbitrage
				</button>
				<button
					class="pill"
					class:active={sortBy === 'price'}
					onclick={() => (sortBy = 'price')}
				>
					Price
				</button>
				<button
					class="pill"
					class:active={sortBy === 'confidence'}
					onclick={() => (sortBy = 'confidence')}
				>
					Confidence
				</button>
				<button
					class="pill"
					class:active={sortBy === 'favorites'}
					onclick={() => (sortBy = 'favorites')}
				>
					♥ Favs
				</button>
			</div>
		</div>

		<div class="filter-group">
			<span class="control-label">Min arb.</span>
			<input
				type="range"
				min="-0.2"
				max="0.4"
				step="0.01"
				bind:value={filterMinArbitrage}
				class="arb-slider"
			/>
			<span class="filter-value">{FormatPct(filterMinArbitrage)}</span>
		</div>

		<button class="btn refresh-btn" onclick={refreshListings} disabled={loading}>
			{#if loading}
				<span class="spinner"></span>
			{:else}
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
					<path d="M3 3v5h5" />
					<path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
					<path d="M16 16h5v5" />
				</svg>
			{/if}
			Refresh
		</button>
	</div>

	<!-- Listing grid -->
	<div class="listings">
		{#each sortListings() as listing (listing.Product.ItemID)}
			{@const p = listing.Product}
			{@const s = listing.Seller}

			<a
				class="listing-card"
				href={p.LotURL}
				target="_blank"
				rel="noopener noreferrer"
				onclick={(e) => {
					e.preventDefault();
					openListing(p.LotURL);
				}}
			>	
				<!-- Image -->
				<div class="thumb">
					{#if p.ImageURLs?.length}
						<img src={p.ImageURLs[0]} alt={p.Title} loading="lazy" />
					{:else}
						<div class="thumb-placeholder">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
							</svg>
						</div>
					{/if}
					<span class="arb-badge {ArbitrageClass(p.ArbitrageRatio)}">
						{FormatPct(p.ArbitrageRatio)}
					</span>
				</div>

				<!-- Body -->
				<div class="card-body">
					<div class="card-top">
						<h2 class="card-title">{p.Title}</h2>
						{#if s.IsPro}
							<span class="pro-badge">PRO</span>
						{/if}
					</div>

					<div class="meta-row">
						<span class="seller">@{s.Username}</span>
						<span class="dot-sep">·</span>
						<span class="favs">♥ {p.FavoriteCount}</span>
						<span class="dot-sep">·</span>
						<span class="weight">{p.Weight} g</span>
					</div>

					<div class="price-grid">
						<div class="price-block">
							<span class="price-label">Asking</span>
							<span class="price-value">{FormatEUR(p.AskingPrice)}</span>
						</div>
						<div class="price-block">
							<span class="price-label">Melt</span>
							<span class="price-value melt">{FormatEUR(p.MeltValue)}</span>
						</div>
						<div class="price-block">
							<span class="price-label">AI est.</span>
							<span class="price-value">{FormatEUR(p.AiEstimate)}</span>
						</div>
						<div class="price-block">
							<span class="price-label">Ship</span>
							<span class="price-value">{FormatEUR(p.DeliveryFeeEUR)}</span>
						</div>
					</div>

					<div class="confidence-bar">
						<div class="confidence-fill" style="width: {p.Confidence * 100}%"></div>
						<span class="confidence-label">{Math.round(p.Confidence * 100)}% conf.</span>
					</div>
				</div>
			</a>
		{:else}
			<div class="empty-state">
				<p>No listings match the current filter.</p>
			</div>
		{/each}
	</div>
</div>

<style>
	/* Local styles that extend the shared neon-night tokens */

	/* Temporarily widen the shared container for this page only */
	:global(.container.wide) {
		max-width: 1080px;
	}

	.page-header {
		margin-bottom: var(--s-4);
	}

	.title-sub {
		display: block;
		font-size: 0.55em;
		font-weight: 600;
		letter-spacing: 0.04em;
		color: var(--muted);
		margin-top: 0.15em;
	}

	/* ----- Layout: controls left / list right on wide screens ------------- */
	.content-layout {
		display: flex;
		flex-direction: column;
		gap: var(--s-4);
		width: 100%;
	}

	/* ----- Controls -------------------------------------------------------- */
	.controls {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		gap: var(--s-3);
		width: 100%;
		max-width: 720px;
	}

	.control-label {
		display: block;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--muted);
		margin-bottom: 6px;
	}

	.pill-row {
		display: flex;
		gap: 6px;
		flex-wrap: wrap;
	}

	.pill {
		padding: 5px 12px;
		font-size: 13px;
		font-weight: 600;
		border-radius: 999px;
		border: 1px solid var(--glass-border);
		background: var(--glass);
		color: var(--muted);
		cursor: pointer;
		transition: all 0.15s ease;
		--wails-draggable: no-drag;
	}

	.pill:hover {
		border-color: rgba(255, 77, 77, 0.45);
		color: var(--text);
	}

	.pill.active {
		background: rgba(255, 45, 114, 0.18);
		border-color: rgba(255, 45, 114, 0.55);
		color: #ff7a8f;
	}

	.filter-group {
		display: flex;
		flex-direction: column;
	}

	.arb-slider {
		width: 140px;
		accent-color: var(--accent);
		cursor: pointer;
	}

	.filter-value {
		font-size: 13px;
		font-weight: 600;
		color: var(--text);
		margin-top: 2px;
	}

	.refresh-btn {
		margin-left: auto;
		padding: 7px 16px;
		font-size: 13px;
	}

	.refresh-btn:disabled {
		opacity: 0.6;
		cursor: wait;
	}

	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: #fff;
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* ----- Listing cards --------------------------------------------------- */
	.listings {
		display: flex;
		flex-direction: column;
		gap: 14px;
		width: 100%;
		max-width: 720px;
		padding-bottom: var(--s-5);
	}

	.listing-card {
		display: flex;
		gap: 16px;
		padding: 14px;
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius);
		cursor: pointer;
		-webkit-backdrop-filter: blur(12px);
		backdrop-filter: blur(12px);
		transition: border-color 0.2s ease, background 0.2s ease, transform 0.15s ease,
			box-shadow 0.2s ease;
		--wails-draggable: no-drag;
	}

	.listing-card:hover {
		border-color: rgba(255, 77, 77, 0.5);
		background: var(--glass-strong);
		box-shadow: 0 10px 28px rgba(255, 45, 114, 0.14);
		transform: translateY(-2px);
	}

	.thumb {
		position: relative;
		flex: 0 0 110px;
		height: 110px;
		border-radius: 9px;
		overflow: hidden;
		background: rgba(0, 0, 0, 0.35);
	}

	.thumb img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.thumb-placeholder {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		color: var(--muted);
	}

	.thumb-placeholder svg {
		width: 36px;
		height: 36px;
		opacity: 0.5;
	}

	.arb-badge {
		position: absolute;
		top: 6px;
		left: 6px;
		padding: 3px 8px;
		font-size: 12px;
		font-weight: 800;
		border-radius: 6px;
		letter-spacing: -0.01em;
		backdrop-filter: blur(8px);
	}

	.arb-great {
		background: rgba(34, 197, 94, 0.85);
		color: #fff;
	}
	.arb-good {
		background: rgba(74, 222, 128, 0.75);
		color: #052e16;
	}
	.arb-ok {
		background: rgba(250, 204, 21, 0.8);
		color: #422006;
	}
	.arb-bad {
		background: rgba(239, 68, 68, 0.8);
		color: #fff;
	}

	.card-body {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.card-top {
		display: flex;
		align-items: flex-start;
		gap: 8px;
	}

	.card-title {
		margin: 0;
		font-size: 15.5px;
		font-weight: 700;
		line-height: 1.3;
		letter-spacing: -0.015em;
		color: var(--text);
		display: -webkit-box;
		line-clamp: 2;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.pro-badge {
		flex-shrink: 0;
		padding: 2px 7px;
		font-size: 10px;
		font-weight: 800;
		letter-spacing: 0.06em;
		color: #ff7a8f;
		background: rgba(255, 45, 114, 0.16);
		border: 1px solid rgba(255, 45, 114, 0.35);
		border-radius: 999px;
	}

	.meta-row {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 12.5px;
		color: var(--muted);
	}

	.seller {
		font-weight: 600;
		color: #b8c0d4;
	}

	.dot-sep {
		opacity: 0.45;
	}

	.price-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 6px;
		margin-top: 2px;
	}

	.price-block {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.price-label {
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--muted);
	}

	.price-value {
		font-size: 13.5px;
		font-weight: 700;
		color: var(--text);
		letter-spacing: -0.02em;
	}

	.price-value.melt {
		color: #4ade80;
	}

	.confidence-bar {
		position: relative;
		height: 4px;
		background: rgba(255, 255, 255, 0.08);
		border-radius: 999px;
		overflow: hidden;
		margin-top: 4px;
	}

	.confidence-fill {
		height: 100%;
		background: linear-gradient(90deg, #ff5a4d, #ff2d72);
		border-radius: 999px;
		transition: width 0.4s ease;
	}

	.confidence-label {
		position: absolute;
		right: 0;
		top: -16px;
		font-size: 10px;
		font-weight: 600;
		color: var(--muted);
	}

	.empty-state {
		padding: 48px 24px;
		text-align: center;
		color: var(--muted);
		font-size: 15px;
	}

	/* ----- Desktop: controls left, list right ------------------------------ */
	@media (min-width: 768px) {
		.content-layout {
			flex-direction: row;
			align-items: flex-start;
			gap: var(--s-4);
		}

		.controls {
			flex: 0 0 210px;
			max-width: 220px;
			flex-direction: column;
			align-items: stretch;
			flex-wrap: nowrap;
			position: sticky;
			top: max(var(--s-3), env(safe-area-inset-top));
			align-self: flex-start;
		}

		.arb-slider {
			width: 100%;
		}

		.refresh-btn {
			margin-left: 0;
			width: 100%;
			justify-content: center;
		}

		.listings {
			flex: 1 1 0;
			min-width: 0;
			max-width: none;
		}
	}

	/* ----- Mobile: stacked + sticky controls header ------------------------ */
	@media (max-width: 767px) {
		.controls {
			flex-direction: column;
			align-items: stretch;
			position: sticky;
			top: 0;
			z-index: 15;
			background: rgba(6, 7, 15, 0.84);
			-webkit-backdrop-filter: blur(16px);
			backdrop-filter: blur(16px);
			padding: 10px 0 14px;
			margin-bottom: var(--s-3);
			border-bottom: 1px solid var(--glass-border);
			max-width: none;
		}

		.refresh-btn {
			margin-left: 0;
			align-self: flex-start;
		}

		.listing-card {
			flex-direction: column;
			gap: 12px;
		}

		.thumb {
			flex: none;
			width: 100%;
			height: 180px;
		}

		.price-grid {
			grid-template-columns: repeat(2, 1fr);
			gap: 10px;
		}
	}
</style>