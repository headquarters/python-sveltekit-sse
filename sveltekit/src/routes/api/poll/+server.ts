import type { RequestHandler } from './$types';

const PYTHON_SSE_URL = 'http://localhost:8000/poll';

export const GET: RequestHandler = async ({ request }) => {
	// Create a ReadableStream to forward SSE events
	const stream = new ReadableStream({
		async start(controller) {
			const encoder = new TextEncoder();
			let abortController: AbortController | null = null;

			try {
				// Subscribe to the Python SSE endpoint
				abortController = new AbortController();
				const response = await fetch(PYTHON_SSE_URL, {
					signal: abortController.signal,
					headers: {
						Accept: 'text/event-stream',
						'Cache-Control': 'no-cache'
					}
				});

				if (!response.ok) {
					throw new Error(`Python SSE endpoint returned ${response.status}`);
				}

				const reader = response.body?.getReader();
				if (!reader) {
					throw new Error('No response body');
				}

				const decoder = new TextDecoder();

				// Read from Python SSE and forward to client
				while (true) {
					const { done, value } = await reader.read();

					if (done) {
						console.log('Python SSE stream ended');
						break;
					}

					// Decode the chunk and forward it
					const chunk = decoder.decode(value, { stream: true });
					controller.enqueue(encoder.encode(chunk));
				}
			} catch (error) {
				if (error instanceof Error && error.name === 'AbortError') {
					console.log('SSE connection aborted by client');
				} else {
					console.error('Error in SSE forwarding:', error);
				}
			} finally {
				controller.close();
				abortController?.abort();
			}
		},

		cancel() {
			console.log('Client disconnected from SvelteKit SSE endpoint');
		}
	});

	return new Response(stream, {
		headers: {
			'Content-Type': 'text/event-stream',
			'Cache-Control': 'no-cache',
			Connection: 'keep-alive'
		}
	});
};
