<script lang="ts">
  import Sidebar from '$lib/components/Sidebar.svelte'
  export let data: {
    user?: {
      email: string;
      name: string;
      user_id: string;
      auth: string;
    };
    rooms:any
  };

</script>
{#if data.user}

  <Sidebar auth={data.user?.auth} username={data.user?.name} rooms={data.rooms}/>

{/if}