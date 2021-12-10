<script lang="ts">
  import { torcstatus, machinfo, machstats, Send, fileSize, hasMachinfo, isAdmin } from './core';
  import { onMount } from 'svelte';
  import slocation from 'slocation';

  let deviceinfoOpen = false;
  let devicestatsOpen = false;
  let torcstatsOpen = false;

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
    Send({ command: 'machstats', aop: 1 });
    Send({ command: 'torcstatus', aop: 1 });
  });

  let devinfoaction = () => {
    deviceinfoOpen = !deviceinfoOpen;
    if ($hasMachinfo === false) {
      Send({ command: 'machinfo', aop: 1 });
    }
  };

  let devicestatsaction = () => {
    if (devicestatsOpen === false) {
      Send({ command: 'machstats', aop: 1 });
      devicestatsOpen = true;
    } else {
      devicestatsOpen = false;
    }
  };

  let torcstatsaction = () => {
    if (torcstatsOpen === false) {
      Send({ command: 'torcstatus', aop: 1 });
      torcstatsOpen = true;
    } else {
      torcstatsOpen = false;
    }
  };
</script>

<div class="mx-auto max-w-xl">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={devinfoaction}>
        <p class="ml-3 font-medium  truncate">Device Info</p>
      </div>

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={devinfoaction}>
        {#if deviceinfoOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col">
      {#if deviceinfoOpen === true}
        {#if $hasMachinfo === true}
          <div class="m-1 p-1 break-all">Arch: {$machinfo?.arch}</div>
          <div class="m-1 p-1 break-all">CPU Model: {$machinfo?.cpumodel}</div>
          <div class="m-1 p-1 break-all">Go Version: {$machinfo?.goversion}</div>
          <div class="m-1 p-1 break-all">Hostname: {$machinfo?.hostname}</div>
          <div class="m-1 p-1 break-all">CPU No: {$machinfo?.numbercpu}</div>
          <div class="m-1 p-1 break-all">OS: {$machinfo?.os}</div>
          <div class="m-1 p-1 break-all">Platform: {$machinfo?.platform}</div>
          <div class="m-1 p-1 break-all">Started at {new Date($machinfo?.startedat)?.toLocaleString()}</div>
          <div class="m-1 p-1 break-all">Total Mem: {fileSize($machinfo?.totalmem)}</div>
        {/if}
      {/if}
    </div>
  </div>
</div>

<div class="mx-auto max-w-xl">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={devicestatsaction}>
        <p class="ml-3 font-medium  truncate">Device Stats</p>
      </div>

      {#if devicestatsOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({ command: 'machstats', aop: 1 });
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={devicestatsaction}>
        {#if devicestatsOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col">
      {#if devicestatsOpen === true}
        <div class="m-1 p-1 break-all">CPU Cycles: {$machstats?.cpucycles}</div>
        <div class="m-1 p-1 break-all">Disk Free: {fileSize($machstats?.diskfree)}</div>
        <div class="m-1 p-1 break-all">Disk Free(Bytes): {$machstats?.diskfree}</div>
        <div class="m-1 p-1 break-all">Disk Percent: {$machstats?.diskpercent} %</div>
        <div class="m-1 p-1 break-all">Memory Percent: {$machstats?.mempercent} %</div>
        <div class="m-1 p-1 break-all">Go Mem: {fileSize($machstats?.gomem)}</div>
        <div class="m-1 p-1 break-all">Go Mem(sys): {fileSize($machstats?.gomemsys)}</div>
        <div class="m-1 p-1 break-all">Goroutines: {$machstats?.goroutines}</div>
      {/if}
    </div>
  </div>
</div>

<div class="mx-auto max-w-xl">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={torcstatsaction}>
        <p class="ml-3 font-medium  truncate">Torrent Client Status</p>
      </div>

      {#if torcstatsOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({ command: 'torcstatus', aop: 1 });
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={torcstatsaction}>
        {#if torcstatsOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col overflow-x-auto bg-neutral-800 rounded-md">
      {#if torcstatsOpen === true}
        <pre class="whitespace-pre mb-2 p-2">{$torcstatus} </pre>
      {/if}
    </div>
  </div>
</div>
