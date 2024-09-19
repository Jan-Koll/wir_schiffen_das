<script lang="ts">
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Button } from '$lib/components/ui/button';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import Square from 'lucide-svelte/icons/square';
	import Repeat2 from 'lucide-svelte/icons/repeat-2';
	import ListRestart from 'lucide-svelte/icons/list-restart';
	import ExternalLink from 'lucide-svelte/icons/external-link';
	export let id: string;
</script>

 <Tooltip.Root>
	<Tooltip.Trigger asChild let:builder>
		<Button builders={[builder]} href={`/analyze/configuration/${id}`} size="icon" variant="secondary">
			<ExternalLink class="h-4 w-4" />
			<span class="sr-only">View</span>
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Open</p>
	</Tooltip.Content>
</Tooltip.Root> 
<Tooltip.Root>
	<Tooltip.Trigger asChild let:builder>
		<Button builders={[builder]} size="icon" variant="ghost">
			<Repeat2 class="h-4 w-4" />
			<span class="sr-only">Retry</span>
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Retry</p>
	</Tooltip.Content>
</Tooltip.Root>
<!-- <Tooltip.Root>
	<Tooltip.Trigger asChild let:builder>
		<Button builders={[builder]} size="icon" variant="ghost">
			<Square class="h-4 w-4" />
			<span class="sr-only">Cancel</span>
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Cancel</p>
	</Tooltip.Content>
</Tooltip.Root> -->
<Tooltip.Root>
	<Tooltip.Trigger asChild let:builder>
		<Button builders={[builder]} size="icon" variant="ghost">
			<ListRestart class="h-4 w-4" />
			<span class="sr-only">Recreate</span>
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Recreate</p>
	</Tooltip.Content>
</Tooltip.Root>
<!-- <DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Button variant="ghost" builders={[builder]} size="icon" class="relative h-8 w-8 p-0">
			<span class="sr-only">Open menu</span>
			<Ellipsis class="h-4 w-4" />
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Group>
			<DropdownMenu.Label>Actions</DropdownMenu.Label>
			<DropdownMenu.Item on:click={() => navigator.clipboard.writeText(id)}>
				Copy payment ID
			</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<DropdownMenu.Item>View customer</DropdownMenu.Item>
		<DropdownMenu.Item>View payment details</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root> -->
