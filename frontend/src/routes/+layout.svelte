<script lang="ts">
	import '../app.css';
	import { ModeWatcher } from 'mode-watcher';
	import Theme_switch from '$lib/components/ui/Theme_switch.svelte';
	import CircleUser from 'lucide-svelte/icons/circle-user';
	import Menu from 'lucide-svelte/icons/menu';
	import Sailboat from 'lucide-svelte/icons/sailboat';
	import Search from 'lucide-svelte/icons/search';

	import { browser } from '$app/environment';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { page, navigating } from '$app/stores';
	import {capitalize} from '$lib/strings';
	import { Toaster } from '$lib/components/ui/sonner';
	import { toast } from 'svelte-sonner';
	import { messages, clientKey } from '$lib/state/sse';
	import { onMount } from 'svelte';
	import { source } from 'sveltekit-sse';
	import type { Readable } from 'svelte/store';
	// Trigger different toast types
	/*$: if ($messages) {
		if ($messages.startsWith('info:')) {
			toast.info("all:"+$messages);
		} else if ($messages.startsWith('success:')) {
			toast.success($messages);
		} else if ($messages.startsWith('error:')) {
			toast.error($messages);
		} else {
			toast($messages);
		}
	}*/

	let breadcrumbs: string[] = [];
	$: breadcrumbs = $page.url.pathname.slice(1).split('/');
	type navItem = { text: string; path: string };
	let navItems: navItem[] = [
		{ text: 'Home', path: '/' },
		{ text: 'Orders', path: '/orders' },
		{ text: 'Products', path: '/products' },
		{ text: 'Customers', path: '/customers' },
		{ text: 'Analyze', path: '/analyze' }
	];
	let open = false;
	$: if ($navigating) open = false;


	// Pass messages to store to make them globally accessible without props
	let messageListener: Readable<string> | null;
	$: if ($messageListener) {
		$messages = $messageListener;
	}

	onMount(async () => {
		if (browser) {
			const getKeyResponse = await fetch('/analyze/sse/request-key', {
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/json'
					// Authorization: `Bearer ${token}`
				}
			});

			if (getKeyResponse.status >= 300) {
				console.log(getKeyResponse.status, 'Failed reading stream.');
			} else {
				const { key }: { key: string } = await getKeyResponse.json();
				clientKey.set(key);
				messageListener = source(`/analyze/sse/${$clientKey}/read`).select('message');
				console.log('Reading stream successfully.');
			}
		}
	});
</script>

<svelte:head>
	<title>
		{$page?.url?.pathname != '/' ? `${capitalize($page?.url?.pathname?.slice(1))}` : 'Home'} | WirSchiffenDas
	</title>
</svelte:head>
<ModeWatcher />
<div class="flex min-h-screen w-full flex-col">
	<header class="sticky top-0 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
		<nav
			class="hidden flex-col gap-6 text-lg font-medium md:flex md:flex-row md:items-center md:gap-5 md:text-sm lg:gap-6"
		>
			<a href="/" class="flex items-center gap-2 text-lg font-semibold md:text-base">
				<Sailboat class="h-6 w-6" />
				<span class="hidden lg:block">WirSchiffenDas</span>
			</a>
			{#each navItems as navItem}
				<a
					href={navItem.path}
					class="{$page.url.pathname === navItem.path
						? 'text-foreground'
						: 'text-muted-foreground'} transition-colors hover:text-foreground"
				>
					{navItem.text}
				</a>
			{/each}
		</nav>
		<Sheet.Root bind:open>
			<Sheet.Trigger asChild let:builder>
				<Button variant="outline" size="icon" class="shrink-0 md:hidden" builders={[builder]}>
					<Menu class="h-5 w-5" />
					<span class="sr-only">Toggle navigation menu</span>
				</Button>
			</Sheet.Trigger>
			<Sheet.Content side="left">
				<nav class="grid gap-6 text-lg font-medium">
					<a href="/" class="flex items-center gap-2 text-lg font-semibold">
						<Sailboat class="h-6 w-6" />
						<span>WirSchiffenDas</span>
					</a>
					{#each navItems as navItem}
						<a
							href={navItem.path}
							class="{$page.url.pathname === navItem.path
								? 'text-foreground'
								: 'text-muted-foreground'} hover:text-foreground"
						>
							{navItem.text}
						</a>
					{/each}
				</nav>
			</Sheet.Content>
		</Sheet.Root>
		<div class="flex w-full items-center gap-4 md:ml-auto md:gap-2 lg:gap-4">
			<form class="ml-auto flex-1 sm:flex-initial">
				<div class="relative">
					<Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
					<Input
						type="search"
						placeholder="Search products..."
						class="pl-8 sm:w-[300px] md:w-[200px] lg:w-[300px]"
					/>
				</div>
			</form>
			<Theme_switch></Theme_switch>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="secondary" size="icon" class="rounded-full">
						<CircleUser class="h-5 w-5" />
						<span class="sr-only">Toggle user menu</span>
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Label>My Account</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item>Settings</DropdownMenu.Item>
					<DropdownMenu.Item>Support</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item>Logout</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</header>
	<main class="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8">
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/">Home</Breadcrumb.Link>
				</Breadcrumb.Item>
				{#if breadcrumbs[0] != ''}
					{#each breadcrumbs as breadcrumb, i}
						<Breadcrumb.Separator />
						<Breadcrumb.Item>
							<Breadcrumb.Link href="/{breadcrumbs.slice(0, i + 1).join('/')}"
								>{breadcrumb.charAt(0).toUpperCase() + breadcrumb.slice(1)}</Breadcrumb.Link
							>
						</Breadcrumb.Item>
					{/each}
				{/if}
			</Breadcrumb.List>
		</Breadcrumb.Root>
		<slot></slot>
		<Toaster expand={true} richColors></Toaster>
	</main>
</div>

<style></style>
