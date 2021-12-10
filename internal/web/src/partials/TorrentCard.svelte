<script lang="ts">
  import { fileSize, Send, adminmode } from './core';
  import slocation from 'slocation';

  export let state: string = '';
  export let name: string | undefined = '';
  export let infohash: string = '';
  export let bytescompleted: number | undefined = 0;
  export let bytesmissing: number | undefined = 0;
  export let length: number | undefined = 0;
  export let seeding: boolean | undefined = false;
  export let isTorrentPage: boolean = false;
  export let locked: boolean = false;

  let progpercentage = 0;

  let refresh = () => {
    if (isTorrentPage === false) {
      Send({
        command: 'listtorrents'
      });
    } else if (isTorrentPage === true) {
      Send({
        command: 'listtorrentinfo',
        data1: infohash
      });
    }
  };

  $: progpercentage = (bytescompleted / length) * 100;
</script>

<div class="bg-black text-neutral-200 rounded-lg m-3">
  <div class="px-3 py-5">
    <div class="flex items-center space-x-3.5 sm:space-x-5 lg:space-x-3.5 xl:space-x-5">
      <div class="min-w-0 flex-auto space-y-0.5">
        <div class="text-neutral-400 flex justify-between text-sm font-medium tabular-nums">
          <div class="truncate">{infohash}</div>
          <div>
            {#if state === 'active' || state === 'inactive'}{#if seeding === true}(Seeding)
              {/if}{fileSize(bytesmissing == null ? 0 : bytesmissing)} R.{/if}
          </div>
        </div>
        {#if state !== 'removed'}
          <h2
            class="text-white text-base sm:text-xl lg:text-base xl:text-xl font-semibold break-all"
            on:click={() => {
              if (isTorrentPage === false) {
                slocation.goto(`/torrent/${infohash}`);
              }
            }}>
            {name}
          </h2>
        {/if}
      </div>
    </div>

    {#if state === 'active' || state === 'inactive'}
      <div class="space-y-2 mt-1">
        <div class="bg-zinc-900 rounded-full overflow-hidden">
          <div class="bg-blue-700  h-1.5" style="width:{progpercentage ? progpercentage : 0}%" />
        </div>
        <div class="text-neutral-500 dark:text-neutral-400 flex justify-between text-sm font-medium tabular-nums">
          <div>{fileSize(bytescompleted)} / {fileSize(length)}</div>
          <div>
            {progpercentage?.toLocaleString('en-US', {
              maximumFractionDigits: 2,
              minimumFractionDigits: 2
            })} %
          </div>
        </div>
      </div>
    {/if}
  </div>

  <div class="bg-gradient-to-r from-black via-neutral-900 to-black text-white py-2 px-1 grid grid-flow-col grid-cols-3 justify-items-center mt-1 rounded-b-lg">
    {#if state === 'active' || state === 'loading' || state === 'inactive'}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'removetorrent',
            data1: infohash,
            ...($adminmode === true && { aop: 1 })
          });
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-neutral-500 text-sm">Remove</p>
      </button>
    {:else if state === 'removed'}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'deletetorrent',
            data1: infohash,
            ...($adminmode === true && { aop: 1 })
          });
          if (isTorrentPage === true) {
            slocation.goto('/');
          } else if (isTorrentPage === false) {
            refresh();
            Send({ command: 'gettorrents' });
          }
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
        <p class="text-neutral-500 text-sm">Delete</p>
      </button>
    {/if}

    {#if state === 'loading'}
      <button type="button" class="grid">
        <svg class="animate-spin h-6 w-6 mx-auto text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        <p class="text-neutral-500 text-sm">Loading</p>
      </button>
    {:else if state === 'active'}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'stoptorrent',
            data1: infohash,
            ...($adminmode === true && { aop: 1 })
          });
          refresh();
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-neutral-500 text-sm">Stop</p>
      </button>
    {:else if state === 'removed'}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'addinfohash',
            data1: infohash
          });
          refresh();
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-neutral-500 text-sm">Add</p>
      </button>
    {:else if state === 'inactive'}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'starttorrent',
            data1: infohash,
            ...($adminmode === true && { aop: 1 })
          });
          refresh();
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-neutral-500 text-sm">Start</p>
      </button>
    {/if}

    {#if isTorrentPage === false}
      <button
        type="button"
        class="grid"
        on:click={() => {
          slocation.goto(`/torrent/${infohash}`);
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
        <p class="text-neutral-500 text-sm">View</p>
      </button>
    {:else}
      <button
        type="button"
        class="grid"
        on:click={() => {
          Send({
            command: 'toggletorrentlock',
            data1: infohash
          });
          setTimeout(() => {
            Send({ command: 'istorrentlocked', data1: infohash });
          }, 1000);
        }}>
        {#if locked === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          <p class="text-neutral-500 text-sm">Locked</p>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" />
          </svg>
          <p class="text-neutral-500 text-sm">Unlocked</p>
        {/if}
      </button>
    {/if}
  </div>
</div>
