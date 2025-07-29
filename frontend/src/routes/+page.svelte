<script lang="ts">
	import { redirect } from "@sveltejs/kit";
	import { onMount } from "svelte";

    onMount(()=>{
        redirect(302, "/register")
    })
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
