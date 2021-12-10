<script lang="ts">
  import { onMount } from 'svelte';
  import { socket, Connect } from './core';
  import slocation from 'slocation';

  function closeSocket() {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.close();
    }
  }

  onMount(() => {
    if (socket == null || socket == undefined || socket?.readyState === WebSocket.CLOSED) {
      console.log('Attempting to Connect');
      slocation.goto('/');
      Connect();
    }
  });
</script>

<nav class="max-w-xl mx-auto">
  <div class="flex justify-between h-16">
    {#if $slocation.pathname !== '/'}
      <button
        class="flex-shrink-0 focus:outline-none  bg-neutral-800 m-2 px-3 rounded-md"
        on:click={() => {
          history.length > 2 ? history.back() : slocation.goto('/');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>
    {:else}
      <button
        class="flex-shrink-0 focus:outline-none  bg-neutral-800 m-2 px-3 rounded-md"
        on:click={() => {
          slocation.goto('/notifications');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </button>
    {/if}

    <button
      class="flex items-center focus:outline-none p-5 noHL"
      title="exatorrent"
      on:click={() => {
        if ($slocation.pathname === '/') {
          slocation.goto('/about');
        } else {
          slocation.goto('/');
        }
      }}>
      <p class="font-sans text-2xl text-neutral-200">exatorrent</p>
    </button>

    <button class="flex-shrink-0 focus:outline-none  bg-neutral-800 m-2 px-3 rounded-md" on:click={closeSocket} title="Disconnect">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.636 18.364a9 9 0 010-12.728m12.728 0a9 9 0 010 12.728m-9.9-2.829a5 5 0 010-7.07m7.072 0a5 5 0 010 7.07M13 12a1 1 0 11-2 0 1 1 0 012 0z" />
      </svg>
    </button>
  </div>
</nav>
