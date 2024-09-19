<script lang="ts">
	import CircleCheck from 'lucide-svelte/icons/circle-check';
	import CircleX from 'lucide-svelte/icons/circle-x';
	import Circle from 'lucide-svelte/icons/circle';
	import CircleEllpsis from 'lucide-svelte/icons/circle-ellipsis';
	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import CircleSlash from 'lucide-svelte/icons/circle-slash';
	import CircleDot from 'lucide-svelte/icons/circle-dot';
	import * as Tooltip from '$lib/components/ui/tooltip';
	export let status: string | undefined;
	export let success: boolean = false;
	export let side: "top" | "right" | "bottom" | "left" | undefined = "right"
	export let name: string;
	//export let small: boolean = false;
</script>

<Tooltip.Root openDelay={0}>
	<Tooltip.Trigger>
		{#if status == 'queued'}
			<CircleDot class="rounded-full text-gray-500"></CircleDot>
		{:else if status == 'running'}
			<CircleEllpsis class="animate-pulse rounded-full text-sky-500"></CircleEllpsis>
		{:else if status == 'canceled'}
			<CircleSlash class="rounded-full text-yellow-600"></CircleSlash>
		{:else if status == 'failed'}
			<CircleAlert class="rounded-full text-red-600"></CircleAlert>
		{:else if status == 'ready'}
			{#if success}
				<CircleCheck class="rounded-full text-emerald-400"></CircleCheck>
			{:else}
				<CircleX class="rounded-full text-red-400"></CircleX>
			{/if}
		{:else if status === undefined}
			<Circle class="rounded-full dark:text-gray-700 text-gray-300"></Circle>
		{/if}
	</Tooltip.Trigger>
	<Tooltip.Content {side}>
		<div class="mb-3 font-bold">{name}</div>
		<div class="flex flex-col">
			<div><span class="font-bold">Status:</span> 
				{#if status}
					{status}
				{:else}
					Not triggered yet...
				{/if}
			</div>
			{#if status == 'ready'}
				<div><span class="font-bold">Success:</span> {success}</div>
			{/if}
		</div>
	</Tooltip.Content>
</Tooltip.Root>
