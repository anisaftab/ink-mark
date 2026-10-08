<script lang="ts">
	import { Greet } from '#lib/wailsjs/go/main/App';

	let name = $state('');
	let result = $state('');

	async function greet() {
		result = await Greet(name);
	}
</script>

<main class="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-center gap-6">
	<h1 class="text-4xl font-bold tracking-tight">🖊️ ink-mark</h1>

	<div class="flex gap-2">
		<input
			bind:value={name}
			placeholder="Enter your name..."
			class="px-4 py-2 rounded-lg bg-gray-800 border border-gray-600 focus:outline-none focus:border-blue-500 text-white placeholder-gray-400"
		/>
		<button
			onclick={greet}
			class="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg font-medium transition-colors"
		>
			Greet
		</button>
	</div>

	{#if result}
		<p class="text-xl text-green-400">{result}</p>
	{/if}
</main>

