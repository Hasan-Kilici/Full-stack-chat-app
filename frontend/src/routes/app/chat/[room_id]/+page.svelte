<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';

    import { getMessages } from '$lib/utils';
    import Sidebar from '$lib/components/Sidebar.svelte';

    export let data: {
        user?: {
            email: string;
            name: string;
            user_id: string;
            auth: string;
        };
        rooms: any;
    };

    let socket: WebSocket | null = null;
    let messages: any[] = [];
    let msg: string = '';

    onMount(async () => {
        if (data.user && $page.params.room_id) {
            messages = await getMessages(data.user.auth, $page.params.room_id);
            socket = new WebSocket("ws://localhost:5000/ws", data.user.auth);

            socket.onopen = () => {
                socket?.send(JSON.stringify({
                    event: "join",
                    user: data.user?.user_id,
                    room: $page.params.room_id
                }));
            };

            let nextId = 0;

            socket.onmessage = (e: MessageEvent) => {
                const incoming = JSON.parse(e.data);
                if (incoming.event === "message") {
                    console.log(incoming);
                    if (!incoming.id) {
                        incoming.id = `msg-${nextId++}`;
                    }
                    messages = [...messages, incoming];
                }
            };
        }
    });

    function sendMessage() {
        if (data.user && socket && msg.trim()) {
            socket.send(JSON.stringify({
                event: "message",
                user: data.user.auth,
                room: $page.params.room_id,
                data: msg,
            }));
            msg = '';
        }
    }
</script>

{#if data.user}
<div class="flex max-h-[100vh] bg-[#f9fafb] text-gray-900">
  <Sidebar auth={data.user.auth} username={data.user.name} rooms={data.rooms} />

  <main class="flex flex-col flex-1 p-6 max-h-[100vh]">
    <div class="flex-1 flex-col-reverse overflow-y-auto space-y-4 mb-4">
      {#if messages.length === 0}
        <p class="text-center text-gray-500 mt-20">No messages yet.</p>
      {:else}
        {#each messages as message (message.id)}
          <div
            class="flex items-start gap-3 max-w-[80%] mx-auto bg-white rounded-lg shadow-md p-4"
          >
            <img
              src={"https://api.dicebear.com/9.x/bottts/svg?seed=" + message.author}
              alt={message.author}
              class="w-10 h-10 rounded-full border border-gray-300"
              loading="lazy"
            />
            <div>
              <div class="flex items-center gap-2 mb-1">
                <span class="font-semibold">{message.author}</span>
                <span class="text-xs text-gray-400">{new Date(message.timestamp || Date.now()).toLocaleTimeString()}</span>
              </div>
              <p class="whitespace-pre-wrap">{message.data}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>

    <form class="flex gap-3" on:submit|preventDefault={sendMessage}>
      <input
        type="text"
        placeholder="Write your message..."
        bind:value={msg}
        class="flex-1 rounded-md border border-gray-300 bg-white px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
        autocomplete="off"
        spellcheck="false"
      />
      <button
        type="submit"
        class="inline-flex items-center justify-center rounded-md bg-blue-600 hover:bg-blue-700 text-white font-semibold px-5 py-2 shadow-sm transition"
        disabled={!msg.trim()}
      >
        Send
      </button>
    </form>
  </main>
</div>
{:else}
  <p class="text-center text-gray-500 mt-20">No user data or still loading...</p>
{/if}
