<script lang="ts">
  import slocation from 'slocation';
  import { onMount } from 'svelte';
  import { torrentsforuser, Send, adminmode, isAdmin } from './core';

  let username = $slocation.pathname?.split('/').reverse()[0];

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
    torrentsforuser.set([]);
    adminmode.set(true);
    Send({ command: 'listtorrentsforuser', data1: username, aop: 1 });
  });
</script>

<div class="mx-auto max-w-3xl ">
  {#if Array.isArray($torrentsforuser) && $torrentsforuser?.length}
    {#each $torrentsforuser as trnt (trnt)}
      <div
        class="bg-neutral-800 text-neutral-200 px-3 py-5 rounded-lg m-3 noHL"
        on:click={() => {
          if (typeof trnt === 'string' && trnt?.length > 0) {
            slocation.goto(`/torrent/${trnt}`);
          }
        }}>
        <div class="break-all mx-1 mb-1 font-bold">
          {trnt}
        </div>
      </div>
    {/each}
  {:else}
    <p class="text-xl text-center text-red-400 font-sans">User owns no Torrents</p>
  {/if}
</div>
