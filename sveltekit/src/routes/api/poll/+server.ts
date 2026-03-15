import type { RequestHandler } from './$types';

const PYTHON_SSE_URL = 'http://localhost:8000/poll';

export const GET: RequestHandler = async () => {
	const abortController = new AbortController();
	const encoder = new TextEncoder();
	const decoder = new TextDecoder();

	const stream = new ReadableStream({
		async start(controller) {
			try {
				const response = await fetch(PYTHON_SSE_URL, {
					signal: abortController.signal,
					headers: { Accept: 'text/event-stream', 'Cache-Control': 'no-cache' }
				});

				if (!response.ok) throw new Error(`Python SSE endpoint returned ${response.status}`);
				if (!response.body) throw new Error('No response body');

				for await (const chunk of response.body as unknown as AsyncIterable<Uint8Array>) {
					controller.enqueue(encoder.encode(decoder.decode(chunk, { stream: true })));
				}
				console.log('Python SSE stream ended');
			} catch (error) {
				if (error instanceof Error && error.name === 'AbortError') {
					console.log('SSE connection aborted by client');
				} else {
					console.error('Error in SSE forwarding:', error);
				}
			} finally {
				controller.close();
			}
		},

		cancel() {
			console.log('Client disconnected from SvelteKit SSE endpoint');
			abortController.abort();
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
