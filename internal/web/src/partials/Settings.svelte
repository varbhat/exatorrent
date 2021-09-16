<script lang="ts">
  import { dontstart, Send, SignOut, socket, adminmode, isAdmin, engconfig, fileSize, diskstats, nooftrackersintrackerdb } from './core';
  import slocation from 'slocation';

  let ds = false;
  if ($dontstart === 'true') {
    ds = true;
  }

  let toggledontstart = (ds: boolean) => {
    console.log('dontstart changed to', ds);
    ds ? (localStorage.setItem('dontstart', 'true'), ($dontstart = 'true')) : (localStorage.setItem('dontstart', 'false'), ($dontstart = 'false'));
  };

  $: toggledontstart(ds);

  let editmode: boolean = false;
  let newpassword: string = '';

  let diskstatsOpen = false;
  let miscOpen = false;
  let miscfirsttime = true;

  let trdelno: string;
  let srno: string;

  let changepw = () => {
    if (newpassword.length < 5) {
      alert('Size of New Password must be more than 5');
      return;
    }
    Send({
      command: 'updatepw',
      data1: newpassword
    });

    socket?.readyState === WebSocket.OPEN ? socket?.close() : console.log('socket already closed');
    SignOut();
  };

  let engsettingsOpen = false;
  let engsettingsstring = '';

  let whenengconfigChange = (ec: Object) => {
    engsettingsstring = JSON.stringify(ec, null, 2);
  };

  $: whenengconfigChange($engconfig);

  let diskusageaction = () => {
    if (diskstatsOpen === false) {
      Send({
        command: 'diskusage'
      });
      diskstatsOpen = true;
    } else {
      diskstatsOpen = false;
    }
  };

  let miscsettingsaction = () => {
    miscOpen = !miscOpen;
    if (miscfirsttime === true) {
      Send({
        command: 'nooftrackersintrackerdb',
        aop: 1
      });
    }
  };

  let enginesettingsOpen = () => {
    if (engsettingsOpen === false) {
      Send({
        command: 'getconfig',
        aop: 1
      });
      engsettingsOpen = true;
    } else {
      engsettingsOpen = false;
    }
  };
</script>

