<script lang="ts">
  import { onDestroy } from 'svelte';
  import { SignOut, Connect } from './core';

  onDestroy(() => {
    document.title = 'exatorrent';
  });
</script>

<svelte:head>
  <title>exatorrent - Disconnected</title>
</svelte:head>

<div class="mx-auto max-w-xl">
  <div class="flex justify-center p-10">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-11 w-11 text-red-600 animate-bounce" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a4.978 4.978 0 01-1.414-2.83m-1.414 5.658a9 9 0 01-2.167-9.238m7.824 2.167a1 1 0 111.414 1.414m-1.414-1.414L3 3m8.293 8.293l1.414 1.414" />
    </svg>
  </div>
  <p class="text-center text-red-400 font-sans">Disconnected</p>

  <div class="grid grid-flow-col grid-cols-2  mt-3">
    <div
      class="bg-gradient-to-r from-green-900 to-green-700 text-neutral-200 px-5 py-5 rounded-lg m-3"
      title="reconnect"
      on:click={() => {
        Connect();
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.636 18.364a9 9 0 010-12.728m12.728 0a9 9 0 010 12.728m-9.9-2.829a5 5 0 010-7.07m7.072 0a5 5 0 010 7.07M13 12a1 1 0 11-2 0 1 1 0 012 0z" />
      </svg>
      Reconnect
    </div>

    <div
      class="bg-gradient-to-r from-red-700 to-red-500 text-neutral-200 px-5 py-5 rounded-lg m-3"
      title="signout"
      on:click={() => {
        SignOut();
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
      </svg>
      Sign Out
    </div>
  </div>
</div>
