<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	interface DataPoint {
		value: number;
		timestamp: string;
	}

	let dataPoints = $state<DataPoint[]>([]);
	let status = $state<'connecting' | 'connected' | 'error' | 'closed'>('connecting');
	let eventSource: EventSource | null = null;

	function connect() {
		status = 'connecting';
		eventSource = new EventSource('/api/poll');

		eventSource.addEventListener('random_number', (event) => {
			const data = JSON.parse(event.data);
			dataPoints = [data, ...dataPoints];
		});

		eventSource.onopen = () => {
			status = 'connected';
		};

		eventSource.onerror = () => {
			status = 'error';
			console.error('SSE connection error');
		};
	}

	function closeConnection() {
		if (eventSource) {
			eventSource.close();
			status = 'closed';
			eventSource = null;
		}
	}

	onMount(() => {
		connect();
	});

	onDestroy(() => {
		closeConnection();
	});
</script>

<div class="container">
	<div class="header">
		<h1>Random Number Stream</h1>
		<button onclick={closeConnection} disabled={status === 'closed'} class="close-button">
			Close Connection
		</button>
	</div>

	<div class="status status-{status}">
		Status: {status}
	</div>

	<div class="data-list">
		{#if dataPoints.length === 0}
			<div class="waiting">Waiting for data...</div>
		{:else}
			{#each dataPoints as point (point.timestamp)}
				<div class="data-row">
					<div class="value">{point.value}</div>
					<div class="time">{new Date(point.timestamp).toLocaleTimeString()}</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.container {
		max-width: 600px;
		margin: 2rem auto;
		padding: 2rem;
		font-family: system-ui, -apple-system, sans-serif;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	h1 {
		margin: 0;
		color: #333;
	}

	.close-button {
		padding: 0.5rem 1rem;
		background-color: #dc2626;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: 500;
		transition: background-color 0.2s;
	}

	.close-button:hover:not(:disabled) {
		background-color: #b91c1c;
	}

	.close-button:disabled {
		background-color: #9ca3af;
		cursor: not-allowed;
	}

	.status {
		padding: 0.5rem 1rem;
		border-radius: 4px;
		margin-bottom: 2rem;
		font-weight: 500;
	}

	.status-connecting {
		background-color: #fef3c7;
		color: #92400e;
	}

	.status-connected {
		background-color: #d1fae5;
		color: #065f46;
	}

	.status-error {
		background-color: #fee2e2;
		color: #991b1b;
	}

	.status-closed {
		background-color: #e5e7eb;
		color: #374151;
	}

	.data-list {
		background-color: #f9fafb;
		border-radius: 8px;
		padding: 1rem;
		min-height: 200px;
		max-height: 500px;
		overflow-y: auto;
	}

	.data-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		margin-bottom: 0.5rem;
		background-color: white;
		border-radius: 4px;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.value {
		font-size: 2rem;
		font-weight: bold;
		color: #1f2937;
	}

	.time {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.waiting {
		text-align: center;
		padding: 3rem;
		color: #9ca3af;
		font-style: italic;
	}
</style>