<div class="mx-auto max-w-xl">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center">
        <p class="mx-1 font-medium  truncate">
          User Settings {#if localStorage.getItem('exausertype') === 'admin'} <span class="text-xs font-semibold inline-block py-1 px-2  rounded-md text-gray-300 bg-secGray-700 ml-3 last:mr-0 mr-1">admin</span>{/if}
        </p>
      </div>
    </div>

    <div class="flex flex-col">
      <div class="flex bg-gray-800 rounded-md my-2 appearance-none border border-gray-800 w-full">
        {#if editmode === false}
          <input id="username" name="username" type="text" class=" bg-gray-800 appearance-none rounded-md w-full flex-grow px-3 py-2  border-none placeholder-gray-500 text-gray-100  focus:outline-none" placeholder={localStorage.getItem('exausername')} disabled />
          <button
            type="button"
            class="focus:outline-none"
            on:click={() => {
              editmode = true;
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400 my-2 mx-2 flex-grow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
          </button>
        {:else if editmode === true}
          <input id="password" name="password" type="text" class=" bg-gray-800 appearance-none rounded-md w-full flex-grow px-3 py-2  border-none placeholder-gray-500 text-gray-200  focus:outline-none" placeholder="Enter New Password" bind:value={newpassword} />
          <button
            type="button"
            class="focus:outline-none"
            on:click={() => {
              editmode = false;
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-red-400 my-2 mx-2 flex-grow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <button type="button" class="focus:outline-none" on:click={changepw}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-400 my-2 mx-2 flex-grow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        {/if}
      </div>
      <div class="grid grid-flow-col grid-cols-4 pr-2 bg-gray-800 rounded-md my-2 appearance-none border border-gray-800 w-full">
        <div class=" bg-gray-800  col-span-3 appearance-none rounded-md w-full flex-grow px-3 py-2  border-none  text-gray-300  focus:outline-none">Don't Start Torrents on Add</div>
        <div class="flex items-center justify-end w-full my-2 mr-2">
          <label for="dontstarttoggle" class="flex items-center cursor-pointer">
            <div class="relative">
              <input type="checkbox" class="rounded text-indigo-700 bg-gray-800 form-checkbox" bind:checked={ds} />
            </div>
          </label>
        </div>
      </div>
      {#if $isAdmin === true}
        <div class="grid grid-flow-col grid-cols-4 pr-2 bg-gray-800 my-2 appearance-none border border-gray-800 w-full rounded-md">
          <div class=" bg-gray-800  col-span-3 appearance-none  w-full flex-grow px-3 py-2  border-none  text-gray-300  focus:outline-none">Admin Mode</div>
          <div class="flex items-center justify-end w-full my-2 mr-2">
            <label for="dontstarttoggle" class="flex items-center cursor-pointer">
              <div class="relative">
                <input type="checkbox" class="rounded text-indigo-700 bg-gray-800 form-checkbox" bind:checked={$adminmode} />
              </div>
            </label>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<div class="mx-auto max-w-xl">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={diskusageaction}>
        <p class="ml-3 font-medium  truncate">Disk Usage</p>
      </div>

      {#if diskstatsOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({
              command: 'diskusage'
            });
          }}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={diskusageaction}>
        {#if diskstatsOpen === true}
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
      {#if diskstatsOpen === true}
        <div class="m-1 p-1 break-all">
          Total: {fileSize($diskstats?.total)}
        </div>
        <div class="m-1 p-1 break-all">
          Free: {fileSize($diskstats?.free)}
        </div>
        <div class="m-1 p-1 break-all">
          Used: {fileSize($diskstats?.used)} ({$diskstats?.usedPercent} %)
        </div>
      {/if}
    </div>
  </div>
</div>

{#if $isAdmin === true}
  <div class="mx-auto max-w-xl">
    <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
      <div class="flex items-center justify-between flex-wrap py-1 px-3">
        <div class="w-0 flex-1 flex items-center" on:click={miscsettingsaction}>
          <p class="ml-3 font-medium  truncate">Misc Settings</p>
        </div>

        {#if miscOpen === true}
          <button
            type="button"
            class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0"
            on:click={() => {
              Send({
                command: 'nooftrackersintrackerdb',
                aop: 1
              });
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        {/if}

        <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={miscsettingsaction}>
          {#if miscOpen === true}
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
        {#if miscOpen === true}
          <div class="m-1 p-1 break-all">
            Total Number of Trackers in TrackerDB: {$nooftrackersintrackerdb}
          </div>
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'deletetrackersintrackerdb',
                data1: 'all',
                aop: 1
              });
            }}>Delete All Trackers in TrackerDB</button
          >
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'trackerdbrefresh',
                aop: 1
              });
            }}>Refresh TrackerDB</button
          >
          <button
            class="m-2 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center"
            on:click={() => {
              Send({
                command: 'stoponseedratio',
                aop: 1
              });
            }}>Seed Ratio Check</button
          >
          <div class="flex flex-col mt-1">
            <input type="text" bind:value={trdelno} required class="m-1 bg-gray-800 appearance-none rounded-md max-w-3xl px-3 py-2 border border-blue-800 placeholder-gray-500 text-gray-200  focus:outline-none sm:text-sm mx-1" placeholder="Delete these many trackers from TrackerDB" />
            <button
              type="button"
              class="m-1 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center noHL"
              on:click={() => {
                Send({
                  command: 'deletetrackersintrackerdb',
                  data1: trdelno,
                  aop: 1
                });
              }}
            >
              Delete Trackers
            </button>
          </div>
          <div class="flex flex-col mt-3">
            <input type="text" bind:value={srno} required class="mt-3 bg-gray-800 appearance-none rounded-md max-w-3xl  px-3 py-2 border border-blue-800 placeholder-gray-500 text-gray-200  focus:outline-none sm:text-sm mx-1" placeholder="Stop Torrents on Reaching this Seedratio" />
            <button
              type="button"
              class="m-1 bg-secGray-800 py-3 max-w-3xl rounded-md justify-center noHL"
              on:click={() => {
                Send({
                  command: 'stoponseedratio',
                  data1: srno,
                  aop: 1
                });
              }}
            >
              Stop Torrents
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>

  <div class="mx-auto max-w-xl ">
    <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
      <div class="flex items-center justify-between flex-wrap py-1 px-3">
        <div class="w-0 flex-1 flex items-center" on:click={enginesettingsOpen}>
          <p class="ml-3 font-medium  truncate">Engine Settings</p>
        </div>

        {#if engsettingsOpen === true}
          <button
            class="bg-indigo-700 text-white p-2 rounded-md focus:outline-none flex-shrink-0 mx-1"
            type="button"
            on:click={() => {
              if (engsettingsstring?.length === 0) {
                alert('Empty Config!');
              } else {
                let b64config = window.btoa(engsettingsstring);
                Send({
                  command: 'updateconfig',
                  data1: b64config,
                  aop: 1
                });
              }
            }}>Update</button
          >
          <button
            type="button"
            class="flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1"
            on:click={() => {
              Send({
                command: 'getconfig',
                aop: 1
              });
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        {/if}

        <button type="button" class="-mr-1 flex p-2 rounded-md bg-gray-800 focus:outline-none flex-shrink-0 mx-1" on:click={enginesettingsOpen}>
          {#if engsettingsOpen === true}
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

      <div class="flex flex-col gap-1">
        {#if engsettingsOpen === true && $engconfig != null}
          <textarea class="form-textarea bg-gray-900 border-1 border-blue-800 max-w-3xl text-white resize-y appearance-none h-screen rounded-md" bind:value={engsettingsstring} />
        {/if}
      </div>
    </div>
  </div>
{/if}

<div class="mx-auto max-w-xl flex items-center mt-2">
  <button
    type="button"
    class="w-full my-2 mx-3 flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-black focus:outline-none"
    on:click={() => {
      slocation.goto('/about');
    }}>About exatorrent</button
  >
</div>
