<script lang="ts">
  import { Send, fileviewpath, fileviewinfohash, fsfileinfo, fileType, socket, fileSize, filepagediscon } from './core';
  import slocation from 'slocation';
  import { onDestroy, onMount } from 'svelte';

  let stream = false;
  let ft = 'unknown';

  onMount(() => {
    if (socket == null || socket == undefined || socket?.readyState === WebSocket.CLOSED) {
      slocation.goto('/');
    }
    Send({
      command: 'getfsfileinfo',
      data1: $fileviewinfohash,
      data2: $fileviewpath
    });
    filepagediscon.set(true);
    ft = fileType($fileviewpath);
  });

  onDestroy(() => {
    filepagediscon.set(false);
    document.title = 'exatorrent';
  });
</script>

<svelte:head>
  <title>{$fsfileinfo?.name}</title>
</svelte:head>

<div class="mx-auto max-w-xl">
  <div class="flex justify-between h-16">
    <button
      class="flex-shrink-0 focus:outline-none  bg-neutral-800 m-2 px-3 rounded-md"
      on:click={() => {
        history.length > 2 ? history.back() : slocation.goto('/');
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>

    <button
      class="flex items-center focus:outline-none p-5 mt-2"
      on:click={() => {
        slocation.goto('/');
      }}>
      <p class="font-sans text-2xl text-neutral-200">exatorrent</p>
    </button>

    <div class="flex-shrink-0 focus:outline-none m-2 px-3 rounded-md" />
  </div>

  <div class="flex justify-center">
    <!-- svelte-ignore a11y-media-has-caption -->
    {#if ft === 'video'}
      <video controls src="/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}" />
    {:else if ft === 'audio'}
      <audio controls src="/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}" />
    {/if}
  </div>

  <p class="text-center text-neutral-400 font-sans break-all">{$fsfileinfo?.name}</p>
  <p class="text-center text-neutral-400 font-sans truncate">({fileSize($fsfileinfo?.size)})</p>

  <div class="grid grid-flow-col grid-cols-4 bg-black my-2 appearance-none border border-neutral-800 w-full rounded-md">
    <div class="col-span-3 appearance-none  w-full flex-grow px-3 py-2  border-none  text-neutral-300  focus:outline-none mx-1">Use Stream API</div>
    <div class="flex items-center justify-end w-full my-2 mr-2">
      <label for="dontstarttoggle" class="flex items-center cursor-pointer mx-1">
        <div class="relative">
          <input type="checkbox" class="rounded text-indigo-700 bg-neutral-800 form-checkbox mx-1" bind:checked={stream} />
        </div>
      </label>
    </div>
  </div>

  {#if ft === 'video' || ft === 'audio'}
    <a href="vlc://{location.origin}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}?token={localStorage.getItem('exasession')}" target="_blank" rel="noopener noreferrer">
      <button class="bg-yellow-700 w-full my-2 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white focus:outline-none"> Play in VLC </button>
    </a>

    <a href="intent://{location.host}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}?token={localStorage.getItem('exasession')}#Intent;type=video/any;package=is.xyz.mpv;scheme={location.protocol.slice(0, -1)};end;" target="_blank" rel="noopener noreferrer">
      <button class="flex md:hidden  bg-purple-700 w-full my-2 justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white focus:outline-none"> Play in MPV </button>
    </a>
  {/if}

  <input type="text" class="bg-zinc-700 rounded-md  w-full my-2 py-2 px-4 text-sm text-neutral-300 truncate" disabled value="{location.origin}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}" />

  <a href="{location.origin}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}?token={localStorage.getItem('exasession')}" target="_blank" rel="noopener noreferrer" download>
    <button type="button" class="w-full my-3 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-900 focus:outline-none"> Download </button>
  </a>
</div>
