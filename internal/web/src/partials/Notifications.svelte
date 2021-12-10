<script lang="ts">
  import slocation from 'slocation';
  import { resplist } from './core';
</script>

<div class="mx-auto max-w-3xl ">
  {#if $resplist?.has === true}
    <button
      type="button"
      class="w-full my-2 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-zinc-600 focus:outline-none"
      on:click={() => {
        resplist.set({ has: false, data: [] });
      }}>
      Clear All
    </button>
    {#each $resplist?.data as resp}
      <div
        class="bg-neutral-800 text-neutral-200 px-3 py-5 rounded-lg m-3 noHL"
        on:click={() => {
          if (typeof resp?.infohash === 'string' && resp?.infohash.length > 0) {
            slocation.goto(`/torrent/${resp?.infohash}`);
          }
        }}>
        <div class="break-all mx-1 mb-1 font-bold">
          {resp?.message}
        </div>
        <div class="text-neutral-300 flex justify-between text-sm font-medium tabular-nums">
          <div class="break-all mx-1">
            {resp?.type}
            {resp?.state}
            {resp?.infohash ? `( ${resp?.infohash} )` : ''}
          </div>
        </div>
      </div>
    {/each}
  {:else}
    <p class="text-xl text-center text-red-400 font-sans">No Notifications</p>
  {/if}
</div>
