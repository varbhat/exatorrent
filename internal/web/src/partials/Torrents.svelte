<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { downloadslist, Send, adminmode, isAdmin, terrormsg } from './core';
  import TorrentCard from './TorrentCard.svelte';

  onMount(() => {
    if ($adminmode === false) {
      Send({ command: 'listtorrents' });
      Send({ command: 'gettorrents' });
    } else {
      Send({ command: 'listalltorrents', aop: 1 });
      Send({ command: 'getalltorrents', aop: 1 });
    }
  });

  onDestroy(() => {
    console.log('on destroy');
    Send({ command: 'stopstream' });
  });

  let checam = (b: boolean) => {
    b ? (Send({ command: 'listalltorrents', aop: 1 }), Send({ command: 'getalltorrents', aop: 1 })) : Send({ command: 'gettorrents' });
  };

  $: checam($adminmode);
</script>

<div class="mx-auto max-w-3xl ">
  {#if $terrormsg?.has === false}
    {#if $isAdmin === true && $adminmode === true}
      <div class="grid grid-flow-col grid-cols-4 pr-2 bg-neutral-800 my-2 appearance-none border border-neutral-800 w-full">
        <div class=" bg-neutral-800  col-span-3 appearance-none  w-full flex-grow px-3 py-2  border-none  text-neutral-300  focus:outline-none">Admin Mode</div>
        <div class="flex items-center justify-end w-full my-2 mr-2">
          <label for="dontstarttoggle" class="flex items-center cursor-pointer">
            <div class="relative">
              <input type="checkbox" class="rounded text-indigo-700 bg-neutral-800 form-checkbox" bind:checked={$adminmode} />
            </div>
          </label>
        </div>
      </div>
    {/if}

    {#each $downloadslist?.data as dl (dl.infohash)}
      <TorrentCard state={dl.state} name={dl?.name} infohash={dl.infohash} bytescompleted={dl?.bytescompleted} bytesmissing={dl?.bytesmissing} length={dl?.length} seeding={dl?.seeding} />
    {/each}
  {:else}
    <p class="text-xl text-center text-red-400 font-sans">{$terrormsg?.msg}</p>
  {/if}
</div>
