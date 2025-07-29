<script lang="ts">
	import { searchUsers } from '$lib/search';
	import { onDestroy } from 'svelte';

	import {createRoom} from '$lib/utils';

	export let auth: string;
    export let username: string;
    export let rooms: any;

	let search = '';
	let timeout: ReturnType<typeof window.setTimeout>;
	let results: any = [];

	$: {
		if (search.trim() === '') {
			results = [];
			clearTimeout(timeout);
		} else {
			clearTimeout(timeout);
			timeout = window.setTimeout(async () => {
				const allResults = await searchUsers(search, auth);
	            results = allResults.filter((user: any) => user.name !== username);
			}, 1000);
		}
	}

	onDestroy(() => {
		clearTimeout(timeout);
	});
</script>
<!-- Sidebar -->
<div class="p-4 w-[18%] bg-white border-r border-gray-300 h-screen max-h-screen overflow-auto space-y-6">
	<!-- Header -->
	<div>
		<h2 class="text-lg font-semibold text-gray-900">Chat App</h2>
	</div>

	<!-- Explore Users Section -->
	<div class="space-y-3">
		<h3 class="text-sm font-medium text-muted-foreground">Explore Users</h3>
		<input
			bind:value={search}
			placeholder="Search users..."
			class="w-full px-3 py-2 text-sm border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
		/>

		{#if results.length}
			<div class="flex flex-col gap-2">
				{#each results as user}
					<button on:click={()=>{
						createRoom(auth, user.id)
					}} class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-100 transition">
						<img
							src={`https://api.dicebear.com/9.x/bottts/svg?seed=${user.name}`}
							alt={user.name}
							class="w-6 h-6 rounded-full"
						/>
						<span class="text-sm text-gray-800">{user.name}</span>
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Chat Rooms Section -->
	{#if rooms.length}
		<div class="space-y-3 pt-4">
			<h3 class="text-sm font-medium text-muted-foreground">Your Chats</h3>
			<div class="flex flex-col gap-2">
				{#each rooms as room}
					<a href="/app/chat/{room.room_id}" class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-100 transition cursor-pointer">
						<img
							src={`https://api.dicebear.com/9.x/bottts/svg?seed=${room.room_name}`}
							alt={room.room_name}
							class="w-6 h-6 rounded-full"
						/>
						<span class="text-sm text-gray-800">{room.room_name}</span>
                    </a>
				{/each}
			</div>
		</div>
	{/if}
</div>
