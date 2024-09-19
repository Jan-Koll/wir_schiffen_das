<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label/index.js';
	import {clientKey, messages} from '$lib/state/sse'
	let keyOutgoing = '';
	let messageOutgoing = '';

	async function submit() {
		const body = JSON.stringify({ message: messageOutgoing });
		messageOutgoing = '';
		let url = `/analyze/sse/${keyOutgoing}/write`;
		if (keyOutgoing === '') {
			url = `/analyze/sse/all/write`;
		}
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body
		});
		if (response.status >= 300) {
			console.log('Failed submitting message.', response);
		}
	}
	</script>

<!--<div class="flex gap-1 space-x-2">-->
<div class="flex w-full max-w-sm flex-col gap-1.5">
	<Label for="key" class="font-bold">Your key</Label>
	<Input
		class="text-muted-foreground w-[500px]"
		on:click={(e) => e.target.select()}
		type="text"
		readonly
		id="key"
		value={$clientKey}
	/>
</div>
<div class="flex w-full max-w-sm flex-col gap-1.5">
	<Label class="font-bold" for="message">Recipient key</Label>
	<Input
		class="w-[500px]"
		type="text"
		on:click={(e) => e.target.select()}
		placeholder="Recipient"
		id="message"
		bind:value={keyOutgoing}
	/>
	<span class="text-muted-foreground text-sm">If empty, message will be sent to everyone.</span>
</div>
<div class="flex w-full max-w-sm flex-col gap-1.5">
	<Label class="font-bold" for="message">Message</Label>
	<Input
		class="w-[500px]"
		type="text"
		placeholder="Yourr message"
		id="message"
		bind:value={messageOutgoing}
	/>
</div>
<Button class="w-[300px] border border-white " on:click={submit}>Send message</Button>
<hr />
<Label for="recievedMessage" class="font-bold">Last recieved message</Label>
<code id="recievedMessage" class="bg-gray-800 px-4 py-2 text-white">
	{#if !$messages}
		Updates will appear here.
	{:else}
		{$messages}
	{/if}
</code>
