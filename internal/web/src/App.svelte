<script lang="ts">
  import slocation from 'slocation';
  import { isAdmin, isDisConnected } from './partials/core';

  import Index from './partials/Index.svelte';
  import Notifications from './partials/Notifications.svelte';
  import Signin from './partials/Signin.svelte';
  import Top from './partials/Top.svelte';
  import Torrents from './partials/Torrents.svelte';
  import Disconnect from './partials/Disconnect.svelte';
  import Settings from './partials/Settings.svelte';
  import Torrent from './partials/Torrent.svelte';
  import Stats from './partials/Stats.svelte';
  import Users from './partials/Users.svelte';
  import File from './partials/File.svelte';
  import About from './partials/About.svelte';
  import User from './partials/User.svelte';
  import { Toaster } from 'svelte-sonner';
</script>

<svelte:head>
  <title>exatorrent</title>
  <link rel="icon" href="data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' stroke='mediumslateblue'><path stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4' /></svg>" />
</svelte:head>

{#if $isDisConnected === false && $slocation.pathname !== '/signin' && $slocation.pathname !== '/file'}
  <Top />
{/if}

{#if $isDisConnected === true && $slocation.pathname !== '/signin' && $slocation.pathname !== '/file'}
  <Disconnect />
{:else if $slocation.pathname === '/'}
  <Index />
{:else if $slocation.pathname === '/signin'}
  <Signin />
{:else if $slocation.pathname === '/notifications'}
  <Notifications />
{:else if $slocation.pathname === '/torrents'}
  <Torrents />
{:else if $slocation.pathname === '/settings'}
  <Settings />
{:else if $slocation.pathname === '/file'}
  <File />
{:else if $slocation.pathname === '/stats' && $isAdmin === true}
  <Stats />
{:else if $slocation.pathname === '/users' && $isAdmin === true}
  <Users />
{:else if $slocation.pathname.startsWith('/user/') && $isAdmin === true}
  <User />
{:else if $slocation.pathname.startsWith('/torrent/')}
  <Torrent />
{:else if $slocation.pathname === '/about'}
  <About />
{:else}
  <div class="mx-auto max-w-3xl">
    <p class="text-xl text-center text-red-400 font-sans">Not Found</p>
  </div>
{/if}

<Toaster position="top-right" theme="dark" richColors />
