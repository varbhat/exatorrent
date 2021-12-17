<script lang="ts">
  import { dontstart, Send, isAdmin } from './core';
  import slocation from 'slocation';
  let ismetainfo = true;
  let torrentinput = '';
  let trntfilestring = '';
  let trntfileinput: HTMLInputElement;

  let addfunc = () => {
    if (ismetainfo === true) {
      if (torrentinput === '') {
        alert('Empty Input');
        return;
      }
      if (torrentinput.startsWith('magnet:')) {
        console.log('Adding Magnet', torrentinput, $dontstart);
        Send({
          command: 'addmagnet',
          data1: torrentinput,
          data2: $dontstart === 'true' ? 'true' : 'false'
        });
        torrentinput = '';
      } else {
        console.log('Adding Infohash', torrentinput, $dontstart);
        Send({
          command: 'addinfohash',
          data1: torrentinput,
          data2: $dontstart === 'true' ? 'true' : 'false'
        });
        torrentinput = '';
      }
    } else {
      if (trntfilestring === '') {
        alert('File Invalid');
        return;
      }
      console.log('Adding Torrent', trntfilestring, $dontstart);
      Send({
        command: 'addtorrent',
        data1: trntfilestring,
        data2: $dontstart === 'true' ? 'true' : 'false'
      });
      trntfilestring = '';
    }
  };

  let entertoadd = (event: KeyboardEvent) => {
    if (event.code === 'Enter') {
      addfunc();
    }
  };
  function toggleismetainfo() {
    ismetainfo = !ismetainfo;
  }

  function readtrnt(e: Event) {
    let f = (e.target as HTMLInputElement).files[0];
    if (f.size > 20971520) {
      alert('Error: Maximum Torrent File Size is 20MB');
      return;
    }
    let reader = new FileReader();
    reader.onload = (e) => {
      trntfilestring = btoa(e.target.result as string);
    };
    reader.readAsBinaryString(f);
    (e.target as HTMLInputElement).value = null;
  }
</script>

<div class="mt-10  flex items-center justify-center px-4">
  <div class="max-w-md w-full ">
    <div>
      <h2 class=" text-center text-3xl font-extrabold text-neutral-300">
        {#if ismetainfo}Enter Magnet or Infohash{:else}Select Torrent File{/if}
      </h2>
    </div>

    <div class="mt-8">
      <div class="flex bg-neutral-800 rounded-md mb-3 border border-neutral-800">
        {#if ismetainfo}
          <input
            id="torrentinput"
            type="text"
            required
            class=" bg-neutral-800 appearance-none rounded-md w-full flex-grow px-3 py-2 border-none placeholder-neutral-500 text-neutral-200  focus:outline-none sm:text-sm"
            placeholder="Magnet / Infohash"
            bind:value={torrentinput}
            on:keydown={entertoadd} />
        {:else}
          <label class="bg-neutral-800 appearance-none rounded-md w-full flex-grow px-3 py-2  placeholder-neutral-500 text-neutral-200  focus:outline-none sm:text-sm">
            <div class="text-neutral-200 flex">
              <svg xmlns="http://www.w3.org/2000/svg" class="text-neutral-400 h-6 w-6 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
              </svg>
              Select a Torrent File
            </div>
            <input accept=".torrent,application/x-bittorrent" bind:this={trntfileinput} on:change={(e) => readtrnt(e)} id="torrentfile" name="torrentfile" type="file" class="hidden" />
          </label>
        {/if}
        <button type="button" class="focus:outline-none focus:text-green-500" on:click={toggleismetainfo}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-400 my-2 mx-2 flex-grow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            {#if ismetainfo}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
            {:else}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            {/if}
          </svg>
        </button>
      </div>

      <button type="button" class="w-full my-2 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-900 hover:bg-indigo-700 focus:outline-none focus:ring-2" on:click={addfunc}> Add </button>
    </div>
  </div>
</div>

<div class="mx-auto max-w-xl">
  <div class="grid grid-flow-col grid-cols-2  mt-3">
    <div
      class="bg-neutral-800 text-neutral-200 px-5 py-5 rounded-lg m-3 cursor-pointer"
      on:click={() => {
        slocation.goto('/torrents');
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
      </svg>
      Torrents
    </div>
    <div
      class="bg-neutral-800 text-neutral-200 px-5 py-5 rounded-lg m-3 cursor-pointer"
      on:click={() => {
        slocation.goto('/settings');
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      Settings
    </div>
  </div>
  {#if $isAdmin}
    <div class="grid grid-flow-col grid-cols-2  mt-3 cursor-pointer">
      <div
        class="bg-neutral-800 text-neutral-200 px-5 py-5 rounded-lg m-3"
        on:click={() => {
          slocation.goto('/users');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        Users
      </div>
      <div
        class="bg-neutral-800 text-neutral-200 px-5 py-5 rounded-lg m-3 cursor-pointer"
        on:click={() => {
          slocation.goto('/stats');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
        </svg>
        Stats
      </div>
    </div>
  {/if}
</div>
