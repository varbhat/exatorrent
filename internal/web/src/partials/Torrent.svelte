<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Send, adminmode, torrentstats, isAdmin, usersfortorrent, torctime, torrentinfo, fileSize, fsdirinfo, torrentfiles, fileviewpath, fileviewinfohash, istrntlocked, torrentinfostat } from './core';
  import { slocation } from 'slocation';
  import TorrentCard from './TorrentCard.svelte';
  import type { DlObject } from './core';
  import ProgStat from './ProgStat.svelte';

  onMount(() => {
    torrentinfo.set({} as DlObject);
    Send({ command: 'listtorrentinfo', data1: infohash });
    Send({ command: 'gettorrentinfo', data1: infohash });
    Send({ command: 'istorrentlocked', data1: infohash });
  });

  onDestroy(() => {
    Send({ command: 'stopstream' });
    document.title = 'exatorrent';
  });

  let infohash = $slocation.pathname?.split('/').reverse()[0];
  let fileProgressOpen = false;
  let browseFilesOpen = false;
  let trntStatsOpen = false;
  let trntUsersOpen = false;
  let miscOpen = false;
  let firsttimetis = true;
  let trackerfilestring = '';
  let trckrfileinput: HTMLInputElement;

  let readtracker = (e: Event) => {
    trackerfilestring = '';
    let f = (e.target as HTMLInputElement).files[0];
    if (f.size > 20971520) {
      alert('Error: Maximum Tracker File Size is 20MB');
      return;
    }
    let reader = new FileReader();
    reader.onload = (e) => {
      trackerfilestring = btoa(e.target.result as string);
      Send({
        command: 'addtrackerstotorrent',
        data1: infohash,
        data2: trackerfilestring,
        ...($adminmode === true && { aop: 1 })
      });
    };
    reader.readAsBinaryString(f);
    (e.target as HTMLInputElement).value = null;
  };

  let fileProgressaction = () => {
    if (fileProgressOpen === false) {
      Send({
        command: 'gettorrentfiles',
        data1: infohash
      });
      fileProgressOpen = true;
    } else {
      fileProgressOpen = false;
    }
  };

  let browseFilesaction = () => {
    if (browseFilesOpen === false) {
      Send({
        command: 'getfsdirinfo',
        data1: infohash
      });
      browseFilesOpen = true;
    } else {
      browseFilesOpen = false;
    }
  };

  let trntStatsaction = () => {
    if (trntStatsOpen === false) {
      Send({
        command: 'gettorrentstats',
        data1: infohash
      });
      trntStatsOpen = true;
    } else {
      trntStatsOpen = false;
    }
  };

  let uwottaction = () => {
    if (trntUsersOpen === false) {
      Send({
        command: 'listusersfortorrent',
        data1: infohash,
        aop: 1
      });
      trntUsersOpen = true;
    } else {
      trntUsersOpen = false;
    }
  };

  let miscaction = () => {
    miscOpen = !miscOpen;
    if (firsttimetis === true) {
      Send({
        command: 'gettorrentinfostat',
        data1: infohash
      });
      firsttimetis = false;
    }
  };
</script>

<svelte:head>
  <title>{$torrentinfo?.name} ({$torrentinfo?.infohash})</title>
</svelte:head>

