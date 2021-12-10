<script lang="ts">
  import { onMount } from 'svelte';
  import slocation from 'slocation';

  onMount(() => {
    console.log('on mount');
    oldselected = selected;
    utchangeallowed = true;
  });

  import { Send } from './core';

  export let username = '';
  export let timestring = '';
  export let selected: 'user' | 'disabled' | 'admin' = 'disabled';

  let oldselected: 'user' | 'disabled' | 'admin' = 'disabled';
  let usereditmode = false;
  let newpw = '';
  let utchangeallowed = false;

  let updateuser = () => {
    if (newpw?.length > 5) {
      Send({
        command: 'updatepw',
        data1: username,
        data2: newpw,
        aop: 1
      });
    }
    if (utchangeallowed === true) {
      if (selected !== oldselected) {
        Send({
          command: 'changeusertype',
          data1: username,
          data2: selected,
          aop: 1
        });
      } else {
        alert('Same Usertype Selected');
      }
    }
  };
</script>

<div class="text-neutral-200 bg-neutral-900 px-3 py-3 rounded-md w-full my-1">
  <div class="flex flex-col justify-between flex-wrap py-1">
    <div class="flex flex-row justify-between">
      <div class="font-medium break-all text-left mx-1">{username}</div>
      <div class="flex">
        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1 noHL"
          on:click={() => {
            usereditmode = !usereditmode;
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
          </svg>
        </button>

        <button
          type="button"
          class="flex p-2 rounded-md bg-red-900 focus:outline-none flex-shrink-0 mx-1 noHL"
          on:click={() => {
            Send({
              command: 'removeuser',
              data1: username,
              aop: 1
            });
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>

        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1 noHL"
          on:click={() => {
            slocation.goto(`/user/${username}`);
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </button>
      </div>
    </div>
    <p class="text-sm  font-extralight truncate">Created at {timestring}</p>
  </div>
  {#if usereditmode === true}
    <div class="flex mt-1">
      <input id="changepw" type="text" bind:value={newpw} required class=" bg-neutral-800 appearance-none rounded-md w-full flex-grow px-3 py-2 border-none placeholder-neutral-500 text-neutral-200  focus:outline-none sm:text-sm mx-1" placeholder="Enter New Password" />
      <button
        type="button"
        class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1 noHL"
        on:click={() => {
          Send({
            command: 'revoketoken',
            data1: username,
            aop: 1
          });
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </button>
      <select class="bg-neutral-800 text-neutral-200 mx-1 rounded-md" bind:value={selected}>
        <option value="user">User</option>
        <option value="disabled">Disabled</option>
        <option value="admin">Admin</option>
      </select>
      <button class="bg-indigo-700 text-white p-3 rounded-md" on:click={updateuser}>Update</button>
    </div>
  {/if}
</div>