<div class="mx-auto max-w-3xl ">
  <TorrentCard state={$torrentinfo?.state} name={$torrentinfo?.name} infohash={$torrentinfo?.infohash} bytescompleted={$torrentinfo?.bytescompleted} bytesmissing={$torrentinfo?.bytesmissing} length={$torrentinfo?.length} seeding={$torrentinfo?.seeding} locked={$istrntlocked} isTorrentPage={true} />

  <div class="bg-black grid grid-flow-col text-white rounded-lg m-3 p-2">
    {#if $adminmode === false}
      <div class=" text-center">
        <button
          type="button"
          on:click={() => {
            Send({
              command: 'abandontorrent',
              data1: infohash
            });
            slocation.goto('/');
          }}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" />
          </svg>
          <p class="text-gray-500 text-sm">Abandon</p>
        </button>
      </div>
    {/if}

    <div class=" text-center">
      <a href="/api/torrent/{infohash}/?dl=tar" target="_blank" rel="noopener noreferrer" download>
        <button type="button">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
          <p class="text-gray-500 text-sm">Tar</p>
        </button>
      </a>
    </div>

    <div class=" text-center">
      <a href="/api/torrent/{infohash}/?dl=zip" target="_blank" rel="noopener noreferrer" download>
        <button type="button">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" viewBox="0 0 20 20" fill="currentColor">
            <path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
            <path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
          </svg>
          <p class="text-gray-500 text-sm">Zip</p>
        </button>
      </a>
    </div>
  </div>

  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive'}
    <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus:appearance-none active:appearance-none focus-within:appearance-none focus:bg-black noHL">
      <div class="flex items-center justify-between flex-wrap py-1 px-3">
        <div class="w-0 flex-1 flex items-center" on:click={fileProgressaction}>
          <p class="ml-3 font-medium  truncate">File Progress</p>
        </div>

        {#if fileProgressOpen === true}
          <div
            class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 noHL"
            on:click={() => {
              Send({
                command: 'gettorrentfiles',
                data1: infohash
              });
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
        {/if}

        <div class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1 noHL" on:click={fileProgressaction}>
          {#if fileProgressOpen === true}
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            </svg>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          {/if}
        </div>
      </div>

      <div class="flex flex-col">
        {#if fileProgressOpen === true}
          {#each $torrentfiles as file (file.path)}
            <div class="text-gray-200 bg-secGray-900 px-3 py-3 rounded-md w-full my-1">
              <div class="flex items-center justify-between flex-wrap py-1">
                <div
                  class="w-0 flex-1 flex"
                  on:click={() => {
                    fileviewpath.set(file?.path);
                    fileviewinfohash.set(infohash);
                    slocation.goto('/file');
                  }}
                >
                  <p class="font-medium  truncate">
                    {file?.displaypath}
                  </p>
                </div>

                {#if file?.priority === 1}
                  <div
                    class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1"
                    on:click={() => {
                      Send({
                        command: 'stopfile',
                        data1: infohash,
                        data2: file?.path,
                        ...($adminmode === true && {
                          aop: 1
                        })
                      });
                      setTimeout(() => {
                        Send({
                          command: 'gettorrentfiles',
                          data1: infohash
                        });
                      }, 1000);
                    }}
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                {:else if file?.priority === 0}
                  <div
                    class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1"
                    on:click={() => {
                      Send({
                        command: 'startfile',
                        data1: infohash,
                        data2: file?.path,
                        ...($adminmode === true && {
                          aop: 1
                        })
                      });
                      setTimeout(() => {
                        Send({
                          command: 'gettorrentfiles',
                          data1: infohash
                        });
                      }, 1000);
                    }}
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                {/if}
              </div>

              <ProgStat bytescompleted={file?.bytescompleted} length={file?.length} offset={file?.offset} />
            </div>
          {/each}
        {/if}
      </div>
    </div>
  {/if}

  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={browseFilesaction}>
        <p class="ml-3 font-medium  truncate">Browse Files</p>
      </div>

      {#if browseFilesOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({
              command: 'getfsdirinfo',
              data1: infohash
            });
          }}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={browseFilesaction}>
        {#if browseFilesOpen === true}
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
      {#if browseFilesOpen === true}
        {#each $fsdirinfo as file (file.path)}
          <div class="text-gray-200 bg-gray-900 px-3 py-3 rounded-md w-full my-1">
            <div class="flex flex-col  justify-between flex-wrap py-1">
              <div
                class="w-full"
                on:click={() => {
                  if (file?.isdir === true) {
                    Send({
                      command: 'getfsdirinfo',
                      data1: infohash,
                      data2: file?.path
                    });
                  } else if (file?.isdir === false) {
                    fileviewpath.set(file?.path);
                    fileviewinfohash.set(infohash);
                    slocation.goto('/file');
                  }
                }}
              >
                <p class="font-medium break-all text-left">
                  {file?.name}
                </p>
              </div>

              <div class="grid grid-flow-col">
                <p class="py-2 font-extralight break-all">
                  {#if file?.isdir === false}{fileSize(file?.size)}{:else}Directory{/if}
                </p>

                <div class="flex flex-row justify-end gap-1 flex-wrap py-1">
                  {#if file?.isdir === true}
                    <a href="/api/torrent/{infohash}/{file?.path}/?dl=zip" target="_blank" rel="noopener noreferrer" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1 noHL" download>
                      <button type="button">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" viewBox="0 0 20 20" fill="currentColor">
                          <path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
                          <path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </a>

                    <a href="/api/torrent/{infohash}/{file?.path}/?dl=tar" target="_blank" rel="noopener noreferrer" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1 noHL" download>
                      <button type="button">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                        </svg>
                      </button>
                    </a>
                  {/if}

                  {#if $torrentinfo?.state === 'removed'}
                    <button
                      type="button"
                      class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1 noHL"
                      on:click={() => {
                        Send({
                          command: 'deletefilepath',
                          data1: infohash,
                          data2: file?.path,
                          ...($adminmode === true && {
                            aop: 1
                          })
                        });
                      }}
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'loading'}
    <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
      <div class="flex items-center justify-between flex-wrap py-1 px-3">
        <div class="w-0 flex-1 flex items-center" on:click={trntStatsaction}>
          <p class="ml-3 font-medium  truncate">Stats</p>
        </div>

        {#if trntStatsOpen === true}
          <button
            type="button"
            class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0"
            on:click={() => {
              Send({
                command: 'gettorrentstats',
                data1: infohash
              });
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        {/if}

        <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={trntStatsaction}>
          {#if trntStatsOpen === true}
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
        {#if trntStatsOpen === true}
          <div class="m-1 p-1 break-all">
            Total Peers: {$torrentstats?.TotalPeers}
          </div>
          <div class="m-1 p-1 break-all">
            Pending Peers: {$torrentstats?.PendingPeers}
          </div>
          <div class="m-1 p-1 break-all">
            Active Peers: {$torrentstats?.ActivePeers}
          </div>
          <div class="m-1 p-1 break-all">
            Connected Seeders: {$torrentstats?.ConnectedSeeders}
          </div>
          <div class="m-1 p-1 break-all">
            Half Open Peers: {$torrentstats?.HalfOpenPeers}
          </div>
          <div class="m-1 p-1 break-all">
            Written: {fileSize($torrentstats?.BytesWritten)}
          </div>
          <div class="m-1 p-1 break-all">
            WrittenData: {fileSize($torrentstats?.BytesWrittenData)}
          </div>
          <div class="m-1 p-1 break-all">
            Read: {fileSize($torrentstats?.BytesRead)}
          </div>
          <div class="m-1 p-1 break-all">
            ReadData: {fileSize($torrentstats?.BytesReadData)}
          </div>
          <div class="m-1 p-1 break-all">
            ReadUsefulData: {fileSize($torrentstats?.BytesReadUsefulData)}
          </div>
          <div class="m-1 p-1 break-all">
            Seed Ratio: {($torrentstats?.BytesWrittenData / $torrentstats.BytesReadData).toLocaleString('en-US', {
              maximumFractionDigits: 5,
              minimumFractionDigits: 5
            })}
          </div>

          <div class="m-1 p-1 break-all">
            Chunks Written: {$torrentstats?.ChunksWritten}
          </div>
          <div class="m-1 p-1 break-all">
            Chunks Read: {$torrentstats?.ChunksRead}
          </div>
          <div class="m-1 p-1 break-all">
            Chunks Read Useful: {$torrentstats?.ChunksReadUseful}
          </div>
          <div class="m-1 p-1 break-all">
            Chunks Read Wasted: {$torrentstats?.ChunksReadWasted}
          </div>
          <div class="m-1 p-1 break-all">
            Metadata Chunks Read: {$torrentstats?.MetadataChunksRead}
          </div>
          <div class="m-1 p-1 break-all">
            Piece Dirtied Good: {$torrentstats?.PiecesDirtiedGood}
          </div>
          <div class="m-1 p-1 break-all">
            Piece Dirtied Bad: {$torrentstats?.PiecesDirtiedBad}
          </div>
        {/if}
      </div>
    </div>
  {/if}

  {#if $isAdmin === true}
    <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
      <div class="flex items-center justify-between flex-wrap py-1 px-3">
        <div class="w-0 flex-1 flex items-center" on:click={uwottaction}>
          <p class="ml-3 font-medium  truncate">Users who Own this Torrent</p>
        </div>

        {#if trntUsersOpen === true}
          <button
            type="button"
            class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0"
            on:click={() => {
              Send({
                command: 'listusersfortorrent',
                data1: infohash,
                aop: 1
              });
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        {/if}

        <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={uwottaction}>
          {#if trntUsersOpen === true}
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
        {#if trntUsersOpen === true}
          {#each $usersfortorrent as eachuser (eachuser)}
            <div class="text-gray-200 bg-secGray-900 px-3 py-3 rounded-md w-full my-1">
              <div class="flex items-center justify-between flex-wrap py-1">
                <div class="w-0 flex-1 flex">
                  <p class="font-medium  break-all">
                    {eachuser}
                  </p>
                </div>

                <div
                  class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1"
                  on:click={() => {
                    Send({
                      command: 'abandontorrent',
                      data1: infohash,
                      data2: eachuser,
                      aop: 1
                    });
                  }}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  {/if}

  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus:appearance-none active:appearance-none focus-within:appearance-none focus:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={miscaction}>
        <p class="ml-3 font-medium  truncate">Misc</p>
      </div>

      {#if miscOpen === true}
        <div
          class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 noHL"
          on:click={() => {
            Send({
              command: 'gettorrentinfostat',
              data1: infohash
            });
          }}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </div>
      {/if}

      <div class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1 noHL" on:click={miscaction}>
        {#if miscOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </div>
    </div>

    <div class="flex flex-col">
      {#if miscOpen === true}
        {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'loading'}
          <div class="m-1 p-1 break-all text-center">
            Added at {$torctime.addedat}
          </div>
          <div class="m-1 p-1 break-all text-center">
            {#if $torrentinfo?.state === 'active'}Started At {$torctime.startedat}{/if}
          </div>
          <label class="appearance-none">
            <div class="text-gray-200 flex m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="text-gray-400 h-6 w-6 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
              </svg>
              Add Trackers (List)
            </div>
            <input accept=".txt" bind:this={trckrfileinput} on:change={(e) => readtracker(e)} id="torrentfile" name="torrentfile" type="file" class="hidden" />
          </label>
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'gettorrentmetainfo',
                data1: infohash
              });
            }}>Download Torrent File</button
          >
        {/if}
        {#if $isAdmin === true && $adminmode === true}
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'changedataload',
                data1: infohash,
                data2: 'upload',
                data3: 'allow',
                aop: 1
              });
            }}>Allow Data Upload</button
          >
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'changedataload',
                data1: infohash,
                data2: 'upload',
                data3: 'disallow',
                aop: 1
              });
            }}>Disallow Data Upload</button
          >
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'changedataload',
                data1: infohash,
                data2: 'download',
                data3: 'allow',
                aop: 1
              });
            }}>Allow Data Download</button
          >
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'changedataload',
                data1: infohash,
                data2: 'download',
                data3: 'disallow',
                aop: 1
              });
            }}>Disallow Data Download</button
          >
        {/if}
      {/if}
    </div>
  </div>
</div>
